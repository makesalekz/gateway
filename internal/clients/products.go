package clients

import "context"

type Product struct {
	ID            int64  `json:"id"`
	TenantID      int64  `json:"tenant_id"`
	Name          string `json:"name"`
	Barcode       string `json:"barcode"`
	CategoryID    int64  `json:"category_id"`
	Unit          string `json:"unit"`
	PurchasePrice string `json:"purchase_price"`
	SellingPrice  string `json:"selling_price"`
	Description   string `json:"description"`
	Sku           string `json:"sku"`
}

type ListProductsRequest struct {
	TenantID int64
	Limit    int32
	FromID   int64
}

type ListProductsResponse struct {
	Items []*Product `json:"items"`
	Total int32      `json:"total"`
}

type SearchProductsRequest struct {
	TenantID int64
	Query    string
	Limit    int32
	FromID   int64
}

type ProductsClient interface {
	ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error)
	SearchProducts(ctx context.Context, req *SearchProductsRequest) (*ListProductsResponse, error)
	GetProduct(ctx context.Context, tenantID, id int64) (*Product, error)
}

type productsStub struct{}

func NewProductsClient() ProductsClient {
	return &productsStub{}
}

func (s *productsStub) ListProducts(_ context.Context, _ *ListProductsRequest) (*ListProductsResponse, error) {
	return &ListProductsResponse{Items: []*Product{}, Total: 0}, nil
}

func (s *productsStub) SearchProducts(_ context.Context, _ *SearchProductsRequest) (*ListProductsResponse, error) {
	return &ListProductsResponse{Items: []*Product{}, Total: 0}, nil
}

func (s *productsStub) GetProduct(_ context.Context, _, _ int64) (*Product, error) {
	return &Product{}, nil
}
