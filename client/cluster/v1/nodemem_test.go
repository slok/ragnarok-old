package v1_test

import (
	"errors"
	"sort"
	"testing"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/validator"
	cliclusterv1 "github.com/slok/ragnarok/client/cluster/v1"
	mvalidator "github.com/slok/ragnarok/mocks/apimachinery/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMemNodeCliCreate(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		registry   map[string]*clusterv1.Node
		newID      string
		invalidObj bool
		expErr     bool
	}{
		{
			name:       "One new node should be created without error.",
			registry:   map[string]*clusterv1.Node{},
			newID:      "node1",
			invalidObj: false,
			expErr:     false,
		},
		{
			name:       "Creating an invalid node should return an error.",
			registry:   map[string]*clusterv1.Node{},
			newID:      "node1",
			invalidObj: true,
			expErr:     true,
		},
		{
			name:       "Storing a node that already was stored should return an error.",
			registry:   map[string]*clusterv1.Node{"node1": nil},
			newID:      "node1",
			invalidObj: false,
			expErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			validationErrs := []error{}
			if test.invalidObj {
				validationErrs = append(validationErrs, errors.New("wanted error"))
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))

			cli := cliclusterv1.NewNodeMem(mv, test.registry)
			n := &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     test.newID,
					Labels: map[string]string{"id": test.newID},
				},
			}
			_, err := cli.Create(n)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				// Check stored node is ok
				nGot, ok := test.registry[test.newID]
				if assert.True(ok) {
					assert.Equal(n, nGot)
				}
			}
		})
	}
}

func TestMemNodeCliUpdate(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name        string
		registry    map[string]*clusterv1.Node
		updateNode  *clusterv1.Node
		invalidObj  bool
		expRegistry map[string]*clusterv1.Node
		expErr      bool
	}{
		{
			name:     "When updating a new node sould return an error.",
			registry: map[string]*clusterv1.Node{},
			updateNode: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "node1",
					Labels: map[string]string{"test": "wrong"},
				},
			},
			invalidObj:  false,
			expRegistry: map[string]*clusterv1.Node{},
			expErr:      true,
		},
		{
			name:     "Updating an invalid node should return an error.",
			registry: map[string]*clusterv1.Node{},
			updateNode: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "node1",
					Labels: map[string]string{"test": "wrong"},
				},
			},
			invalidObj:  true,
			expRegistry: map[string]*clusterv1.Node{},
			expErr:      true,
		},
		{
			name: "Storing a node that already was stored shouldupdate ok.",
			registry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node1",
						Labels: map[string]string{"test": "wrong"},
					},
				},
			},
			updateNode: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "node1",
					Labels: map[string]string{"test": "ok"},
				},
			},
			invalidObj: false,
			expRegistry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node1",
						Labels: map[string]string{"test": "ok"},
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			validationErrs := []error{}
			if test.invalidObj {
				validationErrs = append(validationErrs, errors.New("wanted error"))
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))

			cli := cliclusterv1.NewNodeMem(mv, test.registry)
			_, err := cli.Update(test.updateNode)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expRegistry, test.registry)
			}
		})
	}
}

func TestMemNodeCliGet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		registry map[string]*clusterv1.Node
		id       string
		expNode  *clusterv1.Node
		expErr   bool
	}{
		{
			name: "Getting a present node shouldn't return an error.",
			registry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node1",
					},
				},
			},
			id: "node1",
			expNode: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
				},
			},
			expErr: false,
		},
		{
			name:     "Getting a missing node should return an error.",
			registry: map[string]*clusterv1.Node{},
			id:       "node1",
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mv := &mvalidator.ObjectValidator{}

			cli := cliclusterv1.NewNodeMem(mv, test.registry)
			n, err := cli.Get(test.id)
			if test.expErr {
				assert.Error(err)
			} else {
				// Check retrieved node is ok.
				assert.NoError(err)
				assert.Equal(test.expNode, n)
			}
		})
	}
}

