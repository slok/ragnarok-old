package server

import (
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/clock"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	grpcservice "github.com/slok/ragnarok/master/service/grpc"
)

// TODO: set this as configurable.
const failureStatusUpInterval = 15 * time.Second

// GRPCServiceServer is an interface that wraps all the GRPC service need to implement.
type GRPCServiceServer interface {
	pbns.NodeStatusServer
	pbfs.FailureStatusServer

	// Serve will serve the services.
	Serve(addr string) error
}

// MasterGRPCServiceServer is an implementation of the service server using using master
// implementation as logic, it wraps all the services used as.
type MasterGRPCServiceServer struct {
	*grpcservice.NodeStatus
	*grpcservice.FailureStatus
	server   *grpc.Server
	listener net.Listener
	logger   log.Logger
}

// NewMasterGRPCServiceServer returns a new grpc service server with a master as a base.
func NewMasterGRPCServiceServer(fss service.FailureStatusService, nss service.NodeStatusService, listener net.Listener, clock clock.Clock, logger log.Logger) *MasterGRPCServiceServer {

	// Create different grpc services.
	gnss := grpcservice.NewNodeStatus(nss, serializer.PBSerializerDefault, logger)
	gfss := grpcservice.NewFailureStatus(failureStatusUpInterval, serializer.PBSerializerDefault, fss, clock, logger)

	// TODO: Authentication.
	// Create the GRPC server.
	grpcServer := grpc.NewServer()
	m := &MasterGRPCServiceServer{
		NodeStatus:    gnss, // Node status service.
		FailureStatus: gfss, // Failure status service.

		server:   grpcServer,
		logger:   logger,
		listener: listener,
	}

	// Register our services on the grpc server.
	m.registerServices()

	return m
}

// registerServices will register all the services on the grpc server.
func (m *MasterGRPCServiceServer) registerServices() {
	// Register node status service.
	pbns.RegisterNodeStatusServer(m.server, m)

	// Register failure status service.
	pbfs.RegisterFailureStatusServer(m.server, m)
}

// Serve implements the GRPCServiceServiceServer interface.
func (m *MasterGRPCServiceServer) Serve() error {
	m.logger.Infof("ready to listen GRPC service calls on %s", m.listener.Addr().String())

	// Start serving our GRPC service.
	return m.server.Serve(m.listener)
}
