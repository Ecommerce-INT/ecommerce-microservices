package main

import (
	"net/http/httptest"
	"testing"

	"com.ecommerce/product-service/internal/handler"
	"com.ecommerce/product-service/internal/repository"
	"com.ecommerce/product-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	repo := repository.NewProductRepository(nil)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	h.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqProd := httptest.NewRequest("GET", "/product/actuator/health", nil)
	respProd, err := app.Test(reqProd, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, respProd.StatusCode)
}

func TestFavouriteRoutes(t *testing.T) {
	app := fiber.New()
	favRepo := repository.NewFavouriteRepository(nil)
	favSvc := service.NewFavouriteService(favRepo)
	favH := handler.NewFavouriteHandler(favSvc)
	favH.RegisterRoutes(app)

	// Health check
	req := httptest.NewRequest("GET", "/favourite/actuator/health", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// List favourites
	reqList := httptest.NewRequest("GET", "/api/favourites", nil)
	respList, err := app.Test(reqList, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respList.StatusCode)

	// List favourites via context prefix
	reqListCtx := httptest.NewRequest("GET", "/favourite/api/favourites", nil)
	respListCtx, err := app.Test(reqListCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respListCtx.StatusCode)
}

