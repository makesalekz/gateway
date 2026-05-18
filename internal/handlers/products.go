package handlers

import (
	"net/http"

	"github.com/makesalekz/gateway/internal/clients"
	"github.com/makesalekz/gateway/internal/middleware"
)

type ProductsHandler struct {
	products clients.ProductsClient
}

func NewProductsHandler(products clients.ProductsClient) *ProductsHandler {
	return &ProductsHandler{products: products}
}

func (h *ProductsHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.ListProductsRequest{
		TenantID: tenantID,
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.products.ListProducts(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list products")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProductsHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.SearchProductsRequest{
		TenantID: tenantID,
		Query:    queryParam(r, "q"),
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.products.SearchProducts(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to search products")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProductsHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	id, err := pathParamInt64(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.products.GetProduct(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}

	writeJSON(w, http.StatusOK, product)
}
