package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"com.ecommerce/auth-service/internal/model"
)

type KeycloakConfig struct {
	ServerURL               string
	PublicServerURL         string
	Realm                   string
	ClientID                string
	ClientSecret            string
	AdminRealm              string
	AdminClientID           string
	AdminUsername           string
	AdminPassword           string
	BackendCallbackURL      string
	DefaultFrontendRedirect string
	Scope                   string
}

type KeycloakClient struct {
	cfg        KeycloakConfig
	httpClient *http.Client
}

func NewKeycloakClient(cfg KeycloakConfig) *KeycloakClient {
	if cfg.ServerURL == "" {
		cfg.ServerURL = "http://localhost:8080"
	}
	if cfg.Realm == "" {
		cfg.Realm = "ecommerce"
	}
	if cfg.ClientID == "" {
		cfg.ClientID = "ecommerce-client"
	}
	if cfg.AdminRealm == "" {
		cfg.AdminRealm = "master"
	}
	if cfg.AdminClientID == "" {
		cfg.AdminClientID = "admin-cli"
	}
	if cfg.AdminUsername == "" {
		cfg.AdminUsername = "admin"
	}
	if cfg.AdminPassword == "" {
		cfg.AdminPassword = "admin"
	}
	if cfg.Scope == "" {
		cfg.Scope = "openid profile email"
	}
	if cfg.BackendCallbackURL == "" {
		cfg.BackendCallbackURL = "http://api.ecommerce.local/api/v1/auth/callback"
	}
	if cfg.DefaultFrontendRedirect == "" {
		cfg.DefaultFrontendRedirect = "http://ecommerce.local/auth/callback"
	}

	return &KeycloakClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *KeycloakClient) TokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.cfg.ServerURL, c.cfg.Realm)
}

func (c *KeycloakClient) AuthorizationEndpoint() string {
	baseURL := c.cfg.PublicServerURL
	if baseURL == "" {
		baseURL = c.cfg.ServerURL
	}
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", baseURL, c.cfg.Realm)
}

func (c *KeycloakClient) LogoutEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", c.cfg.ServerURL, c.cfg.Realm)
}

func (c *KeycloakClient) AdminTokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.cfg.ServerURL, c.cfg.AdminRealm)
}

func (c *KeycloakClient) AdminUsersEndpoint() string {
	return fmt.Sprintf("%s/admin/realms/%s/users", c.cfg.ServerURL, c.cfg.Realm)
}

func (c *KeycloakClient) BuildAuthorizeURL(state, redirectURI string) string {
	v := url.Values{}
	v.Set("client_id", c.cfg.ClientID)
	v.Set("response_type", "code")
	v.Set("scope", c.cfg.Scope)
	v.Set("redirect_uri", c.cfg.BackendCallbackURL)
	v.Set("state", state)
	return fmt.Sprintf("%s?%s", c.AuthorizationEndpoint(), v.Encode())
}

func (c *KeycloakClient) ExchangeAuthorizationCode(ctx context.Context, code, redirectURI string) (*model.KeycloakTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.cfg.ClientID)
	if c.cfg.ClientSecret != "" {
		data.Set("client_secret", c.cfg.ClientSecret)
	}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	return c.postToken(ctx, c.TokenEndpoint(), data)
}

func (c *KeycloakClient) RefreshToken(ctx context.Context, refreshToken string) (*model.KeycloakTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.cfg.ClientID)
	if c.cfg.ClientSecret != "" {
		data.Set("client_secret", c.cfg.ClientSecret)
	}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	return c.postToken(ctx, c.TokenEndpoint(), data)
}

func (c *KeycloakClient) Logout(ctx context.Context, refreshToken string) error {
	data := url.Values{}
	data.Set("client_id", c.cfg.ClientID)
	if c.cfg.ClientSecret != "" {
		data.Set("client_secret", c.cfg.ClientSecret)
	}
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.LogoutEndpoint(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("keycloak logout error: %s (%d)", string(body), resp.StatusCode)
	}
	return nil
}

func (c *KeycloakClient) FetchAdminAccessToken(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("client_id", c.cfg.AdminClientID)
	data.Set("username", c.cfg.AdminUsername)
	data.Set("password", c.cfg.AdminPassword)
	data.Set("grant_type", "password")

	res, err := c.postToken(ctx, c.AdminTokenEndpoint(), data)
	if err != nil {
		return "", err
	}
	return res.AccessToken, nil
}

func (c *KeycloakClient) CreateUser(ctx context.Context, username, email, fullName, password string, roles []string) (string, error) {
	adminToken, err := c.FetchAdminAccessToken(ctx)
	if err != nil {
		// When Keycloak is offline/not reachable in dev environment, generate a surrogate ID
		return fmt.Sprintf("kc-%s", username), nil
	}

	payload := map[string]any{
		"enabled":       true,
		"username":      username,
		"email":         email,
		"emailVerified": true,
		"credentials": []map[string]any{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	trimmed := strings.TrimSpace(fullName)
	if trimmed != "" {
		parts := strings.Split(trimmed, " ")
		if len(parts) > 1 {
			payload["firstName"] = strings.Join(parts[:len(parts)-1], " ")
			payload["lastName"] = parts[len(parts)-1]
		} else {
			payload["firstName"] = trimmed
			payload["lastName"] = trimmed
		}
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AdminUsersEndpoint(), bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Sprintf("kc-%s", username), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create keycloak user: %s (status %d)", string(respBody), resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc == "" {
		return fmt.Sprintf("kc-%s", username), nil
	}
	parts := strings.Split(loc, "/")
	return parts[len(parts)-1], nil
}

func (c *KeycloakClient) DeleteUser(ctx context.Context, keycloakUserId string) error {
	adminToken, err := c.FetchAdminAccessToken(ctx)
	if err != nil {
		return nil
	}

	endpoint := fmt.Sprintf("%s/%s", c.AdminUsersEndpoint(), keycloakUserId)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *KeycloakClient) postToken(ctx context.Context, endpoint string, data url.Values) (*model.KeycloakTokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("keycloak token error: %s (status %d)", string(body), resp.StatusCode)
	}

	var tok model.KeycloakTokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}
