package handler

import (
	"errors"
	"strconv"

	"com.ecommerce/order-service/internal/model"
	"com.ecommerce/order-service/internal/service"
	"com.ecommerce/pkg/common/response"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/order/actuator/health", h.HealthCheck)
	app.Get("/health", h.HealthCheck)

	// Cart endpoints
	h.registerCartRoutes(app.Group("/api/carts"))
	h.registerCartRoutes(app.Group("/order/api/carts"))

	// Order endpoints
	h.registerOrderRoutes(app.Group("/api/orders"))
	h.registerOrderRoutes(app.Group("/order/api/orders"))
}

func (h *OrderHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

// -------------------------------------------------------------
// Cart Handlers
// -------------------------------------------------------------

func (h *OrderHandler) registerCartRoutes(r fiber.Router) {
	r.Get("/all", h.FindAllCartsPaged)
	r.Get("/:cartId", h.FindCartById)
	r.Get("", h.FindAllCarts)
	r.Post("", h.SaveCart)
	r.Put("/:cartId", h.UpdateCartById)
	r.Put("", h.UpdateCart)
	r.Delete("/:cartId", h.DeleteCartById)
}

func (h *OrderHandler) FindAllCarts(c *fiber.Ctx) error {
	carts, err := h.svc.FindAllCarts(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(carts)
}

func (h *OrderHandler) FindAllCartsPaged(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "0"))
	size, _ := strconv.Atoi(c.Query("size", "10"))
	sortBy := c.Query("sortBy", "cartId")
	sortOrder := c.Query("sortOrder", "asc")

	res, err := h.svc.FindAllCartsPaged(c.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(res)
}

func (h *OrderHandler) FindCartById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("cartId"))
	if err != nil {
		return response.BadRequest(c, "Invalid cartId")
	}

	cart, err := h.svc.FindCartById(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(cart)
}

func (h *OrderHandler) SaveCart(c *fiber.Ctx) error {
	var req model.CartDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	saved, err := h.svc.SaveCart(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(saved)
}

func (h *OrderHandler) UpdateCart(c *fiber.Ctx) error {
	var req model.CartDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if req.CartId == nil {
		return response.BadRequest(c, "cartId is required for update")
	}

	updated, err := h.svc.UpdateCart(c.Context(), *req.CartId, req)
	if err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *OrderHandler) UpdateCartById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("cartId"))
	if err != nil {
		return response.BadRequest(c, "Invalid cartId")
	}

	var req model.CartDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	updated, err := h.svc.UpdateCart(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *OrderHandler) DeleteCartById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("cartId"))
	if err != nil {
		return response.BadRequest(c, "Invalid cartId")
	}

	if err := h.svc.DeleteCartById(c.Context(), id); err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}

// -------------------------------------------------------------
// Order Handlers
// -------------------------------------------------------------

func (h *OrderHandler) registerOrderRoutes(r fiber.Router) {
	r.Get("/existOrderId", h.ExistsByOrderId)
	r.Get("/all", h.FindAllOrdersPaged)
	r.Get("/:orderId", h.FindOrderById)
	r.Get("", h.FindAllOrders)
	r.Post("", h.SaveOrder)
	r.Put("/:orderId", h.UpdateOrderById)
	r.Put("", h.UpdateOrder)
	r.Delete("/:orderId", h.DeleteOrderById)
}

func (h *OrderHandler) FindAllOrders(c *fiber.Ctx) error {
	orders, err := h.svc.FindAllOrders(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(orders)
}

func (h *OrderHandler) FindAllOrdersPaged(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "0"))
	size, _ := strconv.Atoi(c.Query("size", "10"))
	sortBy := c.Query("sortBy", "orderId")
	sortOrder := c.Query("sortOrder", "asc")

	res, err := h.svc.FindAllOrdersPaged(c.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(res)
}

func (h *OrderHandler) FindOrderById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}

	order, err := h.svc.FindOrderById(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(order)
}

func (h *OrderHandler) ExistsByOrderId(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}

	exists, err := h.svc.ExistsByOrderId(c.Context(), id)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(exists)
}

func (h *OrderHandler) SaveOrder(c *fiber.Ctx) error {
	var req model.OrderDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	saved, err := h.svc.SaveOrder(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(saved)
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	var req model.OrderDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if req.OrderId == nil {
		return response.BadRequest(c, "orderId is required for update")
	}

	updated, err := h.svc.UpdateOrder(c.Context(), *req.OrderId, req)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *OrderHandler) UpdateOrderById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}

	var req model.OrderDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	updated, err := h.svc.UpdateOrder(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *OrderHandler) DeleteOrderById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}

	if err := h.svc.DeleteOrderById(c.Context(), id); err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}
