# Validator Package

A comprehensive, type-safe validation library for Go applications with both simple functions and advanced generic validators.

## Features

- **Type-safe generic validators** for numbers, strings, and parsing
- **Simple validation functions** for common checks
- **Composable validation chains** with detailed error messages
- **Support for multiple types** (int, uint, float, string, []byte, etc.)
- **Comprehensive test coverage** (92.1%)

## Quick Start

### Simple Validation

```go
import "github.com/julianstephens/go-utils/validator"

// Basic validation
if err := validator.ValidateNonEmpty("input"); err != nil {
    // handle empty input
}
```

### Advanced Generic Validators

```go
// Number validation with type safety
numValidator := validator.Numbers[int]()
if err := numValidator.ValidateRange(42, 1, 100); err != nil {
    // handle number out of range
}

// String validation
strValidator := validator.Strings[string]()
if err := strValidator.ValidateMinLength("hello", 3); err != nil {
    // handle string too short
}

// Parsing validation (accessed through string validator)
strValidator := validator.Strings[string]()
if err := strValidator.Parse.ValidateURL("https://example.com"); err != nil {
    // handle invalid URL
}
```

## API Reference

### Factory Functions

- `Numbers[T]() *NumberValidator[T]` - Create a number validator for type T
- `Strings[T]() *StringValidator[T]` - Create a string validator for type T (includes Parse field for parsing validation)
- `Parse() *ParseValidator` - Create a standalone parsing validator (advanced usage)

### Number Validators

All number validators work with `Number` types: `int`, `uint`, `float32`, `float64` and their variants.

#### Range Validation
- `ValidateMin(input T, min T) error` - Input must be ≥ minimum
- `ValidateMax(input T, max T) error` - Input must be ≤ maximum
- `ValidateRange(input T, min, max T) error` - Input must be within range

#### Sign Validation
- `ValidatePositive(input T) error` - Input must be > 0
- `ValidateNegative(input T) error` - Input must be < 0
- `ValidateNonNegative(input T) error` - Input must be ≥ 0
- `ValidateNonPositive(input T) error` - Input must be ≤ 0

#### Equality Validation
- `ValidateZero(input T) error` - Input must be exactly 0
- `ValidateNonZero(input T) error` - Input must not be 0
- `ValidateEqual(input T, expected T) error` - Input must equal expected value
- `ValidateNotEqual(input T, notExpected T) error` - Input must not equal value

#### Comparison Validation
- `ValidateGreaterThan(input T, threshold T) error` - Input must be > threshold
- `ValidateGreaterThanOrEqual(input T, threshold T) error` - Input must be ≥ threshold
- `ValidateLessThan(input T, threshold T) error` - Input must be < threshold
- `ValidateLessThanOrEqual(input T, threshold T) error` - Input must be ≤ threshold
- `ValidateBetween(input T, lower, upper T) error` - Input must be between bounds

### String Validators

String validators work with `StringLike` types: `string`, `[]byte`, `[]rune`.

#### Length Validation
- `ValidateMinLength(input T, min int) error` - String length must be ≥ minimum
- `ValidateMaxLength(input T, max int) error` - String length must be ≤ maximum
- `ValidateLengthRange(input T, min, max int) error` - String length must be within range

### Parse Validators

Parsing validators for string format validation.

#### Format Validation
- `ValidateEmail(input string) error` - Valid email address format
- `ValidatePassword(input string) error` - Password with complexity requirements (8+ chars, mixed case, digits)
- `ValidateUUID(input string) error` - Valid UUID format
- `ValidateURL(input string) error` - Valid URL format

#### Type Parsing
- `ValidateBool(input string) error` - Parseable as boolean (true/false/1/0)
- `ValidateInt(input string) error` - Parseable as int64
- `ValidateUint(input string) error` - Parseable as uint64
- `ValidateFloat(input string) error` - Parseable as float64

#### Constrained Parsing
- `ValidatePositiveInt(input string) error` - Parseable positive integer
- `ValidateNonNegativeInt(input string) error` - Parseable non-negative integer
- `ValidatePositiveFloat(input string) error` - Parseable positive float

#### Network Validation
- `ValidateIPAddress(input string) error` - Valid IPv4 or IPv6 address
- `ValidateIPv4(input string) error` - Valid IPv4 address
- `ValidateIPv6(input string) error` - Valid IPv6 address

### Utility Functions

- `ValidateNonEmpty[T](input T) error` - Generic emptiness check for strings, bytes, runes, maps, and slices
- `New() *Validator` - Create a new validator instance (legacy)
- `Parse() *ParseValidator` - Standalone parsing validator (typically accessed via StringValidator.Parse)

## Type Constraints

The package uses Go generics with these type constraints:

- `Number`: `~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64`
- `StringLike`: `~string | ~[]byte | ~[]rune`
- `Emptyable`: `~string | ~[]byte | ~[]rune | ~map[string]interface{} | ~[]interface{}`

## Error Handling

All validators return `*ValidationError` with detailed context:

```go
if err := numValidator.ValidateRange(-5, 0, 100); err != nil {
    var valErr *validator.ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Validation failed: %s (expected: %v, got: %v)",
            valErr.Cause, valErr.Want, valErr.Have)
    }
}
```

## Examples

### User Registration Validation

```go
func validateUser(name, email, ageStr string) error {
    // String validation
    strVal := validator.Strings[string]()
    if err := strVal.ValidateMinLength(name, 2); err != nil {
        return fmt.Errorf("name validation: %w", err)
    }

    // Parse validation (accessed through string validator)
    if err := strVal.Parse.ValidateEmail(email); err != nil {
        return fmt.Errorf("email validation: %w", err)
    }

    // Number validation after parsing
    age, err := strconv.Atoi(ageStr)
    if err != nil {
        return fmt.Errorf("age parsing: %w", err)
    }

    numVal := validator.Numbers[int]()
    if err := numVal.ValidateRange(age, 13, 120); err != nil {
        return fmt.Errorf("age validation: %w", err)
    }

    return nil
}
```

### Configuration Validation

```go
func validateConfig(port int, host string, timeout float64) error {
    // Port validation
    portVal := validator.Numbers[int]()
    if err := portVal.ValidateRange(port, 1024, 65535); err != nil {
        return fmt.Errorf("invalid port: %w", err)
    }

    // Host validation
    hostVal := validator.Strings[string]()
    if err := hostVal.ValidateMinLength(host, 1); err != nil {
        return fmt.Errorf("invalid host: %w", err)
    }

    // Timeout validation
    timeoutVal := validator.Numbers[float64]()
    if err := timeoutVal.ValidateRange(timeout, 0.1, 300.0); err != nil {
        return fmt.Errorf("invalid timeout: %w", err)
    }

    return nil
}
```

## Testing

The package includes comprehensive tests with 92.1% coverage:

```bash
go test -v ./validator
go test -cover ./validator
```

Tests are organized by validator type:
- `validator_test.go` - General utility functions
- `number_test.go` - Number validator tests
- `string_test.go` - String validator tests
- `parse_test.go` - Parse validator tests
