package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const HeaderCorrelationID = "X-Correlation-Id"

// CorrelationID ensures every request has a correlation ID for distributed tracing
func CorrelationID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		corID := c.Get(HeaderCorrelationID)
		if corID == "" {
			bytes := make([]byte, 16)
			_, _ = rand.Read(bytes)
			corID = hex.EncodeToString(bytes)
		}
		c.Set(HeaderCorrelationID, corID)
		c.Locals("correlationId", corID)
		return c.Next()
	}
}

// UserClaims extracts claims from JWT if present (without blocking unauthenticated routes if APISIX already verified)
func UserClaims() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			parser := jwt.NewParser()
			token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if sub, exists := claims["sub"]; exists {
						c.Locals("userId", sub)
					}
					if email, exists := claims["email"]; exists {
						c.Locals("email", email)
					}
					if preferredUsername, exists := claims["preferred_username"]; exists {
						c.Locals("username", preferredUsername)
					}
				}
			}
		}
		return c.Next()
	}
}
