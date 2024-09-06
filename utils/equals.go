package goutils

// Equals compares two values of any type
func Equals[T comparable](a, b T) bool {
	return a == b
}

// Equals for integers
func EqualsInt(a, b int) bool {
	return a == b
}

// Equals for strings
func EqualsString(a, b string) bool {
	return a == b
}

// Equals for booleans
func EqualsBool(a, b bool) bool {
	return a == b
}

// EqualsSlice compares two slices of any type
func EqualsSlice[T any](a, b []T, equals func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equals(a[i], b[i]) {
			return false
		}
	}
	return true
}

// Equals compares two maps
func EqualsMap(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
