package validator

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

// ValidationFunc is a function type for input validation
type ValidationFunc func(string) error

// ValidateNonEmpty validates that input is not empty
func ValidateNonEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

// ValidateEmail validates basic email format
func ValidateEmail(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return err
	}
	_, err := mail.ParseAddress(input)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}
	// Require domain to contain a dot (simple TLD check) to reject addresses like a@b
	parts := strings.Split(input, "@")
	if len(parts) != 2 || !strings.Contains(parts[1], ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword validates that a password meets basic criteria
// (at least 8 characters, contains upper and lower case letters, and a digit)
func ValidatePassword(input string) error {
	if len(input) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
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
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	return nil
}

func ValidateUUID(input string) error {
	if err := ValidateNonEmpty(input); err != nil {
		return err
	}
	_, err := uuid.Parse(input)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	return nil
}
