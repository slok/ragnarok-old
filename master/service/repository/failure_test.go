package repository_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/master/service/repository"
)

func TestMemStoreFailure(t *testing.T) {
	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		t.Run(string(test.quantity), func(t *testing.T) {
			assert := assert.New(t)
			r := repository.NewMemFailure()
			for i := 0; i < test.quantity; i++ {
				f := &v1.Failure{
					Metadata: api.ObjectMeta{
						ID: fmt.Sprintf("id-%d", i),
					},
					Spec: v1.FailureSpec{},
					Status: v1.FailureStatus{
						CurrentState:  v1.UnknownFailureState,
						ExpectedState: v1.EnabledFailureState,
					},
				}
				// Store the failures.
				err := r.Store(f)
				assert.NoError(err)

				// Check.
				fGot, ok := r.Get(f.Metadata.ID)
				if assert.True(ok) {
					assert.Equal(f, fGot)
				}
			}
		})
	}
}

func TestMemDeleteFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the repository.
	r := repository.NewMemFailure()

	// Store a failure and check is there.
	f := &v1.Failure{Metadata: api.ObjectMeta{ID: "test"}}
	err := r.Store(f)
	require.NoError(err)
	_, ok := r.Get(f.Metadata.ID)
	require.True(ok)

	// Delete and check is missing
	r.Delete(f.Metadata.ID)
	_, ok = r.Get(f.Metadata.ID)
	assert.False(ok)
}

func TestMemGetFailureMissing(t *testing.T) {
	assert := assert.New(t)
	r := repository.NewMemFailure()
	fGot, ok := r.Get("wrong-id")
	if assert.False(ok) {
		assert.Nil(fGot)
	}
}

func TestMemGetAllFailures(t *testing.T) {
	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		t.Run(string(test.quantity), func(t *testing.T) {
			assert := assert.New(t)
			r := repository.NewMemFailure()
			for i := 0; i < test.quantity; i++ {
				f := &v1.Failure{
					Metadata: api.ObjectMeta{
						ID: fmt.Sprintf("id-%d", i),
					},
					Spec: v1.FailureSpec{},
					Status: v1.FailureStatus{
						CurrentState:  v1.UnknownFailureState,
						ExpectedState: v1.EnabledFailureState,
					},
				}
				// Store the failures.
				err := r.Store(f)
				assert.NoError(err)

			}
			// Check.
			fsGot := r.GetAll()
			assert.Len(fsGot, test.quantity)
		})
	}
}

/*
func TestMemGetNotStaleByNodeFailures(t *testing.T) {
	tests := []struct {
		name              string
		nodeFailures      map[string]int
		nodeStaleFailures map[string]int
	}{
		{
			name: "Get a single failure a node and ignore stale failures.",
			nodeFailures: map[string]int{
				"node1": 1,
			},
			nodeStaleFailures: map[string]int{},
		},
		{
			name: "Get multiple failures in multiple nodes and ignore stale failures.",
			nodeFailures: map[string]int{
				"node2": 2,
				"node3": 4,
			},
			nodeStaleFailures: map[string]int{
				"node2": 2,
				"node4": 4,
			},
		},
		{
			name: "Get multiple failures in a lot of nodes and and ignore stale failures.",
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
			nodeStaleFailures: map[string]int{
				"node3": 2,
				"node4": 4,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			r := repository.NewMemFailure()
			// For each node.
			for nID, q := range test.nodeFailures {
				// For each failure per node.
				for i := 0; i < q; i++ {
					f := &v1.Failure{
						Metadata: v1.FailureMetadata{
							ID:     fmt.Sprintf("id-%d", i),
							NodeID: nID,
						},
						Spec: v1.FailureSpec{},
						Status: v1.FailureStatus{
							CurrentState:  v1.EnabledFailureState,
							ExpectedState: v1.EnabledFailureState,
						},
					}
					// Store the failures.
					err := r.Store(f)
					require.NoError(err)
				}
			}

			// For each stale failure per node.
			for nID, q := range test.nodeStaleFailures {
				// For each failure per node.
				for i := 0; i < q; i++ {
					f := &v1.Failure{
						Metadata: v1.FailureMetadata{
							ID:     fmt.Sprintf("id-st-%d", i),
							NodeID: nID,
						},
						Spec: v1.FailureSpec{},
						Status: v1.FailureStatus{
							CurrentState:  v1.StaleFailureState,
							ExpectedState: v1.UnknownFailureState,
						},
					}
					// Store the failures.
					err := r.Store(f)
					require.NoError(err)
				}
			}

			// Check.
			for nID, q := range test.nodeFailures {
				fsGot := r.GetNotStaleByNode(nID)
				assert.Len(fsGot, q)
				if assert.Len(fsGot, q) {
					for _, f := range fsGot {
						assert.Equal(nID, f.Metadata.NodeID)
					}
				}
			}
		})
	}
}

func TestMemGetAllByNodeFailures(t *testing.T) {
	tests := []struct {
		name         string
		nodeFailures map[string]int
	}{
		{
			name: "Get a single failure in a node",
			nodeFailures: map[string]int{
				"node1": 1,
			},
		},
		{
			name: "Get multiple failures in multiple nodes",
			nodeFailures: map[string]int{
				"node2": 2,
				"node3": 4,
			},
		},
		{
			name: "Get a multiple failures in a lot of nodes",
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
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			r := repository.NewMemFailure()
			// For each node.
			for nID, q := range test.nodeFailures {
				// For each failure per node.
				for i := 0; i < q; i++ {
					f := &v1.Failure{
						Metadata: v1.FailureMetadata{
							ID:     fmt.Sprintf("id-%d", i),
							NodeID: nID,
						},
						Spec: v1.FailureSpec{},
						Status: v1.FailureStatus{
							CurrentState:  v1.UnknownFailureState,
							ExpectedState: v1.EnabledFailureState,
						},
					}
					// Store the failures.
					err := r.Store(f)
					require.NoError(err)
				}
				// Check.
				fsGot := r.GetAllByNode(nID)
				if assert.Len(fsGot, q) {
					for _, f := range fsGot {
						assert.Equal(nID, f.Metadata.NodeID)
					}
				}
			}
		})
	}
}

func TestMemGetAllByNodeFailuresMissing(t *testing.T) {
	assert := assert.New(t)

	r := repository.NewMemFailure()
	fsGot := r.GetAllByNode("wrongID")
	assert.Empty(fsGot)
}

func TestDeleteFailureByNode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the repository.
	r := repository.NewMemFailure()

	// Store failures on different nodes.
	f11 := &v1.Failure{Metadata: v1.FailureMetadata{ID: "test1", NodeID: "nid1"}}
	f21 := &v1.Failure{Metadata: v1.FailureMetadata{ID: "test2", NodeID: "nid2"}}
	f22 := &v1.Failure{Metadata: v1.FailureMetadata{ID: "test3", NodeID: "nid2"}}
	require.NoError(r.Store(f11))
	require.NoError(r.Store(f21))
	require.NoError(r.Store(f22))

	fsGot := r.GetAllByNode(f11.Metadata.NodeID)
	require.Len(fsGot, 1)
	fsGot = r.GetAllByNode(f21.Metadata.NodeID)
	require.Len(fsGot, 2)

	// Delete one and check nodes length.
	r.Delete(f21.Metadata.ID)
	fsGot = r.GetAllByNode(f11.Metadata.NodeID)
	assert.Len(fsGot, 1)
	fsGot = r.GetAllByNode(f21.Metadata.NodeID)
	assert.Len(fsGot, 1)
}
*/
