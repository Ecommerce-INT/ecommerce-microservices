package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"com.ecommerce/shipping-service/internal/handler"
	"com.ecommerce/shipping-service/internal/repository"
	"com.ecommerce/shipping-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupShippingApp() *fiber.App {
	app := fiber.New()
	repo := repository.NewShippingRepository(nil)
	svc := service.NewShippingService(repo)
	h := handler.NewShippingHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupShippingApp()

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqShip := httptest.NewRequest("GET", "/shipping/actuator/health", nil)
	respShip, err := app.Test(reqShip, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respShip.StatusCode)
}

func TestShippingEndpoints(t *testing.T) {
	app := setupShippingApp()

	// List all
	reqList := httptest.NewRequest("GET", "/api/shippings", nil)
	respList, err := app.Test(reqList, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respList.StatusCode)

	// List all with context path
	reqCtx := httptest.NewRequest("GET", "/shipping/api/shippings", nil)
	respCtx, err := app.Test(reqCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respCtx.StatusCode)

	// Save shipping item
	body := []byte(`{"orderId": 101, "productId": 202, "orderedQuantity": 3}`)
	reqSave := httptest.NewRequest("POST", "/api/shippings", bytes.NewReader(body))
	reqSave.Header.Set("Content-Type", "application/json")
	respSave, err := app.Test(reqSave, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respSave.StatusCode)

	// Update shipping item
	reqUp := httptest.NewRequest("PUT", "/api/shippings", bytes.NewReader(body))
	reqUp.Header.Set("Content-Type", "application/json")
	respUp, err := app.Test(reqUp, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respUp.StatusCode)

	// Delete shipping item by path
	reqDel := httptest.NewRequest("DELETE", "/api/shippings/101/202", nil)
	respDel, err := app.Test(reqDel, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)

	// Invalid ID validation
	reqInvalid := httptest.NewRequest("GET", "/api/shippings/abc/def", nil)
	respInvalid, err := app.Test(reqInvalid, -1)
	assert.NoError(t, err)
	assert.Equal(t, 400, respInvalid.StatusCode)
}
