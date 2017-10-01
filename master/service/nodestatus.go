package service

import (
	"fmt"
	"sync"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/master/model"
	"github.com/slok/ragnarok/types"
)

// NodeStatusService is how the master manages the status of the nodes.
type NodeStatusService interface {
	// Register registers a new node on the master.
	Register(id string, tags map[string]string) error

	// Heartbeat sets the node state after its heartbeat.
	Heartbeat(id string, state types.NodeState) error
}

// NodeStatus is the implementation of node status service.
type NodeStatus struct {
	repo   NodeRepository // Repo is the repository where all the nodes will be stored.
	logger log.Logger

	nodeLock sync.Mutex
}

// NewNodeStatus returns a new node status service.
func NewNodeStatus(_ config.Config, repository NodeRepository, logger log.Logger) *NodeStatus {
	return &NodeStatus{
		repo:   repository,
		logger: logger,
	}
}

// Register implements NodeStatusService interface.
func (f *NodeStatus) Register(id string, tags map[string]string) error {
	f.logger.WithField("nodeID", id).Infof("node registered on master")
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	n := &model.Node{
		ID:    id,
		Tags:  tags,
		State: types.UnknownNodeState,
	}

	return f.repo.StoreNode(id, n)
}

// Heartbeat sets the node state after its heartbeat.
func (f *NodeStatus) Heartbeat(id string, state types.NodeState) error {
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	// Get the node.
	n, ok := f.repo.GetNode(id)
	if !ok {
		return fmt.Errorf("node '%s' not registered", id)
	}

	// Set state and save.
	n.State = state
	if err := f.repo.StoreNode(id, n); err != nil {
		return fmt.Errorf("could not set the current state: %v", err)
	}

	f.logger.WithField("nodeID", id).Infof("node in state %s", state)
	return nil
}
