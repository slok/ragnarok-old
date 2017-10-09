package model

import (
	"github.com/slok/ragnarok/types"
)

// Labels is a key value pair map.
type Labels map[string]string

// Node is an internal and simplified representation of a failure node on the masters
// TODO: Rethink the reuse of node
type Node struct {
	ID     string          // ID is the id of the node
	Labels Labels          // Labels are the tags related with the node
	State  types.NodeState // State is the state of the Node
}
