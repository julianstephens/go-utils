package validator

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ParseValidator struct{}

func (pv *ParseValidator) ValidateEmail(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, err)
	}
	_, err := mail.ParseAddress(input)
	if err != nil {
		return pv.Errorf(
			"invalid email format",
			"valid email address",
			input,
			fmt.Errorf("%w: %v", ErrInvalidEmail, err),
		)
	}
	// Require domain to contain a dot (simple TLD check) to reject addresses like a@b
	parts := strings.Split(input, "@")
	if len(parts) != 2 || !strings.Contains(parts[1], ".") {
		return pv.Errorf("invalid email format", "valid email address", input, ErrInvalidEmail)
	}
	return nil
}

func (pv *ParseValidator) ValidatePassword(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, err)
	}
	if len(input) < 8 {
		return pv.Errorf("password must be at least 8 characters long", "password length >= 8", len(input), ErrTooShort)
	}
	var hasUpper, hasLower, hasDigit bool
	for _, char := range input {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}
	if !hasUpper {
		return pv.Errorf(
			"password must contain at least one uppercase letter",
			"at least one uppercase letter",
			input,
			ErrMissingUppercase,
		)
	}
	if !hasLower {
		return pv.Errorf(
			"password must contain at least one lowercase letter",
			"at least one lowercase letter",
			input,
			ErrMissingLowercase,
		)
	}
	if !hasDigit {
		return pv.Errorf("password must contain at least one digit", "at least one digit", input, ErrMissingDigit)
	}
	return nil
}

func (pv *ParseValidator) ValidateUUID(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, err)
	}
	_, err := uuid.Parse(input)
	if err != nil {
		return pv.Errorf("invalid UUID format", "valid UUID", input, fmt.Errorf("%w: %v", ErrInvalidUUID, err))
	}
	return nil
}

// ValidateBool validates that input can be parsed as a boolean
func (pv *ParseValidator) ValidateBool(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := strconv.ParseBool(input)
	if err != nil {
		return pv.Errorf("invalid boolean format", "valid boolean", input, ErrInvalidBoolean)
	}
	return nil
}

// ValidateInt validates that input can be parsed as an integer
func (pv *ParseValidator) ValidateInt(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return pv.Errorf("invalid integer format", "valid integer", input, ErrInvalidInteger)
	}
	return nil
}

// ValidateUint validates that input can be parsed as an unsigned integer
func (pv *ParseValidator) ValidateUint(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return pv.Errorf("invalid unsigned integer format", "valid unsigned integer", input, ErrInvalidInteger)
	}
	return nil
}

// ValidateFloat validates that input can be parsed as a float
func (pv *ParseValidator) ValidateFloat(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return pv.Errorf("invalid float format", "valid float", input, ErrInvalidFloat)
	}
	return nil
}

// ValidatePositiveInt validates that input is a positive integer
func (pv *ParseValidator) ValidatePositiveInt(input string) error {
	if err := pv.ValidateInt(input); err != nil {
		return err
	}
	val, _ := strconv.ParseInt(input, 10, 64)
	if val <= 0 {
		return pv.Errorf("integer must be positive", "integer > 0", val, ErrNotPositive)
	}
	return nil
}

// ValidateNonNegativeInt validates that input is a non-negative integer
func (pv *ParseValidator) ValidateNonNegativeInt(input string) error {
	if err := pv.ValidateInt(input); err != nil {
		return err
	}
	val, _ := strconv.ParseInt(input, 10, 64)
	if val < 0 {
		return pv.Errorf("integer must be non-negative", "integer >= 0", val, ErrNumberNegative)
	}
	return nil
}

// ValidatePositiveFloat validates that input is a positive float
func (pv *ParseValidator) ValidatePositiveFloat(input string) error {
	if err := pv.ValidateFloat(input); err != nil {
		return err
	}
	val, _ := strconv.ParseFloat(input, 64)
	if val <= 0 {
		return pv.Errorf("float must be positive", "float > 0", val, ErrNotPositive)
	}
	return nil
}

// ValidateURL validates that input is a valid URL
func (pv *ParseValidator) ValidateURL(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return pv.Errorf("invalid URL format", "valid URL", input, ErrInvalidURL)
	}
	return nil
}

// ValidateIPAddress validates that input is a valid IP address (IPv4 or IPv6)
func (pv *ParseValidator) ValidateIPAddress(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	ip := net.ParseIP(input)
	if ip == nil {
		return pv.Errorf("invalid IP address format", "valid IP address", input, ErrInvalidIPAddress)
	}
	return nil
}

// ValidateIPv4 validates that input is a valid IPv4 address
func (pv *ParseValidator) ValidateIPv4(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	ip := net.ParseIP(input)
	if ip == nil || ip.To4() == nil {
		return pv.Errorf("invalid IPv4 address format", "valid IPv4 address", input, ErrInvalidIPv4)
	}
	return nil
}

// ValidateIPv6 validates that input is a valid IPv6 address
func (pv *ParseValidator) ValidateIPv6(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	ip := net.ParseIP(input)
	if ip == nil || ip.To4() != nil {
		return pv.Errorf("invalid IPv6 address format", "valid IPv6 address", input, ErrInvalidIPv6)
	}
	return nil
}

// ValidateDate validates that input can be parsed as a date in the specified format
func (pv *ParseValidator) ValidateDate(input string, format string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := time.Parse(format, input)
	if err != nil {
		return pv.Errorf("invalid date format", format, input, fmt.Errorf("%w: %v", ErrInvalidDate, err))
	}
	return nil
}

// ValidateDuration validates that input can be parsed as a duration (e.g., "5m", "2h", "1s")
func (pv *ParseValidator) ValidateDuration(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	_, err := time.ParseDuration(input)
	if err != nil {
		return pv.Errorf(
			"invalid duration format",
			"valid duration (e.g., 5m, 2h, 1s)",
			input,
			fmt.Errorf("%w: %v", ErrInvalidDuration, err),
		)
	}
	return nil
}

// ValidatePhoneNumber validates that input is a valid phone number (basic format: digits, spaces, hyphens, +)
func (pv *ParseValidator) ValidatePhoneNumber(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return pv.Errorf("input cannot be empty", "non-empty string", input, ErrEmptyInput)
	}
	// Simple regex allowing digits, spaces, hyphens, parentheses, and +
	reg := regexp.MustCompile(`^\+?[0-9\s\-\(\)]+$`)
	if !reg.MatchString(input) {
		return pv.Errorf("invalid phone number format", "valid phone number", input, ErrInvalidPhone)
	}
	return nil
}

func (pv *ParseValidator) Errorf(cause string, want, have any, err error) error {
	return NewValidationError(ModuleParse, cause, want, have, err)
}
