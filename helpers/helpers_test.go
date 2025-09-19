package helpers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/julianstephens/go-utils/helpers"
)

func TestContainsAll(t *testing.T) {
	tests := []struct {
		name      string
		mainSlice []int
		subset    []int
		expected  bool
	}{
		{
			name:      "empty subset should return true",
			mainSlice: []int{1, 2, 3},
			subset:    []int{},
			expected:  true,
		},
		{
			name:      "subset fully contained",
			mainSlice: []int{1, 2, 3, 4, 5},
			subset:    []int{2, 4},
			expected:  true,
		},
		{
			name:      "subset not fully contained",
			mainSlice: []int{1, 2, 3},
			subset:    []int{2, 4},
			expected:  false,
		},
		{
			name:      "identical slices",
			mainSlice: []int{1, 2, 3},
			subset:    []int{1, 2, 3},
			expected:  true,
		},
		{
			name:      "empty main slice with non-empty subset",
			mainSlice: []int{},
			subset:    []int{1},
			expected:  false,
		},
		{
			name:      "both slices empty",
			mainSlice: []int{},
			subset:    []int{},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.ContainsAll(tt.mainSlice, tt.subset)
			if result != tt.expected {
				t.Errorf("ContainsAll(%v, %v) = %v; expected %v", tt.mainSlice, tt.subset, result, tt.expected)
			}
		})
	}
}

func TestContainsAllString(t *testing.T) {
	mainSlice := []string{"apple", "banana", "cherry"}
	subset := []string{"apple", "cherry"}

	result := helpers.ContainsAll(mainSlice, subset)
	if !result {
		t.Errorf("ContainsAll with strings failed: expected true, got %v", result)
	}

	subset = []string{"apple", "grape"}
	result = helpers.ContainsAll(mainSlice, subset)
	if result {
		t.Errorf("ContainsAll with strings failed: expected false, got %v", result)
	}
}

func TestIf(t *testing.T) {
	// Test with integers
	result := helpers.If(true, 10, 20)
	if result != 10 {
		t.Errorf("If(true, 10, 20) = %v; expected 10", result)
	}

	result = helpers.If(false, 10, 20)
	if result != 20 {
		t.Errorf("If(false, 10, 20) = %v; expected 20", result)
	}

	// Test with strings
	strResult := helpers.If(true, "yes", "no")
	if strResult != "yes" {
		t.Errorf("If(true, 'yes', 'no') = %v; expected 'yes'", strResult)
	}

	strResult = helpers.If(false, "yes", "no")
	if strResult != "no" {
		t.Errorf("If(false, 'yes', 'no') = %v; expected 'no'", strResult)
	}

	// Test with custom types
	type Person struct {
		Name string
	}
	person1 := Person{Name: "Alice"}
	person2 := Person{Name: "Bob"}

	personResult := helpers.If(true, person1, person2)
	if personResult != person1 {
		t.Errorf("If with custom type failed: expected %v, got %v", person1, personResult)
	}
}

func TestDefault(t *testing.T) {
	// Test with strings
	result := helpers.Default("", "default")
	if result != "default" {
		t.Errorf("Default('', 'default') = %v; expected 'default'", result)
	}

	result = helpers.Default("value", "default")
	if result != "value" {
		t.Errorf("Default('value', 'default') = %v; expected 'value'", result)
	}

	// Test with integers
	intResult := helpers.Default(0, 42)
	if intResult != 42 {
		t.Errorf("Default(0, 42) = %v; expected 42", intResult)
	}

	intResult = helpers.Default(10, 42)
	if intResult != 10 {
		t.Errorf("Default(10, 42) = %v; expected 10", intResult)
	}

	// Test with slices
	var nilSlice []int
	sliceResult := helpers.Default(nilSlice, []int{1, 2, 3})
	if !reflect.DeepEqual(sliceResult, []int{1, 2, 3}) {
		t.Errorf("Default with slice failed: expected [1 2 3], got %v", sliceResult)
	}

	nonNilSlice := []int{4, 5, 6}
	sliceResult = helpers.Default(nonNilSlice, []int{1, 2, 3})
	if !reflect.DeepEqual(sliceResult, []int{4, 5, 6}) {
		t.Errorf("Default with slice failed: expected [4 5 6], got %v", sliceResult)
	}

	nilResult := helpers.Default(nil, []int{1, 2, 3})
	if !reflect.DeepEqual(nilResult, []int{1, 2, 3}) {
		t.Errorf("Default with nil slice failed: expected [1 2 3], got %v", nilResult)
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected []string
	}{
		{
			name:     "basic difference",
			a:        []string{"a", "b", "c", "d"},
			b:        []string{"b", "d"},
			expected: []string{"a", "c"},
		},
		{
			name:     "no common elements",
			a:        []string{"a", "b", "c"},
			b:        []string{"x", "y", "z"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all elements in b",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: []string{},
		},
		{
			name:     "empty a slice",
			a:        []string{},
			b:        []string{"a", "b"},
			expected: []string{},
		},
		{
			name:     "empty b slice",
			a:        []string{"a", "b", "c"},
			b:        []string{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "both empty",
			a:        []string{},
			b:        []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.Difference(tt.a, tt.b)
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Difference(%v, %v) = %v; expected %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDeleteElement(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		index    int
		expected []int
	}{
		{
			name:     "delete first element",
			slice:    []int{1, 2, 3, 4},
			index:    0,
			expected: []int{2, 3, 4},
		},
		{
			name:     "delete middle element",
			slice:    []int{1, 2, 3, 4},
			index:    2,
			expected: []int{1, 2, 4},
		},
		{
			name:     "delete last element",
			slice:    []int{1, 2, 3, 4},
			index:    3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "single element slice",
			slice:    []int{42},
			index:    0,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.DeleteElement(tt.slice, tt.index)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeleteElement(%v, %d) = %v; expected %v", tt.slice, tt.index, result, tt.expected)
			}
		})
	}
}

func TestDeleteElementString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	result := helpers.DeleteElement(slice, 1)
	expected := []string{"apple", "cherry"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteElement with strings failed: expected %v, got %v", expected, result)
	}
}

func TestEnsure(t *testing.T) {
	// Test directory creation
	t.Run("create directory", func(t *testing.T) {
		tempDir := filepath.Join(os.TempDir(), "test_ensure_dir")
		defer os.RemoveAll(tempDir)

		err := helpers.Ensure(tempDir, true)
		if err != nil {
			t.Errorf("Ensure directory failed: %v", err)
		}

		// Check if directory exists
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			t.Error("Directory was not created")
		}

		// Call again to ensure it doesn't fail when directory exists
		err = helpers.Ensure(tempDir, true)
		if err != nil {
			t.Errorf("Ensure existing directory failed: %v", err)
		}
	})

	// Test file creation
	t.Run("create file", func(t *testing.T) {
		tempFile := filepath.Join(os.TempDir(), "test_ensure_file.txt")
		defer os.Remove(tempFile)

		err := helpers.Ensure(tempFile, false)
		if err != nil {
			t.Errorf("Ensure file failed: %v", err)
		}

		// Check if file exists
		if _, err := os.Stat(tempFile); os.IsNotExist(err) {
			t.Error("File was not created")
		}

		// Call again to ensure it doesn't fail when file exists
		err = helpers.Ensure(tempFile, false)
		if err != nil {
			t.Errorf("Ensure existing file failed: %v", err)
		}
	})

	// Test nested directory creation
	t.Run("create nested directory", func(t *testing.T) {
		tempDir := filepath.Join(os.TempDir(), "test_ensure", "nested", "dir")
		defer os.RemoveAll(filepath.Join(os.TempDir(), "test_ensure"))

		err := helpers.Ensure(tempDir, true)
		if err != nil {
			t.Errorf("Ensure nested directory failed: %v", err)
		}

		// Check if directory exists
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			t.Error("Nested directory was not created")
		}
	})
}

