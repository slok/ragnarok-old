package node

import (
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node/config"
	"github.com/slok/ragnarok/node/service"
)

// TODO: move node logic to services. IMPORTANT!!

// Node is the interface that a node needs to implement to be a failure node.
type Node interface {
	// Initialize will initialize the node, create connection, register on master, etc.
	Initialize() error
	// Start will start the node and all of its components.
	Start() error
	// Stop will stop the node and all of its components.
	Stop() error
	// GetID Gets the unique ID of the node.
	GetID() string
}

// FailureNode is a kind of node that injects failure on the host.
type FailureNode struct {
	id         string
	cfg        config.Config
	log        log.Logger
	statusSrv  service.Status       // the service that reports the status of the node to the master.
	failureSrv service.FailureState // the service that handle the failure status from the master.

	stopHBHandler chan struct{} // used to stop the background handling of heartbeat errors.
}

// NewFailureNode returns a new FailureNode instance.
func NewFailureNode(id string, cfg config.Config, statusSrv service.Status, failureSrv service.FailureState, logger log.Logger) *FailureNode {
	f := &FailureNode{
		id:            id,
		cfg:           cfg,
		log:           logger,
		statusSrv:     statusSrv,
		failureSrv:    failureSrv,
		stopHBHandler: make(chan struct{}),
	}

	if f.cfg.DryRun {
		logger.Warn("System failure node will run in dry run mode")
	}

	return f
}

// GetID satisfies FailureNode interface.
func (f *FailureNode) GetID() string {
	return f.id
}

// Initialize satisfies FailureNode interface.
func (f *FailureNode) Initialize() error {
	if err := f.statusSrv.RegisterOnMaster(); err != nil {
		return err
	}

	f.log.Info("system failure node ready")
	return nil
}

// handleHeartbeatErrors will handle the errors when the heartbeats to the master fail.
func (f *FailureNode) handleHeartbeatErrors(c chan error) {
	for {
		select {
		case <-f.stopHBHandler:
			return
		case err := <-c:
			f.log.Error(err)
		}
	}
}

// Start satisfies FailureNode interface.
func (f *FailureNode) Start() error {
	// Start heartbeating.
	hbErrC, err := f.statusSrv.StartHeartbeat(f.cfg.HeartbeatInterval)
	if err != nil {
		return err
	}

	// handle errors in background.
	go f.handleHeartbeatErrors(hbErrC)

	if err := f.failureSrv.StartHandling(); err != nil {
		return err
	}

	return nil
}

// Stop satisfies FailureNode interface.
func (f *FailureNode) Stop() error {
	f.log.Info("stopping node...")
	// Stop failure status  handler.
	if err := f.failureSrv.StopHandling(); err != nil {
		return err
	}

	// Stop heartbeating and heartbet error handler.
	if err := f.statusSrv.StopHeartbeat(); err != nil {
		return err
	}
	go func() {
		f.stopHBHandler <- struct{}{}
	}()

	return nil
}
