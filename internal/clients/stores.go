package clients

import "context"

type Store struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type ListStoresRequest struct {
	TenantID int64
	Limit    int32
	FromID   int64
}

type ListStoresResponse struct {
	Items []*Store `json:"items"`
	Total int32    `json:"total"`
}

type StoresClient interface {
	ListStores(ctx context.Context, req *ListStoresRequest) (*ListStoresResponse, error)
	GetStore(ctx context.Context, tenantID, id int64) (*Store, error)
}

type storesStub struct{}

func NewStoresClient() StoresClient {
	return &storesStub{}
}

func (s *storesStub) ListStores(_ context.Context, _ *ListStoresRequest) (*ListStoresResponse, error) {
	return &ListStoresResponse{Items: []*Store{}, Total: 0}, nil
}

func (s *storesStub) GetStore(_ context.Context, _, _ int64) (*Store, error) {
	return &Store{}, nil
}
