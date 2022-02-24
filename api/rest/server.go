package rest

import (
	"context"
	"gophermart/api/rest/handler"
	"gophermart/pkg/logger"
	"net/http"
)

// APIServer base struct.
type APIServer struct {
	server  *http.Server
	handler *handler.Handler
}

// NewAPIServer APIServer constructor.
func NewAPIServer(handler *handler.Handler, addr string) *APIServer {
	server := &APIServer{
		handler: handler,
		server:  &http.Server{Addr: addr, Handler: handler.GetRouter()},
	}

	return server
}

// Run HTTP server.
func (srv *APIServer) Run(ctx context.Context) error {
	go func() {
		logger.GetLogger().Infow("server started", "addr", srv.server.Addr)
		_ = srv.server.ListenAndServe()
	}()

	return nil
}
