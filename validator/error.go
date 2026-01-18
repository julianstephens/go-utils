package validator

import "fmt"

type ValidationModule int

const (
	ModuleUnknown ValidationModule = iota
	ModuleParse
	ModuleString
	ModuleNumber
)

var (
	ErrEmptyInput = fmt.Errorf("input cannot be empty")

	ErrInvalidInput  = fmt.Errorf("invalid input")
	ErrInvalidFormat = fmt.Errorf("invalid format")

	ErrTooShort = fmt.Errorf("input is too short")
	ErrTooLong  = fmt.Errorf("input is too long")

	ErrMissingUppercase   = fmt.Errorf("missing uppercase letter")
	ErrMissingLowercase   = fmt.Errorf("missing lowercase letter")
	ErrMissingDigit       = fmt.Errorf("missing digit")
	ErrMissingSpecialChar = fmt.Errorf("missing special character")

	ErrStringTooShort         = fmt.Errorf("string is too short")
	ErrStringTooLong          = fmt.Errorf("string is too long")
	ErrStringLengthOutOfRange = fmt.Errorf("string length out of range")
	ErrInvalidPattern         = fmt.Errorf("string does not match pattern")
	ErrNotAlphanumeric        = fmt.Errorf("string contains non-alphanumeric characters")
	ErrNotAlpha               = fmt.Errorf("string contains non-alphabetic characters")
	ErrNotNumeric             = fmt.Errorf("string contains non-numeric characters")
	ErrInvalidSlug            = fmt.Errorf("string is not a valid slug")
	ErrNotLowercase           = fmt.Errorf("string contains uppercase characters")
	ErrNotUppercase           = fmt.Errorf("string contains lowercase characters")
	ErrNotContains            = fmt.Errorf("string does not contain substring")
	ErrContains               = fmt.Errorf("string contains substring")
	ErrInvalidPrefix          = fmt.Errorf("string does not have expected prefix")
	ErrInvalidSuffix          = fmt.Errorf("string does not have expected suffix")

	ErrInvalidInteger   = fmt.Errorf("invalid integer")
	ErrInvalidFloat     = fmt.Errorf("invalid float")
	ErrInvalidBoolean   = fmt.Errorf("invalid boolean")
	ErrInvalidUUID      = fmt.Errorf("invalid UUID")
	ErrInvalidEmail     = fmt.Errorf("invalid email")
	ErrInvalidURL       = fmt.Errorf("invalid URL")
	ErrInvalidIPAddress = fmt.Errorf("invalid IP address")
	ErrInvalidIPv4      = fmt.Errorf("invalid IPv4 address")
	ErrInvalidIPv6      = fmt.Errorf("invalid IPv6 address")
	ErrInvalidDate      = fmt.Errorf("invalid date")
	ErrInvalidDuration  = fmt.Errorf("invalid duration")
	ErrInvalidPhone     = fmt.Errorf("invalid phone number")
	ErrNotInSet         = fmt.Errorf("value not in allowed set")
	ErrSliceTooShort    = fmt.Errorf("slice is too short")
	ErrSliceTooLong     = fmt.Errorf("slice is too long")
	ErrFieldMismatch    = fmt.Errorf("field values do not match")

	ErrNumberTooSmall   = fmt.Errorf("number is too small")
	ErrNumberTooLarge   = fmt.Errorf("number is too large")
	ErrNotPositive      = fmt.Errorf("number is not positive")
	ErrNotNegative      = fmt.Errorf("number is not negative")
	ErrNotZero          = fmt.Errorf("number is not zero")
	ErrNumberOutOfRange = fmt.Errorf("number out of range")
	ErrNotEqual         = fmt.Errorf("number not equal")
	ErrNumberZero       = fmt.Errorf("number is zero")
	ErrNumberPositive   = fmt.Errorf("number is positive")
	ErrNumberNegative   = fmt.Errorf("number is negative")
	ErrNotGreaterThan   = fmt.Errorf("number not greater than threshold")
	ErrNotLessThan      = fmt.Errorf("number not less than threshold")

	ErrUnsupportedType = fmt.Errorf("unsupported type")
)

type ValidationError struct {
	Module ValidationModule
	Cause  string
	Want   any
	Have   any
	Err    error
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation error in module %d: %s (want: %v, have: %v): %v",
		ve.Module,
		ve.Cause,
		ve.Want,
		ve.Have,
		ve.Err,
	)
}

func NewValidationError(module ValidationModule, cause string, want, have any, err error) *ValidationError {
	return &ValidationError{
		Module: module,
		Cause:  cause,
		Want:   want,
		Have:   have,
		Err:    err,
	}
}
