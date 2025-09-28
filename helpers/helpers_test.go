package helpers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/julianstephens/go-utils/helpers"
	tst "github.com/julianstephens/go-utils/tests"
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
			tst.AssertDeepEqual(t, result, tt.expected)
		})
	}
}

func TestContainsAllString(t *testing.T) {
	mainSlice := []string{"apple", "banana", "cherry"}
	subset := []string{"apple", "cherry"}

	result := helpers.ContainsAll(mainSlice, subset)
	tst.AssertTrue(t, result, "ContainsAll with strings should be true for existing subset")

	subset = []string{"apple", "grape"}
	result = helpers.ContainsAll(mainSlice, subset)
	tst.AssertFalse(t, result, "ContainsAll with strings should be false for missing element")
}

func TestIf(t *testing.T) {
	// Test with integers
	result := helpers.If(true, 10, 20)
	tst.AssertDeepEqual(t, result, 10)

	result = helpers.If(false, 10, 20)
	tst.AssertDeepEqual(t, result, 20)

	// Test with strings
	strResult := helpers.If(true, "yes", "no")
	tst.AssertDeepEqual(t, strResult, "yes")

	strResult = helpers.If(false, "yes", "no")
	tst.AssertDeepEqual(t, strResult, "no")

	// Test with custom types
	type Person struct {
		Name string
	}
	person1 := Person{Name: "Alice"}
	person2 := Person{Name: "Bob"}

	personResult := helpers.If(true, person1, person2)
	tst.AssertDeepEqual(t, personResult, person1)
}

func TestDefault(t *testing.T) {
	// Test with strings
	result := helpers.Default("", "default")
	tst.AssertDeepEqual(t, result, "default")

	result = helpers.Default("value", "default")
	tst.AssertDeepEqual(t, result, "value")

	// Test with integers
	intResult := helpers.Default(0, 42)
	tst.AssertDeepEqual(t, intResult, 42)

	intResult = helpers.Default(10, 42)
	tst.AssertDeepEqual(t, intResult, 10)

	// Test with slices
	var nilSlice []int
	sliceResult := helpers.Default(nilSlice, []int{1, 2, 3})
	tst.AssertDeepEqual(t, sliceResult, []int{1, 2, 3})

	nonNilSlice := []int{4, 5, 6}
	sliceResult = helpers.Default(nonNilSlice, []int{1, 2, 3})
	tst.AssertDeepEqual(t, sliceResult, []int{4, 5, 6})

	nilResult := helpers.Default(nil, []int{1, 2, 3})
	tst.AssertDeepEqual(t, nilResult, []int{1, 2, 3})
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
			// Treat nil and empty slices as equal for this test suite
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			tst.AssertDeepEqual(t, result, tt.expected)
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
			tst.AssertDeepEqual(t, result, tt.expected)
		})
	}
}

func TestDeleteElementString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	result := helpers.DeleteElement(slice, 1)
	expected := []string{"apple", "cherry"}

	tst.AssertDeepEqual(t, result, expected)
}

func TestExists(t *testing.T) {
	// Test with existing file
	t.Run("existing file", func(t *testing.T) {
		tempFile := filepath.Join(os.TempDir(), "test_exists_file.txt")
		defer os.Remove(tempFile)

		// Create the file
		file, err := os.Create(tempFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()

		// Test that it exists
		if !helpers.Exists(tempFile) {
			t.Error("Exists should return true for existing file")
		}
	})

	// Test with existing directory
	t.Run("existing directory", func(t *testing.T) {
		tempDir := filepath.Join(os.TempDir(), "test_exists_dir")
		defer os.RemoveAll(tempDir)

		// Create the directory
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Test that it exists
		if !helpers.Exists(tempDir) {
			t.Error("Exists should return true for existing directory")
		}
	})

	// Test with non-existing path
	t.Run("non-existing path", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "non_existent_path_12345")

		// Test that it doesn't exist
		if helpers.Exists(nonExistentPath) {
			t.Error("Exists should return false for non-existing path")
		}
	})

	// Test with empty string
	t.Run("empty string", func(t *testing.T) {
		if helpers.Exists("") {
			t.Error("Exists should return false for empty string")
		}
	})
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
		tst.AssertNotNil(t, result, "StructToMap returned nil")

		expected := map[string]any{
			"Name": "Alice",
			"Age":  30,
			"City": "New York",
		}

		tst.AssertDeepEqual(t, result, expected)
	})

	// Test with pointer to struct
	t.Run("pointer to struct conversion", func(t *testing.T) {
		person := &Person{
			Name: "Bob",
			Age:  25,
			City: "Boston",
		}

		result := helpers.StructToMap(person)
		tst.AssertNotNil(t, result, "StructToMap returned nil for pointer")

		expected := map[string]any{
			"Name": "Bob",
			"Age":  25,
			"City": "Boston",
		}

		tst.AssertDeepEqual(t, result, expected)
	})

	// Test with non-struct
	t.Run("non-struct returns nil", func(t *testing.T) {
		result := helpers.StructToMap("not a struct")
		tst.AssertNil(t, result, "StructToMap should return nil for non-struct")

		result = helpers.StructToMap(42)
		tst.AssertNil(t, result, "StructToMap should return nil for non-struct")

		result = helpers.StructToMap([]int{1, 2, 3})
		tst.AssertNil(t, result, "StructToMap should return nil for non-struct")
	})

	// Test with empty struct
	t.Run("empty struct", func(t *testing.T) {
		type Empty struct{}
		empty := Empty{}

		result := helpers.StructToMap(empty)
		tst.AssertNotNil(t, result, "StructToMap returned nil for empty struct")

		tst.AssertDeepEqual(t, len(result), 0)
	})
}
