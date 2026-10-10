package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/order-service/internal/handler"
	"com.ecommerce/order-service/internal/repository"
	"com.ecommerce/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() http.Handler {
	repo := repository.NewOrderRepository(nil)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

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
	assert.Equal(t, 200, performRequest(app, "GET", "/order/actuator/health", nil).Code)
}

func TestCartsEndpoints(t *testing.T) {
	app := setupTestApp()

	// GET all carts
	assert.Equal(t, 200, performRequest(app, "GET", "/api/carts", nil).Code)

	// Save cart
	body := []byte(`{"userId": 10}`)
	assert.Equal(t, 200, performRequest(app, "POST", "/api/carts", body).Code)

	// Delete cart
	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/carts/1", nil).Code)
}

func TestOrdersEndpoints(t *testing.T) {
	app := setupTestApp()

	// GET all orders
	assert.Equal(t, 200, performRequest(app, "GET", "/api/orders", nil).Code)

	// Save order
	body := []byte(`{"orderDesc": "test order", "orderFee": 25.5, "productId": 1}`)
	assert.Equal(t, 200, performRequest(app, "POST", "/api/orders", body).Code)

	// Exists order
	assert.Equal(t, 200, performRequest(app, "GET", "/api/orders/existOrderId?orderId=1", nil).Code)

	// Delete order
	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/orders/1", nil).Code)
}
