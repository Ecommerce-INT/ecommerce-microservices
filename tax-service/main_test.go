package main

import (
	"net/http/httptest"
	"testing"

	"com.ecommerce/tax-service/internal/handler"
	"com.ecommerce/tax-service/internal/repository"
	"com.ecommerce/tax-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	repo := repository.NewTaxRepository(nil)
	svc := service.NewTaxService(repo)
	h := handler.NewTaxHandler(svc)
	h.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqTax := httptest.NewRequest("GET", "/tax/actuator/health", nil)
	respTax, err := app.Test(reqTax, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, respTax.StatusCode)
}
