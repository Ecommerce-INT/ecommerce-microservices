package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/payment-service/internal/handler"
	"com.ecommerce/payment-service/internal/repository"
	"com.ecommerce/payment-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() http.Handler {
	repo := repository.NewPaymentRepository(nil)
	producer := service.NewEventProducer(nil)
	svc := service.NewPaymentService(repo, producer)
	h := handler.NewPaymentHandler(svc)

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
	assert.Equal(t, 200, performRequest(app, "GET", "/payment/actuator/health", nil).Code)
}

func TestPaymentsEndpoints(t *testing.T) {
	app := setupTestApp()

	// GET all payments
	assert.Equal(t, 200, performRequest(app, "GET", "/api/payments", nil).Code)

	// Save payment
	body := []byte(`{"orderId": 1, "userId": 10, "isPayed": true, "paymentStatus": "COMPLETED"}`)
	assert.Equal(t, 200, performRequest(app, "POST", "/api/payments", body).Code)

	// Get order dto
	assert.Equal(t, 200, performRequest(app, "GET", "/api/payments/getOrder/1", nil).Code)

	// Delete payment
	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/payments/1", nil).Code)
}
