package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/makesalekz/gateway/internal/clients"
	"github.com/makesalekz/gateway/internal/middleware"
)

type SalesHandler struct {
	sales clients.SalesClient
}

func NewSalesHandler(sales clients.SalesClient) *SalesHandler {
	return &SalesHandler{sales: sales}
}

func (h *SalesHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	sale, err := h.sales.CreateSale(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create sale")
		return
	}

	writeJSON(w, http.StatusCreated, sale)
}

func (h *SalesHandler) CreateReturn(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.CreateReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	sale, err := h.sales.CreateReturn(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create return")
		return
	}

	writeJSON(w, http.StatusCreated, sale)
}

func (h *SalesHandler) OpenShift(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.OpenShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	shift, err := h.sales.OpenShift(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to open shift")
		return
	}

	writeJSON(w, http.StatusCreated, shift)
}

func (h *SalesHandler) CloseShift(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.CloseShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	shift, err := h.sales.CloseShift(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to close shift")
		return
	}

	writeJSON(w, http.StatusOK, shift)
}

func (h *SalesHandler) SyncOperations(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	resp, err := h.sales.SyncOperations(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "sync failed")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
