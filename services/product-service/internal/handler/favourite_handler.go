package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"com.ecommerce/product-service/internal/model"
	"com.ecommerce/product-service/internal/service"
	"com.ecommerce/shared/response"
	"github.com/go-chi/chi/v5"
)

type FavouriteHandler struct {
	svc *service.FavouriteService
}

func NewFavouriteHandler(svc *service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{svc: svc}
}

func (h *FavouriteHandler) RegisterRoutes(r chi.Router) {
	// Health check aliases for favourite service compatibility
	r.Get("/favourite/actuator/health", h.HealthCheck)

	// Context roots: /favourite/api/favourites, /product/api/favourites and direct root
	for _, base := range []string{"/favourite/api/favourites", "/product/api/favourites", "/api/favourites"} {
		h.registerEndpoints(r, base)
	}
}

func (h *FavouriteHandler) registerEndpoints(r chi.Router, base string) {
	r.Get(base, h.FindAll)
	r.Get(base+"/find", h.Find)
	r.Get(base+"/{userId}/{productId}/{likeDate}", h.FindByID)
	r.Post(base, h.Save)
	r.Put(base, h.Update)
	r.Delete(base+"/delete", h.Delete)
	r.Delete(base+"/{userId}/{productId}/{likeDate}", h.DeleteByID)
}

func (h *FavouriteHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

func (h *FavouriteHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.FindAll(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, model.FavouriteCollectionResponse{Collection: list})
}

func (h *FavouriteHandler) Find(w http.ResponseWriter, r *http.Request) {
	var req model.FavouriteID
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request payload")
		return
	}

	fav, err := h.svc.FindByID(r.Context(), req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	if fav == nil {
		response.NotFound(w, r, "Favourite not found")
		return
	}
	response.WriteJSON(w, http.StatusOK, fav)
}

func (h *FavouriteHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid userId")
		return
	}
	productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}
	likeDate := chi.URLParam(r, "likeDate")

	fav, err := h.svc.FindByID(r.Context(), model.FavouriteID{
		UserID:    userID,
		ProductID: productID,
		LikeDate:  likeDate,
	})
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	if fav == nil {
		response.NotFound(w, r, "Favourite not found")
		return
	}
	response.WriteJSON(w, http.StatusOK, fav)
}

func (h *FavouriteHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req model.FavouriteDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request payload")
		return
	}

	saved, err := h.svc.Save(r.Context(), req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, saved)
}

func (h *FavouriteHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.Save(w, r)
}

func (h *FavouriteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var req model.FavouriteID
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request payload")
		return
	}

	ok, err := h.svc.DeleteByID(r.Context(), req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, ok)
}

func (h *FavouriteHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid userId")
		return
	}
	productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}
	likeDate := chi.URLParam(r, "likeDate")

	ok, err := h.svc.DeleteByID(r.Context(), model.FavouriteID{
		UserID:    userID,
		ProductID: productID,
		LikeDate:  likeDate,
	})
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, ok)
}
