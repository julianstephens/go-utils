package generic_test

import (
	"reflect"
	"testing"

	"github.com/julianstephens/go-utils/generic"
)

func TestContains(t *testing.T) {
	// Test with integers
	slice := []int{1, 2, 3, 4, 5}
	if !generic.Contains(slice, 3) {
		t.Error("Contains should return true for existing element")
	}
	if generic.Contains(slice, 6) {
		t.Error("Contains should return false for non-existing element")
	}

	// Test with strings
	stringSlice := []string{"apple", "banana", "cherry"}
	if !generic.Contains(stringSlice, "banana") {
		t.Error("Contains should return true for existing string")
	}
	if generic.Contains(stringSlice, "grape") {
		t.Error("Contains should return false for non-existing string")
	}

	// Test with empty slice
	if generic.Contains([]int{}, 1) {
		t.Error("Contains should return false for empty slice")
	}
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generic.ContainsAll(tt.mainSlice, tt.subset)
			if result != tt.expected {
				t.Errorf("ContainsAll(%v, %v) = %v; expected %v",
					tt.mainSlice, tt.subset, result, tt.expected)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	// Test with integers
	slice := []int{10, 20, 30, 40, 50}
	if index := generic.IndexOf(slice, 30); index != 2 {
		t.Errorf("IndexOf failed: expected 2, got %v", index)
	}
	if index := generic.IndexOf(slice, 60); index != -1 {
		t.Errorf("IndexOf should return -1 for non-existing element, got %v", index)
	}

	// Test with strings
	stringSlice := []string{"a", "b", "c", "d"}
	if index := generic.IndexOf(stringSlice, "c"); index != 2 {
		t.Errorf("IndexOf with strings failed: expected 2, got %v", index)
	}
}

func TestUnique(t *testing.T) {
	// Test with duplicates
	input := []int{1, 2, 2, 3, 3, 3, 4}
	result := generic.Unique(input)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unique failed: expected %v, got %v", expected, result)
	}

	// Test with no duplicates
	noDupsInput := []int{1, 2, 3, 4}
	noDupsResult := generic.Unique(noDupsInput)
	if !reflect.DeepEqual(noDupsResult, noDupsInput) {
		t.Errorf("Unique with no duplicates failed: expected %v, got %v", noDupsInput, noDupsResult)
	}

	// Test with nil slice
	if result := generic.Unique[int](nil); result != nil {
		t.Errorf("Unique with nil slice should return nil, got %v", result)
	}

	// Test with strings
	stringInput := []string{"apple", "banana", "apple", "cherry", "banana"}
	stringResult := generic.Unique(stringInput)
	expectedStrings := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(stringResult, expectedStrings) {
		t.Errorf("Unique with strings failed: expected %v, got %v", expectedStrings, stringResult)
	}
}

func TestReverse(t *testing.T) {
	// Test with integers
	input := []int{1, 2, 3, 4, 5}
	result := generic.Reverse(input)
	expected := []int{5, 4, 3, 2, 1}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Reverse failed: expected %v, got %v", expected, result)
	}

	// Test with strings
	stringInput := []string{"a", "b", "c"}
	stringResult := generic.Reverse(stringInput)
	expectedStrings := []string{"c", "b", "a"}
	if !reflect.DeepEqual(stringResult, expectedStrings) {
		t.Errorf("Reverse with strings failed: expected %v, got %v", expectedStrings, stringResult)
	}

	// Test with nil slice
	if result := generic.Reverse[int](nil); result != nil {
		t.Errorf("Reverse with nil slice should return nil, got %v", result)
	}

	// Test with empty slice
	emptyResult := generic.Reverse([]int{})
	if len(emptyResult) != 0 {
		t.Errorf("Reverse with empty slice should return empty slice, got %v", emptyResult)
	}
}

