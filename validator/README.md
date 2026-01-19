# Validator Package

A comprehensive, type-safe validation library for Go applications with both simple functions and advanced generic validators.

## Features

- **Type-safe generic validators** for numbers, strings, and parsing
- **Composite validators** with AND/OR logic and field matching
- **Collection validators** for slices and maps
- **Integer-specific validators** (even/odd, divisibility, powers, Fibonacci)
- **String validators** (regex, character types, substring matching)
- **Date and duration parsing** with custom formats
- **Custom validator builder** for fluent chaining
- **Comprehensive test coverage**

## Quick Start

```go
import "github.com/julianstephens/go-utils/validator"

// Simple validation
validator.ValidateNonEmpty("input")
validator.OneOf("active", "active", "inactive", "pending")

// Number validation
numValidator := validator.Numbers[int]()
numValidator.ValidateRange(42, 1, 100)

// String validation
strValidator := validator.Strings[string]()
strValidator.ValidateMinLength("hello", 3)
strValidator.Parse.ValidateEmail("user@example.com")

// Composite validation (AND/OR logic)
validator.All(
    func() error { return numValidator.ValidateRange(age, 18, 65) },
    func() error { return strValidator.ValidateMinLength(email, 5) },
)

// Field matching validation
validator.ValidateMatchesField(password, confirmPassword, "password")

// Custom validator builder
customVal := validator.NewCustomValidator().
    Add(func() error { return numValidator.ValidateRange(age, 18, 65) }).
    Add(func() error { return strValidator.ValidateMinLength(email, 5) })
if err := customVal.Validate(); err != nil {
    // handle error
}

// Collection validation
validator.ValidateSliceLength(items, 5)
validator.ValidateMapHasKey(config, "api_key")
```

## API Reference

### Factory Functions

- `Numbers[T]() *NumberValidator[T]` - Create a number validator for type T
- `Strings[T]() *StringValidator[T]` - Create a string validator for type T (includes Parse field for parsing validation)
- `Parse() *ParseValidator` - Create a standalone parsing validator (advanced usage)

### Number Validators

All number validators work with `Number` types: `int`, `uint`, `float32`, `float64` and their variants.

#### Range Validation
- `ValidateMin(input T, min T) error` - Value ≥ minimum
- `ValidateMax(input T, max T) error` - Value ≤ maximum
- `ValidateRange(input T, min, max T) error` - Value within range

#### Sign Validation
- `ValidatePositive(input T) error` - Value > 0
- `ValidateNegative(input T) error` - Value < 0
- `ValidateNonNegative(input T) error` - Value ≥ 0
- `ValidateNonPositive(input T) error` - Value ≤ 0

#### Equality Validation
- `ValidateZero(input T) error` - Value == 0
- `ValidateNonZero(input T) error` - Value != 0
- `ValidateEqual(input T, expected T) error` - Value equals expected
- `ValidateNotEqual(input T, notExpected T) error` - Value not equal

#### Comparison Validation
- `ValidateGreaterThan(input T, threshold T) error` - Value > threshold
- `ValidateGreaterThanOrEqual(input T, threshold T) error` - Value ≥ threshold
- `ValidateLessThan(input T, threshold T) error` - Value < threshold
- `ValidateLessThanOrEqual(input T, threshold T) error` - Value ≤ threshold
- `ValidateBetween(input T, lower, upper T) error` - Value between bounds
- `ValidateConsecutive(input1, input2 T) error` - input2 == input1 + 1

#### Integer-Only Validation
- `ValidateEven(input T) error` - Value is even
- `ValidateOdd(input T) error` - Value is odd
- `ValidateDivisibleBy(input T, divisor T) error` - Evenly divisible
- `ValidatePowerOf(input T, base T) error` - Power of base
- `ValidateFibonacci(input T) error` - Is Fibonacci number

### String Validators

String validators work with `StringLike` types: `string`, `[]byte`, `[]rune`.

#### Length Validation
- `ValidateMinLength(input T, min int) error` - Length ≥ minimum
- `ValidateMaxLength(input T, max int) error` - Length ≤ maximum
- `ValidateLengthRange(input T, min, max int) error` - Length within range

#### Character Validation
- `ValidatePattern(input T, pattern string) error` - Matches regex pattern
- `ValidateAlphanumeric(input T) error` - Only alphanumeric
- `ValidateAlpha(input T) error` - Only alphabetic
- `ValidateNumeric(input T) error` - Only numeric
- `ValidateSlug(input T) error` - Valid URL slug
- `ValidateLowercase(input T) error` - Only lowercase
- `ValidateUppercase(input T) error` - Only uppercase

#### Substring Validation
- `ValidateContains(input T, substring string) error` - Contains substring
- `ValidateNotContains(input T, substring string) error` - Doesn't contain substring
- `ValidatePrefix(input T, prefix string) error` - Starts with prefix
- `ValidateSuffix(input T, suffix string) error` - Ends with suffix

### Parse Validators

Parsing validators for string format validation.

#### Format Validation
- `ValidateEmail(input string) error` - Valid email format
- `ValidatePassword(input string) error` - 8+ chars, mixed case, digits
- `ValidateUUID(input string) error` - Valid UUID format
- `ValidateURL(input string) error` - Valid URL format
- `ValidateDate(input string, format string) error` - Valid date
- `ValidateDuration(input string) error` - Valid duration (e.g., "5m", "2h")
- `ValidatePhoneNumber(input string) error` - Valid phone format

