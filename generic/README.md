# Generic Package

The `generic` package provides comprehensive generic utilities leveraging Go's type parameters for functional programming, slice operations, map utilities, and type-safe helpers.

## Features

- **Functional Programming**: Map, filter, reduce, find, and other functional operations
- **Slice Utilities**: Unique, reverse, chunk, set operations (union, intersection, difference)
- **Map Operations**: Keys, values, filtering, and transformation functions
- **General Utilities**: Conditional helpers, pointer utilities, and type-safe operations

## Installation

```bash
go get github.com/julianstephens/go-utils/generic
```

## Usage

### Functional Programming

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Map: Convert integers to strings
    stringNumbers := generic.Map(numbers, func(x int) string {
        return strconv.Itoa(x)
    })
    fmt.Printf("Map (int -> string): %v\n", stringNumbers)

    // Filter: Get even numbers
    evens := generic.Filter(numbers, func(x int) bool {
        return x%2 == 0
    })
    fmt.Printf("Filter (evens): %v\n", evens)

    // Reduce: Sum all numbers
    sum := generic.Reduce(numbers, 0, func(acc, x int) int {
        return acc + x
    })
    fmt.Printf("Reduce (sum): %d\n", sum)

    // Find: First number greater than 5
    first, found := generic.Find(numbers, func(x int) bool {
        return x > 5
    })
    fmt.Printf("Find (first > 5): %d (found: %t)\n", first, found)

    // Any: Check if any number is greater than 8
    hasLarge := generic.Any(numbers, func(x int) bool {
        return x > 8
    })
    fmt.Printf("Any (> 8): %t\n", hasLarge)

    // All: Check if all numbers are positive
    allPositive := generic.All(numbers, func(x int) bool {
        return x > 0
    })
    fmt.Printf("All (positive): %t\n", allPositive)
}
```

### Slice Operations

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    duplicates := []string{"apple", "banana", "apple", "cherry", "banana", "date"}

    // Unique: Remove duplicates
    unique := generic.Unique(duplicates)
    fmt.Printf("Unique: %v\n", unique)

    // Contains: Check membership
    hasApple := generic.Contains(duplicates, "apple")
    fmt.Printf("Contains 'apple': %t\n", hasApple)

    // Reverse: Reverse order
    reversed := generic.Reverse(unique)
    fmt.Printf("Reverse: %v\n", reversed)

    // Chunk: Split into chunks
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    chunks := generic.Chunk(numbers, 3)
    fmt.Printf("Chunk (size 3): %v\n", chunks)

    // Set operations
    set1 := []int{1, 2, 3, 4, 5}
    set2 := []int{4, 5, 6, 7, 8}

    difference := generic.Difference(set1, set2)
    intersection := generic.Intersection(set1, set2)
    union := generic.Union(set1, set2)

    fmt.Printf("Difference (set1 - set2): %v\n", difference)
    fmt.Printf("Intersection: %v\n", intersection)
    fmt.Printf("Union: %v\n", union)
}
```

### Map Operations

```go
package main

import (
    "fmt"
    "strings"
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    fruitColors := map[string]string{
        "apple":  "red",
        "banana": "yellow",
        "grape":  "purple",
        "orange": "orange",
    }

    // Keys and Values
    fruits := generic.Keys(fruitColors)
    colors := generic.Values(fruitColors)
    fmt.Printf("Keys: %v\n", fruits)
    fmt.Printf("Values: %v\n", colors)

    // Filter map: Only fruits with colors longer than 3 characters
    filtered := generic.FilterMap(fruitColors, func(fruit, color string) bool {
        return len(color) > 3
    })
    fmt.Printf("Filtered (color len > 3): %v\n", filtered)

    // Map transformation: Uppercase fruit names
    uppercased := generic.MapMap(fruitColors, func(fruit, color string) (string, string) {
        return strings.ToUpper(fruit), color
    })
    fmt.Printf("Transformed (uppercase keys): %v\n", uppercased)

    // Convert to slice of key-value pairs
    pairs := generic.MapToSlice(fruitColors, func(fruit, color string) string {
        return fruit + ":" + color
    })
    fmt.Printf("Map to slice: %v\n", pairs)
}
```

