package web

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/slok/ragnarok/apimachinery/serializer"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/web/handler"
	clusterv1 "github.com/slok/ragnarok/master/web/handler/api/cluster/v1"
)

// Server is the server that will serve all the web app including the http API.
type Server interface {
	// HandleResource will register the resource handler on the server.
	HandleResource(rh handler.ResourceHandler) error
	// Serve will serve the API.
	Serve() error
}

// HTTPServer serves all the application handlers and routers
type HTTPServer struct {
	handler     *http.ServeMux
	listener    net.Listener
	dispatchers map[string]ResourceHandlerDispatcherInterface
	lock        sync.Mutex

	logger log.Logger
}

// NewDefaultHTTPServer returns a new http handler with all the required resources registered.
func NewDefaultHTTPServer(
	serializer serializer.Serializer,
	nodeCli cliclusterv1.NodeClientInterface,
	listener net.Listener,
	logger log.Logger) (*HTTPServer, error) {

	server := NewHTTPServer(listener, logger)

	// Register handlers.
	nodeh := clusterv1.NewNodeHandler(serializer, nodeCli)
	if err := server.HandleResource(nodeh); err != nil {
		return nil, err
	}

	return server, nil
}

// NewHTTPServer will return a new http handler.
func NewHTTPServer(listener net.Listener, logger log.Logger) *HTTPServer {
	sm := http.NewServeMux()

	server := &HTTPServer{
		handler:     sm,
		listener:    listener,
		dispatchers: map[string]ResourceHandlerDispatcherInterface{},
		logger:      logger,
	}
	return server
}

// HandleResource registers a resource handler. Satisfies Server interface.
func (h *HTTPServer) HandleResource(rh handler.ResourceHandler) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	route := rh.GetRoute()
	if _, ok := h.dispatchers[route]; ok {
		return fmt.Errorf("already handler registered on %s", route)
	}

	// Create a dispatcher for the resource handler.
	d := NewResourceHandlerDispatcher(rh)
	h.dispatchers[route] = d
	h.handler.Handle(route, http.HandlerFunc(d.Dispatch))
	h.logger.Infof("registered %s resource handler", route)
	return nil
}

// Serve satisfies Server interface.
func (h *HTTPServer) Serve() error {
	h.logger.Infof("ready to listen HTTP requests on %s", h.listener.Addr())
	return http.Serve(h.listener, h.handler)
}
