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

func TestMemExperimentCliCreate(t *testing.T) {
	tests := []struct {
		name        string
		invalidObj  bool
		createError bool
		expErr      bool
	}{
		{
			name:        "One new experiment should be created without error.",
			invalidObj:  false,
			createError: false,
			expErr:      false,
		},
		{
			name:        "Creating an invalid experiment should return an error.",
			invalidObj:  true,
			createError: false,
			expErr:      true,
		},
		{
			name:        "Creation experiment error on repository should return an error.",
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

			e := &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))
			mr := &mrepository.Client{}
			mr.On("Create", mock.Anything).Return(e, createError)

			// Create our client.
			cli := clichaosv1.NewExperimentClient(mv, mr)

			// Create the experiment and check.
			_, err := cli.Create(e)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestMemExperimentCliUpdate(t *testing.T) {
	tests := []struct {
		name        string
		invalidObj  bool
		updateError bool
		expErr      bool
	}{
		{
			name:        "Updating experiment should be updated without error.",
			invalidObj:  false,
			updateError: false,
			expErr:      false,
		},
		{
			name:        "Updating an invalid experiment should return an error.",
			invalidObj:  true,
			updateError: false,
			expErr:      true,
		},
		{
			name:        "Updating experiment error on repository should return an error.",
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

			e := &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mv.On("Validate", mock.Anything).Return(validator.ErrorList(validationErrs))
			mr := &mrepository.Client{}
			mr.On("Update", mock.Anything).Return(e, updateError)

			// Create our client.
			cli := clichaosv1.NewExperimentClient(mv, mr)

			// Create the experiment and check.
			_, err := cli.Update(e)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestMemExperimentCliDelete(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		deleteError bool
		expFullID   string
		expErr      bool
	}{
		{
			name:        "Deleting experiment should be deleted without error.",
			id:          "test1",
			expFullID:   "chaos/v1/experiment/test1",
			deleteError: false,
			expErr:      false,
		},
		{
			name:        "Deleting experiment with an error on the repository should return an error.",
			id:          "test1",
			expFullID:   "chaos/v1/experiment/test1",
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
			cli := clichaosv1.NewExperimentClient(mv, mr)

			// Create the experiment and check.
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

func TestMemExperimentCliGet(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		getError  bool
		expFullID string
		expErr    bool
	}{
		{
			name:      "Getting experiment should retrieive without error.",
			id:        "test1",
			expFullID: "chaos/v1/experiment/test1",
			getError:  false,
			expErr:    false,
		},
		{
			name:      "Getting a experiment with an error from the repository should return an error.",
			id:        "test1",
			expFullID: "chaos/v1/experiment/test1",
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

			e := &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID:     "test",
					Labels: map[string]string{"id": "test"},
				},
			}

			// Mocks.
			mv := &mvalidator.ObjectValidator{}
			mr := &mrepository.Client{}
			mr.On("Get", test.expFullID).Once().Return(e, getError)

			// Create our client.
			cli := clichaosv1.NewExperimentClient(mv, mr)

			// Create the experiment and check.
			gotN, err := cli.Get(test.id)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
				assert.Equal(e, gotN)
			}
		})
	}
}

func TestMemExperimentCliList(t *testing.T) {
	tests := []struct {
		name              string
		objList           []api.Object
		expExperimentList []*chaosv1.Experiment
		listError         bool
		expErr            bool
	}{
		{
			name: "Getting experiment List should retrieive without error.",
			objList: []api.Object{
				&chaosv1.Experiment{Metadata: api.ObjectMeta{ID: "experiment1"}},
				&chaosv1.Experiment{Metadata: api.ObjectMeta{ID: "experiment2"}},
				&chaosv1.Experiment{Metadata: api.ObjectMeta{ID: "experiment3"}},
			},
			expExperimentList: []*chaosv1.Experiment{
				&chaosv1.Experiment{Metadata: api.ObjectMeta{ID: "experiment1"}},
				&chaosv1.Experiment{Metadata: api.ObjectMeta{ID: "experiment2"}},
				&chaosv1.Experiment{Metadata: api.ObjectMeta{ID: "experiment3"}},
			},
			listError: false,
			expErr:    false,
		},
		{
			name:      "Getting a experiment list with an error from the repository should return an error.",
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
			cli := clichaosv1.NewExperimentClient(mv, mr)

			// Create the experiment and check.
			gotExperiments, err := cli.List(api.ListOptions{})
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				mr.AssertExpectations(t)
				assert.Equal(test.expExperimentList, gotExperiments)
			}
		})
	}
}
