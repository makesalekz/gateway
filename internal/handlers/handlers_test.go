package handlers_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/makesalekz/gateway/internal/clients"
	"github.com/makesalekz/gateway/internal/handlers"
	"github.com/makesalekz/gateway/internal/middleware"
)

// --- Mock clients ---

type mockIAM struct {
	lastPhone string
	lastEmail string
	lastCode  string
}

func (m *mockIAM) AuthByPhone(_ context.Context, req *clients.AuthByPhoneRequest) (*clients.AuthResponse, error) {
	m.lastPhone = req.Phone
	return &clients.AuthResponse{AccessToken: "tok-" + req.Phone, RefreshToken: "ref", ExpiresIn: 3600}, nil
}

func (m *mockIAM) AuthByEmail(_ context.Context, req *clients.AuthByEmailRequest) (*clients.AuthResponse, error) {
	m.lastEmail = req.Email
	return &clients.AuthResponse{AccessToken: "tok-email", RefreshToken: "ref", ExpiresIn: 3600}, nil
}

func (m *mockIAM) AuthByCode(_ context.Context, req *clients.AuthByCodeRequest) (*clients.AuthResponse, error) {
	m.lastCode = req.Code
	return &clients.AuthResponse{AccessToken: "tok-code", RefreshToken: "ref", ExpiresIn: 3600}, nil
}

func (m *mockIAM) RefreshToken(_ context.Context, _ *clients.RefreshTokenRequest) (*clients.AuthResponse, error) {
	return &clients.AuthResponse{AccessToken: "tok-refreshed", RefreshToken: "ref2", ExpiresIn: 3600}, nil
}

type mockProducts struct {
	lastTenantID int64
	lastID       int64
}

func (m *mockProducts) ListProducts(_ context.Context, req *clients.ListProductsRequest) (*clients.ListProductsResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.ListProductsResponse{
		Items: []*clients.Product{{ID: 1, Name: "Test Product", TenantID: req.TenantID}},
		Total: 1,
	}, nil
}

func (m *mockProducts) SearchProducts(_ context.Context, req *clients.SearchProductsRequest) (*clients.ListProductsResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.ListProductsResponse{
		Items: []*clients.Product{{ID: 1, Name: "Found", TenantID: req.TenantID}},
		Total: 1,
	}, nil
}

func (m *mockProducts) GetProduct(_ context.Context, tenantID, id int64) (*clients.Product, error) {
	m.lastTenantID = tenantID
	m.lastID = id
	return &clients.Product{ID: id, TenantID: tenantID, Name: "Product"}, nil
}

type mockSales struct {
	lastTenantID int64
	lastActorID  int64
}

func (m *mockSales) CreateSale(_ context.Context, req *clients.CreateSaleRequest) (*clients.Sale, error) {
	m.lastTenantID = req.TenantID
	m.lastActorID = req.ActorID
	return &clients.Sale{ID: 1, TenantID: req.TenantID, Status: "completed"}, nil
}

func (m *mockSales) CreateReturn(_ context.Context, req *clients.CreateReturnRequest) (*clients.Sale, error) {
	m.lastTenantID = req.TenantID
	return &clients.Sale{ID: 2, Status: "returned"}, nil
}

func (m *mockSales) OpenShift(_ context.Context, req *clients.OpenShiftRequest) (*clients.Shift, error) {
	m.lastTenantID = req.TenantID
	return &clients.Shift{ID: 1, Status: "open"}, nil
}

func (m *mockSales) CloseShift(_ context.Context, req *clients.CloseShiftRequest) (*clients.Shift, error) {
	m.lastTenantID = req.TenantID
	return &clients.Shift{ID: 1, Status: "closed"}, nil
}

func (m *mockSales) SyncOperations(_ context.Context, req *clients.SyncRequest) (*clients.SyncResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.SyncResponse{Processed: int32(len(req.Operations))}, nil
}

func (m *mockSales) ListShifts(_ context.Context, req *clients.ListShiftsRequest) (*clients.ListShiftsResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.ListShiftsResponse{Items: []*clients.Shift{}, Total: 0}, nil
}

func (m *mockSales) GetDashboardSales(_ context.Context, req *clients.DashboardSalesRequest) (*clients.DashboardSalesResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.DashboardSalesResponse{TotalSales: "1000"}, nil
}

func (m *mockSales) GetCashierLog(_ context.Context, req *clients.CashierLogRequest) (*clients.CashierLogResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.CashierLogResponse{Items: []*clients.CashierLogEntry{}, Total: 0}, nil
}

type mockWarehouse struct{}

