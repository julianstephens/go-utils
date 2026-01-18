package validator

type StringValidator[T StringLike] struct {
	Parse *ParseValidator
}

func (sv *StringValidator[T]) ValidateMinLength(input T, min int) error {
	if len(input) < min {
		return sv.Errorf("string below minimum length", min, len(input), ErrStringTooShort)
	}
	return nil
}

func (sv *StringValidator[T]) ValidateMaxLength(input T, max int) error {
	if len(input) > max {
		return sv.Errorf("string above maximum length", max, len(input), ErrStringTooLong)
	}
	return nil
}

func (sv *StringValidator[T]) ValidateLengthRange(input T, min, max int) error {
	length := len(input)
	if length < min || length > max {
		return sv.Errorf("string length out of range", map[string]int{"min": min, "max": max}, length, ErrStringLengthOutOfRange)
	}
	return nil
}

func (sv *StringValidator[T]) Errorf(cause string, want, have any, err error) *ValidationError {
	return NewValidationError(ModuleString, cause, want, have, err)
}
