package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"com.ecommerce/order-service/internal/model"
	"com.ecommerce/order-service/internal/service"
	"com.ecommerce/shared/response"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Get("/actuator/health", h.HealthCheck)
	r.Get("/order/actuator/health", h.HealthCheck)
	r.Get("/health", h.HealthCheck)

	// Cart endpoints
	for _, base := range []string{"/api/carts", "/order/api/carts"} {
		h.registerCartRoutes(r, base)
	}

	// Order endpoints
	for _, base := range []string{"/api/orders", "/order/api/orders"} {
		h.registerOrderRoutes(r, base)
	}
}

func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

func queryDefault(r *http.Request, key, fallback string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return fallback
}

// -------------------------------------------------------------
// Cart Handlers
// -------------------------------------------------------------

func (h *OrderHandler) registerCartRoutes(r chi.Router, base string) {
	r.Get(base+"/all", h.FindAllCartsPaged)
	r.Get(base+"/{cartId}", h.FindCartById)
	r.Get(base, h.FindAllCarts)
	r.Post(base, h.SaveCart)
	r.Put(base+"/{cartId}", h.UpdateCartById)
	r.Put(base, h.UpdateCart)
	r.Delete(base+"/{cartId}", h.DeleteCartById)
}

func (h *OrderHandler) FindAllCarts(w http.ResponseWriter, r *http.Request) {
	carts, err := h.svc.FindAllCarts(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, carts)
}

func (h *OrderHandler) FindAllCartsPaged(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(queryDefault(r, "page", "0"))
	size, _ := strconv.Atoi(queryDefault(r, "size", "10"))
	sortBy := queryDefault(r, "sortBy", "cartId")
	sortOrder := queryDefault(r, "sortOrder", "asc")

	res, err := h.svc.FindAllCartsPaged(r.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, res)
}

func (h *OrderHandler) FindCartById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "cartId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid cartId")
		return
	}

	cart, err := h.svc.FindCartById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, cart)
}

func (h *OrderHandler) SaveCart(w http.ResponseWriter, r *http.Request) {
	var req model.CartDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	saved, err := h.svc.SaveCart(r.Context(), req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, saved)
}

func (h *OrderHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	var req model.CartDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	if req.CartId == nil {
		response.BadRequest(w, r, "cartId is required for update")
		return
	}

	updated, err := h.svc.UpdateCart(r.Context(), *req.CartId, req)
	if err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *OrderHandler) UpdateCartById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "cartId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid cartId")
		return
	}

	var req model.CartDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	updated, err := h.svc.UpdateCart(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrCartNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *OrderHandler) DeleteCartById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "cartId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid cartId")
		return
	}

	if err := h.svc.DeleteCartById(r.Context(), id); err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}

// -------------------------------------------------------------
// Order Handlers
// -------------------------------------------------------------

func (h *OrderHandler) registerOrderRoutes(r chi.Router, base string) {
	r.Get(base+"/existOrderId", h.ExistsByOrderId)
	r.Get(base+"/all", h.FindAllOrdersPaged)
	r.Get(base+"/{orderId}", h.FindOrderById)
	r.Get(base, h.FindAllOrders)
	r.Post(base, h.SaveOrder)
	r.Put(base+"/{orderId}", h.UpdateOrderById)
	r.Put(base, h.UpdateOrder)
	r.Delete(base+"/{orderId}", h.DeleteOrderById)
}

func (h *OrderHandler) FindAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.FindAllOrders(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) FindAllOrdersPaged(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(queryDefault(r, "page", "0"))
	size, _ := strconv.Atoi(queryDefault(r, "size", "10"))
	sortBy := queryDefault(r, "sortBy", "orderId")
	sortOrder := queryDefault(r, "sortOrder", "asc")

	res, err := h.svc.FindAllOrdersPaged(r.Context(), page, size, sortBy, sortOrder)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, res)
}

func (h *OrderHandler) FindOrderById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}

	order, err := h.svc.FindOrderById(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) ExistsByOrderId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}

	exists, err := h.svc.ExistsByOrderId(r.Context(), id)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, exists)
}

func (h *OrderHandler) SaveOrder(w http.ResponseWriter, r *http.Request) {
	var req model.OrderDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	saved, err := h.svc.SaveOrder(r.Context(), req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, saved)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var req model.OrderDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	if req.OrderId == nil {
		response.BadRequest(w, r, "orderId is required for update")
		return
	}

	updated, err := h.svc.UpdateOrder(r.Context(), *req.OrderId, req)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *OrderHandler) UpdateOrderById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}

	var req model.OrderDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}

	updated, err := h.svc.UpdateOrder(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *OrderHandler) DeleteOrderById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid orderId")
		return
	}

	if err := h.svc.DeleteOrderById(r.Context(), id); err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}
