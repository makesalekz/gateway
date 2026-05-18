package handlers

import (
	"net/http"

	"github.com/makesalekz/gateway/internal/clients"
	"github.com/makesalekz/gateway/internal/middleware"
)

type DashboardHandler struct {
	sales     clients.SalesClient
	warehouse clients.WarehouseClient
	stores    clients.StoresClient
	billing   clients.BillingClient
}

func NewDashboardHandler(
	sales clients.SalesClient,
	warehouse clients.WarehouseClient,
	stores clients.StoresClient,
	billing clients.BillingClient,
) *DashboardHandler {
	return &DashboardHandler{
		sales:     sales,
		warehouse: warehouse,
		stores:    stores,
		billing:   billing,
	}
}

func (h *DashboardHandler) GetDashboardSales(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.DashboardSalesRequest{
		TenantID: tenantID,
		StoreID:  queryParamInt64(r, "store_id"),
		DateFrom: queryParam(r, "date_from"),
		DateTo:   queryParam(r, "date_to"),
	}

	resp, err := h.sales.GetDashboardSales(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get dashboard sales")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DashboardHandler) GetDashboardStock(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.GetStockRequest{
		TenantID:    tenantID,
		WarehouseID: queryParamInt64(r, "warehouse_id"),
		Limit:       queryParamInt32(r, "limit"),
		FromID:      queryParamInt64(r, "from_id"),
	}

	resp, err := h.warehouse.GetStockItems(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get stock")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DashboardHandler) ListStores(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.ListStoresRequest{
		TenantID: tenantID,
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.stores.ListStores(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list stores")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DashboardHandler) GetStore(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	id, err := pathParamInt64(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid store id")
		return
	}

	store, err := h.stores.GetStore(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "store not found")
		return
	}

	writeJSON(w, http.StatusOK, store)
}

func (h *DashboardHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.ListShiftsRequest{
		TenantID: tenantID,
		StoreID:  queryParamInt64(r, "store_id"),
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.sales.ListShifts(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list shifts")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DashboardHandler) GetCashierLog(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.CashierLogRequest{
		TenantID: tenantID,
		StoreID:  queryParamInt64(r, "store_id"),
		ActorID:  queryParamInt64(r, "actor_id"),
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.sales.GetCashierLog(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get cashier log")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DashboardHandler) GetBillingReport(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.GetBillingReportRequest{
		TenantID: tenantID,
		Period:   queryParam(r, "period"),
	}

	resp, err := h.billing.GetBillingReport(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get billing report")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
