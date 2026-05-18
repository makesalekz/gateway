package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/makesalekz/gateway/internal/conf"
)

// HTTPServer wraps net/http.Server and implements kratos transport.Server interface.
type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(c *conf.Bootstrap, router *mux.Router) *HTTPServer {
	addr := c.Server.HTTP.Addr
	if addr == "" {
		addr = "0.0.0.0:8080"
	}

	return &HTTPServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *HTTPServer) Start(_ context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
