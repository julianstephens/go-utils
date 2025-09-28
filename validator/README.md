# Validator package

Common input validators extracted from `cliutil`.

## Functions

- `ValidateNonEmpty(input string) error` - returns an error if input is empty or whitespace.
- `ValidateEmail(input string) error` - basic email format validation; returns an error on invalid format.
- `ValidatePassword(input string) error` - basic password format validation; returns an error on invalid format.
- `ValidateUUID(input string) error` - UUID string validation; returns an error on invalid format.

## Usage

```go
import "github.com/julianstephens/go-utils/validator"

if err := validator.ValidateEmail("user@example.com"); err != nil {
    // handle invalid email
}
```
