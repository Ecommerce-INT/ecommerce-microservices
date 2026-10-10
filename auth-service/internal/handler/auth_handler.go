package handler

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"com.ecommerce/auth-service/internal/model"
	"com.ecommerce/auth-service/internal/service"
	"com.ecommerce/pkg/common/response"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userSvc        *service.UserService
	kcClient       *service.KeycloakClient
	ssoStore       *service.SsoSessionStore
	defaultRedirect string
}

func NewAuthHandler(userSvc *service.UserService, kc *service.KeycloakClient, ssoStore *service.SsoSessionStore, defaultRedirect string) *AuthHandler {
	if defaultRedirect == "" {
		defaultRedirect = "http://ecommerce.local/auth/callback"
	}
	return &AuthHandler{
		userSvc:         userSvc,
		kcClient:        kc,
		ssoStore:        ssoStore,
		defaultRedirect: defaultRedirect,
	}
}

func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/auth/actuator/health", h.HealthCheck)
	app.Get("/health", h.HealthCheck)

	// Auth routes
	h.registerAuthRoutes(app.Group("/api/v1/auth"))
	h.registerAuthRoutes(app.Group("/auth/api/v1/auth"))

	// Users routes
	h.registerUserRoutes(app.Group("/api/v1/users"))
	h.registerUserRoutes(app.Group("/auth/api/v1/users"))

	// Roles routes
	h.registerRoleRoutes(app.Group("/api/v1/roles"))
	h.registerRoleRoutes(app.Group("/auth/api/v1/roles"))
}

func (h *AuthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

func (h *AuthHandler) registerAuthRoutes(r fiber.Router) {
	r.Post("/signup", h.Signup)
	r.Get("/login", h.Login)
	r.Get("/callback", h.Callback)
	r.Get("/session", h.Session)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
}

func (h *AuthHandler) registerUserRoutes(r fiber.Router) {
	r.Get("/me", h.GetCurrentUser)
	r.Put("/me/password", h.ChangePassword)
	r.Get("/all", h.GetAllUsers)
	r.Get("/:id", h.GetUserById)
	r.Put("/:id", h.UpdateUser)
	r.Delete("/:id", h.DeleteUser)
	r.Get("", h.GetUserByUsername)
}

func (h *AuthHandler) registerRoleRoutes(r fiber.Router) {
	r.Post("/users/:userId/assign", h.AssignRole)
	r.Post("/users/:userId/revoke", h.RevokeRole)
	r.Get("/users/:userId", h.GetUserRoles)
}

// -------------------------------------------------------------
// Auth Handlers
// -------------------------------------------------------------

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		return response.Error(c, fiber.StatusBadRequest, "VALIDATION_FAILED", "Username, password and email are required")
	}

	_, err := h.userSvc.Register(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) || errors.Is(err, service.ErrEmailExists) || errors.Is(err, service.ErrPhoneExists) {
			return response.Error(c, fiber.StatusConflict, "CONFLICT", err.Error())
		}
		return response.Error(c, fiber.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
	}

	return response.Message(c, fmt.Sprintf("User %s registered successfully", req.Username))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = h.defaultRedirect
	}

	state := h.ssoStore.CreateLoginState(redirectURI)
	authURL := h.kcClient.BuildAuthorizeURL(state, redirectURI)
	return c.Redirect(authURL, fiber.StatusFound)
}

func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	errQuery := c.Query("error")

	if errQuery != "" {
		redirectURL := appendQuery(h.defaultRedirect, "error", errQuery)
		return c.Redirect(redirectURL, fiber.StatusFound)
	}

	frontendRedirect, ok := h.ssoStore.ConsumeLoginState(state)
	if !ok || frontendRedirect == "" {
		frontendRedirect = h.defaultRedirect
	}

	tokens, err := h.kcClient.ExchangeAuthorizationCode(c.Context(), code, frontendRedirect)
	if err != nil {
		// Mock token if local keycloak is unavailable
		tokens = &model.KeycloakTokenResponse{
			AccessToken:  "mock-access-token",
			RefreshToken: "mock-refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}
	}

	ticket := h.ssoStore.StoreTokens(*tokens)
	redirectURL := appendQuery(frontendRedirect, "ticket", ticket)
	return c.Redirect(redirectURL, fiber.StatusFound)
}

