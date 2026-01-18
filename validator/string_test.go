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

func TestStringValidator_Pattern(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidatePattern - should pass
	if err := v.ValidatePattern("hello123", `^[a-z0-9]+$`); err != nil {
		t.Errorf("ValidatePattern('hello123', '^[a-z0-9]+$') should pass, got error: %v", err)
	}

	// Test ValidatePattern - should fail
	if err := v.ValidatePattern("hello_123", `^[a-z0-9]+$`); err == nil {
		t.Error("ValidatePattern('hello_123', '^[a-z0-9]+$') should fail")
	}

	// Test ValidatePattern - invalid regex
	if err := v.ValidatePattern("test", "[invalid"); err == nil {
		t.Error("ValidatePattern with invalid regex should fail")
	}
}

func TestStringValidator_Alphanumeric(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateAlphanumeric - should pass
	if err := v.ValidateAlphanumeric("abc123"); err != nil {
		t.Errorf("ValidateAlphanumeric('abc123') should pass, got error: %v", err)
	}

	// Test ValidateAlphanumeric - with special char
	if err := v.ValidateAlphanumeric("abc_123"); err == nil {
		t.Error("ValidateAlphanumeric('abc_123') should fail")
	}

	// Test ValidateAlphanumeric - with space
	if err := v.ValidateAlphanumeric("abc 123"); err == nil {
		t.Error("ValidateAlphanumeric('abc 123') should fail")
	}
}

func TestStringValidator_Alpha(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateAlpha - should pass
	if err := v.ValidateAlpha("abcXYZ"); err != nil {
		t.Errorf("ValidateAlpha('abcXYZ') should pass, got error: %v", err)
	}

	// Test ValidateAlpha - with digit
	if err := v.ValidateAlpha("abc123"); err == nil {
		t.Error("ValidateAlpha('abc123') should fail")
	}

	// Test ValidateAlpha - with special char
	if err := v.ValidateAlpha("abc_"); err == nil {
		t.Error("ValidateAlpha('abc_') should fail")
	}
}

func TestStringValidator_Numeric(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateNumeric - should pass
	if err := v.ValidateNumeric("12345"); err != nil {
		t.Errorf("ValidateNumeric('12345') should pass, got error: %v", err)
	}

	// Test ValidateNumeric - with letter
	if err := v.ValidateNumeric("12345a"); err == nil {
		t.Error("ValidateNumeric('12345a') should fail")
	}

	// Test ValidateNumeric - with special char
	if err := v.ValidateNumeric("123-45"); err == nil {
		t.Error("ValidateNumeric('123-45') should fail")
	}
}

func TestStringValidator_Slug(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateSlug - should pass
	if err := v.ValidateSlug("my-slug_123"); err != nil {
		t.Errorf("ValidateSlug('my-slug_123') should pass, got error: %v", err)
	}

	// Test ValidateSlug - with uppercase
	if err := v.ValidateSlug("My-Slug"); err == nil {
		t.Error("ValidateSlug('My-Slug') should fail")
	}

	// Test ValidateSlug - with space
	if err := v.ValidateSlug("my slug"); err == nil {
		t.Error("ValidateSlug('my slug') should fail")
	}
}

func TestStringValidator_Lowercase(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateLowercase - should pass
	if err := v.ValidateLowercase("abc123_"); err != nil {
		t.Errorf("ValidateLowercase('abc123_') should pass, got error: %v", err)
	}

	// Test ValidateLowercase - with uppercase
	if err := v.ValidateLowercase("Abc"); err == nil {
		t.Error("ValidateLowercase('Abc') should fail")
	}
}

func TestStringValidator_Uppercase(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateUppercase - should pass
	if err := v.ValidateUppercase("ABC123_"); err != nil {
		t.Errorf("ValidateUppercase('ABC123_') should pass, got error: %v", err)
	}

	// Test ValidateUppercase - with lowercase
	if err := v.ValidateUppercase("ABc"); err == nil {
		t.Error("ValidateUppercase('ABc') should fail")
	}
}

func TestStringValidator_Contains(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateContains - should pass
	if err := v.ValidateContains("hello world", "world"); err != nil {
		t.Errorf("ValidateContains('hello world', 'world') should pass, got error: %v", err)
	}

	// Test ValidateContains - should fail
	if err := v.ValidateContains("hello world", "foo"); err == nil {
		t.Error("ValidateContains('hello world', 'foo') should fail")
	}
}

func TestStringValidator_NotContains(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateNotContains - should pass
	if err := v.ValidateNotContains("hello world", "foo"); err != nil {
		t.Errorf("ValidateNotContains('hello world', 'foo') should pass, got error: %v", err)
	}

	// Test ValidateNotContains - should fail
	if err := v.ValidateNotContains("hello world", "world"); err == nil {
		t.Error("ValidateNotContains('hello world', 'world') should fail")
	}
}

func TestStringValidator_Prefix(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidatePrefix - should pass
	if err := v.ValidatePrefix("hello world", "hello"); err != nil {
		t.Errorf("ValidatePrefix('hello world', 'hello') should pass, got error: %v", err)
	}

	// Test ValidatePrefix - should fail
	if err := v.ValidatePrefix("hello world", "world"); err == nil {
		t.Error("ValidatePrefix('hello world', 'world') should fail")
	}
}

func TestStringValidator_Suffix(t *testing.T) {
	v := validator.Strings[string]()

	// Test ValidateSuffix - should pass
	if err := v.ValidateSuffix("hello world", "world"); err != nil {
		t.Errorf("ValidateSuffix('hello world', 'world') should pass, got error: %v", err)
	}

	// Test ValidateSuffix - should fail
	if err := v.ValidateSuffix("hello world", "hello"); err == nil {
		t.Error("ValidateSuffix('hello world', 'hello') should fail")
	}
}
