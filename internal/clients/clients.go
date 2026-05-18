// Package clients defines interfaces for all backend gRPC services.
// Current implementations are stubs. When replacing with real gRPC clients:
// TODO(grpc-metadata): inject x-md-global-tenant-id, x-md-global-actor-id,
// x-md-global-app-id into outbound gRPC metadata per project conventions.
// TODO(grpc-clients): accept *conf.Bootstrap to read discovery addresses.
package clients

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewIAMClient,
	NewProductsClient,
	NewWarehouseClient,
	NewSalesClient,
	NewOrdersClient,
	NewAgentsClient,
	NewStoresClient,
	NewBillingClient,
)
