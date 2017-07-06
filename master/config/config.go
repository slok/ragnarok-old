package config

import (
	"fmt"
)

// Config is the configuration of a node
type Config struct {
	// HTTPListenAddress is the address where the http server will be listening
	HTTPListenAddress string
	// RPCListenAddress is the address where the rpc server will be listening
	RPCListenAddress string
	// Debug will set the node in debug mode
	Debug bool
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if len(c.HTTPListenAddress) == 0 {
		return fmt.Errorf("%s is not a valid listen address", c.HTTPListenAddress)
	}

	if len(c.RPCListenAddress) == 0 {
		return fmt.Errorf("%s is not a valid listen address", c.RPCListenAddress)
	}

	if c.RPCListenAddress == c.HTTPListenAddress {
		return fmt.Errorf("HTTP and RPC listen addresses can't be the same")
	}

	return nil
}
