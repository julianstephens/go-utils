package generic

// Map applies a function to each element of a slice and returns a new slice with the results.
// The function f is applied to each element of type T and produces an element of type U.
func Map[T, U any](slice []T, f func(T) U) []U {
	if slice == nil {
		return nil
	}
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate function.
// The predicate function should return true for elements to be included in the result.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	if slice == nil {
		return nil
	}
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce applies a reduction function to the elements of a slice from left to right
// to reduce the slice to a single value of type U.
// The initial value serves as the starting accumulator value.
func Reduce[T, U any](slice []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range slice {
		acc = f(acc, v)
	}
	return acc
}

// Find returns the first element in the slice that satisfies the predicate function.
// It returns the element and true if found, or the zero value and false if not found.
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Any returns true if any element in the slice satisfies the predicate function.
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements in the slice satisfy the predicate function.
// Returns true for empty slices.
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// ForEach applies a function to each element of the slice.
// This is useful for side effects like logging or modifying external state.
func ForEach[T any](slice []T, f func(T)) {
	for _, v := range slice {
		f(v)
	}
}