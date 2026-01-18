package validator_test

import (
	"fmt"
	"testing"

	"github.com/julianstephens/go-utils/validator"
)

func TestParseValidator_Date(t *testing.T) {
	pv := validator.Parse()

	// Test ValidateDate - should pass
	if err := pv.ValidateDate("2024-01-15", "2006-01-02"); err != nil {
		t.Errorf("ValidateDate('2024-01-15', '2006-01-02') should pass, got error: %v", err)
	}

	// Test ValidateDate - should fail (wrong format)
	if err := pv.ValidateDate("01/15/2024", "2006-01-02"); err == nil {
		t.Error("ValidateDate('01/15/2024', '2006-01-02') should fail")
	}

	// Test ValidateDate - empty input
	if err := pv.ValidateDate("", "2006-01-02"); err == nil {
		t.Error("ValidateDate('', '2006-01-02') should fail")
	}
}

func TestParseValidator_Duration(t *testing.T) {
	pv := validator.Parse()

	// Test ValidateDuration - should pass with various formats
	testCases := []string{"5m", "2h", "1s", "1h30m", "500ms", "1.5s"}
	for _, tc := range testCases {
		if err := pv.ValidateDuration(tc); err != nil {
			t.Errorf("ValidateDuration('%s') should pass, got error: %v", tc, err)
		}
	}

	// Test ValidateDuration - should fail
	if err := pv.ValidateDuration("invalid"); err == nil {
		t.Error("ValidateDuration('invalid') should fail")
	}

	// Test ValidateDuration - empty input
	if err := pv.ValidateDuration(""); err == nil {
		t.Error("ValidateDuration('') should fail")
	}
}

func TestParseValidator_PhoneNumber(t *testing.T) {
	pv := validator.Parse()

	// Test ValidatePhoneNumber - should pass with various formats
	testCases := []string{
		"1234567890",
		"(123) 456-7890",
		"+1-123-456-7890",
		"123-456-7890",
		"+1 123 456 7890",
	}
	for _, tc := range testCases {
		if err := pv.ValidatePhoneNumber(tc); err != nil {
			t.Errorf("ValidatePhoneNumber('%s') should pass, got error: %v", tc, err)
		}
	}

	// Test ValidatePhoneNumber - should fail (with invalid chars)
	if err := pv.ValidatePhoneNumber("123-abc-7890"); err == nil {
		t.Error("ValidatePhoneNumber('123-abc-7890') should fail")
	}

	// Test ValidatePhoneNumber - empty input
	if err := pv.ValidatePhoneNumber(""); err == nil {
		t.Error("ValidatePhoneNumber('') should fail")
	}
}

func TestOneOf(t *testing.T) {
	// Test OneOf - should pass
	if err := validator.OneOf(2, 1, 2, 3); err != nil {
		t.Errorf("OneOf(2, 1, 2, 3) should pass, got error: %v", err)
	}

	// Test OneOf - should fail
	if err := validator.OneOf(5, 1, 2, 3); err == nil {
		t.Error("OneOf(5, 1, 2, 3) should fail")
	}

	// Test OneOf with strings
	if err := validator.OneOf("active", "active", "inactive"); err != nil {
		t.Errorf("OneOf('active', 'active', 'inactive') should pass, got error: %v", err)
	}

	if err := validator.OneOf("pending", "active", "inactive"); err == nil {
		t.Error("OneOf('pending', 'active', 'inactive') should fail")
	}
}

func TestAll(t *testing.T) {
	// Test All - should pass when all pass
	err := validator.All(
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
	)
	if err != nil {
		t.Errorf("All with all passing validators should pass, got error: %v", err)
	}

	// Test All - should fail when any fail
	err = validator.All(
		func() error { return nil },
		func() error { return fmt.Errorf("test error") },
		func() error { return nil },
	)
	if err == nil {
		t.Error("All with failing validator should fail")
	}
}

func TestAny(t *testing.T) {
	// Test Any - should pass when at least one passes
	err := validator.Any(
		func() error { return fmt.Errorf("test error") },
		func() error { return nil },
		func() error { return fmt.Errorf("test error") },
	)
	if err != nil {
		t.Errorf("Any with at least one passing validator should pass, got error: %v", err)
	}

	// Test Any - should fail when all fail
	err = validator.Any(
		func() error { return fmt.Errorf("test error") },
		func() error { return fmt.Errorf("test error") },
	)
	if err == nil {
		t.Error("Any with all failing validators should fail")
	}

	// Test Any - should fail with no validators
	err = validator.Any()
	if err == nil {
		t.Error("Any with no validators should fail")
	}
}

func TestValidateSliceLength(t *testing.T) {
	// Test ValidateSliceLength - should pass
	if err := validator.ValidateSliceLength([]int{1, 2, 3}, 3); err != nil {
		t.Errorf("ValidateSliceLength([1,2,3], 3) should pass, got error: %v", err)
	}

	// Test ValidateSliceLength - should fail
	if err := validator.ValidateSliceLength([]int{1, 2, 3}, 2); err == nil {
		t.Error("ValidateSliceLength([1,2,3], 2) should fail")
	}
}

