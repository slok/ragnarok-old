package v1_test

import (
	"errors"
	"testing"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/validator"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	mvalidator "github.com/slok/ragnarok/mocks/apimachinery/validator"
	mrepository "github.com/slok/ragnarok/mocks/client/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMemNodeCliCreate(t *testing.T) {
	tests := []struct {
		name        string
		invalidObj  bool
		createError bool
		expErr      bool
	}{
		{
			name:        "One new node should be created without error.",
			invalidObj:  false,
			createError: false,
			expErr:      false,
		},
		{
			name:        "Creating an invalid node should return an error.",
			invalidObj:  true,
			createError: false,
			expErr:      true,
		},
		{
			name:        "Creation node error on repository should return an error.",
			invalidObj:  false,
			createError: true,
			expErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			validationErrs := []error{}
			var createError error

			if test.invalidObj {
				validationErrs = append(validationErrs, errors.New("wanted error"))
			}
			if test.createError {
				createError = errors.New("wanted error")
			}

			n := &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))
			mr := &mrepository.Client{}
			mr.On("Create", mock.Anything).Return(n, createError)

			// Create our client.
			cli := cliclusterv1.NewNodeClient(mv, mr)

			// Create the node and check.
			_, err := cli.Create(n)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestMemNodeCliUpdate(t *testing.T) {
	tests := []struct {
		name        string
		invalidObj  bool
		updateError bool
		expErr      bool
	}{
		{
			name:        "Updating node should be updated without error.",
			invalidObj:  false,
			updateError: false,
			expErr:      false,
		},
		{
			name:        "Updating an invalid node should return an error.",
			invalidObj:  true,
			updateError: false,
			expErr:      true,
		},
		{
			name:        "Updating node error on repository should return an error.",
			invalidObj:  false,
			updateError: true,
			expErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			validationErrs := []error{}
			var updateError error

			if test.invalidObj {
				validationErrs = append(validationErrs, errors.New("wanted error"))
			}
			if test.updateError {
				updateError = errors.New("wanted error")
			}

			n := &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))
			mr := &mrepository.Client{}
			mr.On("Update", mock.Anything).Return(n, updateError)

			// Create our client.
			cli := cliclusterv1.NewNodeClient(mv, mr)

			// Create the node and check.
			_, err := cli.Update(n)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestMemNodeCliDelete(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		deleteError bool
		expFullID   string
		expErr      bool
	}{
		{
			name:        "Deleting node should be deleted without error.",
			id:          "test1",
			expFullID:   "cluster/v1/node/test1",
			deleteError: false,
			expErr:      false,
		},
		{
			name:        "Deleting node with an error on the repository should return an error.",
			id:          "test1",
			expFullID:   "cluster/v1/node/test1",
			deleteError: true,
			expErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var deleteError error

			if test.deleteError {
				deleteError = errors.New("wanted error")
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mr := &mrepository.Client{}
			mr.On("Delete", test.expFullID).Once().Return(deleteError)

			// Create our client.
			cli := cliclusterv1.NewNodeClient(mv, mr)

			// Create the node and check.
			err := cli.Delete(test.id)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestMemNodeCliGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		getError  bool
		expFullID string
		expErr    bool
	}{
		{
			name:      "Getting node should retrieive without error.",
			id:        "test1",
			expFullID: "cluster/v1/node/test1",
			getError:  false,
			expErr:    false,
		},
		{
			name:      "Getting a node with an error from the repository should return an error.",
			id:        "test1",
			expFullID: "cluster/v1/node/test1",
			getError:  true,
			expErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var getError error

			if test.getError {
				getError = errors.New("wanted error")
			}

			n := &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mr := &mrepository.Client{}
			mr.On("Get", test.expFullID).Once().Return(n, getError)

			// Create our client.
			cli := cliclusterv1.NewNodeClient(mv, mr)

			// Create the node and check.
			gotN, err := cli.Get(test.id)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
				assert.Equal(n, gotN)
			}
		})
	}
}

func TestMemNodeCliList(t *testing.T) {
	tests := []struct {
		name        string
		objList     []api.Object
		expNodeList []*clusterv1.Node
		listError   bool
		expErr      bool
	}{
		{
			name: "Getting node List should retrieive without error.",
			objList: []api.Object{
				&clusterv1.Node{Metadata: api.ObjectMeta{ID: "node1"}},
				&clusterv1.Node{Metadata: api.ObjectMeta{ID: "node2"}},
				&clusterv1.Node{Metadata: api.ObjectMeta{ID: "node3"}},
			},
			expNodeList: []*clusterv1.Node{
				&clusterv1.Node{Metadata: api.ObjectMeta{ID: "node1"}},
				&clusterv1.Node{Metadata: api.ObjectMeta{ID: "node2"}},
				&clusterv1.Node{Metadata: api.ObjectMeta{ID: "node3"}},
			},
			listError: false,
			expErr:    false,
		},
		{
			name:      "Getting a node list with an error from the repository should return an error.",
			listError: true,
			expErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var listError error

			if test.listError {
				listError = errors.New("wanted error")
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mr := &mrepository.Client{}
			mr.On("List", mock.Anything).Return(test.objList, listError)

			// Create our client.
			cli := cliclusterv1.NewNodeClient(mv, mr)

			// Create the node and check.
			gotNodes, err := cli.List(api.ListOptions{})
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
				assert.Equal(test.expNodeList, gotNodes)
			}
		})
	}
}
