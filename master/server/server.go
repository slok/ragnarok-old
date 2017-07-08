package server

import (
	"net"

	"google.golang.org/grpc"

	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master"
	"github.com/slok/ragnarok/master/server/service"
)

// GRPCServiceServiceServer is an interface that wraps all the GRPC service need to implement.
type GRPCServiceServiceServer interface {
	pbns.NodeStatusServer
	// Serve will serve the services.
	Serve(addr string) error
}

// MasterGRPCServiceServer is an implementation of the service server using using master
// implementation as logic, it wraps all the services used as.
type MasterGRPCServiceServer struct {
	*service.NodeStatusGRPC
	server   *grpc.Server
	listener net.Listener
	logger   log.Logger
}

// NewMasterGRPCServiceServer returns a new grpc service server with a master as a base.
func NewMasterGRPCServiceServer(master master.Master, listener net.Listener, logger log.Logger) *MasterGRPCServiceServer {

	// TODO: Authentication.
	// Create the GRPC server.
	grpcServer := grpc.NewServer()
	m := &MasterGRPCServiceServer{
		// Node status service.
		NodeStatusGRPC: service.NewNodeStatusGRPC(master, logger),
		server:         grpcServer,
		logger:         logger,
		listener:       listener,
	}

	// Register our services on the grpc server.
	m.registerServices()

	return m
}

// registerServices will register all the services on the grpc server.
func (m *MasterGRPCServiceServer) registerServices() {
	// Register node status service
	pbns.RegisterNodeStatusServer(m.server, m.NodeStatusGRPC)
}

// Serve implements the GRPCServiceServiceServer interface.
func (m *MasterGRPCServiceServer) Serve() error {
	m.logger.Infof("ready to listen GRPC service calls on %s", m.listener.Addr().String)

	// Start serving our GRPC service.
	return m.server.Serve(m.listener)
}
