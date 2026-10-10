package model

type User struct {
	ID             int64    `json:"id"`
	FullName       string   `json:"fullname"`
	Username       string   `json:"username"`
	Email          string   `json:"email"`
	Gender         string   `json:"gender"`
	Phone          string   `json:"phone"`
	Avatar         string   `json:"avatar"`
	KeycloakUserID string   `json:"keycloakUserId,omitempty"`
	Roles          []string `json:"roles,omitempty"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

type RegisterRequest struct {
	FullName string   `json:"fullName"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Gender   string   `json:"gender"`
	Phone    string   `json:"phone"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
}

type UpdateUserRequest struct {
	FullName *string `json:"fullName"`
	Email    *string `json:"email"`
	Gender   *string `json:"gender"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type KeycloakTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}

type PageResponse[T any] struct {
	Content          []T   `json:"content"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int   `json:"totalPages"`
	Size             int   `json:"size"`
	Number           int   `json:"number"`
	NumberOfElements int   `json:"numberOfElements"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	Empty            bool  `json:"empty"`
}
