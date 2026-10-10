package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"com.ecommerce/product-service/internal/handler"
	"com.ecommerce/product-service/internal/repository"
	"com.ecommerce/product-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupProductApp() *fiber.App {
	app := fiber.New()
	repo := repository.NewProductRepository(nil)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	h.RegisterRoutes(app)

	favRepo := repository.NewFavouriteRepository(nil)
	favSvc := service.NewFavouriteService(favRepo)
	favH := handler.NewFavouriteHandler(favSvc)
	favH.RegisterRoutes(app)

	return app
}

func TestHealthCheck(t *testing.T) {
	app := setupProductApp()

	req := httptest.NewRequest("GET", "/actuator/health", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	reqProd := httptest.NewRequest("GET", "/product/actuator/health", nil)
	respProd, err := app.Test(reqProd, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respProd.StatusCode)

	reqFav := httptest.NewRequest("GET", "/favourite/actuator/health", nil)
	respFav, err := app.Test(reqFav, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respFav.StatusCode)
}

func TestCategoryEndpoints(t *testing.T) {
	app := setupProductApp()

	// List categories
	req := httptest.NewRequest("GET", "/api/categories", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// List categories via context prefix
	reqCtx := httptest.NewRequest("GET", "/product/api/categories", nil)
	respCtx, err := app.Test(reqCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respCtx.StatusCode)

	// Create category
	body := []byte(`{"categoryTitle": "Electronics", "imageUrl": "http://img.png"}`)
	reqPost := httptest.NewRequest("POST", "/api/categories", bytes.NewReader(body))
	reqPost.Header.Set("Content-Type", "application/json")
	respPost, err := app.Test(reqPost, -1)
	assert.NoError(t, err)
	assert.Equal(t, 201, respPost.StatusCode)

	// Delete category
	reqDel := httptest.NewRequest("DELETE", "/api/categories/1", nil)
	respDel, err := app.Test(reqDel, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)
}

func TestProductEndpoints(t *testing.T) {
	app := setupProductApp()

	// List products
	req := httptest.NewRequest("GET", "/api/products", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// List products via context prefix
	reqCtx := httptest.NewRequest("GET", "/product/api/products", nil)
	respCtx, err := app.Test(reqCtx, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respCtx.StatusCode)

	// Create product
	body := []byte(`{"productTitle": "iPhone 15 Pro", "priceUnit": 999.99, "quantity": 50}`)
	reqPost := httptest.NewRequest("POST", "/api/products", bytes.NewReader(body))
	reqPost.Header.Set("Content-Type", "application/json")
	respPost, err := app.Test(reqPost, -1)
	assert.NoError(t, err)
	assert.Equal(t, 201, respPost.StatusCode)

	// Delete product
	reqDel := httptest.NewRequest("DELETE", "/api/products/1", nil)
	respDel, err := app.Test(reqDel, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)
}

func TestFavouriteRoutes(t *testing.T) {
	app := setupProductApp()

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

	// Save favourite
	favBody := []byte(`{"userId": 1, "productId": 100, "likeDate": "2026-01-01T10:00:00"}`)
	reqSave := httptest.NewRequest("POST", "/api/favourites", bytes.NewReader(favBody))
	reqSave.Header.Set("Content-Type", "application/json")
	respSave, err := app.Test(reqSave, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respSave.StatusCode)

	// Delete favourite
	reqDel := httptest.NewRequest("DELETE", "/api/favourites/1/100/2026-01-01T10:00:00", nil)
	respDel, err := app.Test(reqDel, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, respDel.StatusCode)
}
