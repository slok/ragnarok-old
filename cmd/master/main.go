package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/cmd/master/flags"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/master/server"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/master/service/repository"
	"github.com/slok/ragnarok/master/service/scheduler"
	"github.com/slok/ragnarok/master/web"
	webapiv1 "github.com/slok/ragnarok/master/web/handler/api/v1"
)

// master dependencies is a helper object to group all the app dependencies
type masterDependencies struct {
	scheduler     scheduler.Scheduler
	nodeRepo      repository.Node
	failureRepo   repository.Failure
	logger        log.Logger
	nodeStatus    service.NodeStatusService
	failureStatus service.FailureStatusService
}

func createGRPCServer(cfg config.Config, deps masterDependencies, logger log.Logger) (*server.MasterGRPCServiceServer, error) {

	// Create the GRPC service server
	l, err := net.Listen("tcp", cfg.RPCListenAddress)
	if err != nil {
		return nil, err
	}
	srvServer := server.NewMasterGRPCServiceServer(deps.failureStatus, deps.nodeStatus, l, clock.Base(), logger)
	return srvServer, nil
}

func createHTTPServer(cfg config.Config, deps masterDependencies, logger log.Logger) (web.Server, error) {
	// Create the GRPC service server
	l, err := net.Listen("tcp", cfg.HTTPListenAddress)
	if err != nil {
		return nil, err
	}

	handler := struct {
		webapiv1.Handler
	}{
		Handler: webapiv1.NewJSONHandler(deps.scheduler, logger),
	}
	server := web.NewHTTPServer(web.DefaultHTTPRoutes, handler, l, logger)
	return server, nil
}

// Main run main logic.
func Main() error {
	logger := log.Base()

	// Get the command line arguments.
	cfg, err := flags.GetMasterConfig(os.Args[1:])
	if err != nil {
		logger.Error(err)
		return err
	}

	// Set debug mode
	if cfg.Debug {
		logger.Set("debug")
	}

	// Create dependencies
	nodeRepo := repository.NewMemNode()
	failureRepo := repository.NewMemFailure()
	deps := masterDependencies{
		nodeRepo:      nodeRepo,
		failureRepo:   failureRepo,
		nodeStatus:    service.NewNodeStatus(*cfg, nodeRepo, logger),
		failureStatus: service.NewFailureStatus(failureRepo, logger),
		scheduler:     scheduler.NewNodeLabels(failureRepo, nodeRepo, logger),
	}

	// TODO: Autoregister this node as a master node.
	grpcServer, err := createGRPCServer(*cfg, deps, logger)
	if err != nil {
		return err
	}
	go func() {
		grpcServer.Serve()
	}()

	httpServer, err := createHTTPServer(*cfg, deps, logger)
	if err != nil {
		return err
	}
	httpServer.Serve()

	return nil
}

func clean() {
	log.Debug("Cleaning...")
}

func main() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	errC := make(chan error)

	// Run main program
	go func() {
		if err := Main(); err != nil {
			errC <- err
		}
		return
	}()

	// Wait until signal (ctr+c, SIGTERM...).
	var exitCode int

Waiter:
	for {
		select {
		// Wait for errors
		case err := <-errC:
			if err != nil {
				exitCode = 1
				break Waiter
			}
			// Wait for signal
		case <-sigC:
			break Waiter
		}
	}

	clean()
	os.Exit(exitCode)
}
