package v1_test

import (
	"errors"
	"testing"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/apimachinery/validator"
	clichaosv1 "github.com/slok/ragnarok/client/api/chaos/v1"
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
			name:        "One new failure should be created without error.",
			invalidObj:  false,
			createError: false,
			expErr:      false,
		},
		{
			name:        "Creating an invalid failure should return an error.",
			invalidObj:  true,
			createError: false,
			expErr:      true,
		},
		{
			name:        "Creation failure error on repository should return an error.",
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

			f := &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))
			mr := &mrepository.Client{}
			mr.On("Create", mock.Anything).Return(f, createError)

			// Create our client.
			cli := clichaosv1.NewFailureClient(mv, mr)

			// Create the failure and check.
			_, err := cli.Create(f)
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
			name:        "Updating failure should be updated without error.",
			invalidObj:  false,
			updateError: false,
			expErr:      false,
		},
		{
			name:        "Updating an invalid failure should return an error.",
			invalidObj:  true,
			updateError: false,
			expErr:      true,
		},
		{
			name:        "Updating failure error on repository should return an error.",
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

			f := &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))
			mr := &mrepository.Client{}
			mr.On("Update", mock.Anything).Return(f, updateError)

			// Create our client.
			cli := clichaosv1.NewFailureClient(mv, mr)

			// Create the failure and check.
			_, err := cli.Update(f)
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
			name:        "Deleting failure should be deleted without error.",
			id:          "test1",
			expFullID:   "chaos/v1/failure/test1",
			deleteError: false,
			expErr:      false,
		},
		{
			name:        "Deleting failure with an error on the repository should return an error.",
			id:          "test1",
			expFullID:   "chaos/v1/failure/test1",
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
			cli := clichaosv1.NewFailureClient(mv, mr)

			// Create the failure and check.
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
			name:      "Getting failure should retrieive without error.",
			id:        "test1",
			expFullID: "chaos/v1/failure/test1",
			getError:  false,
			expErr:    false,
		},
		{
			name:      "Getting a failure with an error from the repository should return an error.",
			id:        "test1",
			expFullID: "chaos/v1/failure/test1",
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

			f := &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mr := &mrepository.Client{}
			mr.On("Get", test.expFullID).Once().Return(f, getError)

			// Create our client.
			cli := clichaosv1.NewFailureClient(mv, mr)

			// Create the failure and check.
			gotN, err := cli.Get(test.id)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
				assert.Equal(f, gotN)
			}
		})
	}
}

func TestMemNodeCliList(t *testing.T) {
	tests := []struct {
		name        string
		objList     []api.Object
		expNodeList []*chaosv1.Failure
		listError   bool
		expErr      bool
	}{
		{
			name: "Getting failure List should retrieive without error.",
			objList: []api.Object{
				&chaosv1.Failure{Metadata: api.ObjectMeta{ID: "failure1"}},
				&chaosv1.Failure{Metadata: api.ObjectMeta{ID: "failure2"}},
				&chaosv1.Failure{Metadata: api.ObjectMeta{ID: "failure3"}},
			},
			expNodeList: []*chaosv1.Failure{
				&chaosv1.Failure{Metadata: api.ObjectMeta{ID: "failure1"}},
				&chaosv1.Failure{Metadata: api.ObjectMeta{ID: "failure2"}},
				&chaosv1.Failure{Metadata: api.ObjectMeta{ID: "failure3"}},
			},
			listError: false,
			expErr:    false,
		},
		{
			name:      "Getting a failure list with an error from the repository should return an error.",
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
			cli := clichaosv1.NewFailureClient(mv, mr)

			// Create the failure and check.
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
