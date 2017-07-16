package config

import (
	"fmt"
	"time"
)

// Config is the configuration of a node.
type Config struct {
	// Debug will set the node in debug mode.
	Debug bool
	// DryRun will set the node in dry run mode.
	DryRun bool
	// MasterAddress is the address where the master server is listening.
	MasterAddress string
	// HeartbeatInterval is the interval the node will send a heartbeat to the master
	HeartbeatInterval time.Duration
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.MasterAddress == "" {
		return fmt.Errorf("master address is required")
	}

	if c.HeartbeatInterval == 0 {
		return fmt.Errorf("heartbeat interval can't be 0")
	}

	return nil
}
