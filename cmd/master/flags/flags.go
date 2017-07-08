package flags

import (
	"flag"
	"fmt"
	"os"

	masterconfig "github.com/slok/ragnarok/master/config"
)

const (
	// defaults values
	defaultHTTPListenAddress = ":10444"
	defaultRPCListenAddress  = ":50444"
	defaultDebug             = false
)

type config struct {
	fs                *flag.FlagSet
	httpListenAddress string
	rpcListenAddress  string
	debug             bool
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

// GetMasterConfig will return a new master configuration from the cmd flags
func GetMasterConfig(args []string) (*masterconfig.Config, error) {
	cfg := new()

	if err := cfg.parse(args); err != nil {
		return nil, err
	}

	nodeCfg := &masterconfig.Config{
		HTTPListenAddress: cfg.httpListenAddress,
		RPCListenAddress:  cfg.rpcListenAddress,
		Debug:             cfg.debug,
	}

	if err := nodeCfg.Validate(); err != nil {
		return nil, err
	}

	return nodeCfg, nil
}
