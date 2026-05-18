package router

import (
	"github.com/gorilla/mux"

	"github.com/makesalekz/gateway/internal/handlers"
	"github.com/makesalekz/gateway/internal/middleware"
)

func New(
	auth *handlers.AuthHandler,
	products *handlers.ProductsHandler,
	sales *handlers.SalesHandler,
	dashboard *handlers.DashboardHandler,
	agents *handlers.AgentsHandler,
	docs *handlers.DocsHandler,
) *mux.Router {
	r := mux.NewRouter()

	// Apply JWT middleware to all routes; public paths are bypassed inside the middleware.
	r.Use(middleware.JWTAuth)

	api := r.PathPrefix("/api/v1").Subrouter()

	// OpenAPI docs (public)
	api.HandleFunc("/docs", docs.ServeSpec).Methods("GET")

	// Auth (public)
	api.HandleFunc("/auth/phone", auth.AuthByPhone).Methods("POST")
	api.HandleFunc("/auth/email", auth.AuthByEmail).Methods("POST")
	api.HandleFunc("/auth/code", auth.AuthByCode).Methods("POST")
	api.HandleFunc("/auth/refresh", auth.RefreshToken).Methods("POST")

	// Cashier (POS app)
	api.HandleFunc("/sales", sales.CreateSale).Methods("POST")
	api.HandleFunc("/sales/return", sales.CreateReturn).Methods("POST")
	api.HandleFunc("/shifts/open", sales.OpenShift).Methods("POST")
	api.HandleFunc("/shifts/close", sales.CloseShift).Methods("POST")
	api.HandleFunc("/sync", sales.SyncOperations).Methods("POST")
	api.HandleFunc("/products", products.ListProducts).Methods("GET")
	api.HandleFunc("/products/search", products.SearchProducts).Methods("GET")
	api.HandleFunc("/products/{id:[0-9]+}", products.GetProduct).Methods("GET")

	// Owner (dashboard app)
	api.HandleFunc("/dashboard/sales", dashboard.GetDashboardSales).Methods("GET")
	api.HandleFunc("/dashboard/stock", dashboard.GetDashboardStock).Methods("GET")
	api.HandleFunc("/stores", dashboard.ListStores).Methods("GET")
	api.HandleFunc("/stores/{id:[0-9]+}", dashboard.GetStore).Methods("GET")
	api.HandleFunc("/shifts", dashboard.ListShifts).Methods("GET")
	api.HandleFunc("/cashier-log", dashboard.GetCashierLog).Methods("GET")
	api.HandleFunc("/billing/report", dashboard.GetBillingReport).Methods("GET")

	// Agent (SFA app)
	api.HandleFunc("/routes", agents.ListRoutes).Methods("GET")
	api.HandleFunc("/routes/{id:[0-9]+}", agents.GetRoute).Methods("GET")
	api.HandleFunc("/visits/checkin", agents.CheckIn).Methods("POST")
	api.HandleFunc("/visits/checkout", agents.CheckOut).Methods("POST")
	api.HandleFunc("/visits/photo", agents.AddVisitPhoto).Methods("POST")
	api.HandleFunc("/orders", agents.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", agents.ListOrders).Methods("GET")
	api.HandleFunc("/onboard", agents.OnboardStore).Methods("POST")

	return r
}
