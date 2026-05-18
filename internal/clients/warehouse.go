package clients

import "context"

type StockItem struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	WarehouseID int64  `json:"warehouse_id"`
	Quantity    int32  `json:"quantity"`
}

type GetStockRequest struct {
	TenantID    int64
	WarehouseID int64
	Limit       int32
	FromID      int64
}

type GetStockResponse struct {
	Items []*StockItem `json:"items"`
	Total int32        `json:"total"`
}

type WarehouseClient interface {
	GetStockItems(ctx context.Context, req *GetStockRequest) (*GetStockResponse, error)
}

type warehouseStub struct{}

func NewWarehouseClient() WarehouseClient {
	return &warehouseStub{}
}

func (s *warehouseStub) GetStockItems(_ context.Context, _ *GetStockRequest) (*GetStockResponse, error) {
	return &GetStockResponse{Items: []*StockItem{}, Total: 0}, nil
}
