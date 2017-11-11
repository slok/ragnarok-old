package v1

import (
	"time"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/attack"
)

const (
	// FailureKind is the kind a failure.
	FailureKind = "failure"
	// FailureVersion is the version of the failure
	FailureVersion = "chaos/v1"
)

// FailureState is the state a failure can be.
type FailureState int

const (
	// UnknownFailureState is an unknown status.
	UnknownFailureState FailureState = iota
	// EnabledFailureState is when the failure should be making stuff.
	EnabledFailureState
	// ExecutingFailureState is when the failure its making stuff.
	ExecutingFailureState
	// RevertingFailureState is when the failure is being reverted.
	RevertingFailureState
	// DisabledFailureState is when the failure is should be not making stuff.
	DisabledFailureState
	// StaleFailureState is when the failure has go through alll the lifecycle and should be archived ((reverted already).
	StaleFailureState
	// ErroredFailureState is when the failure is not making stuff (due to an error).
	ErroredFailureState
	// ErroredRevertingFailureState is when the failure is not making stuff (due to an error reverting).
	ErroredRevertingFailureState
)

func (f FailureState) String() string {
	switch f {
	case EnabledFailureState:
		return "enabled"
	case ExecutingFailureState:
		return "executing"
	case RevertingFailureState:
		return "reverting"
	case DisabledFailureState:
		return "disabled"
	case StaleFailureState:
		return "stale"
	case ErroredFailureState:
		return "errored"
	case ErroredRevertingFailureState:
		return "erroredreverting"

	default:
		return "unknown"
	}
}

// AttackMap is a type that defines a list of map of attackers.
type AttackMap map[string]attack.Opts

// FailureMetadata has information about the object.
type FailureMetadata struct {
	ID     string `json:"id,omitempty" yaml:"id,omitempty"`         // ID is the id of the Failure.
	NodeID string `json:"nodeid,omitempty" yaml:"nodeid,omitempty"` // NodeID is the id of the Node.
}

// FailureStatus has all the information of a failure to create an injection
type FailureStatus struct {
	CurrentState  FailureState `json:"currentState,omitempty" yaml:"currentState,omitempty"`   // CurrentState is the state of the failure.
	ExpectedState FailureState `json:"expectedState,omitempty" yaml:"expectedState,omitempty"` // ExpectedState is the state the failure should be.
	Creation      time.Time    `json:"creation,omitempty" yaml:"creation,omitempty"`           // Creation is when the failure injection was created.
	Executed      time.Time    `json:"executed,omitempty" yaml:"executed,omitempty"`           // Executed is when the failure injectionwas executed.
	Finished      time.Time    `json:"finished,omitempty" yaml:"finished,omitempty"`           // Finished is when the failure injection was reverted.
}

// FailureSpec is the specification that has the information to it can be created and applied.
type FailureSpec struct {
	// Timeout is
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	// Attacks used an array so the no repeated elements of map limitation can be bypassed.
	Attacks []AttackMap `json:"attacks,omitempty" yaml:"attacks,omitempty"`
	// TODO: accuracy
}

// Failure is the way a failure is defined.
type Failure struct {
	api.TypeMeta `json:",inline" yaml:",inline"`

	// Metadta is additional data of a failure object.
	Metadata FailureMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// Spec has all the required data to create a Failure and use it.
	Spec FailureSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	// Status is the current information and status of the Failure.
	Status FailureStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// NewFailure is a plain Failure object contructor.
func NewFailure() Failure {
	return Failure{
		TypeMeta: api.TypeMeta{
			Kind:    FailureKind,
			Version: FailureVersion,
		},
	}
}
