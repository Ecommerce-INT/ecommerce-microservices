package handler

import (
	"strings"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/inventory-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/inventory/actuator/health", h.HealthCheck)

	// Context root: /inventory
	invGroup := app.Group("/inventory/api/inventory")
	h.registerEndpoints(invGroup)

	// Direct root
	directGroup := app.Group("/api/inventory")
	h.registerEndpoints(directGroup)
}

func (h *InventoryHandler) registerEndpoints(r fiber.Router) {
	r.Get("", h.IsInStock)
	r.Get("/:products", h.IsInStockPath)
}

func (h *InventoryHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

func (h *InventoryHandler) IsInStock(c *fiber.Ctx) error {
	// Support both ?productName=a&productName=b and ?productName=a,b
	queryValues := c.Context().QueryArgs().PeekMulti("productName")
	var names []string

	if len(queryValues) > 0 {
		for _, qv := range queryValues {
			parts := strings.Split(string(qv), ",")
			for _, p := range parts {
				trimmed := strings.TrimSpace(p)
				if trimmed != "" {
					names = append(names, trimmed)
				}
			}
		}
	} else if pParam := c.Query("productName"); pParam != "" {
		parts := strings.Split(pParam, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				names = append(names, trimmed)
			}
		}
	}

	result, err := h.svc.IsInStock(c.Context(), names)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(result)
}

func (h *InventoryHandler) IsInStockPath(c *fiber.Ctx) error {
	products := c.Params("products")
	parts := strings.Split(products, ",")
	var names []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			names = append(names, trimmed)
		}
	}

	result, err := h.svc.IsInStock(c.Context(), names)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(result)
}
