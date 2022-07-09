package utils

func SliceContains[T comparable](s []T, x T) bool {
	for _, el := range s {
		if el == x {
			return true
		}
	}
	return false
}
