package clients

import "context"

type OrderItem struct {
	ProductID int64  `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	Price     string `json:"price"`
}

type CreateOrderRequest struct {
	TenantID int64       `json:"-"`
	ActorID  int64       `json:"-"`
	StoreID  int64       `json:"store_id"`
	Items    []OrderItem `json:"items"`
	Notes    string      `json:"notes"`
}

type Order struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	StoreID  int64  `json:"store_id"`
	Status   string `json:"status"`
	Total    string `json:"total"`
}

type ListOrdersRequest struct {
	TenantID int64
	StoreID  int64
	Limit    int32
	FromID   int64
}

type ListOrdersResponse struct {
	Items []*Order `json:"items"`
	Total int32    `json:"total"`
}

type OrdersClient interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error)
	ListOrders(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error)
}

type ordersStub struct{}

func NewOrdersClient() OrdersClient {
	return &ordersStub{}
}

func (s *ordersStub) CreateOrder(_ context.Context, _ *CreateOrderRequest) (*Order, error) {
	return &Order{ID: 1, Status: "created"}, nil
}

func (s *ordersStub) ListOrders(_ context.Context, _ *ListOrdersRequest) (*ListOrdersResponse, error) {
	return &ListOrdersResponse{Items: []*Order{}, Total: 0}, nil
}
