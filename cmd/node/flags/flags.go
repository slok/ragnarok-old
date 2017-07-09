package flags

import (
	"flag"
	"fmt"
	"os"

	nodeconfig "github.com/slok/ragnarok/node/config"
)

const (
	// Default values.
	defaultDebug  = false
	defaultDryRun = false
)

type config struct {
	fs            *flag.FlagSet
	masterAddress string
	debug         bool
	dryRun        bool
}

func new() *config {
	cfg := &config{
		fs: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}

	cfg.fs.StringVar(
		&cfg.masterAddress, "master.address", "",
		"Address where the master is listening",
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
		err = fmt.Errorf("invalid command line arguments. Help: %s -h", os.Args[0])
	}

	// Check master address is set.
	if c.masterAddress == "" {
		err = fmt.Errorf("master address is required")
	}

	return err
}

// GetNodeConfig will return a new node configuration from the cmd flags.
func GetNodeConfig(args []string) (*nodeconfig.Config, error) {
	cfg := new()

	if err := cfg.parse(args); err != nil {
		return nil, err
	}

	nodeCfg := &nodeconfig.Config{
		MasterAddress: cfg.masterAddress,
		Debug:         cfg.debug,
		DryRun:        cfg.dryRun,
	}

	if err := nodeCfg.Validate(); err != nil {
		return nil, err
	}

	return nodeCfg, nil
}
