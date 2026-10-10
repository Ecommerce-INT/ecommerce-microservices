package main

import (
	"net/http/httptest"
	"testing"

	"com.ecommerce/search-service/internal/handler"
	"com.ecommerce/search-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	svc := service.NewSearchService(nil)
	h := handler.NewSearchHandler(svc)
	h.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqSearch := httptest.NewRequest("GET", "/search/actuator/health", nil)
	respSearch, err := app.Test(reqSearch, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, respSearch.StatusCode)
}
