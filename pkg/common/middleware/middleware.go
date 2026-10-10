package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const HeaderCorrelationID = "X-Correlation-Id"

type contextKey string

const (
	correlationIDKey contextKey = "correlationId"
	userIDKey        contextKey = "userId"
	emailKey         contextKey = "email"
	usernameKey      contextKey = "username"
)

// CorrelationID ensures every request has a correlation ID for distributed tracing
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corID := r.Header.Get(HeaderCorrelationID)
		if corID == "" {
			bytes := make([]byte, 16)
			_, _ = rand.Read(bytes)
			corID = hex.EncodeToString(bytes)
		}
		w.Header().Set(HeaderCorrelationID, corID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), correlationIDKey, corID)))
	})
}

// GetCorrelationID returns the correlation ID attached to the request context.
func GetCorrelationID(ctx context.Context) string {
	corID, _ := ctx.Value(correlationIDKey).(string)
	return corID
}

// UserClaims extracts claims from JWT if present (without blocking unauthenticated routes if APISIX already verified)
func UserClaims(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			parser := jwt.NewParser()
			token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if sub, ok := claims["sub"].(string); ok {
						ctx = context.WithValue(ctx, userIDKey, sub)
					}
					if email, ok := claims["email"].(string); ok {
						ctx = context.WithValue(ctx, emailKey, email)
					}
					if username, ok := claims["preferred_username"].(string); ok {
						ctx = context.WithValue(ctx, usernameKey, username)
					}
				}
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserID returns the "sub" claim attached to the request context, if any.
func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}

// Email returns the "email" claim attached to the request context, if any.
func Email(ctx context.Context) string {
	email, _ := ctx.Value(emailKey).(string)
	return email
}

// Username returns the "preferred_username" claim attached to the request context, if any.
func Username(ctx context.Context) string {
	username, _ := ctx.Value(usernameKey).(string)
	return username
}
