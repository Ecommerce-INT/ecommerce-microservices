package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"com.ecommerce/product-service/internal/model"
	"com.ecommerce/product-service/internal/service"
	"com.ecommerce/shared/response"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.Get("/actuator/health", h.HealthCheck)
	r.Get("/product/actuator/health", h.HealthCheck)

	// Context root: /product and direct root
	for _, prefix := range []string{"/product", ""} {
		h.registerCategoryRoutes(r, prefix+"/api/categories")
		h.registerProductRoutes(r, prefix+"/api/products")
	}
}

func (h *ProductHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Health(w, "db")
}

// =================== Categories ===================

func (h *ProductHandler) registerCategoryRoutes(r chi.Router, base string) {
	r.Get(base, h.ListCategories)
	r.Get(base+"/paging", h.ListCategories)
	r.Get(base+"/paging-and-sorting", h.ListCategories)
	r.Get(base+"/{categoryId}", h.GetCategoryByID)
	r.Post(base, h.CreateCategory)
	r.Put(base, h.UpdateCategory)
	r.Put(base+"/{categoryId}", h.UpdateCategoryByID)
	r.Delete(base+"/{categoryId}", h.DeleteCategory)
}

func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.FindAllCategories(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, list)
}

func (h *ProductHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "categoryId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid categoryId")
		return
	}
	cat, err := h.svc.FindCategoryByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, cat)
}

func (h *ProductHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req model.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	created, err := h.svc.SaveCategory(r.Context(), &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusCreated, created)
}

func (h *ProductHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req model.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	updated, err := h.svc.UpdateCategory(r.Context(), req.CategoryId, &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *ProductHandler) UpdateCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "categoryId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid categoryId")
		return
	}
	var req model.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	updated, err := h.svc.UpdateCategory(r.Context(), id, &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *ProductHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "categoryId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid categoryId")
		return
	}
	err = h.svc.DeleteCategory(r.Context(), id)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}

// =================== Products ===================

func (h *ProductHandler) registerProductRoutes(r chi.Router, base string) {
	r.Get(base, h.ListProducts)
	r.Get(base+"/{productId}", h.GetProductByID)
	r.Post(base, h.CreateProduct)
	r.Put(base, h.UpdateProduct)
	r.Put(base+"/{productId}", h.UpdateProductByID)
	r.Delete(base+"/{productId}", h.DeleteProduct)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.FindAllProducts(r.Context())
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, list)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}
	p, err := h.svc.FindProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w, r, err.Error())
			return
		}
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req model.ProductDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	created, err := h.svc.SaveProduct(r.Context(), &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusCreated, created)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var req model.ProductDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	updated, err := h.svc.UpdateProduct(r.Context(), req.ProductId, &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *ProductHandler) UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}
	var req model.ProductDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, "Invalid request body")
		return
	}
	updated, err := h.svc.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, updated)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "productId"))
	if err != nil {
		response.BadRequest(w, r, "Invalid productId")
		return
	}
	err = h.svc.DeleteProduct(r.Context(), id)
	if err != nil {
		response.InternalError(w, r, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, true)
}
