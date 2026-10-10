package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/auth-service/internal/handler"
	"com.ecommerce/auth-service/internal/model"
	"com.ecommerce/auth-service/internal/repository"
	"com.ecommerce/auth-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() http.Handler {
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

	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func performRequest(app http.Handler, method, target string, body []byte) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func TestHealthCheck(t *testing.T) {
	app := setupTestApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/actuator/health", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/auth/actuator/health", nil).Code)
}

func TestSSOFlow(t *testing.T) {
	app := setupTestApp()

	// 1. Login redirects to Keycloak
	respLogin := performRequest(app, "GET", "/api/v1/auth/login", nil)
	assert.Equal(t, 302, respLogin.Code)
	loc := respLogin.Header().Get("Location")
	assert.Contains(t, loc, "protocol/openid-connect/auth")
	assert.Contains(t, loc, "response_type=code")

	// 2. Callback handles code and redirects with ticket
	respCallback := performRequest(app, "GET", "/api/v1/auth/callback?code=mock-code&state=mock-state", nil)
	assert.Equal(t, 302, respCallback.Code)
	locCallback := respCallback.Header().Get("Location")
	assert.Contains(t, locCallback, "ticket=")

	// 3. User change password endpoint
	respPwd := performRequest(app, "PUT", "/api/v1/users/me/password", []byte("{}"))
	assert.Equal(t, 200, respPwd.Code)
}

func TestSignupValidation(t *testing.T) {
	app := setupTestApp()

	// Empty body should return validation error
	resp := performRequest(app, "POST", "/api/v1/auth/signup", []byte("{}"))
	assert.Equal(t, 400, resp.Code)

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
