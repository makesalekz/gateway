# CLAUDE.md — Gateway BFF

This is the REST API gateway for UMAG mobile apps. It proxies REST requests to backend gRPC services.

## Project Documentation

@docs/index.md
@docs/api-contracts.md
@docs/data-models.md
@docs/integration-architecture.md
@docs/development-guide.md
@docs/project-overview.md

## Architecture

- Go + gorilla/mux (REST, not gRPC)
- JWT auth middleware (validates tokens from IAM)
- Stub gRPC clients (replace with real ones as services become available)

## Backend Services (gRPC)

All backend services use Go/Kratos/Ent/PostgreSQL/NATS/gRPC pattern.

### IAM (:9000)
Auth: AuthByPhone, AuthByEmail, AuthByCode, RefreshToken
Users: GetUser, GetUsers, ListUsers, UpdateOwnProfile, DeleteOwnProfile
Settings: GetSettings, UpdateSettings
Privacy: GetPrivacy, UpdatePrivacy

### Products (:9000)
CRUD: CreateProduct, UpdateProduct, DeleteProduct, GetProduct, ListProducts
Barcodes: AddBarcode, RemoveBarcode, GetProductByBarcode, SearchProducts
Categories: CreateCategory, UpdateCategory, DeleteCategory, ListCategories
Pricing: SetPrice, GetPriceHistory
Import: ImportProducts, ImportFromUMAG

### Warehouse (:9000)
Warehouses: CreateWarehouse, GetWarehouse, ListWarehouses
Stock: GetStockItems, GetStockByProduct, SetMinQuantity, GetLowStockItems
Movements: CreateReceipt, CreateTransfer, CreateWriteOff, CreateGift, ListMovements
Inventory: StartInventory, SetInventoryItem, CompleteInventory, GetInventory

### Sales (:9000)
Sales: CreateSale (publishes NATS "sales.sale.completed")
Returns: CreateReturn (publishes NATS "sales.return.completed")
Shifts: OpenShift, CloseShift, GetShift, ListShifts
Sync: SyncOperations (batch offline sync with UUID idempotency)
CashierLog: GetCashierLog

### Orders (:9000)
CRUD: CreateOrder, GetOrder, ListOrders
Status: ConfirmOrder, ShipOrder, AcceptOrder, RejectOrder
Suggestions: GetOrderSuggestions

### Agents (:9000)
Routes: CreateRoute, GetRoute, ListRoutes
Visits: CheckIn (GPS), CheckOut, ListVisits
Photos: AddVisitPhoto, GetVisitPhotos
Onboard: OnboardStore
Reports: GetAgentReport

### Stores (:9000)
CRUD: CreateStore, UpdateStore, DeleteStore, GetStore, ListStores
SetStoreResponsible, GetStoresByCoordinates

### Platform Billing (:9000)
Commissions: GetBillingReport, ClosePeriod
Exclusions: CreateExclusion, ListExclusions, DeleteExclusion

## Conventions

- All routes prefixed with /api/v1/
- JWT required (except /api/v1/auth/* and /api/v1/docs)
- tenant_id and actor_id extracted from JWT claims into context
- JSON request/response
- Error format: {"error": "message", "code": "ERROR_CODE"}
- Pagination: ?page=1&limit=20 or cursor-based ?from_id=123

## When implementing new endpoints

1. Add route in internal/router/router.go
2. Create handler in internal/handlers/
3. Use client interface from internal/clients/
4. Write test
5. Update api/openapi.yaml
