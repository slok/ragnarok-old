package client

// Status interface will implement the required methods to be able to communicate
// with a node status server
type Status interface {
	// RegisterNode registers a node as available on the server
	RegisterNode(id string, tags map[string]string) error
}
