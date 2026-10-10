package main

import (
	"net/http/httptest"
	"testing"

	"com.ecommerce/shipping-service/internal/handler"
	"com.ecommerce/shipping-service/internal/repository"
	"com.ecommerce/shipping-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	repo := repository.NewShippingRepository(nil)
	svc := service.NewShippingService(repo)
	h := handler.NewShippingHandler(svc)
	h.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqShip := httptest.NewRequest("GET", "/shipping/actuator/health", nil)
	respShip, err := app.Test(reqShip, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, respShip.StatusCode)
}
