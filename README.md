# go-utils

A collection of reusable Go utilities and helper functions designed to simplify common programming tasks.

## Available Packages

| Package               | Description                                                                                                                             |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| `generic`             | Comprehensive generic utilities leveraging Go's type parameters for functional programming, slice operations, map utilities, and type-safe helpers |
| `config`              | Reusable and idiomatic configuration management with support for environment variables, YAML/JSON files, validation, and default values |
| `logger`              | Unified structured logger wrapping logrus with log level control, custom formatting, and contextual logging support                     |
| `slices`              | Generic slice utility functions for conditional selection, set operations, and element manipulation                                     |
| `helpers`             | General utility functions including slice operations, conditional helpers, file system utilities, and struct manipulation               |
| `httputil/auth`       | JWT token creation, validation, and management with role-based access control and custom claims support                                 |
| `httputil/middleware` | Common, reusable HTTP middleware for logging, recovery, CORS, and request ID injection                                                  |
| `httputil/request`    | HTTP request parsing utilities for JSON, form data, query parameters, and URL values                                                    |
| `httputil/response`   | Structured HTTP response handling with extensible encoders, hooks, and status code helpers                                              |
## Generic Package

The `generic` package provides comprehensive type-safe utilities leveraging Go's generics (type parameters). It's designed to be the go-to package for functional programming patterns and common data structure operations.

### Features

#### Functional Programming Utilities
- **Map**: Transform elements from one type to another
- **Filter**: Select elements based on a predicate
- **Reduce**: Aggregate elements into a single value
- **Find**: Locate the first matching element
- **Any/All**: Boolean operations on collections
- **ForEach**: Execute side effects on elements

#### Slice Operations
- **Contains/ContainsAll**: Membership testing
- **Unique**: Remove duplicates
- **Reverse**: Reverse element order
- **Difference/Intersection/Union**: Set operations
- **Chunk**: Split slices into smaller pieces
- **IndexOf**: Find element positions

#### Map Utilities
- **Keys/Values**: Extract keys or values
- **MapToSlice/SliceToMap**: Convert between maps and slices
- **FilterMap/MapMap**: Transform and filter maps
- **MergeMap/CopyMap**: Map operations

#### General Utilities
- **If**: Type-safe ternary operator
- **Default**: Provide fallback values for zero values
- **Ptr/Deref**: Pointer utilities
- **Zero**: Get zero values for any type

### Usage Examples

```go
import "github.com/julianstephens/go-utils/generic"

// Functional programming
numbers := []int{1, 2, 3, 4, 5}
doubled := generic.Map(numbers, func(x int) int { return x * 2 })
evens := generic.Filter(numbers, func(x int) bool { return x%2 == 0 })
sum := generic.Reduce(numbers, 0, func(acc, x int) int { return acc + x })

// Slice operations
unique := generic.Unique([]int{1, 2, 2, 3, 3})
contains := generic.Contains(numbers, 3)
chunks := generic.Chunk(numbers, 2)

// Map operations
m := map[string]int{"a": 1, "b": 2, "c": 3}
keys := generic.Keys(m)
values := generic.Values(m)

// Utilities
result := generic.If(condition, "yes", "no")
ptr := generic.Ptr(42)
value := generic.Deref(ptr)
```