func TestDeleteElement(t *testing.T) {
	// Test normal deletion
	input := []int{1, 2, 3, 4, 5}
	result := generic.DeleteElement(input, 2)
	expected := []int{1, 2, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteElement failed: expected %v, got %v", expected, result)
	}

	// Test deletion at beginning
	input2 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input2, 0)
	expected = []int{2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteElement at beginning failed: expected %v, got %v", expected, result)
	}

	// Test deletion at end
	input3 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input3, len(input3)-1)
	expected = []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteElement at end failed: expected %v, got %v", expected, result)
	}

	// Test out of bounds index
	input4 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input4, 10)
	if !reflect.DeepEqual(result, input4) {
		t.Errorf("DeleteElement with out of bounds index should return original slice")
	}

	// Test negative index
	input5 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input5, -1)
	if !reflect.DeepEqual(result, input5) {
		t.Errorf("DeleteElement with negative index should return original slice")
	}
}

func TestInsertElement(t *testing.T) {
	// Test normal insertion
	input := []int{1, 2, 4, 5}
	result := generic.InsertElement(input, 2, 3)
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InsertElement failed: expected %v, got %v", expected, result)
	}

	// Test insertion at beginning
	result = generic.InsertElement(input, 0, 0)
	expected = []int{0, 1, 2, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InsertElement at beginning failed: expected %v, got %v", expected, result)
	}

	// Test insertion out of bounds (should append)
	result = generic.InsertElement(input, 10, 6)
	expected = []int{1, 2, 4, 5, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InsertElement out of bounds should append: expected %v, got %v", expected, result)
	}
}

func TestDifference(t *testing.T) {
	// Test basic difference
	a := []int{1, 2, 3, 4, 5}
	b := []int{2, 4}
	result := generic.Difference(a, b)
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Difference failed: expected %v, got %v", expected, result)
	}

	// Test with strings
	stringA := []string{"apple", "banana", "cherry", "date"}
	stringB := []string{"banana", "date"}
	stringResult := generic.Difference(stringA, stringB)
	expectedStrings := []string{"apple", "cherry"}
	if !reflect.DeepEqual(stringResult, expectedStrings) {
		t.Errorf("Difference with strings failed: expected %v, got %v", expectedStrings, stringResult)
	}

	// Test with nil slice
	if result := generic.Difference[int](nil, []int{1, 2}); result != nil {
		t.Errorf("Difference with nil slice should return nil, got %v", result)
	}
}

func TestIntersection(t *testing.T) {
	// Test basic intersection
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	result := generic.Intersection(a, b)
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersection failed: expected %v, got %v", expected, result)
	}

	// Test with no common elements
	c := []int{8, 9, 10}
	result = generic.Intersection(a, c)
	if len(result) != 0 {
		t.Errorf("Intersection with no common elements should return empty slice, got %v", result)
	}

	// Test with duplicates
	d := []int{1, 1, 2, 2, 3}
	e := []int{1, 2, 2, 4}
	result = generic.Intersection(d, e)
	expected = []int{1, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersection with duplicates failed: expected %v, got %v", expected, result)
	}
}

func TestUnion(t *testing.T) {
	// Test basic union
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	result := generic.Union(a, b)
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Union failed: expected %v, got %v", expected, result)
	}

	// Test with duplicates within slices
	c := []int{1, 1, 2}
	d := []int{2, 2, 3}
	result = generic.Union(c, d)
	expected = []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Union with duplicates failed: expected %v, got %v", expected, result)
	}
}

func TestChunk(t *testing.T) {
	// Test normal chunking
	input := []int{1, 2, 3, 4, 5, 6, 7}
	result := generic.Chunk(input, 3)
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chunk failed: expected %v, got %v", expected, result)
	}

	// Test with exact division
	result = generic.Chunk([]int{1, 2, 3, 4}, 2)
	expected = [][]int{{1, 2}, {3, 4}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chunk with exact division failed: expected %v, got %v", expected, result)
	}

	// Test with chunk size larger than slice
	result = generic.Chunk([]int{1, 2}, 5)
	expected = [][]int{{1, 2}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Chunk with large chunk size failed: expected %v, got %v", expected, result)
	}

	// Test with zero or negative chunk size
	if result := generic.Chunk([]int{1, 2, 3}, 0); result != nil {
		t.Errorf("Chunk with zero size should return nil, got %v", result)
	}

	// Test with empty slice
	if result := generic.Chunk([]int{}, 2); result != nil {
		t.Errorf("Chunk with empty slice should return nil, got %v", result)
	}
}
