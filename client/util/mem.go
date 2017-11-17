package util

// SelectorMatchesLabels will return true if the selector matches the labels.
func SelectorMatchesLabels(labels map[string]string, selector map[string]string) bool {
	for lk, lv := range selector {
		// If one label does not match then we are done.
		if nv, ok := labels[lk]; !ok || (ok && nv != lv) {
			return false
		}
	}
	// All labels matched.
	return true
}