#### Type Parsing
- `ValidateBool(input string) error` - Parseable as bool
- `ValidateInt(input string) error` - Parseable as int64
- `ValidateUint(input string) error` - Parseable as uint64
- `ValidateFloat(input string) error` - Parseable as float64
- `ValidatePositiveInt(input string) error` - Parseable positive int
- `ValidateNonNegativeInt(input string) error` - Parseable non-negative int
- `ValidatePositiveFloat(input string) error` - Parseable positive float

#### Network Validation
- `ValidateIPAddress(input string) error` - Valid IPv4 or IPv6 address
- `ValidateIPv4(input string) error` - Valid IPv4 address
- `ValidateIPv6(input string) error` - Valid IPv6 address

### Composite Validators

Composite validators allow combining and chaining multiple validation functions.

#### Generic Validators
- `OneOf[T comparable](input T, allowed ...T) error` - Value in allowed set
- `All(validators ...func() error) error` - All pass (AND logic)
- `Any(validators ...func() error) error` - At least one passes (OR logic)
- `ValidateMatchesField[T comparable](value1, value2 T, fieldName string) error` - Two values match (e.g., password confirmation)

### Collection Validators

Validators for slices, arrays, and maps.

#### Slice Validators
- `ValidateSliceLength[T any](input []T, length int) error` - Length equals specified value
- `ValidateSliceMinLength[T any](input []T, min int) error` - Length ≥ minimum
- `ValidateSliceMaxLength[T any](input []T, max int) error` - Length ≤ maximum
- `ValidateSliceLengthRange[T any](input []T, min, max int) error` - Length within range
- `ValidateSliceContains[T comparable](input []T, element T) error` - Contains element
- `ValidateSliceNotContains[T comparable](input []T, element T) error` - Doesn't contain element
- `ValidateSliceUnique[T comparable](input []T) error` - All elements unique

#### Map Validators
- `ValidateMapHasKey[K comparable, V any](input map[K]V, key K) error` - Contains key
- `ValidateMapNotHasKey[K comparable, V any](input map[K]V, key K) error` - Doesn't contain key
- `ValidateMapMinLength[K comparable, V any](input map[K]V, min int) error` - Entries ≥ minimum
- `ValidateMapMaxLength[K comparable, V any](input map[K]V, max int) error` - Entries ≤ maximum

### Utility Functions

- `ValidateNonEmpty[T](input T) error` - Generic emptiness check for strings, bytes, runes, maps, and slices
- `NewCustomValidator() *CustomValidator` - Create a custom validator with fluent chaining
- `Parse() *ParseValidator` - Standalone parsing validator (typically accessed via StringValidator.Parse)

### Custom Validator Builder

The `CustomValidator` type provides a fluent interface for composing validators:

```go
cv := validator.NewCustomValidator().
    Add(func() error { return numVal.ValidateRange(age, 18, 65) }).
    Add(func() error { return strVal.ValidateMinLength(email, 5) }).
    Add(func() error { return strVal.ValidatePattern(phone, `^\+?[0-9\s\-\(\)]+$`) })

if err := cv.Validate(); err != nil {
    // All validators passed if no error
}

// For critical code paths, use SafeValidate() to recover from panics
if err := cv.SafeValidate(); err != nil {
    // Handle validation error or panic recovery
}
```

- `NewCustomValidator() *CustomValidator` - Create a custom validator with fluent chaining
- `Add(validator func() error) *CustomValidator` - Append a validator function
- `Validate() error` - Run all validators in sequence (stops at first error)
- `SafeValidate() error` - Run all validators with panic recovery

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

### Safe Validation with Panic Recovery

```go
// Use SafeValidate() to prevent validation panics from crashing the application
cv := validator.NewCustomValidator().
    Add(func() error { return numVal.ValidateRange(age, 18, 65) }).
    Add(func() error { return strVal.ValidateMinLength(email, 5) })

// SafeValidate recovers from panics and returns error instead
if err := cv.SafeValidate(); err != nil {
    // Handle validation error (may be from panic recovery)
    log.Printf("Validation failed: %v", err)
}
```

## Performance

### Efficiency

The validator package is designed for performance:

- **Zero Allocations** for most validators (numbers, basic strings)
- **Lazy Validation** - stops at first error in composite validators
- **Minimal Overhead** - generic validators compile to native code
- **Type-Safe** - no runtime type assertions for generics

### Benchmarks

Typical performance characteristics (on modern hardware):

- Number validation (range): ~10-50 ns/op
- String validation (length): ~20-100 ns/op  
- Pattern validation (regex): ~200-1000 ns/op
- Composite validation (3 validators): ~30-150 ns/op

### Performance Tips

1. **Prefer simple validators** over regex for basic patterns
2. **Order validators** by cost - cheap checks before expensive ones
3. **Reuse validators** - create once, use many times
4. **Use composite validators** rather than nested if statements
5. **Profile first** - use benchmarks to identify bottlenecks

Example optimization:

```go
// Good: cheap validation first
cv := validator.NewCustomValidator().
    Add(func() error { return strVal.ValidateMinLength(email, 5) }).   // ~30 ns
    Add(func() error { return strVal.Parse.ValidateEmail(email) })     // ~500 ns

// Less ideal: expensive validation first  
cv := validator.NewCustomValidator().
    Add(func() error { return strVal.Parse.ValidateEmail(email) }).    // ~500 ns
    Add(func() error { return strVal.ValidateMinLength(email, 5) })    // ~30 ns (never reached if email invalid)
```

## Testing

```bash
go test -v ./validator
go test -cover ./validator
```
