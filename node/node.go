package node

import (
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node/config"
	"github.com/slok/ragnarok/node/service"
)

// TODO: move node logic to services. IMPORTANT!!

// Node is the interface that a node needs to implement to be a failure node.
type Node interface {
	// StartHandlingFailureStates will start handling failures received from the master.
	StartHandlingFailureStates() error
	// StopHandlingFailureStates will stop handling failures received from the master.
	StopHandlingFailureStates() error
	// RegisterOnMaster registers the node on the master.
	RegisterOnMaster() error
	// DeregisterOnMaster deregisters the node on the master.
	DeregisterOnMaster() error
	// Serve serves the RPC and HTTP services.
	Serve() error
	// GetID Gets the unique ID of the node.
	GetID() string
	// StartHeartbeat starts a heartbeat interval to the master.
	StartHeartbeat() error
	// StopHeartbeat stops a heartbeat interval.
	StopHeartbeat() error
}

// FailureNode is a kind of node that injects failure on the host.
type FailureNode struct {
	id         string
	cfg        config.Config
	log        log.Logger
	statusSrv  service.Status       // the service that reports the status of the node to the master.
	failureSrv service.FailureState // the service that handle the failure status from the master.
}

// NewFailureNode returns a new FailureNode instance.
func NewFailureNode(id string, cfg config.Config, statusSrv service.Status, failureSrv service.FailureState, logger log.Logger) *FailureNode {
	f := &FailureNode{
		id:         id,
		cfg:        cfg,
		log:        logger,
		statusSrv:  statusSrv,
		failureSrv: failureSrv,
	}

	logger.Info("System failure node ready")

	if f.cfg.DryRun {
		logger.Warn("System failure node in dry run mode")
	}

	return f
}

// GetID satisfies FailureNode interface.
func (f *FailureNode) GetID() string {
	return f.id
}

// RegisterOnMaster satisfies FailureNode interface.
func (f *FailureNode) RegisterOnMaster() error {
	return f.statusSrv.RegisterOnMaster()
}

// StartHeartbeat satisfies FailureNode interface.
func (f *FailureNode) StartHeartbeat() error {
	// TODO: Handle heartbeat errors.
	_, err := f.statusSrv.StartHeartbeat(f.cfg.HeartbeatInterval)
	return err
}

// StopHeartbeat satisfies FailureNode interface.
func (f *FailureNode) StopHeartbeat() error {
	return f.statusSrv.StopHeartbeat()
}

// StartHandlingFailureStates satisfies FailureNode interface.
func (f *FailureNode) StartHandlingFailureStates() error {
	return f.failureSrv.StartHandling()
}

// StopHandlingFailureStates satisfies FailureNode interface.
func (f *FailureNode) StopHandlingFailureStates() error {
	return f.failureSrv.StopHandling()
}