func (m *mockWarehouse) GetStockItems(_ context.Context, _ *clients.GetStockRequest) (*clients.GetStockResponse, error) {
	return &clients.GetStockResponse{Items: []*clients.StockItem{}, Total: 0}, nil
}

type mockOrders struct {
	lastTenantID int64
}

func (m *mockOrders) CreateOrder(_ context.Context, req *clients.CreateOrderRequest) (*clients.Order, error) {
	m.lastTenantID = req.TenantID
	return &clients.Order{ID: 1, Status: "created"}, nil
}

func (m *mockOrders) ListOrders(_ context.Context, req *clients.ListOrdersRequest) (*clients.ListOrdersResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.ListOrdersResponse{Items: []*clients.Order{}, Total: 0}, nil
}

type mockAgents struct {
	lastTenantID int64
}

func (m *mockAgents) ListRoutes(_ context.Context, req *clients.ListRoutesRequest) (*clients.ListRoutesResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.ListRoutesResponse{Items: []*clients.Route{}, Total: 0}, nil
}

func (m *mockAgents) GetRoute(_ context.Context, tenantID, id int64) (*clients.Route, error) {
	m.lastTenantID = tenantID
	return &clients.Route{ID: id, TenantID: tenantID, Name: "Route 1"}, nil
}

func (m *mockAgents) CheckIn(_ context.Context, req *clients.CheckInRequest) (*clients.Visit, error) {
	m.lastTenantID = req.TenantID
	return &clients.Visit{ID: 1, Status: "checked_in"}, nil
}

func (m *mockAgents) CheckOut(_ context.Context, req *clients.CheckOutRequest) (*clients.Visit, error) {
	m.lastTenantID = req.TenantID
	return &clients.Visit{ID: 1, Status: "checked_out"}, nil
}

func (m *mockAgents) AddVisitPhoto(_ context.Context, req *clients.AddVisitPhotoRequest) (*clients.VisitPhoto, error) {
	m.lastTenantID = req.TenantID
	return &clients.VisitPhoto{ID: 1, VisitID: req.VisitID, URL: "https://example.com/photo.jpg"}, nil
}

func (m *mockAgents) OnboardStore(_ context.Context, req *clients.OnboardStoreRequest) (*clients.OnboardStoreResponse, error) {
	m.lastTenantID = req.TenantID
	return &clients.OnboardStoreResponse{StoreID: 1, Status: "created"}, nil
}

type mockStores struct{}

func (m *mockStores) ListStores(_ context.Context, _ *clients.ListStoresRequest) (*clients.ListStoresResponse, error) {
	return &clients.ListStoresResponse{Items: []*clients.Store{}, Total: 0}, nil
}

func (m *mockStores) GetStore(_ context.Context, tenantID, id int64) (*clients.Store, error) {
	return &clients.Store{ID: id, TenantID: tenantID, Name: "Store"}, nil
}

type mockBilling struct{}

func (m *mockBilling) GetBillingReport(_ context.Context, _ *clients.GetBillingReportRequest) (*clients.BillingReport, error) {
	return &clients.BillingReport{Status: "ok"}, nil
}

// --- Test helpers ---

func buildTestJWT(claims map[string]interface{}) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(claims)
	payloadEnc := base64.RawURLEncoding.EncodeToString(payload)
	sig := base64.RawURLEncoding.EncodeToString([]byte("stub-signature"))
	return fmt.Sprintf("%s.%s.%s", header, payloadEnc, sig)
}

