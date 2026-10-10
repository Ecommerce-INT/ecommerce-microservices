package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/shipping-service/internal/handler"
	"com.ecommerce/shipping-service/internal/repository"
	"com.ecommerce/shipping-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupShippingApp() http.Handler {
	r := chi.NewRouter()
	repo := repository.NewShippingRepository(nil)
	svc := service.NewShippingService(repo)
	h := handler.NewShippingHandler(svc)
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
	app := setupShippingApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/actuator/health", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/shipping/actuator/health", nil).Code)
}

func TestShippingEndpoints(t *testing.T) {
	app := setupShippingApp()

	// List all
	assert.Equal(t, 200, performRequest(app, "GET", "/api/shippings", nil).Code)

	// List all with context path
	assert.Equal(t, 200, performRequest(app, "GET", "/shipping/api/shippings", nil).Code)

	// Save shipping item
	body := []byte(`{"orderId": 101, "productId": 202, "orderedQuantity": 3}`)
	assert.Equal(t, 200, performRequest(app, "POST", "/api/shippings", body).Code)

	// Update shipping item
	assert.Equal(t, 200, performRequest(app, "PUT", "/api/shippings", body).Code)

	// Delete shipping item by path
	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/shippings/101/202", nil).Code)

	// Invalid ID validation
	assert.Equal(t, 400, performRequest(app, "GET", "/api/shippings/abc/def", nil).Code)
}
