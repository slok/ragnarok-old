package service

import (
	"fmt"
	"sync"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
)

// NodeStatusService is how the master manages the status of the nodes.
type NodeStatusService interface {
	// Register registers a new node on the master.
	Register(id string, labels map[string]string) error

	// Heartbeat sets the node state after its heartbeat.
	Heartbeat(id string, state clusterv1.NodeState) error
}

// NodeStatus is the implementation of node status service.
type NodeStatus struct {
	client cliclusterv1.NodeClientInterface // Client will manage the node object operations.
	logger log.Logger

	nodeLock sync.Mutex
}

// NewNodeStatus returns a new node status service.
func NewNodeStatus(_ config.Config, client cliclusterv1.NodeClientInterface, logger log.Logger) *NodeStatus {
	return &NodeStatus{
		client: client,
		logger: logger,
	}
}

// Register implements NodeStatusService interface.
func (f *NodeStatus) Register(id string, labels map[string]string) error {
	f.logger.WithField("nodeID", id).Infof("node registered on master")
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	n := clusterv1.NewNode()
	n.Metadata = api.ObjectMeta{
		ID:     id,
		Labels: labels,
	}
	n.Status = clusterv1.NodeStatus{
		State: clusterv1.UnknownNodeState,
	}
	_, err := f.client.Create(&n)
	return err
}

// Heartbeat sets the node state after its heartbeat.
func (f *NodeStatus) Heartbeat(id string, state clusterv1.NodeState) error {
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	// Get the node.
	n, err := f.client.Get(id)
	if err != nil {
		return fmt.Errorf("node '%s' not registered", id)
	}

	// Set state and save.
	n.Status.State = state
	if _, err := f.client.Update(n); err != nil {
		return err
	}

	f.logger.WithField("nodeID", id).Infof("node in state %s", state)
	return nil
}
