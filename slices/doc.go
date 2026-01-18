// Package slices provides basic generic slice utility functions.
//
// Deprecated: Use the generic package instead, which provides a comprehensive
// set of slice operations with better organization. The generic package includes
// all functionality from slices plus additional utilities like functional programming,
// map operations, and more advanced slice manipulations.
//
// TODO: Remove this package in v0.6.0 release
//
// Migration Guide:
//   - slices.If → generic.If
//   - slices.Difference → generic.Difference
//   - slices.DeleteElement → generic.DeleteElement
//   - slices.ContainsAll → generic.ContainsAll
//
// Example:
//
//	import "github.com/julianstephens/go-utils/generic"
//
//	// Old way
//	// result := slices.Difference(a, b)
//
//	// New way
//	result := generic.Difference(a, b)
package slices
