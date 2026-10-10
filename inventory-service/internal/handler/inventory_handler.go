package handler

import (
	"net/http"
	"strings"

	"com.ecommerce/inventory-service/internal/service"
	"com.ecommerce/pkg/common/response"
	"github.com/go-chi/chi/v5"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) RegisterRoutes(r chi.Router) {
	r.Get("/actuator/health", h.HealthCheck)
	r.Get("/inventory/actuator/health", h.HealthCheck)

	// Context root: /inventory and direct root
	for _, base := range []string{"/inventory/api/inventory", "/api/inventory"} {
		h.registerEndpoints(r, base)
	}
}

func (h *InventoryHandler) registerEndpoints(r chi.Router, base string) {
	r.Get(base, h.IsInStock)
	r.Get(base+"/{products}", h.IsInStockPath)
}

func (h *InventoryHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

func splitProductNames(values []string) []string {
	var names []string
	for _, value := range values {
		parts := strings.Split(value, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				names = append(names, trimmed)
			}
		}
	}
	return names
}

func (h *InventoryHandler) IsInStock(w http.ResponseWriter, r *http.Request) {
	// Support both ?productName=a&productName=b and ?productName=a,b
	names := splitProductNames(r.URL.Query()["productName"])

	result, err := h.svc.IsInStock(r.Context(), names)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) IsInStockPath(w http.ResponseWriter, r *http.Request) {
	names := splitProductNames([]string{chi.URLParam(r, "products")})

	result, err := h.svc.IsInStock(r.Context(), names)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, result)
}
