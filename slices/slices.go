// Package slices provides generic slice utility functions commonly duplicated across Go projects.
//
// This package consolidates slice operations that are frequently needed, offering
// type-safe generic implementations that work with any comparable or arbitrary types
// as appropriate.
//
// The package includes utilities for:
//   - Conditional selection (ternary-like operations)
//   - Set operations (difference, subset checking)
//   - Element manipulation (deletion by index)
//
// Example usage:
//
//	import "github.com/julianstephens/go-utils/slices"
//
//	// Ternary-like conditional selection
//	result := slices.If(condition, "yes", "no")
//
//	// Find elements in slice a but not in slice b
//	diff := slices.Difference([]string{"a", "b", "c"}, []string{"b"})
//	// Returns: []string{"a", "c"}
//
//	// Remove element at index
//	items := slices.DeleteElement([]int{1, 2, 3, 4}, 1)
//	// Returns: []int{1, 3, 4}
//
//	// Check if all elements of subset are in main slice
//	hasAll := slices.ContainsAll([]int{1, 2, 3, 4}, []int{2, 3})
//	// Returns: true
package slices

// If mimics the ternary operator such that: cond ? vtrue : vfalse.
//
// This function provides a concise way to select between two values based on a boolean condition.
// It accepts any type T and returns the first value if the condition is true, otherwise the second value.
//
// Example:
//
//	// Basic usage with strings
//	message := slices.If(isLoggedIn, "Welcome back!", "Please log in")
//
//	// With numbers
//	max := slices.If(a > b, a, b)
//
//	// With complex types
//	config := slices.If(isProd, prodConfig, devConfig)
func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Difference implements slice subtraction, returning elements in slice a that are not in slice b.
//
// This function returns a new slice containing all elements from a that do not appear in b.
// The order of elements from the original slice a is preserved in the result.
// Empty slices are handled gracefully.
//
// Time complexity: O(len(a) + len(b))
// Space complexity: O(len(b)) for the lookup map + O(result) for output
//
// Example:
//
//	a := []string{"apple", "banana", "cherry", "date"}
//	b := []string{"banana", "date"}
//	result := slices.Difference(a, b)
//	// Returns: []string{"apple", "cherry"}
//
//	// Edge cases
//	slices.Difference([]string{}, []string{"a"})        // Returns: []string{}
//	slices.Difference([]string{"a"}, []string{})        // Returns: []string{"a"}
//	slices.Difference([]string{"a"}, []string{"a"})     // Returns: []string{}
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
//
// This function creates a new slice with the element at the given index removed.
// The function does not perform bounds checking - providing an invalid index will panic.
// For safe deletion with bounds checking, validate the index before calling this function.
//
// Time complexity: O(n) where n is len(slice)
// Space complexity: O(n) for the new slice
//
// Example:
//
//	numbers := []int{10, 20, 30, 40, 50}
//	result := slices.DeleteElement(numbers, 2)
//	// Returns: []int{10, 20, 40, 50}
//
//	// Works with any type
//	words := []string{"hello", "world", "foo", "bar"}
//	result := slices.DeleteElement(words, 0)
//	// Returns: []string{"world", "foo", "bar"}
//
// Panics if index < 0 or index >= len(slice).
func DeleteElement[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

// ContainsAll returns true if all elements in subset are present in mainSlice.
//
// This function checks whether mainSlice contains every element that appears in subset.
// An empty subset always returns true (vacuous truth).
// The function uses a map for efficient lookup, providing good performance even with large slices.
//
// Time complexity: O(len(mainSlice) + len(subset))
// Space complexity: O(len(mainSlice)) for the lookup map
//
// Example:
//
//	main := []int{1, 2, 3, 4, 5}
//	subset := []int{2, 4, 5}
//	result := slices.ContainsAll(main, subset)
//	// Returns: true
//
//	subset2 := []int{2, 6}
//	result2 := slices.ContainsAll(main, subset2)
//	// Returns: false (6 is not in main)
//
//	// Edge cases
//	slices.ContainsAll([]int{1, 2}, []int{})           // Returns: true (empty subset)
//	slices.ContainsAll([]int{}, []int{1})              // Returns: false
//	slices.ContainsAll([]int{1, 2}, []int{1, 2})       // Returns: true (identical)
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
