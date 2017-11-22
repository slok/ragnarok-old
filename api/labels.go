package api

// Common reserved labels that are standarized and used by the system, they are reserved
// and can only be used by the system.
const (
	LabelNode       = "node"
	LabelFailure    = "failure"
	LabelExperiment = "experiment"
	LabelID         = "id"
)

// GetAllReserverLabels will return all the reserved labels on the system.
func GetAllReserverLabels() map[string]struct{} {
	return map[string]struct{}{
		LabelNode:       struct{}{},
		LabelFailure:    struct{}{},
		LabelExperiment: struct{}{},
		LabelID:         struct{}{},
	}
}
