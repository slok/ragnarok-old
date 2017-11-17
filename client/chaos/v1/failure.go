package v1

import (
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
)

// FailureWatch is the watch object that the Watch will return.
type FailureWatch <-chan *chaosv1.Failure

// FailureListOptions are the options required to list failures
type FailureListOptions struct{}

// FailureInterface has the required logic to manage Failures.
type FailureInterface interface {
	Create(*chaosv1.Failure) (*chaosv1.Failure, error)
	Update(*chaosv1.Failure) (*chaosv1.Failure, error)
	Delete(id string) error
	Get(id string) (*chaosv1.Failure, error)
	List(opts FailureListOptions) ([]*chaosv1.Failure, error)
	Watch(opts FailureListOptions) (FailureWatch, error)
}
