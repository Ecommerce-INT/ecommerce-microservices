package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// ApiResponse represents the platform unified response envelope
type ApiResponse[T any] struct {
	Success   bool      `json:"success"`
	Code      string    `json:"code,omitempty"`
	Message   string    `json:"message,omitempty"`
	Data      T         `json:"data,omitempty"`
	Errors    []string  `json:"errors,omitempty"`
	Path      string    `json:"path,omitempty"`
	TraceID   string    `json:"traceId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func OK[T any](c *fiber.Ctx, data T, message ...string) error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	traceID, _ := c.Locals("correlation_id").(string)
	return c.JSON(ApiResponse[T]{
		Success:   true,
		Code:      "OK",
		Message:   msg,
		Data:      data,
		TraceID:   traceID,
		Timestamp: time.Now().UTC(),
	})
}

func Message(c *fiber.Ctx, message string) error {
	traceID, _ := c.Locals("correlation_id").(string)
	return c.JSON(ApiResponse[any]{
		Success:   true,
		Code:      "OK",
		Message:   message,
		TraceID:   traceID,
		Timestamp: time.Now().UTC(),
	})
}

func Error(c *fiber.Ctx, status int, code, message string, errors ...string) error {
	traceID, _ := c.Locals("correlation_id").(string)
	return c.Status(status).JSON(ApiResponse[any]{
		Success:   false,
		Code:      code,
		Message:   message,
		Errors:    errors,
		Path:      c.Path(),
		TraceID:   traceID,
		Timestamp: time.Now().UTC(),
	})
}
