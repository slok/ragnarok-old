package types

import (
	"fmt"
	"strings"

	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
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
	// DisabledFailureState is when the failure is not making stuff (reverted already).
	DisabledFailureState
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
	case ErroredFailureState:
		return "errored"
	case ErroredRevertingFailureState:
		return "erroredreverting"

	default:
		return "unknown"
	}
}

// FailureStateParser has the required methods to transform failure state.
type FailureStateParser interface {
	// StrToFailureState transforms an string to a failure state.
	StrToFailureState(state string) (FailureState, error)
	// PBToFailureState transforms a GRPC failure staet to a failure state.
	PBToFailureState(state pbfs.State) (FailureState, error)
	// FailureState transforms a failure state to a GRPC failure state
	FailureStateToPB(state FailureState) (pbfs.State, error)
}

// failureStateParser will convert failure state in different types.
type failureStateParser struct{}

// FailureStateTransformer is the util to transform failure status kinds.
var FailureStateTransformer = &failureStateParser{}

// StrToFailureState implements FailureStateParser interface.
func (f *failureStateParser) StrToFailureState(state string) (FailureState, error) {
	switch strings.ToLower(state) {
	case "enabled":
		return EnabledFailureState, nil
	case "executing":
		return ExecutingFailureState, nil
	case "reverting":
		return RevertingFailureState, nil
	case "disabled":
		return DisabledFailureState, nil
	case "errored":
		return ErroredFailureState, nil
	case "erroredreverting":
		return ErroredRevertingFailureState, nil
	default:
		return UnknownFailureState, fmt.Errorf("invalid failure state: %s", state)
	}
}

// PBToFailureState implements FailureStateParser interface.
func (f *failureStateParser) PBToFailureState(state pbfs.State) (FailureState, error) {
	switch state {
	case pbfs.State_ENABLED:
		return EnabledFailureState, nil
	case pbfs.State_EXECUTING:
		return ExecutingFailureState, nil
	case pbfs.State_REVERTING:
		return RevertingFailureState, nil
	case pbfs.State_DISABLED:
		return DisabledFailureState, nil
	case pbfs.State_ERRORED:
		return ErroredFailureState, nil
	case pbfs.State_ERRORED_REVERTING:
		return ErroredRevertingFailureState, nil
	default:
		return UnknownFailureState, fmt.Errorf("invalid failure state: %s", state)
	}
}

// FailureStateToPB implements FailureStateParser interface.
func (f *failureStateParser) FailureStateToPB(state FailureState) (pbfs.State, error) {
	switch state {
	case EnabledFailureState:
		return pbfs.State_ENABLED, nil
	case ExecutingFailureState:
		return pbfs.State_EXECUTING, nil
	case RevertingFailureState:
		return pbfs.State_REVERTING, nil
	case DisabledFailureState:
		return pbfs.State_DISABLED, nil
	case ErroredFailureState:
		return pbfs.State_ERRORED, nil
	case ErroredRevertingFailureState:
		return pbfs.State_ERRORED_REVERTING, nil
	default:
		return pbfs.State_UNKNOWN, fmt.Errorf("invalid failure state: %s", state)
	}
}
