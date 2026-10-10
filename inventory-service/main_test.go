package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/inventory-service/internal/handler"
	"com.ecommerce/inventory-service/internal/repository"
	"com.ecommerce/inventory-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupInventoryApp() http.Handler {
	r := chi.NewRouter()
	repo := repository.NewInventoryRepository(nil)
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryHandler(svc)
	h.RegisterRoutes(r)
	return r
}

func performRequest(app http.Handler, method, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func TestHealthCheck(t *testing.T) {
	app := setupInventoryApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/actuator/health").Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/inventory/actuator/health").Code)
}

func TestIsInStockQuery(t *testing.T) {
	app := setupInventoryApp()

	// Query with single productName
	assert.Equal(t, 200, performRequest(app, "GET", "/api/inventory?productName=iPhone%2015").Code)

	// Query with comma-separated productNames
	assert.Equal(t, 200, performRequest(app, "GET", "/inventory/api/inventory?productName=iPhone,MacBook").Code)

	// Query with empty param
	assert.Equal(t, 200, performRequest(app, "GET", "/api/inventory").Code)
}

func TestIsInStockPath(t *testing.T) {
	app := setupInventoryApp()

	// Path with comma-separated products
	assert.Equal(t, 200, performRequest(app, "GET", "/api/inventory/iPhone,iPad,MacBook").Code)

	// Context prefix
	assert.Equal(t, 200, performRequest(app, "GET", "/inventory/api/inventory/iPhone").Code)
}
