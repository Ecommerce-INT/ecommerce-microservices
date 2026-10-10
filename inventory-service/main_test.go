package main

import (
	"net/http/httptest"
	"testing"

	"com.ecommerce/inventory-service/internal/handler"
	"com.ecommerce/inventory-service/internal/repository"
	"com.ecommerce/inventory-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupInventoryApp() *fiber.App {
	app := fiber.New()
	repo := repository.NewInventoryRepository(nil)
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupInventoryApp()

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqInv := httptest.NewRequest("GET", "/inventory/actuator/health", nil)
	respInv, err := app.Test(reqInv, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respInv.StatusCode)
}

func TestIsInStockQuery(t *testing.T) {
	app := setupInventoryApp()

	// Query with single productName
	req := httptest.NewRequest("GET", "/api/inventory?productName=iPhone%2015", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Query with comma-separated productNames
	reqMulti := httptest.NewRequest("GET", "/inventory/api/inventory?productName=iPhone,MacBook", nil)
	respMulti, err := app.Test(reqMulti, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respMulti.StatusCode)

	// Query with empty param
	reqEmpty := httptest.NewRequest("GET", "/api/inventory", nil)
	respEmpty, err := app.Test(reqEmpty, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respEmpty.StatusCode)
}

func TestIsInStockPath(t *testing.T) {
	app := setupInventoryApp()

	// Path with comma-separated products
	req := httptest.NewRequest("GET", "/api/inventory/iPhone,iPad,MacBook", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Context prefix
	reqCtx := httptest.NewRequest("GET", "/inventory/api/inventory/iPhone", nil)
	respCtx, err := app.Test(reqCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respCtx.StatusCode)
}
