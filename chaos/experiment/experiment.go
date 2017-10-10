package experiment

import (
	"github.com/slok/ragnarok/chaos/failure"
)

// Status is the status after the creation of the Experiment.
type Status struct {
	// FailureIDs are the IDs of the failures that have been created.
	FailureIDs []string `yaml:"failureIDs,omitempty"`
}

// Experiment is only a simple group of failures that are being injected in
// the targets that have been selected by the experiment using selectors.
type Experiment struct {
	// ID is the id of the experiment
	ID string `yaml:"id,omitempty"`

	// Name is the name of the experiment.
	Name string `yaml:"name,omitempty"`

	// Description is the description of the experiment.
	Description string `yaml:"description,omitempty"`

	// Selector is the map of key-value pairs that will match the desired nodes where the attacks
	// will be injected.
	Selector map[string]string `yaml:"selector,omitempty"`

	// Definition is the definition of a Failure.
	Spec failure.Failure `yaml:"spec,omitempty"`

	// Status is the status of the experiment.
	Status Status `yaml:"status,omitempty"`
}
