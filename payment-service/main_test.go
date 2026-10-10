package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"com.ecommerce/payment-service/internal/handler"
	"com.ecommerce/payment-service/internal/repository"
	"com.ecommerce/payment-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	repo := repository.NewPaymentRepository(nil)
	producer := service.NewEventProducer(nil)
	svc := service.NewPaymentService(repo, producer)
	h := handler.NewPaymentHandler(svc)

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

	reqPayment := httptest.NewRequest("GET", "/payment/actuator/health", nil)
	respPayment, err := app.Test(reqPayment)
	assert.NoError(t, err)
	assert.Equal(t, 200, respPayment.StatusCode)
}

func TestPaymentsEndpoints(t *testing.T) {
	app := setupTestApp()

	// GET all payments
	req := httptest.NewRequest("GET", "/api/payments", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Save payment
	body := []byte(`{"orderId": 1, "userId": 10, "isPayed": true, "paymentStatus": "COMPLETED"}`)
	reqPost := httptest.NewRequest("POST", "/api/payments", bytes.NewReader(body))
	reqPost.Header.Set("Content-Type", "application/json")
	respPost, err := app.Test(reqPost)
	assert.NoError(t, err)
	assert.Equal(t, 200, respPost.StatusCode)

	// Get order dto
	reqOrder := httptest.NewRequest("GET", "/api/payments/getOrder/1", nil)
	respOrder, err := app.Test(reqOrder)
	assert.NoError(t, err)
	assert.Equal(t, 200, respOrder.StatusCode)

	// Delete payment
	reqDel := httptest.NewRequest("DELETE", "/api/payments/1", nil)
	respDel, err := app.Test(reqDel)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)
}
