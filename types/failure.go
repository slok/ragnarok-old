package types

import (
	"github.com/slok/ragnarok/api/chaos/v1"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
)

// FailureParser will transform failures from and to different formats.
type FailureParser interface {
	// FailureToPB transforms a failure.Failure to Protobufs format.
	FailureToPB(fl *v1.Failure) (*pbfs.Failure, error)
	// PBToFailure transforms a protobuf Failure to failure.Failure.
	PBToFailure(fl *pbfs.Failure) (*v1.Failure, error)
}

// FailureTransformer is the Failure transformer.
var FailureTransformer = &failureTransformer{
	stateParser: FailureStateTransformer,
}

// failureTransformer implements the logic of a failure parser.
type failureTransformer struct {
	stateParser FailureStateParser
}

// FailureToPB implements FailureParser interface.
func (f *failureTransformer) FailureToPB(fl *v1.Failure) (*pbfs.Failure, error) {
	bs, err := fl.Spec.Render()
	if err != nil {
		return nil, err
	}

	cs, err := f.stateParser.FailureStateToPB(fl.Status.CurrentState)
	if err != nil {
		return nil, err
	}

	es, err := f.stateParser.FailureStateToPB(fl.Status.ExpectedState)
	if err != nil {
		return nil, err
	}

	return &pbfs.Failure{
		Id:            fl.Metadata.ID,
		NodeID:        fl.Metadata.NodeID,
		Definition:    string(bs),
		CurrentState:  cs,
		ExpectedState: es,
	}, nil
}

// PBToFailure implements FailureParser interface.
func (f *failureTransformer) PBToFailure(fl *pbfs.Failure) (*v1.Failure, error) {

	spec, err := v1.ReadFailureSpec([]byte(fl.Definition))
	if err != nil {
		return nil, err
	}

	cs, err := f.stateParser.PBToFailureState(fl.GetCurrentState())
	if err != nil {
		return nil, err
	}
	es, err := f.stateParser.PBToFailureState(fl.GetExpectedState())
	if err != nil {
		return nil, err
	}

	res := v1.Failure{
		Metadata: v1.FailureMetadata{
			ID:     fl.GetId(),
			NodeID: fl.GetNodeID(),
		},
		Spec: spec,
		Status: v1.FailureStatus{
			CurrentState:  cs,
			ExpectedState: es,
		},
	}

	return &res, nil
}
