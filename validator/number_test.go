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

	// Test ValidateConsecutive
	if err := v.ValidateConsecutive(5, 6); err != nil {
		t.Errorf("ValidateConsecutive(5, 6) should pass, got error: %v", err)
	}
	if err := v.ValidateConsecutive(5, 7); err == nil {
		t.Error("ValidateConsecutive(5, 7) should fail")
	}
	if err := v.ValidateConsecutive(5, 5); err == nil {
		t.Error("ValidateConsecutive(5, 5) should fail")
	}
	if err := v.ValidateConsecutive(-1, 0); err != nil {
		t.Errorf("ValidateConsecutive(-1, 0) should pass, got error: %v", err)
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

	// Test ValidateConsecutive
	if err := v.ValidateConsecutive(5.0, 6.0); err != nil {
		t.Errorf("ValidateConsecutive(5.0, 6.0) should pass, got error: %v", err)
	}
	if err := v.ValidateConsecutive(5.0, 7.0); err == nil {
		t.Error("ValidateConsecutive(5.0, 7.0) should fail")
	}
}

func TestNumberValidator_Even(t *testing.T) {
	v := validator.Numbers[int]()

	// Test ValidateEven
	if err := v.ValidateEven(4); err != nil {
		t.Errorf("ValidateEven(4) should pass, got error: %v", err)
	}
	if err := v.ValidateEven(0); err != nil {
		t.Errorf("ValidateEven(0) should pass, got error: %v", err)
	}
	if err := v.ValidateEven(-2); err != nil {
		t.Errorf("ValidateEven(-2) should pass, got error: %v", err)
	}
	if err := v.ValidateEven(3); err == nil {
		t.Error("ValidateEven(3) should fail")
	}
	if err := v.ValidateEven(-1); err == nil {
		t.Error("ValidateEven(-1) should fail")
	}
}

func TestNumberValidator_Odd(t *testing.T) {
	v := validator.Numbers[int]()

	// Test ValidateOdd
	if err := v.ValidateOdd(3); err != nil {
		t.Errorf("ValidateOdd(3) should pass, got error: %v", err)
	}
	if err := v.ValidateOdd(-1); err != nil {
		t.Errorf("ValidateOdd(-1) should pass, got error: %v", err)
	}
	if err := v.ValidateOdd(1); err != nil {
		t.Errorf("ValidateOdd(1) should pass, got error: %v", err)
	}
	if err := v.ValidateOdd(4); err == nil {
		t.Error("ValidateOdd(4) should fail")
	}
	if err := v.ValidateOdd(0); err == nil {
		t.Error("ValidateOdd(0) should fail")
	}
}

// Test that float types throw an error
func TestNumberValidator_EvenOdd_FloatError(t *testing.T) {
	v := validator.Numbers[float64]()

	if err := v.ValidateEven(4.0); err == nil {
		t.Error("ValidateEven(4.0) should fail with float type")
	}
	if err := v.ValidateOdd(3.0); err == nil {
		t.Error("ValidateOdd(3.0) should fail with float type")
	}
}

func TestValidateDivisibleBy(t *testing.T) {
	v := validator.Numbers[int]()

	// Test ValidateDivisibleBy
	if err := v.ValidateDivisibleBy(10, 2); err != nil {
		t.Errorf("ValidateDivisibleBy(10, 2) should pass, got error: %v", err)
	}
	if err := v.ValidateDivisibleBy(10, 3); err == nil {
		t.Error("ValidateDivisibleBy(10, 3) should fail")
	}
	if err := v.ValidateDivisibleBy(10, 0); err == nil {
		t.Error("ValidateDivisibleBy(10, 0) should fail due to division by zero")
	}
}

func TestValidatePowerOf(t *testing.T) {
	v := validator.Numbers[int]()

	// Test ValidatePowerOf
	if err := v.ValidatePowerOf(8, 2); err != nil {
		t.Errorf("ValidatePowerOf(8, 2) should pass, got error: %v", err)
	}
	if err := v.ValidatePowerOf(9, 3); err != nil {
		t.Errorf("ValidatePowerOf(9, 3) should pass, got error: %v", err)
	}
	if err := v.ValidatePowerOf(10, 2); err == nil {
		t.Error("ValidatePowerOf(10, 2) should fail")
	}
	if err := v.ValidatePowerOf(27, 3); err != nil {
		t.Errorf("ValidatePowerOf(27, 3) should pass, got error: %v", err)
	}
	if err := v.ValidatePowerOf(16, 4); err != nil {
		t.Errorf("ValidatePowerOf(16, 4) should pass, got error: %v", err)
	}
	if err := v.ValidatePowerOf(20, 4); err == nil {
		t.Error("ValidatePowerOf(20, 4) should fail")
	}
	if err := v.ValidatePowerOf(8, 1); err == nil {
		t.Error("ValidatePowerOf(8, 1) should fail due to invalid base")
	}
	if err := v.ValidatePowerOf(8, 0); err == nil {
		t.Error("ValidatePowerOf(8, 0) should fail due to invalid base")
	}
}

// Test that float types throw an error
func TestNumberValidator_DivisibleBy_PowerOf_FloatError(t *testing.T) {
	v := validator.Numbers[float64]()

	if err := v.ValidateDivisibleBy(10.0, 2.0); err == nil {
		t.Error("ValidateDivisibleBy(10.0, 2.0) should fail with float type")
	}
	if err := v.ValidatePowerOf(8.0, 2.0); err == nil {
		t.Error("ValidatePowerOf(8.0, 2.0) should fail with float type")
	}
}

func TestValidateFibonacci(t *testing.T) {
	v := validator.Numbers[int]()

	// Test ValidateFibonacci
	if err := v.ValidateFibonacci(8); err != nil {
		t.Errorf("ValidateFibonacci(8) should pass, got error: %v", err)
	}
	if err := v.ValidateFibonacci(13); err != nil {
		t.Errorf("ValidateFibonacci(13) should pass, got error: %v", err)
	}
	if err := v.ValidateFibonacci(10); err == nil {
		t.Error("ValidateFibonacci(10) should fail")
	}
}

// Test that float types throw an error
func TestNumberValidator_Fibonacci_FloatError(t *testing.T) {
	v := validator.Numbers[float64]()

	if err := v.ValidateFibonacci(8.0); err == nil {
		t.Error("ValidateFibonacci(8.0) should fail with float type")
	}
}
