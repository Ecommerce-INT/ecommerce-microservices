package handler

import (
	"strconv"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/search-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	svc *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/search/actuator/health", h.HealthCheck)

	// Context root: /search
	searchGroup := app.Group("/search")
	h.registerEndpoints(searchGroup)

	// Direct root
	h.registerEndpoints(app)
}

func (h *SearchHandler) registerEndpoints(r fiber.Router) {
	r.Get("/storefront/catalog-search", h.CatalogSearch)
	r.Get("/storefront/search_suggest", h.SearchSuggest)
}

func (h *SearchHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"elasticsearch": fiber.Map{"status": "UP"},
		},
	})
}

func (h *SearchHandler) CatalogSearch(c *fiber.Ctx) error {
	keyword := c.Query("keyword", "")
	page := c.QueryInt("page", 0)
	size := c.QueryInt("size", 12)
	brand := c.Query("brand", "")
	category := c.Query("category", "")
	attribute := c.Query("attribute", "")
	sortType := c.Query("sortType", "DEFAULT")

	var minPrice *float64
	if minP := c.Query("minPrice"); minP != "" {
		if val, err := strconv.ParseFloat(minP, 64); err == nil {
			minPrice = &val
		}
	}

	var maxPrice *float64
	if maxP := c.Query("maxPrice"); maxP != "" {
		if val, err := strconv.ParseFloat(maxP, 64); err == nil {
			maxPrice = &val
		}
	}

	res, err := h.svc.FindProductAdvance(c.Context(), keyword, page, size, brand, category, attribute, minPrice, maxPrice, sortType)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return c.JSON(res)
}

func (h *SearchHandler) SearchSuggest(c *fiber.Ctx) error {
	keyword := c.Query("keyword", "")
	res, err := h.svc.AutoComplete(c.Context(), keyword)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(res)
}
