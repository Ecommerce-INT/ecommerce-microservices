package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"com.ecommerce/payment-service/internal/model"
	"com.ecommerce/payment-service/internal/service"
	"com.ecommerce/shared/response"
	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) RegisterRoutes(r chi.Router) {
	r.Get("/actuator/health", h.HealthCheck)
	r.Get("/payment/actuator/health", h.HealthCheck)
	r.Get("/health", h.HealthCheck)

	for _, base := range []string{"/api/payments", "/payment/api/payments"} {
		h.registerEndpoints(r, base)
	}
}

func (h *PaymentHandler) registerEndpoints(r chi.Router, base string) {
	r.Get(base+"/getOrder/{orderId}", h.GetOrderDto)
	r.Get(base+"/all", h.FindAllPaged)
	r.Get(base+"/{paymentId}", h.FindById)
	r.Get(base, h.FindAll)
	r.Post(base, h.Save)
	r.Put(base+"/{paymentId}", h.UpdateById)
	r.Put(base, h.Update)
	r.Delete(base+"/{paymentId}", h.DeleteById)
}

func (h *PaymentHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

func queryDefault(r *http.Request, key, fallback string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return fallback
}

func (h *PaymentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	payments, err := h.svc.FindAll(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, payments)
}

func (h *PaymentHandler) FindAllPaged(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(queryDefault(r, "page", "0"))
	size, _ := strconv.Atoi(queryDefault(r, "size", "10"))
	sortBy := queryDefault(r, "sortBy", "paymentId")
	sortOrder := queryDefault(r, "sortOrder", "asc")

	res, err := h.svc.FindAllPaged(r.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, res)
}

func (h *PaymentHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "paymentId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid paymentId")
		return
	}

	payment, err := h.svc.FindById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) GetOrderDto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}

	dto := h.svc.GetOrderDto(r.Context(), id)
	response.WriteJSON(w, http.StatusOK, dto)
}

func (h *PaymentHandler) Save(w http.ResponseWriter, r *http.Request) {
	var req model.PaymentDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	saved, err := h.svc.Save(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrOrderAlreadyPaid) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, saved)
}

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.PaymentDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	updated, err := h.svc.Update(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *PaymentHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "paymentId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid paymentId")
		return
	}

	var req model.PaymentDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	updated, err := h.svc.UpdateById(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrPaymentNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *PaymentHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "paymentId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid paymentId")
		return
	}

	if err := h.svc.DeleteById(r.Context(), id); err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}
