package clients

import "context"

type Route struct {
	ID       int64    `json:"id"`
	TenantID int64    `json:"tenant_id"`
	Name     string   `json:"name"`
	Date     string   `json:"date"`
	StoreIDs []int64  `json:"store_ids"`
	Status   string   `json:"status"`
}

type ListRoutesRequest struct {
	TenantID int64
	ActorID  int64
	Date     string
	Limit    int32
	FromID   int64
}

type ListRoutesResponse struct {
	Items []*Route `json:"items"`
	Total int32    `json:"total"`
}

type CheckInRequest struct {
	TenantID  int64   `json:"-"`
	ActorID   int64   `json:"-"`
	RouteID   int64   `json:"route_id"`
	StoreID   int64   `json:"store_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CheckOutRequest struct {
	TenantID  int64   `json:"-"`
	ActorID   int64   `json:"-"`
	VisitID   int64   `json:"visit_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Notes     string  `json:"notes"`
}

type Visit struct {
	ID        int64  `json:"id"`
	RouteID   int64  `json:"route_id"`
	StoreID   int64  `json:"store_id"`
	Status    string `json:"status"`
	CheckInAt string `json:"check_in_at,omitempty"`
	CheckOutAt string `json:"check_out_at,omitempty"`
}

type AddVisitPhotoRequest struct {
	TenantID int64  `json:"-"`
	ActorID  int64  `json:"-"`
	VisitID  int64
	PhotoData []byte
	Filename  string
}

type VisitPhoto struct {
	ID      int64  `json:"id"`
	VisitID int64  `json:"visit_id"`
	URL     string `json:"url"`
}

type OnboardStoreRequest struct {
	TenantID  int64   `json:"-"`
	ActorID   int64   `json:"-"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type OnboardStoreResponse struct {
	StoreID int64  `json:"store_id"`
	Status  string `json:"status"`
}

type AgentsClient interface {
	ListRoutes(ctx context.Context, req *ListRoutesRequest) (*ListRoutesResponse, error)
	GetRoute(ctx context.Context, tenantID, id int64) (*Route, error)
	CheckIn(ctx context.Context, req *CheckInRequest) (*Visit, error)
	CheckOut(ctx context.Context, req *CheckOutRequest) (*Visit, error)
	AddVisitPhoto(ctx context.Context, req *AddVisitPhotoRequest) (*VisitPhoto, error)
	OnboardStore(ctx context.Context, req *OnboardStoreRequest) (*OnboardStoreResponse, error)
}

type agentsStub struct{}

func NewAgentsClient() AgentsClient {
	return &agentsStub{}
}

func (s *agentsStub) ListRoutes(_ context.Context, _ *ListRoutesRequest) (*ListRoutesResponse, error) {
	return &ListRoutesResponse{Items: []*Route{}, Total: 0}, nil
}

func (s *agentsStub) GetRoute(_ context.Context, _, _ int64) (*Route, error) {
	return &Route{}, nil
}

func (s *agentsStub) CheckIn(_ context.Context, _ *CheckInRequest) (*Visit, error) {
	return &Visit{ID: 1, Status: "checked_in"}, nil
}

func (s *agentsStub) CheckOut(_ context.Context, _ *CheckOutRequest) (*Visit, error) {
	return &Visit{ID: 1, Status: "checked_out"}, nil
}

func (s *agentsStub) AddVisitPhoto(_ context.Context, _ *AddVisitPhotoRequest) (*VisitPhoto, error) {
	return &VisitPhoto{ID: 1, URL: "https://stub.example.com/photo.jpg"}, nil
}

func (s *agentsStub) OnboardStore(_ context.Context, _ *OnboardStoreRequest) (*OnboardStoreResponse, error) {
	return &OnboardStoreResponse{StoreID: 1, Status: "created"}, nil
}