// buildRouter creates a mux.Router with all routes wired to mock clients.
// This avoids importing the router package (which would create an import cycle).
func buildRouter(
	authH *handlers.AuthHandler,
	productsH *handlers.ProductsHandler,
	salesH *handlers.SalesHandler,
	dashboardH *handlers.DashboardHandler,
	agentsH *handlers.AgentsHandler,
) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.JWTAuth)

	api := r.PathPrefix("/api/v1").Subrouter()

	// Docs (public)
	docsH := handlers.NewDocsHandler([]byte(`openapi: "3.0.3"`))
	api.HandleFunc("/docs", docsH.ServeSpec).Methods("GET")

	// Auth (public)
	api.HandleFunc("/auth/phone", authH.AuthByPhone).Methods("POST")
	api.HandleFunc("/auth/email", authH.AuthByEmail).Methods("POST")
	api.HandleFunc("/auth/code", authH.AuthByCode).Methods("POST")
	api.HandleFunc("/auth/refresh", authH.RefreshToken).Methods("POST")

	// Cashier
	api.HandleFunc("/sales", salesH.CreateSale).Methods("POST")
	api.HandleFunc("/sales/return", salesH.CreateReturn).Methods("POST")
	api.HandleFunc("/shifts/open", salesH.OpenShift).Methods("POST")
	api.HandleFunc("/shifts/close", salesH.CloseShift).Methods("POST")
	api.HandleFunc("/sync", salesH.SyncOperations).Methods("POST")
	api.HandleFunc("/products", productsH.ListProducts).Methods("GET")
	api.HandleFunc("/products/search", productsH.SearchProducts).Methods("GET")
	api.HandleFunc("/products/{id:[0-9]+}", productsH.GetProduct).Methods("GET")

	// Owner
	api.HandleFunc("/dashboard/sales", dashboardH.GetDashboardSales).Methods("GET")
	api.HandleFunc("/dashboard/stock", dashboardH.GetDashboardStock).Methods("GET")
	api.HandleFunc("/stores", dashboardH.ListStores).Methods("GET")
	api.HandleFunc("/stores/{id:[0-9]+}", dashboardH.GetStore).Methods("GET")
	api.HandleFunc("/shifts", dashboardH.ListShifts).Methods("GET")
	api.HandleFunc("/cashier-log", dashboardH.GetCashierLog).Methods("GET")
	api.HandleFunc("/billing/report", dashboardH.GetBillingReport).Methods("GET")

	// Agent
	api.HandleFunc("/routes", agentsH.ListRoutes).Methods("GET")
	api.HandleFunc("/routes/{id:[0-9]+}", agentsH.GetRoute).Methods("GET")
	api.HandleFunc("/visits/checkin", agentsH.CheckIn).Methods("POST")
	api.HandleFunc("/visits/checkout", agentsH.CheckOut).Methods("POST")
	api.HandleFunc("/visits/photo", agentsH.AddVisitPhoto).Methods("POST")
	api.HandleFunc("/orders", agentsH.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", agentsH.ListOrders).Methods("GET")
	api.HandleFunc("/onboard", agentsH.OnboardStore).Methods("POST")

	return r
}

func setupRouter() (*mux.Router, *mockIAM, *mockProducts, *mockSales, *mockAgents, *mockOrders) {
	iam := &mockIAM{}
	products := &mockProducts{}
	sales := &mockSales{}
	warehouse := &mockWarehouse{}
	orders := &mockOrders{}
	agents := &mockAgents{}
	stores := &mockStores{}
	billing := &mockBilling{}

	authH := handlers.NewAuthHandler(iam)
	productsH := handlers.NewProductsHandler(products)
	salesH := handlers.NewSalesHandler(sales)
	dashboardH := handlers.NewDashboardHandler(sales, warehouse, stores, billing)
	agentsH := handlers.NewAgentsHandler(agents, orders)

	r := buildRouter(authH, productsH, salesH, dashboardH, agentsH)
	return r, iam, products, sales, agents, orders
}

func authedRequest(method, path string, body io.Reader) *http.Request {
	token := buildTestJWT(map[string]interface{}{
		"tenantId": float64(10),
		"memberId": float64(42),
	})
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// --- Auth handler tests ---

func TestAuthByPhone(t *testing.T) {
	r, iam, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"phone":"+77001234567"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/phone", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "+77001234567", iam.lastPhone)

	var resp clients.AuthResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "tok-+77001234567", resp.AccessToken)
}

func TestAuthByPhone_MissingPhone(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/phone", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthByEmail(t *testing.T) {
	r, iam, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"email":"test@example.com","password":"secret"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/email", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test@example.com", iam.lastEmail)
}

func TestAuthByCode(t *testing.T) {
	r, iam, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"phone":"+77001234567","code":"777333"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/code", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "777333", iam.lastCode)
}

func TestAuthByCode_MissingFields(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"phone":"+77001234567"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/code", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRefreshToken(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"refresh_token":"old-token"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp clients.AuthResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "tok-refreshed", resp.AccessToken)
}

func TestRefreshToken_MissingField(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// --- Products handler tests ---

func TestListProducts(t *testing.T) {
	r, _, products, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/products?limit=10&from_id=5", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), products.lastTenantID)

	var resp clients.ListProductsResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, int32(1), resp.Total)
}

