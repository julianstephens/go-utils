package generic_test

import (
	"testing"

	"github.com/julianstephens/go-utils/generic"
)

func TestZero(t *testing.T) {
	// Test with various types
	if result := generic.Zero[int](); result != 0 {
		t.Errorf("Zero[int]() = %v; expected 0", result)
	}

	if result := generic.Zero[string](); result != "" {
		t.Errorf("Zero[string]() = %v; expected empty string", result)
	}

	if result := generic.Zero[bool](); result != false {
		t.Errorf("Zero[bool]() = %v; expected false", result)
	}

	if result := generic.Zero[*int](); result != nil {
		t.Errorf("Zero[*int]() = %v; expected nil", result)
	}

	// Test with slice
	if result := generic.Zero[[]int](); result != nil {
		t.Errorf("Zero[[]int]() = %v; expected nil", result)
	}

	// Test with custom struct
	type Person struct {
		Name string
		Age  int
	}
	expected := Person{}
	if result := generic.Zero[Person](); result != expected {
		t.Errorf("Zero[Person]() = %v; expected %v", result, expected)
	}
}

func TestIf(t *testing.T) {
	// Test with integers
	result := generic.If(true, 10, 20)
	if result != 10 {
		t.Errorf("If(true, 10, 20) = %v; expected 10", result)
	}

	result = generic.If(false, 10, 20)
	if result != 20 {
		t.Errorf("If(false, 10, 20) = %v; expected 20", result)
	}

	// Test with strings
	strResult := generic.If(true, "yes", "no")
	if strResult != "yes" {
		t.Errorf("If(true, 'yes', 'no') = %v; expected 'yes'", strResult)
	}

	strResult = generic.If(false, "yes", "no")
	if strResult != "no" {
		t.Errorf("If(false, 'yes', 'no') = %v; expected 'no'", strResult)
	}

	// Test with custom types
	type Person struct {
		Name string
	}
	person1 := Person{Name: "Alice"}
	person2 := Person{Name: "Bob"}

	personResult := generic.If(true, person1, person2)
	if personResult != person1 {
		t.Errorf("If with custom type failed: expected %v, got %v", person1, personResult)
	}
}

func TestDefault(t *testing.T) {
	// Test with strings
	result := generic.Default("", "default")
	if result != "default" {
		t.Errorf("Default('', 'default') = %v; expected 'default'", result)
	}

	result = generic.Default("value", "default")
	if result != "value" {
		t.Errorf("Default('value', 'default') = %v; expected 'value'", result)
	}

	// Test with integers
	intResult := generic.Default(0, 42)
	if intResult != 42 {
		t.Errorf("Default(0, 42) = %v; expected 42", intResult)
	}

	intResult = generic.Default(10, 42)
	if intResult != 10 {
		t.Errorf("Default(10, 42) = %v; expected 10", intResult)
	}

	// Test with slices
	var nilSlice []int
	sliceResult := generic.Default(nilSlice, []int{1, 2, 3})
	expected := []int{1, 2, 3}
	if len(sliceResult) != len(expected) {
		t.Errorf("Default with nil slice failed: expected %v, got %v", expected, sliceResult)
	}
}

func TestPtr(t *testing.T) {
	// Test with integer
	value := 42
	ptr := generic.Ptr(value)
	if ptr == nil {
		t.Error("Ptr returned nil")
	}
	if *ptr != value {
		t.Errorf("Ptr failed: expected %v, got %v", value, *ptr)
	}

	// Test with string
	str := "hello"
	strPtr := generic.Ptr(str)
	if strPtr == nil {
		t.Error("Ptr returned nil for string")
	}
	if *strPtr != str {
		t.Errorf("Ptr failed for string: expected %v, got %v", str, *strPtr)
	}

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}
	person := Person{Name: "Alice", Age: 30}
	personPtr := generic.Ptr(person)
	if personPtr == nil {
		t.Error("Ptr returned nil for struct")
	}
	if *personPtr != person {
		t.Errorf("Ptr failed for struct: expected %v, got %v", person, *personPtr)
	}
}

func TestDeref(t *testing.T) {
	// Test with valid pointer
	value := 42
	ptr := &value
	result := generic.Deref(ptr)
	if result != value {
		t.Errorf("Deref failed: expected %v, got %v", value, result)
	}

	// Test with nil pointer
	var nilPtr *int
	result = generic.Deref(nilPtr)
	if result != 0 {
		t.Errorf("Deref with nil pointer failed: expected 0, got %v", result)
	}

	// Test with string pointer
	str := "hello"
	strPtr := &str
	strResult := generic.Deref(strPtr)
	if strResult != str {
		t.Errorf("Deref failed for string: expected %v, got %v", str, strResult)
	}

	// Test with nil string pointer
	var nilStrPtr *string
	strResult = generic.Deref(nilStrPtr)
	if strResult != "" {
		t.Errorf("Deref with nil string pointer failed: expected empty string, got %v", strResult)
	}
}
