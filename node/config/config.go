package config

import (
	"fmt"
)

// Config is the configuration of a node.
type Config struct {
	// Debug will set the node in debug mode.
	Debug bool
	// DryRun will set the node in dry run mode.
	DryRun bool
	// MasterAddress is the address where the master server is listening.
	MasterAddress string
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.MasterAddress == "" {
		return fmt.Errorf("master address is required")
	}

	return nil
}
