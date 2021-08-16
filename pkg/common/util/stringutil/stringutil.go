package stringutil

// StrInSlice returns whether the string is in the slice.
func StrInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}
