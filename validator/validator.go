package validator

import (
	"fmt"
	"strings"
)

// Type constraints
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Number interface {
	Integer | ~float32 | ~float64
}

type StringLike interface {
	~string | ~[]byte | ~[]rune
}

type Emptyable interface {
	~string | ~[]byte | ~[]rune | ~map[string]interface{} | ~[]interface{}
}

// Numbers returns a NumberValidator for the specified numeric type.
func Numbers[T Number]() *NumberValidator[T] {
	return &NumberValidator[T]{}
}

// Strings returns a StringValidator for the specified string-like type.
func Strings[T StringLike]() *StringValidator[T] {
	return &StringValidator[T]{
		Parse: &ParseValidator{},
	}
}

// Parse returns a ParseValidator for parsing string inputs.
func Parse() *ParseValidator {
	return &ParseValidator{}
}

// ValidateNonEmpty validates that input is not empty
func ValidateNonEmpty[T Emptyable](input T) error {
	switch v := any(input).(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("input cannot be empty")
		}
	case []byte:
		if len(v) == 0 {
			return fmt.Errorf("input cannot be empty")
		}
	case []rune:
		if len(v) == 0 {
			return fmt.Errorf("input cannot be empty")
		}
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("input cannot be empty")
		}
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("input cannot be empty")
		}
	default:
		return fmt.Errorf("unsupported type for emptiness check")
	}
	return nil
}

// ValidateMatchesField validates that two comparable values are equal (useful for password confirmation)
func ValidateMatchesField[T comparable](value1 T, value2 T, fieldName string) error {
	if value1 != value2 {
		return fmt.Errorf("%w: %s does not match confirmation", ErrFieldMismatch, fieldName)
	}
	return nil
}

// CustomValidator provides a builder pattern for composing multiple validators
type CustomValidator struct {
	validators []func() error
}

// NewCustomValidator creates a new CustomValidator instance
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validators: []func() error{},
	}
}

// Add appends a validator function to the custom validator
func (cv *CustomValidator) Add(validator func() error) *CustomValidator {
	cv.validators = append(cv.validators, validator)
	return cv
}

// Validate runs all accumulated validators in sequence (AND logic)
func (cv *CustomValidator) Validate() error {
	for _, validator := range cv.validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}
