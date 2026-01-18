# Generic Package

The `generic` package provides comprehensive generic utilities leveraging Go's type parameters for functional programming, slice operations, map utilities, and type-safe helpers.

## Features

- **Functional Programming**: Map, filter, reduce, find, and other operations
- **Slice Utilities**: Unique, reverse, chunk, set operations (union, intersection, difference)
- **Map Operations**: Keys, values, filtering, and transformation
- **General Utilities**: Conditional helpers and pointer utilities

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
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Map: Convert to strings
    strings := generic.Map(numbers, func(x int) string {
        return fmt.Sprintf("%d", x)
    })
    _ = strings

    // Filter: Get even numbers
    evens := generic.Filter(numbers, func(x int) bool {
        return x%2 == 0
    })
    _ = evens

    // Reduce: Sum all numbers
    sum := generic.Reduce(numbers, 0, func(acc, x int) int {
        return acc + x
    })
    _ = sum

    // Find: First number greater than 5
    first, found := generic.Find(numbers, func(x int) bool {
        return x > 5
    })
    _ = first
    _ = found
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
    duplicates := []string{"apple", "banana", "apple", "cherry", "banana"}
    unique := generic.Unique(duplicates)
    _ = unique

    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    chunks := generic.Chunk(numbers, 3)
    _ = chunks

    set1 := []int{1, 2, 3, 4, 5}
    set2 := []int{4, 5, 6, 7, 8}

    difference := generic.Difference(set1, set2)
    intersection := generic.Intersection(set1, set2)
    union := generic.Union(set1, set2)

    _ = difference
    _ = intersection
    _ = union
}
```

### Map Operations

```go
package main

import (
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    fruitColors := map[string]string{
        "apple":  "red",
        "banana": "yellow",
        "grape":  "purple",
        "orange": "orange",
    }

    keys := generic.Keys(fruitColors)
    values := generic.Values(fruitColors)
    _ = keys
    _ = values

    // Filter: Colors longer than 3 characters
    filtered := generic.FilterMap(fruitColors, func(_, color string) bool {
        return len(color) > 3
    })
    _ = filtered

    // Convert to slice of pairs
    pairs := generic.MapToSlice(fruitColors, func(fruit, color string) string {
        return fruit + ":" + color
    })
    _ = pairs
}
```

### General Utilities

```go
package main

import (
    "github.com/julianstephens/go-utils/generic"
)

func main() {
    // Provide fallback for zero values
    emptyString := ""
    defaulted := generic.Default(emptyString, "default")
    _ = defaulted

    // Create and dereference pointers
    value := 42
    ptr := generic.Ptr(value)
    dereferenced := generic.Deref(ptr)
    _ = dereferenced

    // Get zero value for type
    zeroInt := generic.Zero[int]()
    zeroString := generic.Zero[string]()
    _ = zeroInt
    _ = zeroString
}
```

### Complex Example: Processing People

```go
package main

import (
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
    _ = adults

    // Map: Extract names
    names := generic.Map(people, func(p Person) string {
        return p.Name
    })
    _ = names

    // Group by city
    cityMap := generic.SliceToMapBy(people, func(p Person) string {
        return p.City
    })
    _ = cityMap
}
```

## Available Functions

### Functional Programming
- `Map[T, U any](slice []T, f func(T) U) []U` - Apply function to each element
- `Filter[T any](slice []T, predicate func(T) bool) []T` - Filter by predicate
- `Reduce[T, U any](slice []T, initial U, f func(U, T) U) U` - Reduce to single value
- `Find[T any](slice []T, predicate func(T) bool) (T, bool)` - Find first match
- `Any[T any](slice []T, predicate func(T) bool) bool` - Check if any matches
- `All[T any](slice []T, predicate func(T) bool) bool` - Check if all match
- `ForEach[T any](slice []T, f func(T))` - Execute for each element

### Slice Operations
- `Contains[T comparable](slice []T, value T) bool` - Check if contains value
- `ContainsAll[T comparable](mainSlice, subset []T) bool` - Check if all subset elements present
- `IndexOf[T comparable](slice []T, value T) int` - Find index (-1 if not found)
- `Unique[T comparable](slice []T) []T` - Remove duplicates
- `Reverse[T any](slice []T) []T` - Reverse order
- `DeleteElement[T any](slice []T, index int) []T` - Remove element at index
- `InsertElement[T any](slice []T, index int, element T) []T` - Insert element at index
- `Chunk[T any](slice []T, size int) [][]T` - Split into chunks
- `Union[T comparable](slice1, slice2 []T) []T` - Union of slices
- `Intersection[T comparable](slice1, slice2 []T) []T` - Intersection of slices
- `Difference[T comparable](slice1, slice2 []T) []T` - Elements in first but not second

### Map Operations
- `Keys[K comparable, V any](m map[K]V) []K` - Extract all keys
- `Values[K comparable, V any](m map[K]V) []V` - Extract all values
- `HasKey[K comparable, V any](m map[K]V, key K) bool` - Check if key exists
- `HasValue[K comparable, V comparable](m map[K]V, value V) bool` - Check if value exists
- `FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V` - Filter entries
- `MapMap[K1 comparable, V1, K2 comparable, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2)) map[K2]V2` - Transform map
- `MapToSlice[K comparable, V, T any](m map[K]V, f func(K, V) T) []T` - Convert to slice
- `SliceToMap[T any, K comparable, V any](slice []T, keyFunc func(T) K, valueFunc func(T) V) map[K]V` - Convert to map
- `SliceToMapBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T` - Convert to map with elements as values
- `MergeMap[K comparable, V any](mergeMaps ...map[K]V) map[K]V` - Merge multiple maps
- `CopyMap[K comparable, V any](m map[K]V) map[K]V` - Shallow copy

### General Utilities
- `Default[T any](val T, defaultVal T) T` - Return default if zero value
- `Zero[T any]() T` - Get zero value for type
- `Ptr[T any](v T) *T` - Create pointer
- `Deref[T any](ptr *T) T` - Safely dereference

## Notes

- All functions are type-safe using Go generics
- Nil slices and maps are handled gracefully
- For ternary-like operations, use `helpers.If[T any]()` from the helpers package
- Functions follow idiomatic Go conventions
- No external dependencies except Go standard library