func TestMemNodeCliDelete(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		registry map[string]*clusterv1.Node
		id       string
		expNodes map[string]*clusterv1.Node
		expErr   bool
	}{
		{
			name: "Deleting a node shouldn't return an error.",
			registry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node1",
					},
				},
				"node2": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node2",
					},
				},
			},
			id: "node1",
			expNodes: map[string]*clusterv1.Node{
				"node2": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node2",
					},
				},
			},
			expErr: false,
		},
		{
			name:     "delering a missing node shouldn't return an error.",
			registry: map[string]*clusterv1.Node{},
			id:       "node1",
			expNodes: map[string]*clusterv1.Node{},
			expErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mv := &mvalidator.ObjectValidator{}

			cli := cliclusterv1.NewNodeMem(mv, test.registry)
			err := cli.Delete(test.id)
			if test.expErr {
				assert.Error(err)
			} else {
				// Check retrieved node is ok.
				assert.NoError(err)
				assert.Equal(test.expNodes, test.registry)
			}
		})
	}
}

func TestMemNodeCliList(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		registry map[string]*clusterv1.Node
		selector map[string]string
		expNodes []*clusterv1.Node
		expErr   bool
	}{
		{
			name:     "No selectors should return all the registry.",
			selector: map[string]string{},
			registry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node1",
						Labels: map[string]string{"kind": "master", "region": "eu-west-1"},
					},
				},
				"node2": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node2",
						Labels: map[string]string{"kind": "node", "region": "eu-central-1"},
					},
				},
				"node3": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node3",
						Labels: map[string]string{"kind": "master", "region": "eu-central-1"},
					},
				},
			},
			expNodes: []*clusterv1.Node{
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node1",
						Labels: map[string]string{"kind": "master", "region": "eu-west-1"},
					},
				},
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node2",
						Labels: map[string]string{"kind": "node", "region": "eu-central-1"},
					},
				},
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node3",
						Labels: map[string]string{"kind": "master", "region": "eu-central-1"},
					},
				},
			},
			expErr: false,
		},
		{
			name:     "One selectors should return only the selector ones.",
			selector: map[string]string{"region": "eu-central-1"},
			registry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node1",
						Labels: map[string]string{"kind": "master", "region": "eu-west-1"},
					},
				},
				"node2": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node2",
						Labels: map[string]string{"kind": "node", "region": "eu-central-1"},
					},
				},
				"node3": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node3",
						Labels: map[string]string{"kind": "master", "region": "eu-central-1"},
					},
				},
			},
			expNodes: []*clusterv1.Node{
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node2",
						Labels: map[string]string{"kind": "node", "region": "eu-central-1"},
					},
				},
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node3",
						Labels: map[string]string{"kind": "master", "region": "eu-central-1"},
					},
				},
			},
			expErr: false,
		},
		{
			name:     "Multiple selectors should return only the selector ones.",
			selector: map[string]string{"kind": "master", "region": "eu-central-1"},
			registry: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node1",
						Labels: map[string]string{"kind": "master", "region": "eu-west-1"},
					},
				},
				"node2": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node2",
						Labels: map[string]string{"kind": "node", "region": "eu-central-1"},
					},
				},
				"node3": &clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node3",
						Labels: map[string]string{"kind": "master", "region": "eu-central-1"},
					},
				},
			},
			expNodes: []*clusterv1.Node{
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID:     "node3",
						Labels: map[string]string{"kind": "master", "region": "eu-central-1"},
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mv := &mvalidator.ObjectValidator{}

			cli := cliclusterv1.NewNodeMem(mv, test.registry)
			listOpts := cliclusterv1.NodeListOptions{Selector: test.selector}
			nodes, err := cli.List(listOpts)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				sort.Slice(nodes, func(i, j int) bool {
					return nodes[i].Metadata.ID < nodes[j].Metadata.ID
				})
				assert.Equal(test.expNodes, nodes)
			}
		})
	}
}
