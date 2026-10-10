package response

import (
	"encoding/json"
	"net/http"
	"time"

	"com.ecommerce/pkg/common/middleware"
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

// WriteJSON writes v as an application/json response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func OK[T any](w http.ResponseWriter, r *http.Request, data T, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	WriteJSON(w, http.StatusOK, ApiResponse[T]{
		Success:   true,
		Code:      "OK",
		Message:   msg,
		Data:      data,
		TraceID:   middleware.GetCorrelationID(r.Context()),
		Timestamp: time.Now().UTC(),
	})
}

func Message(w http.ResponseWriter, r *http.Request, message string) {
	WriteJSON(w, http.StatusOK, ApiResponse[any]{
		Success:   true,
		Code:      "OK",
		Message:   message,
		TraceID:   middleware.GetCorrelationID(r.Context()),
		Timestamp: time.Now().UTC(),
	})
}

func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, errors ...string) {
	WriteJSON(w, status, ApiResponse[any]{
		Success:   false,
		Code:      code,
		Message:   message,
		Errors:    errors,
		Path:      r.URL.Path,
		TraceID:   middleware.GetCorrelationID(r.Context()),
		Timestamp: time.Now().UTC(),
	})
}

// Health writes the liveness payload consumed by the Kubernetes probes.
func Health(w http.ResponseWriter, component string) {
	WriteJSON(w, http.StatusOK, map[string]any{
		"status": "UP",
		"components": map[string]any{
			component: map[string]any{"status": "UP"},
		},
	})
}
