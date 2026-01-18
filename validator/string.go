package validator

import (
	"regexp"
	"strings"
	"unicode"
)

type StringValidator[T StringLike] struct {
	Parse *ParseValidator
}

// ValidateMinLength validates that a string is at least the minimum length
func (sv *StringValidator[T]) ValidateMinLength(input T, min int) error {
	if len(input) < min {
		return sv.Errorf("string below minimum length", min, len(input), ErrStringTooShort)
	}
	return nil
}

// ValidateMaxLength validates that a string is at most the maximum length
func (sv *StringValidator[T]) ValidateMaxLength(input T, max int) error {
	if len(input) > max {
		return sv.Errorf("string above maximum length", max, len(input), ErrStringTooLong)
	}
	return nil
}

// ValidateLengthRange validates that a string length is within the specified range (inclusive)
func (sv *StringValidator[T]) ValidateLengthRange(input T, min, max int) error {
	length := len(input)
	if length < min || length > max {
		return sv.Errorf("string length out of range", map[string]int{"min": min, "max": max}, length, ErrStringLengthOutOfRange)
	}
	return nil
}

// ValidatePattern validates that a string matches the specified regex pattern
func (sv *StringValidator[T]) ValidatePattern(input T, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return sv.Errorf("invalid regex pattern", "valid regex pattern", pattern, err)
	}
	str := toString(input)
	if !regex.MatchString(str) {
		return sv.Errorf("string does not match pattern", pattern, str, ErrInvalidPattern)
	}
	return nil
}

// ValidateAlphanumeric validates that a string contains only alphanumeric characters
func (sv *StringValidator[T]) ValidateAlphanumeric(input T) error {
	str := toString(input)
	for _, ch := range str {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return sv.Errorf("string contains non-alphanumeric characters", "only letters and digits", str, ErrNotAlphanumeric)
		}
	}
	return nil
}

// ValidateAlpha validates that a string contains only alphabetic characters
func (sv *StringValidator[T]) ValidateAlpha(input T) error {
	str := toString(input)
	for _, ch := range str {
		if !unicode.IsLetter(ch) {
			return sv.Errorf("string contains non-alphabetic characters", "only letters", str, ErrNotAlpha)
		}
	}
	return nil
}

// ValidateNumeric validates that a string contains only numeric characters
func (sv *StringValidator[T]) ValidateNumeric(input T) error {
	str := toString(input)
	for _, ch := range str {
		if !unicode.IsDigit(ch) {
			return sv.Errorf("string contains non-numeric characters", "only digits", str, ErrNotNumeric)
		}
	}
	return nil
}

// ValidateSlug validates that a string is a valid URL slug (lowercase alphanumeric with hyphens and underscores)
func (sv *StringValidator[T]) ValidateSlug(input T) error {
	str := toString(input)
	for _, ch := range str {
		if !unicode.IsLower(ch) && !unicode.IsDigit(ch) && ch != '-' && ch != '_' {
			return sv.Errorf("string is not a valid slug", "lowercase alphanumeric with hyphens/underscores", str, ErrInvalidSlug)
		}
	}
	return nil
}

// ValidateLowercase validates that a string contains only lowercase characters (letters only)
func (sv *StringValidator[T]) ValidateLowercase(input T) error {
	str := toString(input)
	for _, ch := range str {
		if unicode.IsLetter(ch) && unicode.IsUpper(ch) {
			return sv.Errorf("string contains uppercase characters", "only lowercase letters", str, ErrNotLowercase)
		}
	}
	return nil
}

// ValidateUppercase validates that a string contains only uppercase characters (letters only)
func (sv *StringValidator[T]) ValidateUppercase(input T) error {
	str := toString(input)
	for _, ch := range str {
		if unicode.IsLetter(ch) && unicode.IsLower(ch) {
			return sv.Errorf("string contains lowercase characters", "only uppercase letters", str, ErrNotUppercase)
		}
	}
	return nil
}

// ValidateContains validates that a string contains the specified substring
func (sv *StringValidator[T]) ValidateContains(input T, substring string) error {
	str := toString(input)
	if !strings.Contains(str, substring) {
		return sv.Errorf("string does not contain substring", substring, str, ErrNotContains)
	}
	return nil
}

// ValidateNotContains validates that a string does not contain the specified substring
func (sv *StringValidator[T]) ValidateNotContains(input T, substring string) error {
	str := toString(input)
	if strings.Contains(str, substring) {
		return sv.Errorf("string contains forbidden substring", "should not contain"+substring, str, ErrContains)
	}
	return nil
}

// ValidatePrefix validates that a string starts with the specified prefix
func (sv *StringValidator[T]) ValidatePrefix(input T, prefix string) error {
	str := toString(input)
	if !strings.HasPrefix(str, prefix) {
		return sv.Errorf("string does not have expected prefix", prefix, str, ErrInvalidPrefix)
	}
	return nil
}

// ValidateSuffix validates that a string ends with the specified suffix
func (sv *StringValidator[T]) ValidateSuffix(input T, suffix string) error {
	str := toString(input)
	if !strings.HasSuffix(str, suffix) {
		return sv.Errorf("string does not have expected suffix", suffix, str, ErrInvalidSuffix)
	}
	return nil
}

func (sv *StringValidator[T]) Errorf(cause string, want, have any, err error) *ValidationError {
	return NewValidationError(ModuleString, cause, want, have, err)
}

// toString converts a StringLike type to a string for processing
func toString[T StringLike](input T) string {
	switch v := any(input).(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case []rune:
		return string(v)
	default:
		return ""
	}
}
