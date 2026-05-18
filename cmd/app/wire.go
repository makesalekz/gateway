//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/makesalekz/gateway/internal/clients"
	"github.com/makesalekz/gateway/internal/conf"
	"github.com/makesalekz/gateway/internal/handlers"
	"github.com/makesalekz/gateway/internal/router"
	"github.com/makesalekz/gateway/internal/server"
)

func wireApp(*conf.Bootstrap, log.Logger, []byte) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		clients.ProviderSet,
		wire.NewSet(handlers.NewAuthHandler, handlers.NewProductsHandler, handlers.NewSalesHandler, handlers.NewDashboardHandler, handlers.NewAgentsHandler, handlers.NewDocsHandler),
		wire.NewSet(router.New),
		newApp,
	))
}
