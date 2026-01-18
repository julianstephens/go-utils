package validator

import (
	"errors"
	"fmt"
	"reflect"
)

type NumberValidator[T Number] struct{}

// ValidateMin validates that a number is greater than or equal to the minimum value
func (nv *NumberValidator[T]) ValidateMin(input T, min T) error {
	if input < min {
		return nv.Errorf("number below minimum", min, input, ErrNumberTooSmall)
	}
	return nil
}

// ValidateMax validates that a number is less than or equal to the maximum value
func (nv *NumberValidator[T]) ValidateMax(input T, max T) error {
	if input > max {
		return nv.Errorf("number above maximum", max, input, ErrNumberTooLarge)
	}
	return nil
}

// ValidateRange validates that a number is within the specified range (inclusive)
func (nv *NumberValidator[T]) ValidateRange(input T, min, max T) error {
	if input < min || input > max {
		return nv.Errorf("number out of range", map[string]T{"min": min, "max": max}, input, ErrNumberOutOfRange)
	}
	return nil
}

// ValidatePositive validates that a number is greater than zero
func (nv *NumberValidator[T]) ValidatePositive(input T) error {
	if input <= 0 {
		return nv.Errorf("number is not positive", "positive number", input, ErrNotPositive)
	}
	return nil
}

// ValidateNegative validates that a number is less than zero
func (nv *NumberValidator[T]) ValidateNegative(input T) error {
	if input >= 0 {
		return nv.Errorf("number is not negative", "negative number", input, ErrNotNegative)
	}
	return nil
}

// ValidateNonNegative validates that a number is greater than or equal to zero
func (nv *NumberValidator[T]) ValidateNonNegative(input T) error {
	if input < 0 {
		return nv.Errorf("number is negative", "non-negative number", input, ErrNumberNegative)
	}
	return nil
}

// ValidateNonPositive validates that a number is less than or equal to zero
func (nv *NumberValidator[T]) ValidateNonPositive(input T) error {
	if input > 0 {
		return nv.Errorf("number is positive", "non-positive number", input, ErrNumberPositive)
	}
	return nil
}

// ValidateZero validates that a number equals zero
func (nv *NumberValidator[T]) ValidateZero(input T) error {
	if input != 0 {
		return nv.Errorf("number is not zero", 0, input, ErrNotZero)
	}
	return nil
}

// ValidateNonZero validates that a number is not equal to zero
func (nv *NumberValidator[T]) ValidateNonZero(input T) error {
	if input == 0 {
		return nv.Errorf("number is zero", "non-zero number", input, ErrNumberZero)
	}
	return nil
}

// ValidateEqual validates that a number equals the expected value
func (nv *NumberValidator[T]) ValidateEqual(input T, expected T) error {
	if input != expected {
		return nv.Errorf("number not equal", expected, input, ErrNotEqual)
	}
	return nil
}

// ValidateNotEqual validates that a number is not equal to the specified value
func (nv *NumberValidator[T]) ValidateNotEqual(input T, notExpected T) error {
	if input == notExpected {
		return nv.Errorf("number equal to disallowed value", notExpected, input, ErrNotEqual)
	}
	return nil
}

// ValidateGreaterThan validates that a number is strictly greater than the threshold
func (nv *NumberValidator[T]) ValidateGreaterThan(input T, threshold T) error {
	if input <= threshold {
		return nv.Errorf("number not greater than threshold", fmt.Sprintf("> %v", threshold), input, ErrNotGreaterThan)
	}
	return nil
}

// ValidateGreaterThanOrEqual validates that a number is greater than or equal to the threshold
func (nv *NumberValidator[T]) ValidateGreaterThanOrEqual(input T, threshold T) error {
	if input < threshold {
		return nv.Errorf(
			"number not greater than or equal to threshold",
			fmt.Sprintf(">= %v", threshold),
			input,
			ErrNotGreaterThan,
		)
	}
	return nil
}

// ValidateLessThan validates that a number is strictly less than the threshold
func (nv *NumberValidator[T]) ValidateLessThan(input T, threshold T) error {
	if input >= threshold {
		return nv.Errorf("number not less than threshold", fmt.Sprintf("< %v", threshold), input, ErrNotLessThan)
	}
	return nil
}

// ValidateLessThanOrEqual validates that a number is less than or equal to the threshold
func (nv *NumberValidator[T]) ValidateLessThanOrEqual(input T, threshold T) error {
	if input > threshold {
		return nv.Errorf(
			"number not less than or equal to threshold",
			fmt.Sprintf("<= %v", threshold),
			input,
			ErrNotLessThan,
		)
	}
	return nil
}

