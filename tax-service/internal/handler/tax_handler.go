package handler

import (
	"errors"
	"strconv"
	"strings"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/tax-service/internal/model"
	"com.ecommerce/tax-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TaxHandler struct {
	svc *service.TaxService
}

func NewTaxHandler(svc *service.TaxService) *TaxHandler {
	return &TaxHandler{svc: svc}
}

func (h *TaxHandler) RegisterRoutes(app *fiber.App) {
	// Liveness & Readiness checks
	app.Get("/actuator/health", h.HealthCheck)
	app.Get("/tax/actuator/health", h.HealthCheck)

	// Context root: /tax
	taxGroup := app.Group("/tax")

	// Backoffice Tax Classes
	classes := taxGroup.Group("/backoffice/tax-classes")
	classes.Get("/paging", h.GetPageableTaxClasses)
	classes.Get("/:id", h.GetTaxClassByID)
	classes.Get("", h.ListTaxClasses)
	classes.Post("", h.CreateTaxClass)
	classes.Put("/:id", h.UpdateTaxClass)
	classes.Delete("/:id", h.DeleteTaxClass)

	// Backoffice Tax Rates
	rates := taxGroup.Group("/backoffice/tax-rates")
	rates.Get("/paging", h.GetPageableTaxRates)
	rates.Get("/tax-percent", h.GetTaxPercent)
	rates.Get("/location-based-batch", h.GetBatchTaxRates)
	rates.Get("/:id", h.GetTaxRateByID)
	rates.Get("", h.ListTaxRates)
	rates.Post("", h.CreateTaxRate)
	rates.Put("/:id", h.UpdateTaxRate)
	rates.Delete("/:id", h.DeleteTaxRate)
}

func (h *TaxHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

// =================== TaxClass Handlers ===================

func (h *TaxHandler) ListTaxClasses(c *fiber.Ctx) error {
	list, err := h.svc.FindAllTaxClasses(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(list)
}

func (h *TaxHandler) GetTaxClassByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid tax class ID")
	}

	tc, err := h.svc.FindTaxClassByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(tc)
}

func (h *TaxHandler) CreateTaxClass(c *fiber.Ctx) error {
	var req model.TaxClassPostVm
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.BadRequest(c, "Tax class name must not be empty")
	}

	created, err := h.svc.CreateTaxClass(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrDuplicated) {
			return response.BadRequest(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	c.Set("Location", "/tax-classes/"+strconv.FormatInt(created.ID, 10))
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *TaxHandler) UpdateTaxClass(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid tax class ID")
	}

	var req model.TaxClassPostVm
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}
	if strings.TrimSpace(req.Name) == "" {
		return response.BadRequest(c, "Tax class name must not be empty")
	}

	err = h.svc.UpdateTaxClass(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		if errors.Is(err, service.ErrDuplicated) {
			return response.BadRequest(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaxHandler) DeleteTaxClass(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid tax class ID")
	}

	err = h.svc.DeleteTaxClass(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaxHandler) GetPageableTaxClasses(c *fiber.Ctx) error {
	pageNo := c.QueryInt("pageNo", 0)
	pageSize := c.QueryInt("pageSize", 10)

	result, err := h.svc.GetPageableTaxClasses(c.Context(), pageNo, pageSize)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(result)
}

// =================== TaxRate Handlers ===================

func (h *TaxHandler) ListTaxRates(c *fiber.Ctx) error {
	list, err := h.svc.FindAllTaxRates(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(list)
}

func (h *TaxHandler) GetTaxRateByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid tax rate ID")
	}

	tr, err := h.svc.FindTaxRateByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return c.JSON(tr)
}

func (h *TaxHandler) CreateTaxRate(c *fiber.Ctx) error {
	var req model.TaxRatePostVm
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	created, err := h.svc.CreateTaxRate(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	c.Set("Location", "/tax-rates/"+strconv.FormatInt(created.ID, 10))
	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *TaxHandler) UpdateTaxRate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid tax rate ID")
	}

	var req model.TaxRatePostVm
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	err = h.svc.UpdateTaxRate(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaxHandler) DeleteTaxRate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid tax rate ID")
	}

	err = h.svc.DeleteTaxRate(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaxHandler) GetPageableTaxRates(c *fiber.Ctx) error {
	pageNo := c.QueryInt("pageNo", 0)
	pageSize := c.QueryInt("pageSize", 10)

	result, err := h.svc.GetPageableTaxRates(c.Context(), pageNo, pageSize)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(result)
}

func (h *TaxHandler) GetTaxPercent(c *fiber.Ctx) error {
	taxClassIDStr := c.Query("taxClassId")
	countryIDStr := c.Query("countryId")

	if taxClassIDStr == "" || countryIDStr == "" {
		return response.BadRequest(c, "taxClassId and countryId are required")
	}

	taxClassID, err := strconv.ParseInt(taxClassIDStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid taxClassId")
	}

	countryID, err := strconv.ParseInt(countryIDStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid countryId")
	}

	var stateOrProvinceID *int64
	if sID := c.Query("stateOrProvinceId"); sID != "" {
		if val, err := strconv.ParseInt(sID, 10, 64); err == nil {
			stateOrProvinceID = &val
		}
	}

	var zipCode *string
	if z := c.Query("zipCode"); z != "" {
		zipCode = &z
	}

	percent, err := h.svc.GetTaxPercent(c.Context(), taxClassID, countryID, stateOrProvinceID, zipCode)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return c.JSON(percent)
}

func (h *TaxHandler) GetBatchTaxRates(c *fiber.Ctx) error {
	taxClassIDsStr := c.Query("taxClassIds")
	countryIDStr := c.Query("countryId")

	if taxClassIDsStr == "" || countryIDStr == "" {
		return response.BadRequest(c, "taxClassIds and countryId are required")
	}

	countryID, err := strconv.ParseInt(countryIDStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid countryId")
	}

	var taxClassIDs []int64
	for _, idStr := range strings.Split(taxClassIDsStr, ",") {
		idStr = strings.TrimSpace(idStr)
		if val, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			taxClassIDs = append(taxClassIDs, val)
		}
	}

	var stateOrProvinceID *int64
	if sID := c.Query("stateOrProvinceId"); sID != "" {
		if val, err := strconv.ParseInt(sID, 10, 64); err == nil {
			stateOrProvinceID = &val
		}
	}

	var zipCode *string
	if z := c.Query("zipCode"); z != "" {
		zipCode = &z
	}

	list, err := h.svc.GetBatchTaxRates(c.Context(), taxClassIDs, countryID, stateOrProvinceID, zipCode)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return c.JSON(list)
}
