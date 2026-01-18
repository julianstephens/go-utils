package validator

import "fmt"

type NumberValidator[T Number] struct{}

func (nv *NumberValidator[T]) ValidateMin(input T, min T) error {
	if input < min {
		return nv.Errorf("number below minimum", min, input, ErrNumberTooSmall)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateMax(input T, max T) error {
	if input > max {
		return nv.Errorf("number above maximum", max, input, ErrNumberTooLarge)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateRange(input T, min, max T) error {
	if input < min || input > max {
		return nv.Errorf("number out of range", map[string]T{"min": min, "max": max}, input, ErrNumberOutOfRange)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidatePositive(input T) error {
	if input <= 0 {
		return nv.Errorf("number is not positive", "positive number", input, ErrNotPositive)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateNegative(input T) error {
	if input >= 0 {
		return nv.Errorf("number is not negative", "negative number", input, ErrNotNegative)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateNonNegative(input T) error {
	if input < 0 {
		return nv.Errorf("number is negative", "non-negative number", input, ErrNumberNegative)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateNonPositive(input T) error {
	if input > 0 {
		return nv.Errorf("number is positive", "non-positive number", input, ErrNumberPositive)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateZero(input T) error {
	if input != 0 {
		return nv.Errorf("number is not zero", 0, input, ErrNotZero)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateNonZero(input T) error {
	if input == 0 {
		return nv.Errorf("number is zero", "non-zero number", input, ErrNumberZero)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateEqual(input T, expected T) error {
	if input != expected {
		return nv.Errorf("number not equal", expected, input, ErrNotEqual)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateNotEqual(input T, notExpected T) error {
	if input == notExpected {
		return nv.Errorf("number equal to disallowed value", notExpected, input, ErrNotEqual)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateGreaterThan(input T, threshold T) error {
	if input <= threshold {
		return nv.Errorf("number not greater than threshold", fmt.Sprintf("> %v", threshold), input, ErrNotGreaterThan)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateGreaterThanOrEqual(input T, threshold T) error {
	if input < threshold {
		return nv.Errorf("number not greater than or equal to threshold", fmt.Sprintf(">= %v", threshold), input, ErrNotGreaterThan)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateLessThan(input T, threshold T) error {
	if input >= threshold {
		return nv.Errorf("number not less than threshold", fmt.Sprintf("< %v", threshold), input, ErrNotLessThan)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateLessThanOrEqual(input T, threshold T) error {
	if input > threshold {
		return nv.Errorf("number not less than or equal to threshold", fmt.Sprintf("<= %v", threshold), input, ErrNotLessThan)
	}
	return nil
}

func (nv *NumberValidator[T]) ValidateBetween(input T, lower T, upper T) error {
	if input < lower || input > upper {
		return nv.Errorf("number not between bounds", fmt.Sprintf("[%v, %v]", lower, upper), input, ErrNumberOutOfRange)
	}
	return nil
}

func (nv *NumberValidator[T]) Errorf(cause string, want, have any, err error) *ValidationError {
	return NewValidationError(ModuleNumber, cause, want, have, err)
}
