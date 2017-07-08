package flags_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/cmd/node/flags"
	"github.com/slok/ragnarok/node/config"
)

func TestCofiguration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		flags       []string
		expectedCfg config.Config
		shouldError bool
	}{
		{
			[]string{
				"--run.debug",
				"-master.address", "127.0.0.1:8080",
			},
			config.Config{
				HTTPListenAddress: ":10443",
				RPCListenAddress:  ":50443",
				MasterAddress:     "127.0.0.1:8080",
				Debug:             true,
				DryRun:            false,
			},
			false,
		},
		{
			[]string{
				"--run.dry-run",
				"-master.address", "127.0.0.1:8080",
			},
			config.Config{
				HTTPListenAddress: ":10443",
				RPCListenAddress:  ":50443",
				MasterAddress:     "127.0.0.1:8080",
				Debug:             false,
				DryRun:            true,
			},
			false,
		},
		{
			[]string{
				"-run.dry-run",
				"-http.listen-address", "127.0.0.1:9999",
				"-master.address", "127.0.0.1:8080",
			},
			config.Config{
				HTTPListenAddress: "127.0.0.1:9999",
				RPCListenAddress:  ":50443",
				MasterAddress:     "127.0.0.1:8080",
				Debug:             false,
				DryRun:            true,
			},
			false,
		},
		{
			[]string{
				"-run.dry-run",
				"-rpc.listen-address", "127.0.0.1:9999",
				"--master.address", "127.0.0.1:8080",
			},
			config.Config{
				HTTPListenAddress: ":10443",
				RPCListenAddress:  "127.0.0.1:9999",
				MasterAddress:     "127.0.0.1:8080",
				Debug:             false,
				DryRun:            true,
			},
			false,
		},
		{
			[]string{
				"-run.dry-run",
				"-http.listen-address", "127.0.0.1:9999",
				"-rpc.listen-address", "127.0.0.1:9999",
				"-master.address", "127.0.0.1:8080",
			},
			config.Config{},
			true,
		},
		{
			[]string{
				"-run.dry-run",
			},
			config.Config{},
			true,
		},
		{
			[]string{"something"},
			config.Config{},
			true,
		},
		{
			[]string{"--not-real"},
			config.Config{},
			true,
		},
	}

	for _, test := range tests {
		cfg, err := flags.GetNodeConfig(test.flags)

		if test.shouldError {
			assert.Error(err)
		} else {
			require.NoError(err)
			assert.Equal(test.expectedCfg, *cfg)
		}
	}
}
