package tests

import (
	"fmt"
	"strings"
)

// contains is a helper for substring search using strings.Contains.
func contains(s, substr string) bool {
	return substr == "" || strings.Contains(s, substr)
}

// Print helper to ease debugging in tests (keeps fmt import used).
func Print(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
