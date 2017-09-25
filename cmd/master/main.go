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
)

func createGRPCServer(cfg config.Config, logger log.Logger) (*server.MasterGRPCServiceServer, error) {
	// Create the services.
	nodeReg := service.NewMemNodeRepository()
	failureReg := service.NewMemFailureRepository()
	nss := service.NewNodeStatus(cfg, nodeReg, logger)
	fss := service.NewFailureStatus(failureReg, logger)

	// Create the GRPC service server
	l, err := net.Listen("tcp", cfg.RPCListenAddress)
	if err != nil {
		return nil, err
	}
	srvServer := server.NewMasterGRPCServiceServer(fss, nss, l, clock.Base(), logger)
	return srvServer, nil
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

	gserver, err := createGRPCServer(*cfg, logger)
	if err != nil {
		return err
	}
	gserver.Serve()

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
