package main

import (
	"net/http/httptest"
	"testing"

	"com.ecommerce/search-service/internal/handler"
	"com.ecommerce/search-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupSearchApp() *fiber.App {
	app := fiber.New()
	svc := service.NewSearchService(nil)
	h := handler.NewSearchHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupSearchApp()

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqSearch := httptest.NewRequest("GET", "/search/actuator/health", nil)
	respSearch, err := app.Test(reqSearch, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respSearch.StatusCode)
}

func TestCatalogSearch(t *testing.T) {
	app := setupSearchApp()

	// Direct endpoint
	req := httptest.NewRequest("GET", "/storefront/catalog-search?keyword=phone&page=0&size=10&minPrice=100&maxPrice=1000", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Context prefix
	reqCtx := httptest.NewRequest("GET", "/search/storefront/catalog-search?keyword=laptop", nil)
	respCtx, err := app.Test(reqCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respCtx.StatusCode)
}

func TestSearchSuggest(t *testing.T) {
	app := setupSearchApp()

	// Direct endpoint
	req := httptest.NewRequest("GET", "/storefront/search_suggest?keyword=iph", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Context prefix
	reqCtx := httptest.NewRequest("GET", "/search/storefront/search_suggest?keyword=sam", nil)
	respCtx, err := app.Test(reqCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respCtx.StatusCode)
}
