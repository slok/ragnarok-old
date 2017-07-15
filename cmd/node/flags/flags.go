package flags

import (
	"flag"
	"fmt"
	"os"

	"time"

	nodeconfig "github.com/slok/ragnarok/node/config"
)

const (
	// Default values.
	defaultDebug  = false
	defaultDryRun = false
)

type config struct {
	fs                *flag.FlagSet
	masterAddress     string
	heartbeatInterval string
	debug             bool
	dryRun            bool
}

func new() *config {
	cfg := &config{
		fs: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}

	cfg.fs.StringVar(
		&cfg.masterAddress, "master.address", "",
		"Address where the master is listening",
	)

	cfg.fs.StringVar(
		&cfg.heartbeatInterval, "heartbeat.interval", "15s",
		"Time interval the node will send a heartbeat to the master",
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

	// Check heartbeat valid timing.
	d, err := time.ParseDuration(c.heartbeatInterval)
	if err != nil || d <= 0 {
		err = fmt.Errorf("invalid heartbeat interval")
	}

	return err
}

// GetNodeConfig will return a new node configuration from the cmd flags.
func GetNodeConfig(args []string) (*nodeconfig.Config, error) {
	cfg := new()

	if err := cfg.parse(args); err != nil {
		return nil, err
	}

	// Parse intervals (parsing error validated on the parse).
	d, _ := time.ParseDuration(cfg.heartbeatInterval)

	nodeCfg := &nodeconfig.Config{
		MasterAddress:     cfg.masterAddress,
		HeartbeatInterval: d,
		Debug:             cfg.debug,
		DryRun:            cfg.dryRun,
	}

	if err := nodeCfg.Validate(); err != nil {
		return nil, err
	}

	return nodeCfg, nil
}
