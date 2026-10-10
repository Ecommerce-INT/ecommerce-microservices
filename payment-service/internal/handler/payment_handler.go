package handler

import (
	"errors"
	"strconv"

	"com.ecommerce/payment-service/internal/model"
	"com.ecommerce/payment-service/internal/service"
	"com.ecommerce/pkg/common/response"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/payment/actuator/health", h.HealthCheck)
	app.Get("/health", h.HealthCheck)

	h.registerEndpoints(app.Group("/api/payments"))
	h.registerEndpoints(app.Group("/payment/api/payments"))
}

func (h *PaymentHandler) registerEndpoints(r fiber.Router) {
	r.Get("/getOrder/:orderId", h.GetOrderDto)
	r.Get("/all", h.FindAllPaged)
	r.Get("/:paymentId", h.FindById)
	r.Get("", h.FindAll)
	r.Post("", h.Save)
	r.Put("/:paymentId", h.UpdateById)
	r.Put("", h.Update)
	r.Delete("/:paymentId", h.DeleteById)
}

func (h *PaymentHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

func (h *PaymentHandler) FindAll(c *fiber.Ctx) error {
	payments, err := h.svc.FindAll(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(payments)
}

func (h *PaymentHandler) FindAllPaged(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "0"))
	size, _ := strconv.Atoi(c.Query("size", "10"))
	sortBy := c.Query("sortBy", "paymentId")
	sortOrder := c.Query("sortOrder", "asc")

	res, err := h.svc.FindAllPaged(c.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(res)
}

func (h *PaymentHandler) FindById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("paymentId"))
	if err != nil {
		return response.BadRequest(c, "Invalid paymentId")
	}

	payment, err := h.svc.FindById(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(payment)
}

func (h *PaymentHandler) GetOrderDto(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return response.BadRequest(c, "Invalid orderId")
	}

	dto := h.svc.GetOrderDto(c.Context(), id)
	return c.JSON(dto)
}

func (h *PaymentHandler) Save(c *fiber.Ctx) error {
	var req model.PaymentDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	saved, err := h.svc.Save(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrOrderAlreadyPaid) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(saved)
}

func (h *PaymentHandler) Update(c *fiber.Ctx) error {
	var req model.PaymentDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	updated, err := h.svc.Update(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *PaymentHandler) UpdateById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("paymentId"))
	if err != nil {
		return response.BadRequest(c, "Invalid paymentId")
	}

	var req model.PaymentDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	updated, err := h.svc.UpdateById(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *PaymentHandler) DeleteById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("paymentId"))
	if err != nil {
		return response.BadRequest(c, "Invalid paymentId")
	}

	if err := h.svc.DeleteById(c.Context(), id); err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}
