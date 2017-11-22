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

var failureObjType = api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion}

// FailureClientInterface has the required logic to manage Failure.
type FailureClientInterface interface {
	Create(failure *chaosv1.Failure) (*chaosv1.Failure, error)
	Update(failure *chaosv1.Failure) (*chaosv1.Failure, error)
	Delete(id string) error
	Get(id string) (*chaosv1.Failure, error)
	List(opts api.ListOptions) ([]*chaosv1.Failure, error)
	Watch(opts api.ListOptions) (watch.Watcher, error)
	// TODO Patch
}

// FailureClient has the required logic to manage Failures.
type FailureClient struct {
	validator validator.ObjectValidator
	repoCli   repository.Client
}

// NewFailureClient returns a new FailureClient.
func NewFailureClient(validator validator.ObjectValidator, repoCli repository.Client) *FailureClient {
	return &FailureClient{
		validator: validator,
		repoCli:   repoCli,
	}
}

func (f *FailureClient) typeAssert(obj api.Object) (*chaosv1.Failure, error) {
	failure, ok := obj.(*chaosv1.Failure)
	if !ok {
		return nil, fmt.Errorf("could not make the type assertion from obj to failure. Wrong type")
	}
	return failure, nil
}

func (f *FailureClient) validate(failure *chaosv1.Failure) error {
	// Check valid object.
	if errs := f.validator.Validate(failure); len(errs) > 0 {
		return fmt.Errorf("error on validation: %s", errs)
	}
	return nil
}

// Create satisfies FailureClientInterface interface.
func (f *FailureClient) Create(failure *chaosv1.Failure) (*chaosv1.Failure, error) {
	// Check valid object.
	if err := f.validate(failure); err != nil {
		return nil, err
	}

	obj, err := f.repoCli.Create(failure)
	if err != nil {
		return nil, err
	}
	return f.typeAssert(obj)
}

// Update satisfies FailureClientInterface interface.
func (f *FailureClient) Update(failure *chaosv1.Failure) (*chaosv1.Failure, error) {
	// Check valid object.
	if err := f.validate(failure); err != nil {
		return nil, err
	}

	obj, err := f.repoCli.Update(failure)
	if err != nil {
		return nil, err
	}
	return f.typeAssert(obj)
}

// Delete satisfies FailureClientInterface interface.
func (f *FailureClient) Delete(id string) error {
	// get the full ID
	fullID := apiutil.GetFullIDFromType(failureObjType, id)
	return f.repoCli.Delete(fullID)
}

// Get satisfies FailureClientInterface interface.
func (f *FailureClient) Get(id string) (*chaosv1.Failure, error) {
	fullID := apiutil.GetFullIDFromType(failureObjType, id)
	obj, err := f.repoCli.Get(fullID)
	if err != nil {
		return nil, err
	}
	return f.typeAssert(obj)
}

// List satisfies FailureClientInterface interface.
func (f *FailureClient) List(opts api.ListOptions) ([]*chaosv1.Failure, error) {
	opts.TypeMeta = chaosv1.FailureTypeMeta
	failures := []*chaosv1.Failure{}

	objs, err := f.repoCli.List(opts)
	if err != nil {
		return failures, err
	}

	failures = make([]*chaosv1.Failure, len(objs))
	for i, obj := range objs {
		failure, err := f.typeAssert(obj)
		if err != nil {
			return failures, err
		}
		failures[i] = failure
	}

	return failures, nil
}

// Watch satisfies FailureClientInterface interface.
func (f *FailureClient) Watch(opts api.ListOptions) (watch.Watcher, error) {
	opts.TypeMeta = chaosv1.FailureTypeMeta
	return f.repoCli.Watch(opts)
}