func TestListProducts_Unauthorized(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSearchProducts(t *testing.T) {
	r, _, products, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/products/search?q=milk&limit=20", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), products.lastTenantID)
}

func TestGetProduct(t *testing.T) {
	r, _, products, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/products/42", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), products.lastTenantID)
	assert.Equal(t, int64(42), products.lastID)
}

func TestGetProduct_InvalidID(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/products/abc", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	// gorilla/mux won't match non-numeric due to [0-9]+ constraint
	assert.True(t, rec.Code == http.StatusNotFound || rec.Code == http.StatusMethodNotAllowed)
}

// --- Sales handler tests ---

func TestCreateSale(t *testing.T) {
	r, _, _, sales, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"store_id":1,"items":[{"product_id":1,"quantity":2,"price":"100"}]}`)
	req := authedRequest("POST", "/api/v1/sales", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, int64(10), sales.lastTenantID)
	assert.Equal(t, int64(42), sales.lastActorID)

	var sale clients.Sale
	json.NewDecoder(rec.Body).Decode(&sale)
	assert.Equal(t, "completed", sale.Status)
}

func TestCreateReturn(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"sale_id":1,"items":[],"reason":"defective"}`)
	req := authedRequest("POST", "/api/v1/sales/return", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestOpenShift(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"store_id":1,"cash_open":"5000"}`)
	req := authedRequest("POST", "/api/v1/shifts/open", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCloseShift(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"shift_id":1,"cash_close":"15000"}`)
	req := authedRequest("POST", "/api/v1/shifts/close", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestSyncOperations(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"operations":[{"type":"sale","payload":{}}]}`)
	req := authedRequest("POST", "/api/v1/sync", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp clients.SyncResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, int32(1), resp.Processed)
}

// --- Dashboard handler tests ---

func TestGetDashboardSales(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/dashboard/sales?date_from=2026-01-01&date_to=2026-01-31", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp clients.DashboardSalesResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "1000", resp.TotalSales)
}

func TestGetDashboardStock(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/dashboard/stock", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestListStores(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/stores", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetStore(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/stores/5", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var store clients.Store
	json.NewDecoder(rec.Body).Decode(&store)
	assert.Equal(t, int64(5), store.ID)
}

func TestListShifts(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/shifts?store_id=1", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetCashierLog(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/cashier-log?store_id=1", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetBillingReport(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/billing/report?period=2026-01", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp clients.BillingReport
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "ok", resp.Status)
}

// --- Agent handler tests ---

func TestListRoutes(t *testing.T) {
	r, _, _, _, agents, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/routes?date=2026-05-18", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), agents.lastTenantID)
}

func TestGetRoute(t *testing.T) {
	r, _, _, _, agents, _ := setupRouter()

	req := authedRequest("GET", "/api/v1/routes/7", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), agents.lastTenantID)
}

func TestCheckIn(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"route_id":1,"store_id":2,"latitude":43.238,"longitude":76.945}`)
	req := authedRequest("POST", "/api/v1/visits/checkin", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCheckOut(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"visit_id":1,"latitude":43.238,"longitude":76.945,"notes":"all good"}`)
	req := authedRequest("POST", "/api/v1/visits/checkout", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAddVisitPhoto(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	// Build multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("visit_id", "5")
	part, _ := writer.CreateFormFile("photo", "test.jpg")
	part.Write([]byte("fake-image-data"))
	writer.Close()

	token := buildTestJWT(map[string]interface{}{
		"tenantId": float64(10),
		"memberId": float64(42),
	})

	req := httptest.NewRequest("POST", "/api/v1/visits/photo", &buf)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var photo clients.VisitPhoto
	json.NewDecoder(rec.Body).Decode(&photo)
	assert.Equal(t, int64(5), photo.VisitID)
}

func TestCreateOrder(t *testing.T) {
	r, _, _, _, _, orders := setupRouter()

	body := bytes.NewBufferString(`{"store_id":1,"items":[{"product_id":1,"quantity":10,"price":"500"}],"notes":"urgent"}`)
	req := authedRequest("POST", "/api/v1/orders", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, int64(10), orders.lastTenantID)
}

func TestListOrders(t *testing.T) {
	r, _, _, _, _, orders := setupRouter()

	req := authedRequest("GET", "/api/v1/orders?store_id=1", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, int64(10), orders.lastTenantID)
}

func TestOnboardStore(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{"name":"New Shop","address":"123 Main St","phone":"+77001234567"}`)
	req := authedRequest("POST", "/api/v1/onboard", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp clients.OnboardStoreResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "created", resp.Status)
}

// --- Error mapping tests ---

func TestProtectedRoute_NoToken_Returns401(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/products"},
		{"GET", "/api/v1/products/1"},
		{"POST", "/api/v1/sales"},
		{"GET", "/api/v1/dashboard/sales"},
		{"GET", "/api/v1/routes"},
		{"POST", "/api/v1/orders"},
	}

	r, _, _, _, _, _ := setupRouter()

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}

func TestDocsEndpoint_NoAuth(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	req := httptest.NewRequest("GET", "/api/v1/docs", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "openapi")
}

func TestInvalidJSON_Returns400(t *testing.T) {
	r, _, _, _, _, _ := setupRouter()

	body := bytes.NewBufferString(`{invalid json}`)
	req := authedRequest("POST", "/api/v1/sales", body)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
