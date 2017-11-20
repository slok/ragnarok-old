package v1

import (
	"fmt"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	apiutil "github.com/slok/ragnarok/api/util"
	"github.com/slok/ragnarok/apimachinery/validator"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/repository"
)

var experimentObjType = api.TypeMeta{Kind: chaosv1.ExperimentKind, Version: chaosv1.ExperimentVersion}

// ExperimentClientInterface has the required logic to manage Experiment.
type ExperimentClientInterface interface {
	Create(experiment *chaosv1.Experiment) (*chaosv1.Experiment, error)
	Update(experiment *chaosv1.Experiment) (*chaosv1.Experiment, error)
	Delete(id string) error
	Get(id string) (*chaosv1.Experiment, error)
	List(opts api.ListOptions) ([]*chaosv1.Experiment, error)
	Watch(opts api.ListOptions) (watch.Watcher, error)
	// TODO Patch
}

// ExperimentClient has the required logic to manage Experiments.
type ExperimentClient struct {
	validator validator.ObjectValidator
	repoCli   repository.Client
}

// NewExperimentClient returns a new ExperimentClient.
func NewExperimentClient(validator validator.ObjectValidator, repoCli repository.Client) *ExperimentClient {
	return &ExperimentClient{
		validator: validator,
		repoCli:   repoCli,
	}
}

func (e *ExperimentClient) typeAssert(obj api.Object) (*chaosv1.Experiment, error) {
	experiment, ok := obj.(*chaosv1.Experiment)
	if !ok {
		return nil, fmt.Errorf("could not make the type assertion from obj to experiment. Wrong type")
	}
	return experiment, nil
}

func (e *ExperimentClient) validate(experiment *chaosv1.Experiment) error {
	// Check valid object.
	if errs := e.validator.Validate(experiment); len(errs) > 0 {
		return fmt.Errorf("error on validation: %s", errs)
	}
	return nil
}

// Create satisfies ExperimentClientInterface interface.
func (e *ExperimentClient) Create(experiment *chaosv1.Experiment) (*chaosv1.Experiment, error) {
	// Check valid object.
	if err := e.validate(experiment); err != nil {
		return nil, err
	}

	obj, err := e.repoCli.Create(experiment)
	if err != nil {
		return nil, err
	}
	return e.typeAssert(obj)
}

// Update satisfies ExperimentClientInterface interface.
func (e *ExperimentClient) Update(experiment *chaosv1.Experiment) (*chaosv1.Experiment, error) {
	// Check valid object.
	if err := e.validate(experiment); err != nil {
		return nil, err
	}

	obj, err := e.repoCli.Update(experiment)
	if err != nil {
		return nil, err
	}
	return e.typeAssert(obj)
}

// Delete satisfies ExperimentClientInterface interface.
func (e *ExperimentClient) Delete(id string) error {
	// get the full ID
	fullID := apiutil.GetFullIDFromType(experimentObjType, id)
	return e.repoCli.Delete(fullID)
}

// Get satisfies ExperimentClientInterface interface.
func (e *ExperimentClient) Get(id string) (*chaosv1.Experiment, error) {
	fullID := apiutil.GetFullIDFromType(experimentObjType, id)
	obj, err := e.repoCli.Get(fullID)
	if err != nil {
		return nil, err
	}
	return e.typeAssert(obj)
}

// List satisfies ExperimentClientInterface interface.
func (e *ExperimentClient) List(opts api.ListOptions) ([]*chaosv1.Experiment, error) {
	experiments := []*chaosv1.Experiment{}

	objs, err := e.repoCli.List(opts)
	if err != nil {
		return experiments, err
	}

	experiments = make([]*chaosv1.Experiment, len(objs))
	for i, obj := range objs {
		experiment, err := e.typeAssert(obj)
		if err != nil {
			return experiments, err
		}
		experiments[i] = experiment
	}

	return experiments, nil
}

// Watch satisfies ExperimentClientInterface interface.
func (e *ExperimentClient) Watch(opts api.ListOptions) (watch.Watcher, error) {
	return nil, nil
}
