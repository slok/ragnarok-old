package flags_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/cmd/master/flags"
	"github.com/slok/ragnarok/master/config"
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
			[]string{"--run.debug"},
			config.Config{
				HTTPListenAddress: ":10444",
				RPCListenAddress:  ":50444",
				Debug:             true,
			},
			false,
		},
		{
			[]string{"-http.listen-address", "127.0.0.1:9999"},
			config.Config{
				HTTPListenAddress: "127.0.0.1:9999",
				RPCListenAddress:  ":50444",
				Debug:             false,
			},
			false,
		},
		{
			[]string{"-rpc.listen-address", "127.0.0.1:9999"},
			config.Config{
				HTTPListenAddress: ":10444",
				RPCListenAddress:  "127.0.0.1:9999",
				Debug:             false,
			},
			false,
		},
		{
			[]string{"-http.listen-address", "127.0.0.1:9999", "-rpc.listen-address", "127.0.0.1:9999"},
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
		cfg, err := flags.GetMasterConfig(test.flags)

		if test.shouldError {
			assert.Error(err)
		} else {
			require.NoError(err)
			assert.Equal(test.expectedCfg, *cfg)
		}
	}
}
