package clients

import "context"

type BillingReport struct {
	TenantID     int64  `json:"tenant_id"`
	Period       string `json:"period"`
	TotalRevenue string `json:"total_revenue"`
	Commission   string `json:"commission"`
	Status       string `json:"status"`
}

type GetBillingReportRequest struct {
	TenantID int64
	Period   string
}

type BillingClient interface {
	GetBillingReport(ctx context.Context, req *GetBillingReportRequest) (*BillingReport, error)
}

type billingStub struct{}

func NewBillingClient() BillingClient {
	return &billingStub{}
}

func (s *billingStub) GetBillingReport(_ context.Context, _ *GetBillingReportRequest) (*BillingReport, error) {
	return &BillingReport{Status: "stub"}, nil
}
