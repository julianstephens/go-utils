package validator

import (
	"fmt"
)

// OneOf validates that a value is one of the allowed values
func OneOf[T comparable](input T, allowed ...T) error {
	for _, v := range allowed {
		if input == v {
			return nil
		}
	}
	return fmt.Errorf("%w: value %v not in allowed set %v", ErrNotInSet, input, allowed)
}

// All validates that all validator functions pass (AND logic)
func All(validators ...func() error) error {
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

// Any validates that at least one validator function passes (OR logic)
func Any(validators ...func() error) error {
	if len(validators) == 0 {
		return fmt.Errorf("at least one validator must be provided")
	}
	var lastErr error
	for _, validator := range validators {
		if err := validator(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}
	return lastErr
}

// ValidateSliceLength validates that a slice has the specified length
func ValidateSliceLength[T any](input []T, length int) error {
	if len(input) != length {
		return fmt.Errorf("slice length mismatch: expected %d, got %d", length, len(input))
	}
	return nil
}

// ValidateSliceMinLength validates that a slice has at least the minimum length
func ValidateSliceMinLength[T any](input []T, min int) error {
	if len(input) < min {
		return fmt.Errorf("%w: slice length is %d, minimum required is %d", ErrSliceTooShort, len(input), min)
	}
	return nil
}

// ValidateSliceMaxLength validates that a slice has at most the maximum length
func ValidateSliceMaxLength[T any](input []T, max int) error {
	if len(input) > max {
		return fmt.Errorf("%w: slice length is %d, maximum allowed is %d", ErrSliceTooLong, len(input), max)
	}
	return nil
}

// ValidateSliceLengthRange validates that a slice length is within the specified range (inclusive)
func ValidateSliceLengthRange[T any](input []T, min, max int) error {
	length := len(input)
	if length < min || length > max {
		return fmt.Errorf("slice length out of range: expected [%d, %d], got %d", min, max, length)
	}
	return nil
}

// ValidateSliceContains validates that a slice contains the specified element
func ValidateSliceContains[T comparable](input []T, element T) error {
	for _, v := range input {
		if v == element {
			return nil
		}
	}
	return fmt.Errorf("slice does not contain element %v", element)
}

// ValidateSliceNotContains validates that a slice does not contain the specified element
func ValidateSliceNotContains[T comparable](input []T, element T) error {
	for _, v := range input {
		if v == element {
			return fmt.Errorf("slice contains forbidden element %v", element)
		}
	}
	return nil
}

// ValidateSliceUnique validates that all elements in a slice are unique
func ValidateSliceUnique[T comparable](input []T) error {
	seen := make(map[any]bool)
	for _, v := range input {
		if seen[v] {
			return fmt.Errorf("slice contains duplicate element %v", v)
		}
		seen[v] = true
	}
	return nil
}

// ValidateMapHasKey validates that a map contains the specified key
func ValidateMapHasKey[K comparable, V any](input map[K]V, key K) error {
	if _, ok := input[key]; !ok {
		return fmt.Errorf("map does not contain key %v", key)
	}
	return nil
}

// ValidateMapNotHasKey validates that a map does not contain the specified key
func ValidateMapNotHasKey[K comparable, V any](input map[K]V, key K) error {
	if _, ok := input[key]; ok {
		return fmt.Errorf("map contains forbidden key %v", key)
	}
	return nil
}

// ValidateMapMinLength validates that a map has at least the minimum number of entries
func ValidateMapMinLength[K comparable, V any](input map[K]V, min int) error {
	if len(input) < min {
		return fmt.Errorf("map length is %d, minimum required is %d", len(input), min)
	}
	return nil
}

// ValidateMapMaxLength validates that a map has at most the maximum number of entries
func ValidateMapMaxLength[K comparable, V any](input map[K]V, max int) error {
	if len(input) > max {
		return fmt.Errorf("map length is %d, maximum allowed is %d", len(input), max)
	}
	return nil
}
