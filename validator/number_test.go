package validator_test

import (
	"testing"

	"github.com/julianstephens/go-utils/validator"
)

func TestNumberValidator_Int(t *testing.T) {
	v := validator.Numbers[int]()

	// Test ValidateMin
	if err := v.ValidateMin(5, 3); err != nil {
		t.Errorf("ValidateMin(5, 3) should pass, got error: %v", err)
	}
	if err := v.ValidateMin(2, 3); err == nil {
		t.Error("ValidateMin(2, 3) should fail")
	}

	// Test ValidateMax
	if err := v.ValidateMax(5, 7); err != nil {
		t.Errorf("ValidateMax(5, 7) should pass, got error: %v", err)
	}
	if err := v.ValidateMax(8, 7); err == nil {
		t.Error("ValidateMax(8, 7) should fail")
	}

	// Test ValidateRange
	if err := v.ValidateRange(5, 3, 7); err != nil {
		t.Errorf("ValidateRange(5, 3, 7) should pass, got error: %v", err)
	}
	if err := v.ValidateRange(2, 3, 7); err == nil {
		t.Error("ValidateRange(2, 3, 7) should fail")
	}
	if err := v.ValidateRange(8, 3, 7); err == nil {
		t.Error("ValidateRange(8, 3, 7) should fail")
	}

	// Test ValidatePositive
	if err := v.ValidatePositive(5); err != nil {
		t.Errorf("ValidatePositive(5) should pass, got error: %v", err)
	}
	if err := v.ValidatePositive(0); err == nil {
		t.Error("ValidatePositive(0) should fail")
	}
	if err := v.ValidatePositive(-1); err == nil {
		t.Error("ValidatePositive(-1) should fail")
	}

	// Test ValidateNegative
	if err := v.ValidateNegative(-5); err != nil {
		t.Errorf("ValidateNegative(-5) should pass, got error: %v", err)
	}
	if err := v.ValidateNegative(0); err == nil {
		t.Error("ValidateNegative(0) should fail")
	}
	if err := v.ValidateNegative(1); err == nil {
		t.Error("ValidateNegative(1) should fail")
	}

	// Test ValidateNonNegative
	if err := v.ValidateNonNegative(5); err != nil {
		t.Errorf("ValidateNonNegative(5) should pass, got error: %v", err)
	}
	if err := v.ValidateNonNegative(0); err != nil {
		t.Errorf("ValidateNonNegative(0) should pass, got error: %v", err)
	}
	if err := v.ValidateNonNegative(-1); err == nil {
		t.Error("ValidateNonNegative(-1) should fail")
	}

	// Test ValidateNonPositive
	if err := v.ValidateNonPositive(-5); err != nil {
		t.Errorf("ValidateNonPositive(-5) should pass, got error: %v", err)
	}
	if err := v.ValidateNonPositive(0); err != nil {
		t.Errorf("ValidateNonPositive(0) should pass, got error: %v", err)
	}
	if err := v.ValidateNonPositive(1); err == nil {
		t.Error("ValidateNonPositive(1) should fail")
	}

	// Test ValidateZero
	if err := v.ValidateZero(0); err != nil {
		t.Errorf("ValidateZero(0) should pass, got error: %v", err)
	}
	if err := v.ValidateZero(1); err == nil {
		t.Error("ValidateZero(1) should fail")
	}

	// Test ValidateNonZero
	if err := v.ValidateNonZero(5); err != nil {
		t.Errorf("ValidateNonZero(5) should pass, got error: %v", err)
	}
	if err := v.ValidateNonZero(0); err == nil {
		t.Error("ValidateNonZero(0) should fail")
	}

	// Test ValidateEqual
	if err := v.ValidateEqual(5, 5); err != nil {
		t.Errorf("ValidateEqual(5, 5) should pass, got error: %v", err)
	}
	if err := v.ValidateEqual(5, 3); err == nil {
		t.Error("ValidateEqual(5, 3) should fail")
	}

	// Test ValidateNotEqual
	if err := v.ValidateNotEqual(5, 3); err != nil {
		t.Errorf("ValidateNotEqual(5, 3) should pass, got error: %v", err)
	}
	if err := v.ValidateNotEqual(5, 5); err == nil {
		t.Error("ValidateNotEqual(5, 5) should fail")
	}

	// Test ValidateGreaterThan
	if err := v.ValidateGreaterThan(5, 3); err != nil {
		t.Errorf("ValidateGreaterThan(5, 3) should pass, got error: %v", err)
	}
	if err := v.ValidateGreaterThan(3, 5); err == nil {
		t.Error("ValidateGreaterThan(3, 5) should fail")
	}
	if err := v.ValidateGreaterThan(5, 5); err == nil {
		t.Error("ValidateGreaterThan(5, 5) should fail")
	}

	// Test ValidateGreaterThanOrEqual
	if err := v.ValidateGreaterThanOrEqual(5, 3); err != nil {
		t.Errorf("ValidateGreaterThanOrEqual(5, 3) should pass, got error: %v", err)
	}
	if err := v.ValidateGreaterThanOrEqual(5, 5); err != nil {
		t.Errorf("ValidateGreaterThanOrEqual(5, 5) should pass, got error: %v", err)
	}
	if err := v.ValidateGreaterThanOrEqual(3, 5); err == nil {
		t.Error("ValidateGreaterThanOrEqual(3, 5) should fail")
	}

	// Test ValidateLessThan
	if err := v.ValidateLessThan(3, 5); err != nil {
		t.Errorf("ValidateLessThan(3, 5) should pass, got error: %v", err)
	}
	if err := v.ValidateLessThan(5, 3); err == nil {
		t.Error("ValidateLessThan(5, 3) should fail")
	}
	if err := v.ValidateLessThan(5, 5); err == nil {
		t.Error("ValidateLessThan(5, 5) should fail")
	}

	// Test ValidateLessThanOrEqual
	if err := v.ValidateLessThanOrEqual(3, 5); err != nil {
		t.Errorf("ValidateLessThanOrEqual(3, 5) should pass, got error: %v", err)
	}
	if err := v.ValidateLessThanOrEqual(5, 5); err != nil {
		t.Errorf("ValidateLessThanOrEqual(5, 5) should pass, got error: %v", err)
	}
	if err := v.ValidateLessThanOrEqual(5, 3); err == nil {
		t.Error("ValidateLessThanOrEqual(5, 3) should fail")
	}

	// Test ValidateBetween
	if err := v.ValidateBetween(5, 3, 7); err != nil {
		t.Errorf("ValidateBetween(5, 3, 7) should pass, got error: %v", err)
	}
	if err := v.ValidateBetween(2, 3, 7); err == nil {
		t.Error("ValidateBetween(2, 3, 7) should fail")
	}
	if err := v.ValidateBetween(8, 3, 7); err == nil {
		t.Error("ValidateBetween(8, 3, 7) should fail")
	}
}

func TestNumberValidator_Float(t *testing.T) {
	v := validator.Numbers[float64]()

	// Test ValidateMin
	if err := v.ValidateMin(5.5, 3.0); err != nil {
		t.Errorf("ValidateMin(5.5, 3.0) should pass, got error: %v", err)
	}
	if err := v.ValidateMin(2.5, 3.0); err == nil {
		t.Error("ValidateMin(2.5, 3.0) should fail")
	}

	// Test ValidatePositive
	if err := v.ValidatePositive(3.14); err != nil {
		t.Errorf("ValidatePositive(3.14) should pass, got error: %v", err)
	}
	if err := v.ValidatePositive(0.0); err == nil {
		t.Error("ValidatePositive(0.0) should fail")
	}
	if err := v.ValidatePositive(-1.5); err == nil {
		t.Error("ValidatePositive(-1.5) should fail")
	}
}
