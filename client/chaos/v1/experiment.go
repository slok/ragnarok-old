package v1

import (
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
)

// ExperimentWatch is the watch object that the Watch will return.
type ExperimentWatch <-chan *chaosv1.Experiment

// ExperimentListOptions are the options required to list experiments
type ExperimentListOptions struct{}

// ExperimentInterface has the required logic to manage Experiments.
type ExperimentInterface interface {
	Create(*chaosv1.Experiment) (*chaosv1.Experiment, error)
	Update(*chaosv1.Experiment) (*chaosv1.Experiment, error)
	Delete(id string) error
	Get(id string) (*chaosv1.Experiment, error)
	List(opts ExperimentListOptions) ([]*chaosv1.Experiment, error)
	Watch(opts ExperimentListOptions) (ExperimentWatch, error)
}