func TestStringPtr(t *testing.T) {
	testString := "hello world"
	ptr := helpers.StringPtr(testString)

	if ptr == nil {
		t.Error("StringPtr returned nil")
	} else {
		if *ptr != testString {
			t.Errorf("StringPtr failed: expected %q, got %q", testString, *ptr)
		}
	}

	// Test that it returns a new pointer each time
	ptr2 := helpers.StringPtr(testString)
	if ptr == ptr2 {
		t.Error("StringPtr should return different pointers for different calls")
	}
}

func TestMustMarshalJson(t *testing.T) {
	// Test successful marshaling
	t.Run("successful marshal", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "test",
			"age":  30,
		}

		result := helpers.MustMarshalJson(data)
		if len(result) == 0 {
			t.Error("MustMarshalJson returned empty result")
		}

		// Verify it's valid JSON by unmarshaling
		var unmarshaled map[string]interface{}
		err := json.Unmarshal(result, &unmarshaled)
		if err != nil {
			t.Errorf("MustMarshalJson produced invalid JSON: %v", err)
		}

		if unmarshaled["name"] != "test" || unmarshaled["age"] != float64(30) {
			t.Errorf("MustMarshalJson produced incorrect JSON: %v", unmarshaled)
		}
	})

	// Test panic with unmarshalable data
	t.Run("panic on unmarshalable data", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustMarshalJson should panic on unmarshalable data")
			}
		}()

		// Channel cannot be marshaled to JSON
		unmarshalable := make(chan int)
		helpers.MustMarshalJson(unmarshalable)
	})
}

func TestStructToMap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	// Test with struct
	t.Run("struct conversion", func(t *testing.T) {
		person := Person{
			Name: "Alice",
			Age:  30,
			City: "New York",
		}

		result := helpers.StructToMap(person)
		if result == nil {
			t.Error("StructToMap returned nil")
		}

		expected := map[string]any{
			"Name": "Alice",
			"Age":  30,
			"City": "New York",
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("StructToMap failed: expected %v, got %v", expected, result)
		}
	})

	// Test with pointer to struct
	t.Run("pointer to struct conversion", func(t *testing.T) {
		person := &Person{
			Name: "Bob",
			Age:  25,
			City: "Boston",
		}

		result := helpers.StructToMap(person)
		if result == nil {
			t.Error("StructToMap returned nil for pointer")
		}

		expected := map[string]any{
			"Name": "Bob",
			"Age":  25,
			"City": "Boston",
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("StructToMap with pointer failed: expected %v, got %v", expected, result)
		}
	})

	// Test with non-struct
	t.Run("non-struct returns nil", func(t *testing.T) {
		result := helpers.StructToMap("not a struct")
		if result != nil {
			t.Errorf("StructToMap should return nil for non-struct, got %v", result)
		}

		result = helpers.StructToMap(42)
		if result != nil {
			t.Errorf("StructToMap should return nil for non-struct, got %v", result)
		}

		result = helpers.StructToMap([]int{1, 2, 3})
		if result != nil {
			t.Errorf("StructToMap should return nil for non-struct, got %v", result)
		}
	})

	// Test with empty struct
	t.Run("empty struct", func(t *testing.T) {
		type Empty struct{}
		empty := Empty{}

		result := helpers.StructToMap(empty)
		if result == nil {
			t.Error("StructToMap returned nil for empty struct")
		}

		if len(result) != 0 {
			t.Errorf("StructToMap should return empty map for empty struct, got %v", result)
		}
	})
}
