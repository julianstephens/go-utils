package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/julianstephens/go-utils/generic"
)

func main() {
	fmt.Println("=== Generic Package Examples ===\n")

	// Functional Programming Examples
	fmt.Println("1. Functional Programming:")
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Map: Convert integers to strings
	stringNumbers := generic.Map(numbers, func(x int) string {
		return strconv.Itoa(x)
	})
	fmt.Printf("   Map (int -> string): %v\n", stringNumbers)

	// Filter: Get even numbers
	evens := generic.Filter(numbers, func(x int) bool {
		return x%2 == 0
	})
	fmt.Printf("   Filter (evens): %v\n", evens)

	// Reduce: Sum all numbers
	sum := generic.Reduce(numbers, 0, func(acc, x int) int {
		return acc + x
	})
	fmt.Printf("   Reduce (sum): %d\n", sum)

	// Find: First number greater than 5
	first, found := generic.Find(numbers, func(x int) bool {
		return x > 5
	})
	fmt.Printf("   Find (first > 5): %d (found: %t)\n", first, found)

	// Any: Check if any number is greater than 8
	hasLarge := generic.Any(numbers, func(x int) bool {
		return x > 8
	})
	fmt.Printf("   Any (> 8): %t\n", hasLarge)

	// All: Check if all numbers are positive
	allPositive := generic.All(numbers, func(x int) bool {
		return x > 0
	})
	fmt.Printf("   All (positive): %t\n", allPositive)

	fmt.Println()

	// Slice Operations Examples
	fmt.Println("2. Slice Operations:")
	duplicates := []string{"apple", "banana", "apple", "cherry", "banana", "date"}

	// Unique: Remove duplicates
	unique := generic.Unique(duplicates)
	fmt.Printf("   Unique: %v\n", unique)

	// Contains: Check membership
	hasApple := generic.Contains(duplicates, "apple")
	fmt.Printf("   Contains 'apple': %t\n", hasApple)

	// Reverse: Reverse order
	reversed := generic.Reverse(unique)
	fmt.Printf("   Reverse: %v\n", reversed)

	// Chunk: Split into chunks of 2
	chunks := generic.Chunk(numbers, 3)
	fmt.Printf("   Chunk (size 3): %v\n", chunks)

	// Set operations
	set1 := []int{1, 2, 3, 4, 5}
	set2 := []int{4, 5, 6, 7, 8}

	difference := generic.Difference(set1, set2)
	fmt.Printf("   Difference (set1 - set2): %v\n", difference)

	intersection := generic.Intersection(set1, set2)
	fmt.Printf("   Intersection: %v\n", intersection)

	union := generic.Union(set1, set2)
	fmt.Printf("   Union: %v\n", union)

	fmt.Println()

	// Map Operations Examples
	fmt.Println("3. Map Operations:")
	fruitColors := map[string]string{
		"apple":  "red",
		"banana": "yellow",
		"grape":  "purple",
		"orange": "orange",
	}

	// Keys and Values
	fruits := generic.Keys(fruitColors)
	colors := generic.Values(fruitColors)
	fmt.Printf("   Keys: %v\n", fruits)
	fmt.Printf("   Values: %v\n", colors)

	// Filter map: Only fruits with colors longer than 3 characters
	filtered := generic.FilterMap(fruitColors, func(fruit, color string) bool {
		return len(color) > 3
	})
	fmt.Printf("   Filtered (color len > 3): %v\n", filtered)

	// Map transformation: Uppercase fruit names
	uppercased := generic.MapMap(fruitColors, func(fruit, color string) (string, string) {
		return strings.ToUpper(fruit), color
	})
	fmt.Printf("   Transformed (uppercase keys): %v\n", uppercased)

	// Convert to slice of key-value pairs
	pairs := generic.MapToSlice(fruitColors, func(fruit, color string) string {
		return fruit + ":" + color
	})
	fmt.Printf("   Map to slice: %v\n", pairs)

	fmt.Println()

	// General Utilities Examples
	fmt.Println("4. General Utilities:")

	// If: Ternary operator
	condition := true
	result := generic.If(condition, "success", "failure")
	fmt.Printf("   If (ternary): %s\n", result)

	// Default: Provide fallback for zero values
	emptyString := ""
	defaulted := generic.Default(emptyString, "default value")
	fmt.Printf("   Default: '%s'\n", defaulted)

	// Ptr: Create pointer
	value := 42
	ptr := generic.Ptr(value)
	fmt.Printf("   Ptr: *%d = %d\n", ptr, *ptr)

	// Deref: Dereference pointer safely
	dereferenced := generic.Deref(ptr)
	fmt.Printf("   Deref: %d\n", dereferenced)

	// Deref with nil pointer
	var nilPtr *int
	safeDeref := generic.Deref(nilPtr)
	fmt.Printf("   Deref (nil): %d\n", safeDeref)

	// Zero: Get zero value
	zeroInt := generic.Zero[int]()
	zeroString := generic.Zero[string]()
	fmt.Printf("   Zero values: int=%d, string='%s'\n", zeroInt, zeroString)

	fmt.Println()

	// Complex Example: Processing a list of people
	fmt.Println("5. Complex Example - Processing People:")

	type Person struct {
		Name string
		Age  int
		City string
	}

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
	fmt.Printf("   Adults over 30: %v\n", adults)

	// Map: Extract names
	names := generic.Map(people, func(p Person) string {
		return p.Name
	})
	fmt.Printf("   Names: %v\n", names)

	// Group by city using SliceToMap
	cityMap := generic.SliceToMapBy(people, func(p Person) string {
		return p.City
	})
	fmt.Printf("   Last person by city: %v\n", cityMap)

	// Check if anyone is from Chicago
	hasChicago := generic.Any(people, func(p Person) bool {
		return p.City == "Chicago"
	})
	fmt.Printf("   Anyone from Chicago: %t\n", hasChicago)

	fmt.Println("\n=== End of Examples ===")
}
