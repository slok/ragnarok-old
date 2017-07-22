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
				ID:            fmt.Sprintf("id-%d", i),
				NodeID:        fmt.Sprintf("nodeid-%d", i),
				Definition:    fmt.Sprintf("definition-%d", i),
				CurrentState:  types.UnknownFailureState,
				ExpectedState: types.EnabledFailureState,
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
				ID:            fmt.Sprintf("id-%d", i),
				NodeID:        fmt.Sprintf("nodeid-%d", i),
				Definition:    fmt.Sprintf("definition-%d", i),
				CurrentState:  types.UnknownFailureState,
				ExpectedState: types.EnabledFailureState,
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

func TestGetAllByNodeFailures(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		nodeFailures map[string]int
	}{
		{
			nodeFailures: map[string]int{
				"node1": 1,
			},
		},
		{

			nodeFailures: map[string]int{
				"node2": 2,
				"node3": 4,
			},
		},
		{
			nodeFailures: map[string]int{
				"node3":  8,
				"node4":  16,
				"node5":  32,
				"node6":  64,
				"node7":  128,
				"node8":  256,
				"node9":  512,
				"node10": 1024,
			},
		},
	}

	for _, test := range tests {
		r := service.NewMemFailureRepository()
		// For each node.
		for nID, q := range test.nodeFailures {
			// For each failure per node.
			for i := 0; i < q; i++ {
				f := &model.Failure{
					ID:            fmt.Sprintf("id-%d", i),
					NodeID:        nID,
					Definition:    fmt.Sprintf("definition-%s-f%d-", nID, i),
					CurrentState:  types.UnknownFailureState,
					ExpectedState: types.EnabledFailureState,
				}
				// Store the failures.
				err := r.Store(f)
				require.NoError(err)
			}
			// Check.
			fsGot := r.GetAllByNode(nID)
			if assert.Len(fsGot, q) {
				for _, f := range fsGot {
					assert.Equal(nID, f.NodeID)
				}
			}
		}
	}
}

func TestGetAllByNodeFailuresMissing(t *testing.T) {
	assert := assert.New(t)

	r := service.NewMemFailureRepository()
	fsGot := r.GetAllByNode("wrongID")
	assert.Empty(fsGot)
}

func TestDeleteFailureByNode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the repository.
	r := service.NewMemFailureRepository()

	// Store failures on different nodes.
	f11 := &model.Failure{ID: "test1", NodeID: "nid1"}
	f21 := &model.Failure{ID: "test2", NodeID: "nid2"}
	f22 := &model.Failure{ID: "test3", NodeID: "nid2"}
	require.NoError(r.Store(f11))
	require.NoError(r.Store(f21))
	require.NoError(r.Store(f22))

	fsGot := r.GetAllByNode(f11.NodeID)
	require.Len(fsGot, 1)
	fsGot = r.GetAllByNode(f21.NodeID)
	require.Len(fsGot, 2)

	// Delete one and check nodes length.
	r.Delete(f21.ID)
	fsGot = r.GetAllByNode(f11.NodeID)
	assert.Len(fsGot, 1)
	fsGot = r.GetAllByNode(f21.NodeID)
	assert.Len(fsGot, 1)
}
