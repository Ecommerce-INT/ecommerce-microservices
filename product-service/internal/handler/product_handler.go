package handler

import (
	"errors"
	"strconv"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/product-service/internal/model"
	"com.ecommerce/product-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/product/actuator/health", h.HealthCheck)

	// Context root: /product
	prodGroup := app.Group("/product")
	h.registerCategoryRoutes(prodGroup.Group("/api/categories"))
	h.registerProductRoutes(prodGroup.Group("/api/products"))

	// Direct root
	h.registerCategoryRoutes(app.Group("/api/categories"))
	h.registerProductRoutes(app.Group("/api/products"))
}

func (h *ProductHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

// =================== Categories ===================

func (h *ProductHandler) registerCategoryRoutes(r fiber.Router) {
	r.Get("", h.ListCategories)
	r.Get("/paging", h.ListCategories)
	r.Get("/paging-and-sorting", h.ListCategories)
	r.Get("/:categoryId", h.GetCategoryByID)
	r.Post("", h.CreateCategory)
	r.Put("", h.UpdateCategory)
	r.Put("/:categoryId", h.UpdateCategoryByID)
	r.Delete("/:categoryId", h.DeleteCategory)
}

func (h *ProductHandler) ListCategories(c *fiber.Ctx) error {
	list, err := h.svc.FindAllCategories(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(list)
}

func (h *ProductHandler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return response.BadRequest(c, "Invalid categoryId")
	}
	cat, err := h.svc.FindCategoryByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(cat)
}

func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	var req model.CategoryDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	created, err := h.svc.SaveCategory(c.Context(), &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *ProductHandler) UpdateCategory(c *fiber.Ctx) error {
	var req model.CategoryDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	updated, err := h.svc.UpdateCategory(c.Context(), req.CategoryId, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *ProductHandler) UpdateCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return response.BadRequest(c, "Invalid categoryId")
	}
	var req model.CategoryDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	updated, err := h.svc.UpdateCategory(c.Context(), id, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *ProductHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return response.BadRequest(c, "Invalid categoryId")
	}
	err = h.svc.DeleteCategory(c.Context(), id)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}

// =================== Products ===================

func (h *ProductHandler) registerProductRoutes(r fiber.Router) {
	r.Get("", h.ListProducts)
	r.Get("/:productId", h.GetProductByID)
	r.Post("", h.CreateProduct)
	r.Put("", h.UpdateProduct)
	r.Put("/:productId", h.UpdateProductByID)
	r.Delete("/:productId", h.DeleteProduct)
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	list, err := h.svc.FindAllProducts(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(list)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}
	p, err := h.svc.FindProductByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(p)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req model.ProductDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	created, err := h.svc.SaveProduct(c.Context(), &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var req model.ProductDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	updated, err := h.svc.UpdateProduct(c.Context(), req.ProductId, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *ProductHandler) UpdateProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}
	var req model.ProductDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	updated, err := h.svc.UpdateProduct(c.Context(), id, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(updated)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}
	err = h.svc.DeleteProduct(c.Context(), id)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(true)
}
