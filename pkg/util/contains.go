package util

// Contains takes a slice and the item to be found, it returns true if successful
func Contains[T comparable](slice []T, predicate T) bool {
	for _, item := range slice {
		if item == predicate {
			return true
		}
	}

	return false
}
