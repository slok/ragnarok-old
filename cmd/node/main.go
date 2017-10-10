package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	"github.com/slok/ragnarok/chaos/failure"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/cmd/node/flags"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/client"
	"github.com/slok/ragnarok/node/service"
	"github.com/slok/ragnarok/types"
)

// Main run main logic.
func Main() error {
	nodeID := uuid.New().String()
	nodeTags := map[string]string{"id": nodeID, "version": "v0.1alpha"}
	logger := log.Base().WithField("id", nodeID)

	// Get the command line arguments.
	cfg, err := flags.GetNodeConfig(os.Args[1:])
	if err != nil {
		logger.Error(err)
		return err
	}

	// Set debug mode.
	if cfg.Debug {
		logger.Set("debug")
	}

	// Create node GRPC clients
	conn, err := grpc.Dial(cfg.MasterAddress, grpc.WithInsecure()) // TODO: secured.
	if err != nil {
		return err
	}
	// TODO: Handle correctly the disconnect, reconnects...
	//defer conn.Close()

	// Create GRPC clients.
	nsCli, err := client.NewStatusGRPCFromConnection(conn, types.NodeStateTransformer, logger)
	if err != nil {
		return err
	}

	fCli, err := client.NewFailureGRPCFromConnection(conn, failure.Transformer, types.FailureStateTransformer, clock.Base(), logger)
	if err != nil {
		return err
	}

	// Create services.
	stSrv := service.NewNodeStatus(nodeID, nodeTags, nsCli, clock.Base(), logger)
	fSrv := service.NewLogFailureState(nodeID, fCli, clock.Base(), logger)

	// Create the node.
	n := node.NewFailureNode(nodeID, *cfg, stSrv, fSrv, logger)

	// Register node & start.
	if err := n.Initialize(); err != nil {
		return fmt.Errorf("node could not inicialize: %v", err)
	}

	if err := n.Start(); err != nil {
		return fmt.Errorf("could not start the node: %v", err)
	}

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

	// Wait until signal (ctr+c, SIGTERM...)
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
