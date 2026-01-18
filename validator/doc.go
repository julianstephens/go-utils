// Package validator provides comprehensive, reusable input validators
// for Go applications.
//
// The package offers both simple validation functions and advanced generic
// validators for type-safe validation of numbers, strings, and parsed values.
// Validators are designed to be composable and reusable across different
// parts of an application.
//
// # Basic Usage
//
// Simple validation functions for common checks:
//
//	import "github.com/julianstephens/go-utils/validator"
//
//	if err := validator.ValidateNonEmpty("input"); err != nil {
//	    // handle empty input
//	}
//
// # Advanced Generic Validators
//
// For comprehensive validation, use the generic validator factories:
//
//	// Number validation
//	numValidator := validator.Numbers[int]()
//	if err := numValidator.ValidateRange(42, 1, 100); err != nil {
//	    // handle invalid number
//	}
//
//	// String validation
//	strValidator := validator.Strings[string]()
//	if err := strValidator.ValidateMinLength("hello", 3); err != nil {
//	    // handle string too short
//	}
//
//	// Parsing validation (accessed through string validator)
//	strValidator := validator.Strings[string]()
//	if err := strValidator.Parse.ValidateEmail("user@example.com"); err != nil {
//	    // handle invalid email
//	}
//
// # Type Constraints
//
// The package uses Go generics with type constraints:
//   - Number: ~int, ~uint, ~float32, ~float64 variants
//   - StringLike: ~string, ~[]byte, ~[]rune
//   - Emptyable: types that can be checked for emptiness
//
// # Error Handling
//
// All validators return detailed ValidationError instances with context
// about what was expected vs. what was received, making debugging easier.
package validator
