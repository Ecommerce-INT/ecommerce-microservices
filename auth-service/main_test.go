package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"com.ecommerce/auth-service/internal/handler"
	"com.ecommerce/auth-service/internal/model"
	"com.ecommerce/auth-service/internal/repository"
	"com.ecommerce/auth-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	repo := repository.NewUserRepository(nil)
	kcCfg := service.KeycloakConfig{
		ServerURL:               "http://localhost:8080",
		Realm:                   "ecommerce",
		ClientID:                "ecommerce-client",
		BackendCallbackURL:      "http://api.ecommerce.local/api/v1/auth/callback",
		DefaultFrontendRedirect: "http://ecommerce.local/auth/callback",
	}
	kcClient := service.NewKeycloakClient(kcCfg)
	ssoStore := service.NewSsoSessionStore()
	userSvc := service.NewUserService(repo, kcClient)
	h := handler.NewAuthHandler(userSvc, kcClient, ssoStore, kcCfg.DefaultFrontendRedirect)

	app := fiber.New()
	h.RegisterRoutes(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqAuth := httptest.NewRequest("GET", "/auth/actuator/health", nil)
	respAuth, err := app.Test(reqAuth)
	assert.NoError(t, err)
	assert.Equal(t, 200, respAuth.StatusCode)
}

func TestSSOFlow(t *testing.T) {
	app := setupTestApp()

	// 1. Login redirects to Keycloak
	reqLogin := httptest.NewRequest("GET", "/api/v1/auth/login", nil)
	respLogin, err := app.Test(reqLogin)
	assert.NoError(t, err)
	assert.Equal(t, 302, respLogin.StatusCode)
	loc := respLogin.Header.Get("Location")
	assert.Contains(t, loc, "protocol/openid-connect/auth")
	assert.Contains(t, loc, "response_type=code")

	// 2. Callback handles code and redirects with ticket
	reqCallback := httptest.NewRequest("GET", "/api/v1/auth/callback?code=mock-code&state=mock-state", nil)
	respCallback, err := app.Test(reqCallback)
	assert.NoError(t, err)
	assert.Equal(t, 302, respCallback.StatusCode)
	locCallback := respCallback.Header.Get("Location")
	assert.Contains(t, locCallback, "ticket=")

	// 3. User change password endpoint
	reqPwd := httptest.NewRequest("PUT", "/api/v1/users/me/password", bytes.NewReader([]byte("{}")))
	reqPwd.Header.Set("Content-Type", "application/json")
	respPwd, err := app.Test(reqPwd)
	assert.NoError(t, err)
	assert.Equal(t, 200, respPwd.StatusCode)
}

func TestSignupValidation(t *testing.T) {
	app := setupTestApp()

	// Empty body should return validation error
	req := httptest.NewRequest("POST", "/api/v1/auth/signup", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var res map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&res)
	assert.Equal(t, false, res["success"])
}

func TestSSOSessionStore(t *testing.T) {
	store := service.NewSsoSessionStore()

	state := store.CreateLoginState("http://example.com/redirect")
	assert.NotEmpty(t, state)

	redirect, ok := store.ConsumeLoginState(state)
	assert.True(t, ok)
	assert.Equal(t, "http://example.com/redirect", redirect)

	// Consumed once -> second time must fail
	_, ok2 := store.ConsumeLoginState(state)
	assert.False(t, ok2)

	tok := model.KeycloakTokenResponse{
		AccessToken: "access-123",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}
	ticket := store.StoreTokens(tok)
	assert.NotEmpty(t, ticket)

	consumedTok, ok := store.ConsumeTokens(ticket)
	assert.True(t, ok)
	assert.Equal(t, "access-123", consumedTok.AccessToken)
}
