package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/product-service/internal/handler"
	"com.ecommerce/product-service/internal/repository"
	"com.ecommerce/product-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupProductApp() http.Handler {
	r := chi.NewRouter()

	repo := repository.NewProductRepository(nil)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	h.RegisterRoutes(r)

	favRepo := repository.NewFavouriteRepository(nil)
	favSvc := service.NewFavouriteService(favRepo)
	favH := handler.NewFavouriteHandler(favSvc)
	favH.RegisterRoutes(r)

	return r
}

func performRequest(app http.Handler, method, target string, body []byte) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func TestHealthCheck(t *testing.T) {
	app := setupProductApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/actuator/health", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/product/actuator/health", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/favourite/actuator/health", nil).Code)
}

func TestCategoryEndpoints(t *testing.T) {
	app := setupProductApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/api/categories", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/product/api/categories", nil).Code)

	body := []byte(`{"categoryTitle": "Electronics", "imageUrl": "http://img.png"}`)
	assert.Equal(t, 201, performRequest(app, "POST", "/api/categories", body).Code)

	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/categories/1", nil).Code)
}

func TestProductEndpoints(t *testing.T) {
	app := setupProductApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/api/products", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/product/api/products", nil).Code)

	body := []byte(`{"productTitle": "iPhone 15 Pro", "priceUnit": 999.99, "quantity": 50}`)
	assert.Equal(t, 201, performRequest(app, "POST", "/api/products", body).Code)

	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/products/1", nil).Code)
}

func TestFavouriteRoutes(t *testing.T) {
	app := setupProductApp()

	assert.Equal(t, 200, performRequest(app, "GET", "/api/favourites", nil).Code)
	assert.Equal(t, 200, performRequest(app, "GET", "/favourite/api/favourites", nil).Code)

	favBody := []byte(`{"userId": 1, "productId": 100, "likeDate": "2026-01-01T10:00:00"}`)
	assert.Equal(t, 200, performRequest(app, "POST", "/api/favourites", favBody).Code)

	assert.Equal(t, 200, performRequest(app, "DELETE", "/api/favourites/1/100/2026-01-01T10:00:00", nil).Code)
}
