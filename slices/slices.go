// Package slices provides generic slice utility functions commonly duplicated across Go projects.
package slices

// If mimics the ternary operator such that: cond ? vtrue : vfalse.
func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Difference implements slice subtraction, returning elements in slice a that are not in slice b.
func Difference(a []string, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// DeleteElement removes an element from a slice at the specified index and returns the modified slice.
func DeleteElement[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

// ContainsAll returns true if all elements in subset are present in mainSlice.
func ContainsAll[T comparable](mainSlice, subset []T) bool {
	if len(subset) == 0 {
		return true
	}

	mainMap := make(map[T]bool)
	for _, item := range mainSlice {
		mainMap[item] = true
	}

	for _, item := range subset {
		if !mainMap[item] {
			return false
		}
	}
	return true
}
