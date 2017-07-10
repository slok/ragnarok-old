package master

import (
	"sync"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
)

// Master is the master node that sends attacks to the nodes interface
type Master interface {
	// RegisterNode registers a new node on the master
	RegisterNode(id string, tags map[string]string) error
	//GetRegisteredNodes() []string
}

// FailureMaster is the implementation of master failure sender
type FailureMaster struct {
	debug  bool
	repo   NodeRepository // repository where all the nodes will be stored
	logger log.Logger

	nodeLock sync.Mutex
}

// NewFailureMaster returns a new failure master
func NewFailureMaster(cfg config.Config, repository NodeRepository, logger log.Logger) *FailureMaster {
	return &FailureMaster{
		debug:  cfg.Debug,
		repo:   repository,
		logger: logger,
	}
}

// RegisterNode implements Master interface
func (f *FailureMaster) RegisterNode(id string, tags map[string]string) error {
	f.logger.WithField("nodeID", id).Infof("node registered on master")
	f.nodeLock.Lock()
	defer f.nodeLock.Unlock()

	n := &Node{
		ID:   id,
		Tags: tags,
	}

	return f.repo.StoreNode(id, n)
}
