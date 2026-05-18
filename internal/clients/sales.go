package clients

import "context"

type CreateSaleRequest struct {
	TenantID int64       `json:"-"`
	ActorID  int64       `json:"-"`
	StoreID  int64       `json:"store_id"`
	Items    []SaleItem  `json:"items"`
}

type SaleItem struct {
	ProductID int64  `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Price     string `json:"price"`
}

type Sale struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	StoreID  int64  `json:"store_id"`
	Total    string `json:"total"`
	Status   string `json:"status"`
}

type CreateReturnRequest struct {
	TenantID int64  `json:"-"`
	ActorID  int64  `json:"-"`
	SaleID   int64  `json:"sale_id"`
	Items    []SaleItem `json:"items"`
	Reason   string `json:"reason"`
}

type OpenShiftRequest struct {
	TenantID int64  `json:"-"`
	ActorID  int64  `json:"-"`
	StoreID  int64  `json:"store_id"`
	CashOpen string `json:"cash_open"`
}

type CloseShiftRequest struct {
	TenantID int64  `json:"-"`
	ActorID  int64  `json:"-"`
	ShiftID  int64  `json:"shift_id"`
	CashClose string `json:"cash_close"`
}

type Shift struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	StoreID  int64  `json:"store_id"`
	ActorID  int64  `json:"actor_id"`
	Status   string `json:"status"`
	CashOpen string `json:"cash_open"`
}

type SyncOperation struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type SyncRequest struct {
	TenantID   int64           `json:"-"`
	ActorID    int64           `json:"-"`
	Operations []SyncOperation `json:"operations"`
}

type SyncResponse struct {
	Processed int32    `json:"processed"`
	Failed    int32    `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
}

type ListShiftsRequest struct {
	TenantID int64
	StoreID  int64
	Limit    int32
	FromID   int64
}

type ListShiftsResponse struct {
	Items []*Shift `json:"items"`
	Total int32    `json:"total"`
}

type DashboardSalesRequest struct {
	TenantID int64
	StoreID  int64
	DateFrom string
	DateTo   string
}

type DashboardSalesResponse struct {
	TotalSales    string `json:"total_sales"`
	TotalReturns  string `json:"total_returns"`
	NetSales      string `json:"net_sales"`
	TransactionsCount int32 `json:"transactions_count"`
}

type CashierLogRequest struct {
	TenantID int64
	StoreID  int64
	ActorID  int64
	Limit    int32
	FromID   int64
}

type CashierLogEntry struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type CashierLogResponse struct {
	Items []*CashierLogEntry `json:"items"`
	Total int32              `json:"total"`
}

type SalesClient interface {
	CreateSale(ctx context.Context, req *CreateSaleRequest) (*Sale, error)
	CreateReturn(ctx context.Context, req *CreateReturnRequest) (*Sale, error)
	OpenShift(ctx context.Context, req *OpenShiftRequest) (*Shift, error)
	CloseShift(ctx context.Context, req *CloseShiftRequest) (*Shift, error)
	SyncOperations(ctx context.Context, req *SyncRequest) (*SyncResponse, error)
	ListShifts(ctx context.Context, req *ListShiftsRequest) (*ListShiftsResponse, error)
	GetDashboardSales(ctx context.Context, req *DashboardSalesRequest) (*DashboardSalesResponse, error)
	GetCashierLog(ctx context.Context, req *CashierLogRequest) (*CashierLogResponse, error)
}

type salesStub struct{}

func NewSalesClient() SalesClient {
	return &salesStub{}
}

func (s *salesStub) CreateSale(_ context.Context, _ *CreateSaleRequest) (*Sale, error) {
	return &Sale{ID: 1, Status: "completed"}, nil
}

func (s *salesStub) CreateReturn(_ context.Context, _ *CreateReturnRequest) (*Sale, error) {
	return &Sale{ID: 1, Status: "returned"}, nil
}

func (s *salesStub) OpenShift(_ context.Context, _ *OpenShiftRequest) (*Shift, error) {
	return &Shift{ID: 1, Status: "open"}, nil
}

func (s *salesStub) CloseShift(_ context.Context, _ *CloseShiftRequest) (*Shift, error) {
	return &Shift{ID: 1, Status: "closed"}, nil
}

func (s *salesStub) SyncOperations(_ context.Context, _ *SyncRequest) (*SyncResponse, error) {
	return &SyncResponse{Processed: 0}, nil
}

func (s *salesStub) ListShifts(_ context.Context, _ *ListShiftsRequest) (*ListShiftsResponse, error) {
	return &ListShiftsResponse{Items: []*Shift{}, Total: 0}, nil
}

func (s *salesStub) GetDashboardSales(_ context.Context, _ *DashboardSalesRequest) (*DashboardSalesResponse, error) {
	return &DashboardSalesResponse{}, nil
}

func (s *salesStub) GetCashierLog(_ context.Context, _ *CashierLogRequest) (*CashierLogResponse, error) {
	return &CashierLogResponse{Items: []*CashierLogEntry{}, Total: 0}, nil
}
