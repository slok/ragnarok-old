package master

import (
	"sync"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
)

// Master is the master node that sends attacks to the nodes interface
type Master interface {
	// RegisterNode registers a new node on the master
	RegisterNode(id string, address string) error
	//GetRegisteredNodes() []string
}

// FailureMaster is the implementation of master failure sender
type FailureMaster struct {
	debug  bool
	reg    NodeRegistry // registry where all the nodes will be stored
	logger log.Logger

	nodeLock sync.Mutex
}

// NewFailureMaster returns a new failure master
func NewFailureMaster(cfg config.Config, registry NodeRegistry, logger log.Logger) *FailureMaster {
	return &FailureMaster{
		debug:  cfg.Debug,
		reg:    registry,
		logger: logger,
	}
}

// RegisterNode implements Master interface
func (f *FailureMaster) RegisterNode(id string, address string) error {
	f.logger.WithField("nodeID", id).Infof("node registered on master")
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	n := &Node{
		ID:      id,
		Address: address,
	}

	return f.reg.AddNode(id, n)
}
