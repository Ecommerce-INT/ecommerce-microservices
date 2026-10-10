package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"com.ecommerce/order-service/internal/handler"
	"com.ecommerce/order-service/internal/repository"
	"com.ecommerce/order-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	repo := repository.NewOrderRepository(nil)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	app := fiber.New()
	h.RegisterRoutes(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqOrder := httptest.NewRequest("GET", "/order/actuator/health", nil)
	respOrder, err := app.Test(reqOrder)
	assert.NoError(t, err)
	assert.Equal(t, 200, respOrder.StatusCode)
}

func TestCartsEndpoints(t *testing.T) {
	app := setupTestApp()

	// GET all carts
	req := httptest.NewRequest("GET", "/api/carts", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Save cart
	body := []byte(`{"userId": 10}`)
	reqPost := httptest.NewRequest("POST", "/api/carts", bytes.NewReader(body))
	reqPost.Header.Set("Content-Type", "application/json")
	respPost, err := app.Test(reqPost)
	assert.NoError(t, err)
	assert.Equal(t, 200, respPost.StatusCode)

	// Delete cart
	reqDel := httptest.NewRequest("DELETE", "/api/carts/1", nil)
	respDel, err := app.Test(reqDel)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)
}

func TestOrdersEndpoints(t *testing.T) {
	app := setupTestApp()

	// GET all orders
	req := httptest.NewRequest("GET", "/api/orders", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Save order
	body := []byte(`{"orderDesc": "test order", "orderFee": 25.5, "productId": 1}`)
	reqPost := httptest.NewRequest("POST", "/api/orders", bytes.NewReader(body))
	reqPost.Header.Set("Content-Type", "application/json")
	respPost, err := app.Test(reqPost)
	assert.NoError(t, err)
	assert.Equal(t, 200, respPost.StatusCode)

	// Exists order
	reqExists := httptest.NewRequest("GET", "/api/orders/existOrderId?orderId=1", nil)
	respExists, err := app.Test(reqExists)
	assert.NoError(t, err)
	assert.Equal(t, 200, respExists.StatusCode)

	// Delete order
	reqDel := httptest.NewRequest("DELETE", "/api/orders/1", nil)
	respDel, err := app.Test(reqDel)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)
}
