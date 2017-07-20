package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/master/model"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/types"
)

func TestStoreFailure(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		r := service.NewMemFailureRepository()
		for i := 0; i < test.quantity; i++ {
			f := &model.Failure{
				ID:         fmt.Sprintf("id-%d", i),
				NodeID:     fmt.Sprintf("nodeid-%d", i),
				Definition: fmt.Sprintf("definition-%d", i),
				State:      types.UnknownFailureState,
			}
			// Store the failures.
			err := r.Store(f)
			assert.NoError(err)

			// Check.
			fGot, ok := r.Get(f.ID)
			if assert.True(ok) {
				assert.Equal(f, fGot)
			}
		}
	}
}

func TestDeleteFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the repository.
	r := service.NewMemFailureRepository()

	// Store a failure and check is there.
	f := &model.Failure{ID: "test"}
	err := r.Store(f)
	require.NoError(err)
	_, ok := r.Get(f.ID)
	require.True(ok)

	// Delete and check is missing
	r.Delete(f.ID)
	_, ok = r.Get(f.ID)
	assert.False(ok)

}

func TestGetFailureMissing(t *testing.T) {
	assert := assert.New(t)
	r := service.NewMemFailureRepository()
	fGot, ok := r.Get("wrong-id")
	if assert.False(ok) {
		assert.Nil(fGot)
	}
}

func TestGetAllFailures(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		r := service.NewMemFailureRepository()
		for i := 0; i < test.quantity; i++ {
			f := &model.Failure{
				ID:         fmt.Sprintf("id-%d", i),
				NodeID:     fmt.Sprintf("nodeid-%d", i),
				Definition: fmt.Sprintf("definition-%d", i),
				State:      types.UnknownFailureState,
			}
			// Store the failures.
			err := r.Store(f)
			assert.NoError(err)

		}
		// Check.
		fsGot := r.GetAll()
		assert.Len(fsGot, test.quantity)
	}
}
