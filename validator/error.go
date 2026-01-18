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

	ErrInvalidInteger   = fmt.Errorf("invalid integer")
	ErrInvalidFloat     = fmt.Errorf("invalid float")
	ErrInvalidBoolean   = fmt.Errorf("invalid boolean")
	ErrInvalidUUID      = fmt.Errorf("invalid UUID")
	ErrInvalidEmail     = fmt.Errorf("invalid email")
	ErrInvalidURL       = fmt.Errorf("invalid URL")
	ErrInvalidIPAddress = fmt.Errorf("invalid IP address")
	ErrInvalidIPv4      = fmt.Errorf("invalid IPv4 address")
	ErrInvalidIPv6      = fmt.Errorf("invalid IPv6 address")

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
)

var (
	ErrStringTooShort         = fmt.Errorf("string is too short")
	ErrStringTooLong          = fmt.Errorf("string is too long")
	ErrStringLengthOutOfRange = fmt.Errorf("string length out of range")
)

type ValidationError struct {
	Module ValidationModule
	Cause  string
	Want   any
	Have   any
	Err    error
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error in module %d: %s (want: %v, have: %v): %v", ve.Module, ve.Cause, ve.Want, ve.Have, ve.Err)
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
