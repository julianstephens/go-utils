package validator_test

import (
	"errors"
	"testing"

	"github.com/julianstephens/go-utils/validator"
)

func TestStringValidator(t *testing.T) {
	v := validator.Strings[string]()

	// Test that StringValidator exists and has Errorf method
	if v == nil {
		t.Fatal("Strings[string]() returned nil")
	}

	// The StringValidator currently only has Errorf method
	// Add more tests when more methods are implemented
}

func TestStringValidator_Errorf(t *testing.T) {
	validator := validator.Strings[string]()

	// Test Errorf method
	err := validator.Errorf("test cause", "expected", "actual", nil)
	if err == nil {
		t.Fatal("Errorf returned nil")
	}

	// Verify it's a ValidationError
	var valErr interface{}
	if !errors.As(err, &valErr) {
		t.Fatal("Errorf should return a ValidationError")
	}

	// Check the error message format
	if err.Error() == "" {
		t.Error("ValidationError should have a non-empty error message")
	}
}

func TestStringValidator_ValidateMinLength(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateMinLength - should pass
	if err := v.ValidateMinLength("hello", 3); err != nil {
		t.Errorf("ValidateMinLength('hello', 3) should pass, got error: %v", err)
	}

	// Test ValidateMinLength - should fail
	if err := v.ValidateMinLength("hi", 3); err == nil {
		t.Error("ValidateMinLength('hi', 3) should fail")
	}

	// Test with empty string
	if err := v.ValidateMinLength("", 1); err == nil {
		t.Error("ValidateMinLength('', 1) should fail")
	}

	// Test with exact length
	if err := v.ValidateMinLength("abc", 3); err != nil {
		t.Errorf("ValidateMinLength('abc', 3) should pass, got error: %v", err)
	}
}

func TestStringValidator_ValidateMaxLength(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateMaxLength - should pass
	if err := v.ValidateMaxLength("hi", 5); err != nil {
		t.Errorf("ValidateMaxLength('hi', 5) should pass, got error: %v", err)
	}

	// Test ValidateMaxLength - should fail
	if err := v.ValidateMaxLength("hello world", 5); err == nil {
		t.Error("ValidateMaxLength('hello world', 5) should fail")
	}

	// Test with empty string
	if err := v.ValidateMaxLength("", 5); err != nil {
		t.Errorf("ValidateMaxLength('', 5) should pass, got error: %v", err)
	}

	// Test with exact length
	if err := v.ValidateMaxLength("abc", 3); err != nil {
		t.Errorf("ValidateMaxLength('abc', 3) should pass, got error: %v", err)
	}
}

func TestStringValidator_ValidateLengthRange(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateLengthRange - should pass
	if err := v.ValidateLengthRange("hello", 2, 10); err != nil {
		t.Errorf("ValidateLengthRange('hello', 2, 10) should pass, got error: %v", err)
	}

	// Test ValidateLengthRange - below minimum
	if err := v.ValidateLengthRange("hi", 5, 10); err == nil {
		t.Error("ValidateLengthRange('hi', 5, 10) should fail (below minimum)")
	}

	// Test ValidateLengthRange - above maximum
	if err := v.ValidateLengthRange("hello world this is long", 5, 10); err == nil {
		t.Error("ValidateLengthRange('hello world this is long', 5, 10) should fail (above maximum)")
	}

	// Test with empty string
	if err := v.ValidateLengthRange("", 1, 5); err == nil {
		t.Error("ValidateLengthRange('', 1, 5) should fail")
	}

	// Test with exact min length
	if err := v.ValidateLengthRange("abc", 3, 5); err != nil {
		t.Errorf("ValidateLengthRange('abc', 3, 5) should pass, got error: %v", err)
	}

	// Test with exact max length
	if err := v.ValidateLengthRange("abcde", 3, 5); err != nil {
		t.Errorf("ValidateLengthRange('abcde', 3, 5) should pass, got error: %v", err)
	}
}

func TestStringValidator_Bytes(t *testing.T) {
	v := validator.Strings[[]byte]()

	// Test ValidateMinLength with bytes
	if err := v.ValidateMinLength([]byte("hello"), 3); err != nil {
		t.Errorf("ValidateMinLength([]byte('hello'), 3) should pass, got error: %v", err)
	}

	if err := v.ValidateMinLength([]byte("hi"), 3); err == nil {
		t.Error("ValidateMinLength([]byte('hi'), 3) should fail")
	}

	// Test ValidateMaxLength with bytes
	if err := v.ValidateMaxLength([]byte("hi"), 5); err != nil {
		t.Errorf("ValidateMaxLength([]byte('hi'), 5) should pass, got error: %v", err)
	}

	if err := v.ValidateMaxLength([]byte("hello world"), 5); err == nil {
		t.Error("ValidateMaxLength([]byte('hello world'), 5) should fail")
	}

	// Test ValidateLengthRange with bytes
	if err := v.ValidateLengthRange([]byte("hello"), 2, 10); err != nil {
		t.Errorf("ValidateLengthRange([]byte('hello'), 2, 10) should pass, got error: %v", err)
	}

	if err := v.ValidateLengthRange([]byte("hi"), 5, 10); err == nil {
		t.Error("ValidateLengthRange([]byte('hi'), 5, 10) should fail")
	}
}
