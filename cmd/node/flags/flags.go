package flags

import (
	"flag"
	"fmt"
	"os"

	nodeconfig "github.com/slok/ragnarok/node/config"
)

const (
	// defaults values
	defaultHTTPListenAddress = ":9413"
	defaultRPCListenAddress  = ":9414"
	defaultDebug             = false
	defaultDryRun            = false
)

type config struct {
	fs                *flag.FlagSet
	httpListenAddress string
	rpcListenAddress  string
	debug             bool
	dryRun            bool
}

func new() *config {
	cfg := &config{
		fs: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}

	// Set flags
	cfg.fs.StringVar(
		&cfg.httpListenAddress, "http.listen-address", defaultHTTPListenAddress,
		"Address to listen for HTTP communication",
	)

	cfg.fs.StringVar(
		&cfg.rpcListenAddress, "rpc.listen-address", defaultRPCListenAddress,
		"Address to listen for RPC communication",
	)

	cfg.fs.BoolVar(
		&cfg.debug, "run.debug", defaultDebug,
		"Run in debug mode",
	)

	cfg.fs.BoolVar(
		&cfg.dryRun, "run.dry-run", defaultDryRun,
		"Run in dry run mode",
	)
	return cfg
}

func (c *config) parse(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return err
	}

	if len(c.fs.Args()) != 0 {
		err = fmt.Errorf("Invalid command line arguments. Help: %s -h", os.Args[0])
	}

	return err
}

// GetNodeConfig will return a new node configuration from the cmd flags
func GetNodeConfig(args []string) (*nodeconfig.Config, error) {
	cfg := new()

	if err := cfg.parse(args); err != nil {
		return nil, err
	}

	nodeCfg := &nodeconfig.Config{
		HTTPListenAddress: cfg.httpListenAddress,
		RPCListenAddress:  cfg.rpcListenAddress,
		Debug:             cfg.debug,
		DryRun:            cfg.dryRun,
	}

	if err := nodeCfg.Validate(); err != nil {
		return nil, err
	}

	return nodeCfg, nil
}
