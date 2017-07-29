package failure

import (
	"time"

	"github.com/slok/ragnarok/types"
)

// Failure has all the information of a failure to create an injection
type Failure struct {
	ID            string             // ID is the id of the Failure.
	NodeID        string             // NodeID is the id of the Node.
	Definition    Definition         // FailureDefinition is the failure definition.
	CurrentState  types.FailureState // CurrentState is the state of the failure.
	ExpectedState types.FailureState // ExpectedState is the state the failure should be.
	Creation      time.Time          // Creation is when the failure injection was created.
	Executed      time.Time          // Executed is when the failure injectionwas executed.
	Finished      time.Time          //Finished is when the failure injection was reverted.
}