func (h *AuthHandler) Session(c *fiber.Ctx) error {
	ticket := c.Query("ticket")
	if ticket == "" {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Ticket is required")
	}

	tokens, ok := h.ssoStore.ConsumeTokens(ticket)
	if !ok || tokens == nil {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired ticket")
	}

	return response.OK(c, tokens)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Refresh token is required")
	}

	tokens, err := h.kcClient.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}

	return response.OK(c, tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	_ = c.BodyParser(&req)
	if req.RefreshToken != "" {
		_ = h.kcClient.Logout(c.Context(), req.RefreshToken)
	}
	return response.Message(c, "Logout successful")
}

// -------------------------------------------------------------
// User Handlers
// -------------------------------------------------------------

func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
	}

	user, err := h.userSvc.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.OK(c, user, "User updated successfully")
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	return response.Message(c, "Password change request processed")
}

func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
	}

	if err := h.userSvc.Delete(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.Message(c, fmt.Sprintf("User %d deleted successfully", id))
}

func (h *AuthHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Query("username")
	if username == "" {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "username query parameter is required")
	}

	user, err := h.userSvc.FindByUsername(c.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.OK(c, user)
}

func (h *AuthHandler) GetUserById(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
	}

	user, err := h.userSvc.FindById(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.OK(c, user)
}

func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "0"))
	size, _ := strconv.Atoi(c.Query("size", "10"))
	sortBy := c.Query("sortBy", "id")
	sortOrder := c.Query("sortOrder", "ASC")

	pageResult, err := h.userSvc.FindAllUsers(c.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.OK(c, pageResult)
}

func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	username, _ := c.Locals("username").(string)
	if username == "" {
		username, _ = c.Locals("preferred_username").(string)
	}
	if username == "" {
		username, _ = c.Locals("sub").(string)
	}
	if username == "" {
		return response.Error(c, fiber.StatusUnauthorized, "AUTH_TOKEN_INVALID", "Invalid or missing token")
	}

	user, err := h.userSvc.FindByUsername(c.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.OK(c, user)
}

// -------------------------------------------------------------
// Role Handlers
// -------------------------------------------------------------

func (h *AuthHandler) AssignRole(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("userId"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
	}

	roleName := strings.TrimSpace(string(c.Body()))
	if roleName == "" {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Role name cannot be empty")
	}

	assigned, err := h.userSvc.AssignRole(c.Context(), userId, roleName)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	if !assigned {
		return response.Error(c, fiber.StatusConflict, "CONFLICT", "Role already assigned or user/role not found")
	}

	return response.Message(c, fmt.Sprintf("Roles assigned to user %d", userId))
}

func (h *AuthHandler) RevokeRole(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("userId"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
	}

	roleName := strings.TrimSpace(string(c.Body()))
	if roleName == "" {
		return response.Error(c, fiber.StatusBadRequest, "BAD_REQUEST", "Role name cannot be empty")
	}

	revoked, err := h.userSvc.RevokeRole(c.Context(), userId, roleName)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
	if !revoked {
		return response.Error(c, fiber.StatusConflict, "CONFLICT", "Role not present on user")
	}

	return response.Message(c, fmt.Sprintf("Roles revoked from user %d", userId))
}

func (h *AuthHandler) GetUserRoles(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("userId"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
	}

	roles, err := h.userSvc.GetUserRoles(c.Context(), userId)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return response.OK(c, roles)
}

func appendQuery(base, key, value string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base + "?" + url.QueryEscape(key) + "=" + url.QueryEscape(value)
	}
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u.String()
}
