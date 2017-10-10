package v1

import (
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

// Parser will transform failures from and to different formats.
type Parser interface {
	// FailureToPB transforms a failure.Failure to Protobufs format.
	FailureToPB(fl *Failure) (*pbfs.Failure, error)
	// PBToFailure transforms a protobuf Failure to failure.Failure.
	PBToFailure(fl *pbfs.Failure) (*Failure, error)
}

// Transformer is the Failure transformer.
var Transformer = &transformer{
	stateParser: types.FailureStateTransformer,
}

// transformer implements the logic of a failure parser.
type transformer struct {
	stateParser types.FailureStateParser
}

// FailureToPB implements Parser interface.
func (t *transformer) FailureToPB(fl *Failure) (*pbfs.Failure, error) {
	bs, err := fl.Spec.Render()
	if err != nil {
		return nil, err
	}

	cs, err := t.stateParser.FailureStateToPB(fl.Status.CurrentState)
	if err != nil {
		return nil, err
	}

	es, err := t.stateParser.FailureStateToPB(fl.Status.ExpectedState)
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

// PBToFailure implements Parser interface.
func (t *transformer) PBToFailure(fl *pbfs.Failure) (*Failure, error) {

	spec, err := ReadFailureSpec([]byte(fl.Definition))
	if err != nil {
		return nil, err
	}

	cs, err := t.stateParser.PBToFailureState(fl.GetCurrentState())
	if err != nil {
		return nil, err
	}
	es, err := t.stateParser.PBToFailureState(fl.GetExpectedState())
	if err != nil {
		return nil, err
	}

	res := Failure{
		Metadata: Metadata{
			ID:     fl.GetId(),
			NodeID: fl.GetNodeID(),
		},
		Spec: spec,
		Status: Status{
			CurrentState:  cs,
			ExpectedState: es,
		},
	}

	return &res, nil
}
