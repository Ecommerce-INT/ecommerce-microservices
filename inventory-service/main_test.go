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

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	repo := repository.NewInventoryRepository(nil)
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryHandler(svc)
	h.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqInv := httptest.NewRequest("GET", "/inventory/actuator/health", nil)
	respInv, err := app.Test(reqInv, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, respInv.StatusCode)
}
