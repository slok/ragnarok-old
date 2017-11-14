package web

import (
	"net"
	"net/http"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/web/handler"
)

// Server is the server that will serve all the web app including the http API.
type Server interface {
	Serve() error
}

// HTTPServer serves all the application handlers and routers
type HTTPServer struct {
	handler  http.Handler
	listener net.Listener

	logger log.Logger
}

func registerRoutes(routes HTTPRoutes, handler handler.Handler, sm *http.ServeMux) {
	// API v1 routes
	sm.Handle(routes.APIV1.Debug, http.HandlerFunc(handler.Debug))
	sm.Handle(routes.APIV1.WriteExperiment, http.HandlerFunc(handler.WriteExperiment))
}

// NewHTTPServer will return a new http handler.
func NewHTTPServer(routes HTTPRoutes, handler handler.Handler, listener net.Listener, logger log.Logger) *HTTPServer {
	sm := http.NewServeMux()
	registerRoutes(routes, handler, sm)

	server := &HTTPServer{
		handler:  sm,
		listener: listener,
		logger:   logger,
	}
	return server
}

// Serve satisfies Server interface.
func (h *HTTPServer) Serve() error {
	h.logger.Infof("ready to listen HTTP requests on %s", h.listener.Addr())
	return http.Serve(h.listener, h.handler)
}
