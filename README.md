# go-utils 

A collection of reusable Go utilities and helper functions designed to simplify common programming tasks.

## Available Packages

| Package                    | Description                                                                |
|----------------------------|----------------------------------------------------------------------------|
| `slices`                   | Generic slice utility functions for conditional selection, set operations, and element manipulation |
| `helpers`                  | General utility functions including slice operations, conditional helpers, file system utilities, and struct manipulation |
| `httputil/auth`            | JWT token creation, validation, and management with role-based access control and custom claims support |
| `httputil/middleware`      | Common, reusable HTTP middleware for logging, recovery, CORS, and request ID injection |
| `httputil/request`         | HTTP request parsing utilities for JSON, form data, query parameters, and URL values |
| `httputil/response`        | Structured HTTP response handling with extensible encoders, hooks, and status code helpers |

## Package Details

### `slices`

The `slices` package provides generic, type-safe utility functions for slice operations commonly needed across Go projects. It consolidates functions that are frequently duplicated, offering clean, idiomatic implementations.

**Features:**
- **Generic ternary operator**: `If[T any](cond bool, vtrue T, vfalse T) T`
- **Set difference**: `Difference(a []string, b []string) []string`
- **Element deletion**: `DeleteElement[T any](slice []T, index int) []T`
- **Subset checking**: `ContainsAll[T comparable](mainSlice, subset []T) bool`

**Example Usage:**
```go
import "github.com/julianstephens/go-utils/slices"

// Conditional selection (ternary-like)
message := slices.If(isLoggedIn, "Welcome!", "Please log in")

// Find elements in a but not in b
diff := slices.Difference([]string{"a", "b", "c"}, []string{"b"})
// Returns: []string{"a", "c"}

// Remove element at index
items := slices.DeleteElement([]int{1, 2, 3, 4}, 1)
// Returns: []int{1, 3, 4}

// Check if subset is contained in main slice
hasAll := slices.ContainsAll([]int{1, 2, 3, 4}, []int{2, 3})
// Returns: true
```
