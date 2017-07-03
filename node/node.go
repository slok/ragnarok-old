package node

import (
	"github.com/google/uuid"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node/client"
	"github.com/slok/ragnarok/node/config"
)

// Node is the interface that a node needs to implement to be a failure node
type Node interface {
	// RegisterOnMaster registers the node on the master
	RegisterOnMaster() error
	// DeregisterOnMaster deregisters the node on the master
	DeregisterOnMaster() error
	// Serve serves the RPC and HTTP services
	Serve() error
	// GetID Gets the unique ID of the node
	GetID() string
}

// FailureNode is a kind of node that injects failure on the host
type FailureNode struct {
	id string

	log          log.Logger
	dryRun       bool
	debug        bool
	statusClient client.Status // client to communicate with node status service
}

// NewFailureNode returns a new FailureNode instnace
func NewFailureNode(cfg config.Config, statusClient client.Status, logger log.Logger) *FailureNode {
	id := uuid.New().String()

	logger = logger.WithField("id", id)

	f := &FailureNode{
		id:           id,
		log:          logger,
		dryRun:       cfg.DryRun,
		debug:        cfg.Debug,
		statusClient: statusClient,
	}

	logger.Info("System failure node ready")

	if f.dryRun {
		logger.Warn("System failure node in dry run mode")
	}

	return f
}

// GetID satisfies FailureNode interface
func (f *FailureNode) GetID() string {
	return f.id
}

// RegisterOnMaster satisfies FailureNode interface
func (f *FailureNode) RegisterOnMaster() error {
	return f.statusClient.RegisterNode(f.id, map[string]string{})
}
