package types

import (
	"fmt"
	"strings"

	"github.com/slok/ragnarok/api/chaos/v1"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
)

// FailureStateParser has the required methods to transform failure state.
type FailureStateParser interface {
	// StrToFailureState transforms an string to a failure state.
	StrToFailureState(state string) (v1.FailureState, error)
	// PBToFailureState transforms a GRPC failure staet to a failure state.
	PBToFailureState(state pbfs.State) (v1.FailureState, error)
	// FailureState transforms a failure state to a GRPC failure state
	FailureStateToPB(state v1.FailureState) (pbfs.State, error)
}

// failureStateParser will convert failure state in different types.
type failureStateParser struct{}

// FailureStateTransformer is the util to transform failure status kinds.
var FailureStateTransformer = &failureStateParser{}

// StrToFailureState implements FailureStateParser interface.
func (f *failureStateParser) StrToFailureState(state string) (v1.FailureState, error) {
	switch strings.ToLower(state) {
	case "enabled":
		return v1.EnabledFailureState, nil
	case "executing":
		return v1.ExecutingFailureState, nil
	case "reverting":
		return v1.RevertingFailureState, nil
	case "disabled":
		return v1.DisabledFailureState, nil
	case "stale":
		return v1.StaleFailureState, nil
	case "errored":
		return v1.ErroredFailureState, nil
	case "erroredreverting":
		return v1.ErroredRevertingFailureState, nil
	default:
		return v1.UnknownFailureState, fmt.Errorf("invalid failure state: %s", state)
	}
}

// PBToFailureState implements FailureStateParser interface.
func (f *failureStateParser) PBToFailureState(state pbfs.State) (v1.FailureState, error) {
	switch state {
	case pbfs.State_ENABLED:
		return v1.EnabledFailureState, nil
	case pbfs.State_EXECUTING:
		return v1.ExecutingFailureState, nil
	case pbfs.State_REVERTING:
		return v1.RevertingFailureState, nil
	case pbfs.State_DISABLED:
		return v1.DisabledFailureState, nil
	case pbfs.State_STALE:
		return v1.StaleFailureState, nil
	case pbfs.State_ERRORED:
		return v1.ErroredFailureState, nil
	case pbfs.State_ERRORED_REVERTING:
		return v1.ErroredRevertingFailureState, nil
	default:
		return v1.UnknownFailureState, fmt.Errorf("invalid failure state: %s", state)
	}
}

// FailureStateToPB implements FailureStateParser interface.
func (f *failureStateParser) FailureStateToPB(state v1.FailureState) (pbfs.State, error) {
	switch state {
	case v1.EnabledFailureState:
		return pbfs.State_ENABLED, nil
	case v1.ExecutingFailureState:
		return pbfs.State_EXECUTING, nil
	case v1.RevertingFailureState:
		return pbfs.State_REVERTING, nil
	case v1.DisabledFailureState:
		return pbfs.State_DISABLED, nil
	case v1.StaleFailureState:
		return pbfs.State_STALE, nil
	case v1.ErroredFailureState:
		return pbfs.State_ERRORED, nil
	case v1.ErroredRevertingFailureState:
		return pbfs.State_ERRORED_REVERTING, nil
	default:
		return pbfs.State_UNKNOWN, fmt.Errorf("invalid failure state: %s", state)
	}
}
