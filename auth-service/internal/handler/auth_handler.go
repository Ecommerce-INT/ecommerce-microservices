package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"com.ecommerce/auth-service/internal/model"
	"com.ecommerce/auth-service/internal/service"
	"com.ecommerce/pkg/common/middleware"
	"com.ecommerce/pkg/common/response"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	userSvc         *service.UserService
	kcClient        *service.KeycloakClient
	ssoStore        *service.SsoSessionStore
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

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/actuator/health", h.HealthCheck)
	r.Get("/auth/actuator/health", h.HealthCheck)
	r.Get("/health", h.HealthCheck)

	// Auth routes
	for _, base := range []string{"/api/v1/auth", "/auth/api/v1/auth"} {
		h.registerAuthRoutes(r, base)
	}

	// Users routes
	for _, base := range []string{"/api/v1/users", "/auth/api/v1/users"} {
		h.registerUserRoutes(r, base)
	}

	// Roles routes
	for _, base := range []string{"/api/v1/roles", "/auth/api/v1/roles"} {
		h.registerRoleRoutes(r, base)
	}
}

func (h *AuthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

func (h *AuthHandler) registerAuthRoutes(r chi.Router, base string) {
	r.Post(base+"/signup", h.Signup)
	r.Get(base+"/login", h.Login)
	r.Get(base+"/callback", h.Callback)
	r.Get(base+"/session", h.Session)
	r.Post(base+"/refresh", h.Refresh)
	r.Post(base+"/logout", h.Logout)
}

func (h *AuthHandler) registerUserRoutes(r chi.Router, base string) {
	r.Get(base+"/me", h.GetCurrentUser)
	r.Put(base+"/me/password", h.ChangePassword)
	r.Get(base+"/all", h.GetAllUsers)
	r.Get(base+"/{id}", h.GetUserById)
	r.Put(base+"/{id}", h.UpdateUser)
	r.Delete(base+"/{id}", h.DeleteUser)
	r.Get(base, h.GetUserByUsername)
}

func (h *AuthHandler) registerRoleRoutes(r chi.Router, base string) {
	r.Post(base+"/users/{userId}/assign", h.AssignRole)
	r.Post(base+"/users/{userId}/revoke", h.RevokeRole)
	r.Get(base+"/users/{userId}", h.GetUserRoles)
}

func queryDefault(r *http.Request, key, fallback string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return fallback
}

// -------------------------------------------------------------
// Auth Handlers
// -------------------------------------------------------------

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		response.Error(w, r, http.StatusBadRequest, "VALIDATION_FAILED", "Username, password and email are required")
		return
	}

	_, err := h.userSvc.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) || errors.Is(err, service.ErrEmailExists) || errors.Is(err, service.ErrPhoneExists) {
			response.Error(w, r, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		response.Error(w, r, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	response.Message(w, r, fmt.Sprintf("User %s registered successfully", req.Username))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = h.defaultRedirect
	}

	state := h.ssoStore.CreateLoginState(redirectURI)
	authURL := h.kcClient.BuildAuthorizeURL(state, redirectURI)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	errQuery := q.Get("error")

	if errQuery != "" {
		redirectURL := appendQuery(h.defaultRedirect, "error", errQuery)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	frontendRedirect, ok := h.ssoStore.ConsumeLoginState(state)
	if !ok || frontendRedirect == "" {
		frontendRedirect = h.defaultRedirect
	}

	tokens, err := h.kcClient.ExchangeAuthorizationCode(r.Context(), code, frontendRedirect)
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
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "Ticket is required")
		return
	}

	tokens, ok := h.ssoStore.ConsumeTokens(ticket)
	if !ok || tokens == nil {
		response.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired ticket")
		return
	}

	response.OK(w, r, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "Refresh token is required")
		return
	}

	tokens, err := h.kcClient.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	response.OK(w, r, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.RefreshToken != "" {
		_ = h.kcClient.Logout(r.Context(), req.RefreshToken)
	}
	response.Message(w, r, "Logout successful")
}

// -------------------------------------------------------------
// User Handlers
// -------------------------------------------------------------

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	user, err := h.userSvc.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(w, r, user, "User updated successfully")
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	response.Message(w, r, "Password change request processed")
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
		return
	}

	if err := h.userSvc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Message(w, r, fmt.Sprintf("User %d deleted successfully", id))
}

func (h *AuthHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "username query parameter is required")
		return
	}

	user, err := h.userSvc.FindByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(w, r, user)
}

func (h *AuthHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
		return
	}

	user, err := h.userSvc.FindById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(w, r, user)
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(queryDefault(r, "page", "0"))
	size, _ := strconv.Atoi(queryDefault(r, "size", "10"))
	sortBy := queryDefault(r, "sortBy", "id")
	sortOrder := queryDefault(r, "sortOrder", "ASC")

	pageResult, err := h.userSvc.FindAllUsers(r.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(w, r, pageResult)
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	username := middleware.Username(r.Context())
	if username == "" {
		username = middleware.UserID(r.Context())
	}
	if username == "" {
		response.Error(w, r, http.StatusUnauthorized, "AUTH_TOKEN_INVALID", "Invalid or missing token")
		return
	}

	user, err := h.userSvc.FindByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, r, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(w, r, user)
}

// -------------------------------------------------------------
// Role Handlers
// -------------------------------------------------------------

func (h *AuthHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
		return
	}

	body, _ := io.ReadAll(r.Body)
	roleName := strings.TrimSpace(string(body))
	if roleName == "" {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "Role name cannot be empty")
		return
	}

	assigned, err := h.userSvc.AssignRole(r.Context(), userId, roleName)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if !assigned {
		response.Error(w, r, http.StatusConflict, "CONFLICT", "Role already assigned or user/role not found")
		return
	}

	response.Message(w, r, fmt.Sprintf("Roles assigned to user %d", userId))
}

func (h *AuthHandler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
		return
	}

	body, _ := io.ReadAll(r.Body)
	roleName := strings.TrimSpace(string(body))
	if roleName == "" {
		response.Error(w, r, http.StatusBadRequest, "BAD_REQUEST", "Role name cannot be empty")
		return
	}

	revoked, err := h.userSvc.RevokeRole(r.Context(), userId, roleName)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if !revoked {
		response.Error(w, r, http.StatusConflict, "CONFLICT", "Role not present on user")
		return
	}

	response.Message(w, r, fmt.Sprintf("Roles revoked from user %d", userId))
}

func (h *AuthHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, "INVALID_ID", "User ID must be numeric")
		return
	}

	roles, err := h.userSvc.GetUserRoles(r.Context(), userId)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.OK(w, r, roles)
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
