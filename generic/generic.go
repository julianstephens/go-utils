// Package generic provides commonly used generic utilities for Go projects,
// leveraging Go's type parameters to create reusable and type-safe code abstractions.
// This package includes functional programming utilities, slice operations, map utilities,
// and general helper functions that work with any type.
package generic

import "reflect"

// Zero returns the zero value for type T.
func Zero[T any]() T {
	var zero T
	return zero
}

// If mimics the ternary operator such that: cond ? vtrue : vfalse.
// It provides a concise way to conditionally select between two values of the same type.
func If[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Default returns defaultVal if val is the zero value for its type, otherwise returns val.
// This is useful for providing fallback values when dealing with potentially empty values.
func Default[T any](val T, defaultVal T) T {
	var zero T
	if reflect.DeepEqual(val, zero) {
		return defaultVal
	}
	return val
}

// Ptr returns a pointer to the given value.
// This is useful for creating pointers to literal values or variables.
func Ptr[T any](v T) *T {
	return &v
}

// Deref dereferences a pointer and returns the value it points to.
// If the pointer is nil, it returns the zero value for type T.
func Deref[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