func TestValidateSliceMinLength(t *testing.T) {
	// Test ValidateSliceMinLength - should pass
	if err := validator.ValidateSliceMinLength([]int{1, 2, 3}, 3); err != nil {
		t.Errorf("ValidateSliceMinLength([1,2,3], 3) should pass, got error: %v", err)
	}

	// Test ValidateSliceMinLength - should fail
	if err := validator.ValidateSliceMinLength([]int{1, 2}, 3); err == nil {
		t.Error("ValidateSliceMinLength([1,2], 3) should fail")
	}
}

func TestValidateSliceMaxLength(t *testing.T) {
	// Test ValidateSliceMaxLength - should pass
	if err := validator.ValidateSliceMaxLength([]int{1, 2, 3}, 3); err != nil {
		t.Errorf("ValidateSliceMaxLength([1,2,3], 3) should pass, got error: %v", err)
	}

	// Test ValidateSliceMaxLength - should fail
	if err := validator.ValidateSliceMaxLength([]int{1, 2, 3}, 2); err == nil {
		t.Error("ValidateSliceMaxLength([1,2,3], 2) should fail")
	}
}

func TestValidateSliceLengthRange(t *testing.T) {
	// Test ValidateSliceLengthRange - should pass
	if err := validator.ValidateSliceLengthRange([]int{1, 2, 3}, 2, 4); err != nil {
		t.Errorf("ValidateSliceLengthRange([1,2,3], 2, 4) should pass, got error: %v", err)
	}

	// Test ValidateSliceLengthRange - below minimum
	if err := validator.ValidateSliceLengthRange([]int{1}, 2, 4); err == nil {
		t.Error("ValidateSliceLengthRange([1], 2, 4) should fail")
	}

	// Test ValidateSliceLengthRange - above maximum
	if err := validator.ValidateSliceLengthRange([]int{1, 2, 3, 4, 5}, 2, 4); err == nil {
		t.Error("ValidateSliceLengthRange([1,2,3,4,5], 2, 4) should fail")
	}
}

func TestValidateSliceContains(t *testing.T) {
	// Test ValidateSliceContains - should pass
	if err := validator.ValidateSliceContains([]int{1, 2, 3}, 2); err != nil {
		t.Errorf("ValidateSliceContains([1,2,3], 2) should pass, got error: %v", err)
	}

	// Test ValidateSliceContains - should fail
	if err := validator.ValidateSliceContains([]int{1, 2, 3}, 5); err == nil {
		t.Error("ValidateSliceContains([1,2,3], 5) should fail")
	}
}

func TestValidateSliceNotContains(t *testing.T) {
	// Test ValidateSliceNotContains - should pass
	if err := validator.ValidateSliceNotContains([]int{1, 2, 3}, 5); err != nil {
		t.Errorf("ValidateSliceNotContains([1,2,3], 5) should pass, got error: %v", err)
	}

	// Test ValidateSliceNotContains - should fail
	if err := validator.ValidateSliceNotContains([]int{1, 2, 3}, 2); err == nil {
		t.Error("ValidateSliceNotContains([1,2,3], 2) should fail")
	}
}

func TestValidateSliceUnique(t *testing.T) {
	// Test ValidateSliceUnique - should pass
	if err := validator.ValidateSliceUnique([]int{1, 2, 3}); err != nil {
		t.Errorf("ValidateSliceUnique([1,2,3]) should pass, got error: %v", err)
	}

	// Test ValidateSliceUnique - should fail (duplicate)
	if err := validator.ValidateSliceUnique([]int{1, 2, 2, 3}); err == nil {
		t.Error("ValidateSliceUnique([1,2,2,3]) should fail")
	}
}

func TestValidateMapHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	// Test ValidateMapHasKey - should pass
	if err := validator.ValidateMapHasKey(m, "a"); err != nil {
		t.Errorf("ValidateMapHasKey with existing key should pass, got error: %v", err)
	}

	// Test ValidateMapHasKey - should fail
	if err := validator.ValidateMapHasKey(m, "c"); err == nil {
		t.Error("ValidateMapHasKey with missing key should fail")
	}
}

func TestValidateMapNotHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	// Test ValidateMapNotHasKey - should pass
	if err := validator.ValidateMapNotHasKey(m, "c"); err != nil {
		t.Errorf("ValidateMapNotHasKey with missing key should pass, got error: %v", err)
	}

	// Test ValidateMapNotHasKey - should fail
	if err := validator.ValidateMapNotHasKey(m, "a"); err == nil {
		t.Error("ValidateMapNotHasKey with existing key should fail")
	}
}

func TestValidateMapMinLength(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	// Test ValidateMapMinLength - should pass
	if err := validator.ValidateMapMinLength(m, 2); err != nil {
		t.Errorf("ValidateMapMinLength with length >= min should pass, got error: %v", err)
	}

	// Test ValidateMapMinLength - should fail
	if err := validator.ValidateMapMinLength(m, 3); err == nil {
		t.Error("ValidateMapMinLength with length < min should fail")
	}
}

func TestValidateMapMaxLength(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	// Test ValidateMapMaxLength - should pass
	if err := validator.ValidateMapMaxLength(m, 2); err != nil {
		t.Errorf("ValidateMapMaxLength with length <= max should pass, got error: %v", err)
	}

	// Test ValidateMapMaxLength - should fail
	if err := validator.ValidateMapMaxLength(m, 1); err == nil {
		t.Error("ValidateMapMaxLength with length > max should fail")
	}
}
