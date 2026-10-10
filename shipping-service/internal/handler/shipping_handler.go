package handler

import (
	"errors"
	"strconv"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/shipping-service/internal/model"
	"com.ecommerce/shipping-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ShippingHandler struct {
	svc *service.ShippingService
}

func NewShippingHandler(svc *service.ShippingService) *ShippingHandler {
	return &ShippingHandler{svc: svc}
}

func (h *ShippingHandler) RegisterRoutes(app *fiber.App) {
	// Health checks
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/shipping/actuator/health", h.HealthCheck)

	// Context root: /shipping/api/shippings
	shippingGroup := app.Group("/shipping/api/shippings")
	h.registerEndpoints(shippingGroup)

	// Fallback root: /api/shippings
	directGroup := app.Group("/api/shippings")
	h.registerEndpoints(directGroup)
}

func (h *ShippingHandler) registerEndpoints(r fiber.Router) {
	r.Get("", h.FindAll)
	r.Get("/find", h.FindByBody)
	r.Get("/:orderId/:productId", h.FindByID)
	r.Post("", h.Save)
	r.Put("", h.Update)
	r.Delete("/delete", h.DeleteByBody)
	r.Delete("/:orderId/:productId", h.DeleteByID)
}

func (h *ShippingHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

func (h *ShippingHandler) FindAll(c *fiber.Ctx) error {
	items, err := h.svc.FindAll(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(model.DtoCollectionResponse[model.OrderItemDto]{
		Collection: items,
	})
}

func (h *ShippingHandler) FindByID(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}
	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}

	item, err := h.svc.FindByID(c.Context(), orderId, productId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(item)
}

func (h *ShippingHandler) FindByBody(c *fiber.Ctx) error {
	var req model.OrderItemId
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	item, err := h.svc.FindByID(c.Context(), req.OrderId, req.ProductId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(item)
}

func (h *ShippingHandler) Save(c *fiber.Ctx) error {
	var req model.OrderItemDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	saved, err := h.svc.Save(c.Context(), &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(saved)
}

func (h *ShippingHandler) Update(c *fiber.Ctx) error {
	var req model.OrderItemDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	updated, err := h.svc.Update(c.Context(), &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *ShippingHandler) DeleteByID(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}
	productId, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}

	err = h.svc.DeleteByID(c.Context(), orderId, productId)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}

func (h *ShippingHandler) DeleteByBody(c *fiber.Ctx) error {
	var req model.OrderItemId
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	err := h.svc.DeleteByID(c.Context(), req.OrderId, req.ProductId)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}
