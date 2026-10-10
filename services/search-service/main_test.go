package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"com.ecommerce/search-service/internal/handler"
	"com.ecommerce/search-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupSearchApp() http.Handler {
	r := chi.NewRouter()
	svc := service.NewSearchService(nil)
	h := handler.NewSearchHandler(svc)
	h.RegisterRoutes(r)
	return r
}

func performRequest(app http.Handler, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", target, nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	return rec
}

func TestHealthCheck(t *testing.T) {
	app := setupSearchApp()

	assert.Equal(t, 200, performRequest(app, "/actuator/health").Code)
	assert.Equal(t, 200, performRequest(app, "/search/actuator/health").Code)
}

func TestCatalogSearch(t *testing.T) {
	app := setupSearchApp()

	// Direct endpoint
	assert.Equal(t, 200, performRequest(app, "/storefront/catalog-search?keyword=phone&page=0&size=10&minPrice=100&maxPrice=1000").Code)

	// Context prefix
	assert.Equal(t, 200, performRequest(app, "/search/storefront/catalog-search?keyword=laptop").Code)
}

func TestSearchSuggest(t *testing.T) {
	app := setupSearchApp()

	// Direct endpoint
	assert.Equal(t, 200, performRequest(app, "/storefront/search_suggest?keyword=iph").Code)

	// Context prefix
	assert.Equal(t, 200, performRequest(app, "/search/storefront/search_suggest?keyword=sam").Code)
}
