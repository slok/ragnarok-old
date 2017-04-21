package memory

import (
	"testing"

	"github.com/slok/ragnarok/attack"
	"github.com/stretchr/testify/assert"
)

func TestCreationWithOpts(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		opts      attack.Opts
		expectErr bool
	}{
		{
			opts: attack.Opts{
				"size": 1000,
			},
			expectErr: false,
		},
		{
			opts: attack.Opts{
				"Size": 1000,
			},
			expectErr: true,
		},
		{
			opts: attack.Opts{
				"size": "1",
			},
			expectErr: true,
		},
		{
			opts: attack.Opts{
				"size": nil,
			},
			expectErr: true,
		},
		{
			opts: attack.Opts{
				"size": 0,
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		m, err := NewMemAllocationOpts(test.opts)
		if !test.expectErr {
			if assert.NoError(err, "Creation of memory allocator shouldn't error") {
				assert.EqualValues(test.opts[sizeKey], m.Size)
			}
		} else {
			assert.Error(err, "Creation of memory allocator shoud error")
		}
	}
}
