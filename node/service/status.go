package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node/client"
	"github.com/slok/ragnarok/types"
)

const (
	hbErrTimeout = time.Millisecond * 500
)

// Status is the interface all the services that want to notify the status
// to the master need to implement
type Status interface {
	// Returns the node status.
	State() types.NodeState

	// RegisterOnMaster registers the node on the master.
	RegisterOnMaster() error

	// DeregisterOnMaster deregisters the node on the master.
	DeregisterOnMaster() error

	// StartHeartbeat starts a heartbeat interval to the master, it will return
	// an error when the heartbeat is already running and a channel that can be used to
	// be notified when the heartbeat start failing.
	StartHeartbeat(interval time.Duration) (hbErr chan error, err error)

	// StopHeartbeat stops a heartbeat interval.
	StopHeartbeat() error
}

// NodeStatus is a service that a node will use to report its status to its master.
type NodeStatus struct {
	nodeID string
	cli    client.Status
	tags   map[string]string
	logger log.Logger
	clock  clock.Clock

	state types.NodeState
	stMu  sync.Mutex // stMu is the node status mutex.

	hbFinishC   chan struct{}
	hearbeating bool
	hbT         *time.Ticker // Ticker of the heartbeat.
	hbMu        sync.Mutex   // hbMu is the node heartbeat mutex.
}

// NewNodeStatus returns a new NodeStatus.
func NewNodeStatus(nodeID string, tags map[string]string, cli client.Status, clock clock.Clock, logger log.Logger) *NodeStatus {
	return &NodeStatus{
		nodeID:    nodeID,
		cli:       cli,
		tags:      tags,
		logger:    logger,
		clock:     clock,
		hbFinishC: make(chan struct{}),
	}
}

// State satisfies Status interface.
func (n *NodeStatus) State() types.NodeState {
	n.stMu.Lock()
	defer n.stMu.Unlock()
	return n.state
}

// RegisterOnMaster satisfies Status interface.
func (n *NodeStatus) RegisterOnMaster() error {
	n.stMu.Lock()
	defer n.stMu.Unlock()

	if err := n.cli.RegisterNode(n.nodeID, n.tags); err != nil {
		return err
	}
	n.state = types.ReadyNodeState

	return nil
}

// DeregisterOnMaster satisfies Status interface.
func (n *NodeStatus) DeregisterOnMaster() error {
	return fmt.Errorf("deregistering on node not implemented")
}

// StartHeartbeat satisfies Status interface.
func (n *NodeStatus) StartHeartbeat(interval time.Duration) (chan error, error) {
	n.stMu.Lock()
	st := n.state
	n.stMu.Unlock()
	if st != types.ReadyNodeState {
		return nil, fmt.Errorf("register the node on the master before start heartbeating")
	}

	n.hbMu.Lock()
	ht := n.hearbeating
	n.hbMu.Unlock()

	if ht {
		return nil, fmt.Errorf("already heartbeating")
	}

	// Create a new ticker in each heartbeat interval and set the control channel.
	n.hbMu.Lock()
	n.hbT = n.clock.NewTicker(interval)
	n.hbFinishC = make(chan struct{})
	n.hbMu.Unlock()

	// Start a heatbeat in a periodic interval.
	n.logger.Infof("heartbeat started every %s", interval)
	n.hearbeating = true
	hbErrC := make(chan error)

	go func() {
		for {
			select {
			case <-n.hbFinishC:
				n.logger.Info("heartbeat stop signal received, stopping heartbeat")
				n.hbMu.Lock()
				n.hearbeating = false
				n.hbMu.Unlock()
				return

			case <-n.hbT.C:
				n.stMu.Lock()
				st := n.state
				n.stMu.Unlock()

				// The show must go on, if heartbeat error it will be notified and the one that started
				// the heartbeat is responsible of stopping it.
				if err := n.cli.NodeHeartbeat(n.nodeID, st); err != nil {
					select {
					case <-n.clock.After(hbErrTimeout):
						n.logger.Errorf("timeout notifying heartbeat error. Heartbeat error: %v", err)
					case hbErrC <- fmt.Errorf("heartbeat failed: %v", err):
					}
				} else {
					n.logger.Debug("heartbeat sent")
				}
			}
		}
	}()

	return hbErrC, nil
}

// StopHeartbeat satisfies Status interface.
func (n *NodeStatus) StopHeartbeat() error {
	n.hbMu.Lock()
	defer n.hbMu.Unlock()

	// Check if theres is a heartbeat.
	if !n.hearbeating {
		return fmt.Errorf("heartbeat not running")
	}

	// Sent heartbeat stop signal.
	close(n.hbFinishC)
	n.hbFinishC = nil

	// Stop heartbeat intervals for the GC.
	n.hbT.Stop()
	n.hbT = nil
	n.hearbeating = false

	return nil
}
