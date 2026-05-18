package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/makesalekz/gateway/internal/clients"
	"github.com/makesalekz/gateway/internal/middleware"
)

type AgentsHandler struct {
	agents clients.AgentsClient
	orders clients.OrdersClient
}

func NewAgentsHandler(agents clients.AgentsClient, orders clients.OrdersClient) *AgentsHandler {
	return &AgentsHandler{agents: agents, orders: orders}
}

func (h *AgentsHandler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	req := &clients.ListRoutesRequest{
		TenantID: tenantID,
		ActorID:  actorID,
		Date:     queryParam(r, "date"),
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.agents.ListRoutes(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list routes")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AgentsHandler) GetRoute(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	id, err := pathParamInt64(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid route id")
		return
	}

	route, err := h.agents.GetRoute(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}

	writeJSON(w, http.StatusOK, route)
}

func (h *AgentsHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.CheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	visit, err := h.agents.CheckIn(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "check-in failed")
		return
	}

	writeJSON(w, http.StatusCreated, visit)
}

func (h *AgentsHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.CheckOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	visit, err := h.agents.CheckOut(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "check-out failed")
		return
	}

	writeJSON(w, http.StatusOK, visit)
}

func (h *AgentsHandler) AddVisitPhoto(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	// Max 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	visitIDStr := r.FormValue("visit_id")
	visitID, err := strconv.ParseInt(visitIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid visit_id")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "photo file is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read photo")
		return
	}

	req := &clients.AddVisitPhotoRequest{
		TenantID:  tenantID,
		ActorID:   actorID,
		VisitID:   visitID,
		PhotoData: data,
		Filename:  header.Filename,
	}

	photo, err := h.agents.AddVisitPhoto(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to upload photo")
		return
	}

	writeJSON(w, http.StatusCreated, photo)
}

func (h *AgentsHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	order, err := h.orders.CreateOrder(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *AgentsHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())

	req := &clients.ListOrdersRequest{
		TenantID: tenantID,
		StoreID:  queryParamInt64(r, "store_id"),
		Limit:    queryParamInt32(r, "limit"),
		FromID:   queryParamInt64(r, "from_id"),
	}

	resp, err := h.orders.ListOrders(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list orders")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AgentsHandler) OnboardStore(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromContext(r.Context())
	actorID := middleware.ActorIDFromContext(r.Context())

	var req clients.OnboardStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.TenantID = tenantID
	req.ActorID = actorID

	resp, err := h.agents.OnboardStore(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to onboard store")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}
