package validator

import (
	"fmt"
	"strings"
)

// Type constraints
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type StringLike interface {
	~string | ~[]byte | ~[]rune
}

type Emptyable interface {
	~string | ~[]byte | ~[]rune | ~map[string]interface{} | ~[]interface{}
}

// Validator is the main validator that provides access to specialized validators.
type Validator struct{}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{}
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
