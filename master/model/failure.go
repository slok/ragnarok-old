package model

import (
	"github.com/slok/ragnarok/types"
)

// Failure is a FeailureDefinition assigned to a node.
type Failure struct {
	ID         string             // ID is the id of the Failure.
	NodeID     string             // NodeID is the id of the Node.
	Definition string             // FailureDefinition is the failure definition.
	State      types.FailureState // State is the state of the failure.
}