// ValidateBetween validates that a number is within the specified bounds (inclusive)
func (nv *NumberValidator[T]) ValidateBetween(input T, lower T, upper T) error {
	if input < lower || input > upper {
		return nv.Errorf("number not between bounds", fmt.Sprintf("[%v, %v]", lower, upper), input, ErrNumberOutOfRange)
	}
	return nil
}

// ValidateConsecutive validates that the second number is exactly one greater than the first
func (nv *NumberValidator[T]) ValidateConsecutive(input1, input2 T) error {
	if input2 != input1+1 {
		return nv.Errorf(
			"numbers are not consecutive",
			fmt.Sprintf("%v followed by %v", input1, input1+1),
			input2,
			fmt.Errorf("%w:%v", ErrNumberOutOfRange, errors.New("not consecutive")),
		)
	}
	return nil
}

// ValidateEven validates that a number is even (integer types only)
func (nv *NumberValidator[T]) ValidateEven(input T) error {
	intVal, err := toInt64(input)
	if err != nil {
		return nv.Errorf("type does not support even/odd validation", "integer type", input, err)
	}
	if intVal%2 != 0 {
		return nv.Errorf("number is not even", "even number", input, fmt.Errorf("%w:%v", ErrNotEqual, "odd number"))
	}
	return nil
}

// ValidateOdd validates that a number is odd (integer types only)
func (nv *NumberValidator[T]) ValidateOdd(input T) error {
	intVal, err := toInt64(input)
	if err != nil {
		return nv.Errorf("type does not support even/odd validation", "integer type", input, err)
	}
	if intVal%2 == 0 {
		return nv.Errorf("number is not odd", "odd number", input, fmt.Errorf("%w:%v", ErrNotEqual, "even number"))
	}
	return nil
}

// ValidateDivisibleBy validates that a number is evenly divisible by the divisor (integer types only)
func (nv *NumberValidator[T]) ValidateDivisibleBy(input T, divisor T) error {
	intInput, err := toInt64(input)
	if err != nil {
		return nv.Errorf("type does not support divisibility validation", "integer type", input, err)
	}
	intDivisor, err := toInt64(divisor)
	if err != nil {
		return nv.Errorf("type does not support divisibility validation", "integer type", divisor, err)
	}
	if intDivisor == 0 {
		return nv.Errorf("divisor cannot be zero", "non-zero divisor", divisor, fmt.Errorf("division by zero"))
	}
	if intInput%intDivisor != 0 {
		return nv.Errorf(
			"number is not divisible by divisor",
			fmt.Sprintf("divisible by %v", divisor),
			input,
			fmt.Errorf("%w:%v", ErrNotEqual, "not divisible"),
		)
	}
	return nil
}

// ValidatePowerOf validates that a number is a power of the specified base (integer types only)
func (nv *NumberValidator[T]) ValidatePowerOf(input T, base T) error {
	intInput, err := toInt64(input)
	if err != nil {
		return nv.Errorf("type does not support power validation", "integer type", input, err)
	}
	intBase, err := toInt64(base)
	if err != nil {
		return nv.Errorf("type does not support power validation", "integer type", base, err)
	}
	if intBase <= 1 {
		return nv.Errorf("base must be greater than 1", "> 1", base, fmt.Errorf("invalid base for power"))
	}
	power := int64(1)
	for power < intInput {
		power *= intBase
	}
	if power != intInput {
		return nv.Errorf(
			"number is not a power of the base",
			fmt.Sprintf("power of %v", base),
			input,
			fmt.Errorf("%w:%v", ErrNotEqual, "not a power"),
		)
	}
	return nil
}

// ValidateFibonacci validates that a number is a Fibonacci number (integer types only)
func (nv *NumberValidator[T]) ValidateFibonacci(input T) error {
	intInput, err := toInt64(input)
	if err != nil {
		return nv.Errorf("type does not support Fibonacci validation", "integer type", input, err)
	}
	if intInput < 0 {
		return nv.Errorf(
			"Fibonacci number cannot be negative",
			"non-negative number",
			input,
			fmt.Errorf("negative number"),
		)
	}
	a, b := int64(0), int64(1)
	for a < intInput {
		a, b = b, a+b
	}
	if a != intInput {
		return nv.Errorf(
			"number is not a Fibonacci number",
			"Fibonacci number",
			input,
			fmt.Errorf("%w:%v", ErrNotEqual, "not a Fibonacci number"),
		)
	}
	return nil
}

func (nv *NumberValidator[T]) Errorf(cause string, want, have any, err error) *ValidationError {
	return NewValidationError(ModuleNumber, cause, want, have, err)
}

// toInt64 converts a Number type to int64, returning an error if it's a float type
func toInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32, float64:
		return 0, fmt.Errorf("%w: float types do not support even/odd validation", ErrUnsupportedType)
	default:
		return 0, fmt.Errorf("%w: %v", ErrUnsupportedType, reflect.TypeOf(v))
	}
}
