package failure

import (
	"time"

	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

// Failure has all the information of a failure to create an injection
type Failure struct {
	ID            string             // ID is the id of the Failure.
	NodeID        string             // NodeID is the id of the Node.
	Definition    Definition         // Definition is the failure definition.
	CurrentState  types.FailureState // CurrentState is the state of the failure.
	ExpectedState types.FailureState // ExpectedState is the state the failure should be.
	Creation      time.Time          // Creation is when the failure injection was created.
	Executed      time.Time          // Executed is when the failure injectionwas executed.
	Finished      time.Time          //Finished is when the failure injection was reverted.
}

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
	bs, err := fl.Definition.Render()
	if err != nil {
		return nil, err
	}

	cs, err := t.stateParser.FailureStateToPB(fl.CurrentState)
	if err != nil {
		return nil, err
	}

	es, err := t.stateParser.FailureStateToPB(fl.ExpectedState)
	if err != nil {
		return nil, err
	}

	return &pbfs.Failure{
		Id:            fl.ID,
		NodeID:        fl.NodeID,
		Definition:    string(bs),
		CurrentState:  cs,
		ExpectedState: es,
	}, nil
}

// PBToFailure implements Parser interface.
func (t *transformer) PBToFailure(fl *pbfs.Failure) (*Failure, error) {

	def, err := ReadDefinition([]byte(fl.Definition))
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

	return &Failure{
		ID:            fl.GetId(),
		NodeID:        fl.GetNodeID(),
		Definition:    def,
		CurrentState:  cs,
		ExpectedState: es,
	}, nil
}