### General Utilities

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    // If: Ternary operator
    condition := true
    result := generic.If(condition, "success", "failure")
    fmt.Printf("If (ternary): %s\n", result)

    // Default: Provide fallback for zero values
    emptyString := ""
    defaulted := generic.Default(emptyString, "default value")
    fmt.Printf("Default: '%s'\n", defaulted)

    // Ptr: Create pointer
    value := 42
    ptr := generic.Ptr(value)
    fmt.Printf("Ptr: *%d = %d\n", ptr, *ptr)

    // Deref: Dereference pointer safely
    dereferenced := generic.Deref(ptr)
    fmt.Printf("Deref: %d\n", dereferenced)

    // Deref with nil pointer (returns zero value)
    var nilPtr *int
    safeDeref := generic.Deref(nilPtr)
    fmt.Printf("Deref (nil): %d\n", safeDeref)

    // Zero: Get zero value
    zeroInt := generic.Zero[int]()
    zeroString := generic.Zero[string]()
    fmt.Printf("Zero values: int=%d, string='%s'\n", zeroInt, zeroString)
}
```

### Complex Example: Processing People

```go
package main

import (
    "fmt"
    "github.com/julianstephens/go-utils/generic"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    people := []Person{
        {Name: "Alice", Age: 30, City: "New York"},
        {Name: "Bob", Age: 25, City: "San Francisco"},
        {Name: "Charlie", Age: 35, City: "New York"},
        {Name: "Diana", Age: 28, City: "Chicago"},
        {Name: "Eve", Age: 32, City: "San Francisco"},
    }

    // Filter: Adults over 30
    adults := generic.Filter(people, func(p Person) bool {
        return p.Age > 30
    })
    fmt.Printf("Adults over 30: %v\n", adults)

    // Map: Extract names
    names := generic.Map(people, func(p Person) string {
        return p.Name
    })
    fmt.Printf("Names: %v\n", names)

    // Group by city using SliceToMapBy
    cityMap := generic.SliceToMapBy(people, func(p Person) string {
        return p.City
    })
    fmt.Printf("Last person by city: %v\n", cityMap)

    // Check if anyone is from Chicago
    hasChicago := generic.Any(people, func(p Person) bool {
        return p.City == "Chicago"
    })
    fmt.Printf("Anyone from Chicago: %t\n", hasChicago)
}
```

## Available Functions

### Functional Programming
- `Map[T, U any](slice []T, f func(T) U) []U` - Apply function to each element
- `Filter[T any](slice []T, predicate func(T) bool) []T` - Filter elements by predicate
- `Reduce[T, U any](slice []T, initial U, f func(U, T) U) U` - Reduce slice to single value
- `Find[T any](slice []T, predicate func(T) bool) (T, bool)` - Find first matching element
- `Any[T any](slice []T, predicate func(T) bool) bool` - Check if any element matches
- `All[T any](slice []T, predicate func(T) bool) bool` - Check if all elements match

### Slice Operations
- `Contains[T comparable](slice []T, value T) bool` - Check if slice contains value
- `Unique[T comparable](slice []T) []T` - Remove duplicate elements
- `Reverse[T any](slice []T) []T` - Reverse slice order
- `Chunk[T any](slice []T, size int) [][]T` - Split slice into chunks
- `Union[T comparable](slice1, slice2 []T) []T` - Union of two slices
- `Intersection[T comparable](slice1, slice2 []T) []T` - Intersection of two slices
- `Difference[T comparable](slice1, slice2 []T) []T` - Elements in slice1 but not slice2

### Map Operations
- `Keys[K comparable, V any](m map[K]V) []K` - Extract all keys
- `Values[K comparable, V any](m map[K]V) []V` - Extract all values
- `FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V` - Filter map entries
- `MapMap[K comparable, V, U any](m map[K]V, f func(K, V) (K, U)) map[K]U` - Transform map
- `MapToSlice[K comparable, V, U any](m map[K]V, f func(K, V) U) []U` - Convert map to slice

### General Utilities
- `If[T any](cond bool, vtrue T, vfalse T) T` - Ternary operator
- `Default[T any](val T, defaultVal T) T` - Return default if zero value
- `Zero[T any]() T` - Get zero value for type
- `Ptr[T any](v T) *T` - Create pointer to value
- `Deref[T any](ptr *T) T` - Safely dereference pointer

## Notes

- All functions are type-safe using Go generics
- Nil slices and maps are handled gracefully
- Functions follow idiomatic Go conventions
- No external dependencies except Go standard library