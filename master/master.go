package master

import (
	"sync"
	"fmt"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/types"
)

// Master is the master node that sends attacks to the nodes interface.
type Master interface {
	// RegisterNode registers a new node on the master.
	RegisterNode(id string, tags map[string]string) error

	// NodeHeartbeat sets the node state after its heartbeat.
	NodeHeartbeat(id string, state types.NodeState) error
}

// FailureMaster is the implementation of master failure sender.
type FailureMaster struct {
	debug  bool
	repo   NodeRepository // repository where all the nodes will be stored.
	logger log.Logger

	nodeLock sync.Mutex
}

// NewFailureMaster returns a new failure master.
func NewFailureMaster(cfg config.Config, repository NodeRepository, logger log.Logger) *FailureMaster {
	return &FailureMaster{
		debug:  cfg.Debug,
		repo:   repository,
		logger: logger,
	}
}

// RegisterNode implements Master interface.
func (f *FailureMaster) RegisterNode(id string, tags map[string]string) error {
	f.logger.WithField("nodeID", id).Infof("node registered on master")
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	n := &Node{
		ID:    id,
		Tags:  tags,
		State: types.UnknownNodeState,
	}

	return f.repo.StoreNode(id, n)
}

// NodeHeartbeat sets the node state after its heartbeat.
func (f *FailureMaster) NodeHeartbeat(id string, state types.NodeState) error {
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
