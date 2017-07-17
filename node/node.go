package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node/client"
	"github.com/slok/ragnarok/node/config"
	"github.com/slok/ragnarok/types"
)

// Node is the interface that a node needs to implement to be a failure node.
type Node interface {
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
	id           string
	cfg          config.Config
	clock        clock.Clock
	log          log.Logger
	dryRun       bool
	debug        bool
	statusClient client.Status   // Client to communicate with node status service.
	state        types.NodeState // The status of the node.

	// Heartbeat
	hbT       *time.Ticker  // Ticker of the heartbeat.
	hbFinishC chan struct{} // Used to notify that the heartbeat has finished

	// Mutexes
	stMu  sync.Mutex // Mutex for the state.
	hbTMu sync.Mutex // Mutex for the heartbeat ticker.
}

// NewFailureNode returns a new FailureNode instance.
func NewFailureNode(cfg config.Config,
	statusClient client.Status,
	clock clock.Clock,
	logger log.Logger) *FailureNode {

	id := uuid.New().String()

	logger = logger.WithField("id", id)

	f := &FailureNode{
		id:           id,
		cfg:          cfg,
		clock:        clock,
		log:          logger,
		dryRun:       cfg.DryRun,
		debug:        cfg.Debug,
		statusClient: statusClient,
		state:        types.UnknownNodeState,
	}

	logger.Info("System failure node ready")

	if f.dryRun {
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
	// TODO: set the node state to ready state
	return f.statusClient.RegisterNode(f.id, map[string]string{})
}

// StartHeartbeat satisfies FailureNode interface.
func (f *FailureNode) StartHeartbeat() error {
	f.hbTMu.Lock()
	t := f.hbT
	f.hbTMu.Unlock()

	if t != nil {
		return fmt.Errorf("seems that the heartbeat is already working")
	}

	// Create a new ticker in each heartbeat interval.
	f.hbTMu.Lock()
	t = f.clock.NewTicker(f.cfg.HeartbeatInterval)
	f.hbT = t
	f.hbTMu.Unlock()

	// Set our heartbeat control channel.
	f.hbFinishC = make(chan struct{})

	// Start a heatbeat in a periodic interval.
	f.log.Info("heartbeat started")
	for {
		select {
		case <-f.hbFinishC:
			f.log.Info("heartbeat stopped")
			return nil

		case <-t.C:
			f.stMu.Lock()
			st := f.state
			f.stMu.Unlock()

			// The show must go on.
			if err := f.statusClient.NodeHeartbeat(f.id, st); err != nil {
				f.log.Errorf("heartbeat failed: %v", err)
			} else {
				f.log.Debug("heartbeat sent")
			}
		}
	}
}

// StopHeartbeat satisfies FailureNode interface.
func (f *FailureNode) StopHeartbeat() error {
	f.hbTMu.Lock()
	defer f.hbTMu.Unlock()

	// Check if theres is a heartbeat.
	if f.hbT == nil {
		return fmt.Errorf("heartbeat not running")
	}

	// Sent heartbeat stop signal.
	close(f.hbFinishC)
	f.hbFinishC = nil

	// Stop heartbeat intervals for the GC.
	f.hbT.Stop()
	f.hbT = nil

	return nil
}
