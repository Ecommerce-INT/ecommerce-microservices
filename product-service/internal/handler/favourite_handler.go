package handler

import (
	"net/url"
	"strconv"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/product-service/internal/model"
	"com.ecommerce/product-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type FavouriteHandler struct {
	svc *service.FavouriteService
}

func NewFavouriteHandler(svc *service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{svc: svc}
}

func (h *FavouriteHandler) RegisterRoutes(app *fiber.App) {
	// Health check aliases for favourite service compatibility
	app.Get("/favourite/actuator/health", h.HealthCheck)

	// Context root: /favourite/api/favourites
	favGroup := app.Group("/favourite/api/favourites")
	h.registerEndpoints(favGroup)

	// Context root: /product/api/favourites
	prodFavGroup := app.Group("/product/api/favourites")
	h.registerEndpoints(prodFavGroup)

	// Direct root: /api/favourites
	directGroup := app.Group("/api/favourites")
	h.registerEndpoints(directGroup)
}

func (h *FavouriteHandler) registerEndpoints(r fiber.Router) {
	r.Get("", h.FindAll)
	r.Get("/find", h.Find)
	r.Get("/:userId/:productId/:likeDate", h.FindByID)
	r.Post("", h.Save)
	r.Put("", h.Update)
	r.Delete("/delete", h.Delete)
	r.Delete("/:userId/:productId/:likeDate", h.DeleteByID)
}

func (h *FavouriteHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
		"components": fiber.Map{
			"db": fiber.Map{"status": "UP"},
		},
	})
}

func (h *FavouriteHandler) FindAll(c *fiber.Ctx) error {
	list, err := h.svc.FindAll(c.Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(model.FavouriteCollectionResponse{Collection: list})
}

func (h *FavouriteHandler) Find(c *fiber.Ctx) error {
	var req model.FavouriteID
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	fav, err := h.svc.FindByID(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	if fav == nil {
		return response.NotFound(c, "Favourite not found")
	}
	return c.JSON(fav)
}

func (h *FavouriteHandler) FindByID(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return response.BadRequest(c, "Invalid userId")
	}
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}
	likeDate, _ := url.QueryUnescape(c.Params("likeDate"))

	fav, err := h.svc.FindByID(c.Context(), model.FavouriteID{
		UserID:    userID,
		ProductID: productID,
		LikeDate:  likeDate,
	})
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	if fav == nil {
		return response.NotFound(c, "Favourite not found")
	}
	return c.JSON(fav)
}

func (h *FavouriteHandler) Save(c *fiber.Ctx) error {
	var req model.FavouriteDto
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	saved, err := h.svc.Save(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(saved)
}

func (h *FavouriteHandler) Update(c *fiber.Ctx) error {
	return h.Save(c)
}

func (h *FavouriteHandler) Delete(c *fiber.Ctx) error {
	var req model.FavouriteID
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	ok, err := h.svc.DeleteByID(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(ok)
}

func (h *FavouriteHandler) DeleteByID(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return response.BadRequest(c, "Invalid userId")
	}
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return response.BadRequest(c, "Invalid productId")
	}
	likeDate, _ := url.QueryUnescape(c.Params("likeDate"))

	ok, err := h.svc.DeleteByID(c.Context(), model.FavouriteID{
		UserID:    userID,
		ProductID: productID,
		LikeDate:  likeDate,
	})
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return c.JSON(ok)
}
