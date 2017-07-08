package main

import (
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"fmt"

	"github.com/slok/ragnarok/cmd/node/flags"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/client"
)

// Main run main logic.
func Main() error {
	logger := log.Base()

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

	// Create node status client
	conn, err := grpc.Dial(cfg.MasterAddress, grpc.WithInsecure()) // TODO: secured.
	if err != nil {
		return err
	}
	defer conn.Close()
	nsCli, err := client.NewStatusGRPCFromConnection(conn, logger)
	if err != nil {
		return err
	}
	// Create the node.
	n := node.NewFailureNode(*cfg, nsCli, logger)

	// Register node.
	if err := n.RegisterOnMaster(); err != nil {
		return fmt.Errorf("node not registered on master: %v", err)
	}

	// TODO: Listen for service calls

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
