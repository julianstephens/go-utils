package generic_test

import (
	"testing"

	"github.com/julianstephens/go-utils/generic"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestZero(t *testing.T) {
	// Test with various types
	tst.AssertDeepEqual(t, generic.Zero[int](), 0)

	tst.AssertDeepEqual(t, generic.Zero[string](), "")

	tst.AssertDeepEqual(t, generic.Zero[bool](), false)

	tst.AssertNil(t, generic.Zero[*int](), "Zero[*int]() should be nil")

	// Test with slice
	tst.AssertNil(t, generic.Zero[[]int](), "Zero[[]int]() should be nil")

	// Test with custom struct
	type Person struct {
		Name string
		Age  int
	}
	expected := Person{}
	tst.AssertDeepEqual(t, generic.Zero[Person](), expected)
}

func TestIf(t *testing.T) {
	// Test with integers
	result := generic.If(true, 10, 20)
	tst.AssertDeepEqual(t, result, 10)

	result = generic.If(false, 10, 20)
	tst.AssertDeepEqual(t, result, 20)

	// Test with strings
	strResult := generic.If(true, "yes", "no")
	tst.AssertDeepEqual(t, strResult, "yes")

	strResult = generic.If(false, "yes", "no")
	tst.AssertDeepEqual(t, strResult, "no")

	// Test with custom types
	type Person struct {
		Name string
	}
	person1 := Person{Name: "Alice"}
	person2 := Person{Name: "Bob"}

	personResult := generic.If(true, person1, person2)
	tst.AssertDeepEqual(t, personResult, person1)
}

func TestDefault(t *testing.T) {
	// Test with strings
	result := generic.Default("", "default")
	tst.AssertDeepEqual(t, result, "default")

	result = generic.Default("value", "default")
	tst.AssertDeepEqual(t, result, "value")

	// Test with integers
	intResult := generic.Default(0, 42)
	tst.AssertDeepEqual(t, intResult, 42)

	intResult = generic.Default(10, 42)
	tst.AssertDeepEqual(t, intResult, 10)

	// Test with slices
	var nilSlice []int
	sliceResult := generic.Default(nilSlice, []int{1, 2, 3})
	expected := []int{1, 2, 3}
	tst.AssertDeepEqual(t, sliceResult, expected)
}

func TestPtr(t *testing.T) {
	// Test with integer
	value := 42
	ptr := generic.Ptr(value)
	tst.AssertNotNil(t, ptr, "Ptr returned nil")
	tst.AssertDeepEqual(t, *ptr, value)

	// Test with string
	str := "hello"
	strPtr := generic.Ptr(str)
	tst.AssertNotNil(t, strPtr, "Ptr returned nil for string")
	tst.AssertDeepEqual(t, *strPtr, str)

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}
	person := Person{Name: "Alice", Age: 30}
	personPtr := generic.Ptr(person)
	tst.AssertNotNil(t, personPtr, "Ptr returned nil for struct")
	tst.AssertDeepEqual(t, *personPtr, person)
}

func TestDeref(t *testing.T) {
	// Test with valid pointer
	value := 42
	ptr := &value
	result := generic.Deref(ptr)
	tst.AssertDeepEqual(t, result, value)

	// Test with nil pointer
	var nilPtr *int
	result = generic.Deref(nilPtr)
	tst.AssertDeepEqual(t, result, 0)

	// Test with string pointer
	str := "hello"
	strPtr := &str
	strResult := generic.Deref(strPtr)
	tst.AssertDeepEqual(t, strResult, str)

	// Test with nil string pointer
	var nilStrPtr *string
	strResult = generic.Deref(nilStrPtr)
	tst.AssertDeepEqual(t, strResult, "")
}
