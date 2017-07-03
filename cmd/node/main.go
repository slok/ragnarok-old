package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/slok/ragnarok/cmd/node/flags"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node"
)

// Main run main logic
func Main() error {
	logger := log.Base()

	// Get the command line arguments
	cfg, err := flags.GetNodeConfig(os.Args[1:])
	if err != nil {
		logger.Error(err)
		return err
	}

	// Set debug mode
	if cfg.Debug {
		logger.Set("debug")
	}

	// Create the node
	// TODO: GRPC client, for now nil
	n := node.NewFailureNode(*cfg, nil, logger)
	n.GetID()

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
