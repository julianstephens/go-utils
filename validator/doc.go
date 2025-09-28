// Package validator provides small, reusable input validators used across
// the go-utils repository.
//
// It includes common checks such as non-empty validation and a basic
// email format validator. Extracting these helpers to their own package
// makes them easy to import from multiple packages (for example, CLI
// prompts, HTTP handlers, and configuration loaders) without creating
// circular dependencies.
//
// Example
//
//	import "github.com/julianstephens/go-utils/validator"
//
//	if err := validator.ValidateEmail("user@example.com"); err != nil {
//	    // handle invalid email
//	}
package validator
