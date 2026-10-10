package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"com.ecommerce/pkg/common/response"
	"com.ecommerce/shipping-service/internal/model"
	"com.ecommerce/shipping-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type ShippingHandler struct {
	svc *service.ShippingService
}

func NewShippingHandler(svc *service.ShippingService) *ShippingHandler {
	return &ShippingHandler{svc: svc}
}

func (h *ShippingHandler) RegisterRoutes(r chi.Router) {
	// Health checks
	r.Get("/actuator/health", h.HealthCheck)
	r.Get("/shipping/actuator/health", h.HealthCheck)

	// Context root: /shipping and fallback root
	for _, base := range []string{"/shipping/api/shippings", "/api/shippings"} {
		h.registerEndpoints(r, base)
	}
}

func (h *ShippingHandler) registerEndpoints(r chi.Router, base string) {
	r.Get(base, h.FindAll)
	r.Get(base+"/find", h.FindByBody)
	r.Get(base+"/{orderId}/{productId}", h.FindByID)
	r.Post(base, h.Save)
	r.Put(base, h.Update)
	r.Delete(base+"/delete", h.DeleteByBody)
	r.Delete(base+"/{orderId}/{productId}", h.DeleteByID)
}

func (h *ShippingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

func (h *ShippingHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.FindAll(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, model.DtoCollectionResponse[model.OrderItemDto]{
		Collection: items,
	})
}

func (h *ShippingHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}
	productId, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}

	item, err := h.svc.FindByID(r.Context(), orderId, productId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, item)
}

func (h *ShippingHandler) FindByBody(w http.ResponseWriter, r *http.Request) {
	var req model.OrderItemId
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	item, err := h.svc.FindByID(r.Context(), req.OrderId, req.ProductId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, item)
}

func (h *ShippingHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req model.OrderItemDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	saved, err := h.svc.Save(r.Context(), &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, saved)
}

func (h *ShippingHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.OrderItemDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	updated, err := h.svc.Update(r.Context(), &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *ShippingHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}
	productId, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}

	err = h.svc.DeleteByID(r.Context(), orderId, productId)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}

func (h *ShippingHandler) DeleteByBody(w http.ResponseWriter, r *http.Request) {
	var req model.OrderItemId
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	err := h.svc.DeleteByID(r.Context(), req.OrderId, req.ProductId)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}
