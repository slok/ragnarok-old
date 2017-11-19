package v1

import (
	"time"

	"github.com/slok/ragnarok/api"
)

const (
	// ExperimentKind is the kind a experiment.
	ExperimentKind = "experiment"
	// ExperimentVersion is the version of an experiment.
	ExperimentVersion = "chaos/v1"
)

// ExperimentStatus is the status after the creation of the Experiment.
type ExperimentStatus struct {
	// FailureIDs are the IDs of the failures that have been created.
	FailureIDs []string  `json:"failureIDs,omitempty"`
	Creation   time.Time `json:"creation,omitempty"` // Creation is when the creation of the node happenned.
}

// ExperimentFailureTemplate is the template of a failure
type ExperimentFailureTemplate struct {
	Spec FailureSpec `json:"spec,omitempty"`
}

// ExperimentSpec is the spec of the experiment
type ExperimentSpec struct {
	// Name is the name of the experiment.
	Name string `json:"name,omitempty"`
	// Description is the description of the experiment.
	Description string `json:"description,omitempty"`
	// Selector is the map of key-value pairs that will match the desired nodes where the attacks
	// will be injected.
	Selector map[string]string         `json:"selector,omitempty"`
	Template ExperimentFailureTemplate `json:"template,omitempty"`
}

// Experiment is only a simple group of failures that are being injected in
// the targets that have been selected by the experiment using selectors.
type Experiment struct {
	api.TypeMeta

	Metadata api.ObjectMeta   `json:"metadata,omitempty"`
	Spec     ExperimentSpec   `json:"spec,omitempty"`
	Status   ExperimentStatus `json:"status,omitempty"`
}

// NewExperiment is a plain Experiment object contructor.
func NewExperiment() Experiment {
	return Experiment{
		TypeMeta: api.TypeMeta{
			Kind:    ExperimentKind,
			Version: ExperimentVersion,
		},
	}
}

// GetObjectMetadata satisfies object interface.
func (e *Experiment) GetObjectMetadata() api.ObjectMeta {
	return e.Metadata
}
