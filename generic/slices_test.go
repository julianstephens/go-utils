package generic_test

import (
	"testing"

	"github.com/julianstephens/go-utils/generic"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestContains(t *testing.T) {
	// Test with integers
	slice := []int{1, 2, 3, 4, 5}
	tst.AssertTrue(t, generic.Contains(slice, 3), "Contains should return true for existing element")
	tst.AssertFalse(t, generic.Contains(slice, 6), "Contains should return false for non-existing element")

	// Test with strings
	stringSlice := []string{"apple", "banana", "cherry"}
	tst.AssertTrue(t, generic.Contains(stringSlice, "banana"), "Contains should return true for existing string")
	tst.AssertFalse(t, generic.Contains(stringSlice, "grape"), "Contains should return false for non-existing string")

	// Test with empty slice
	tst.AssertFalse(t, generic.Contains([]int{}, 1), "Contains should return false for empty slice")
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
			tst.AssertDeepEqual(t, result, tt.expected)
		})
	}
}

func TestIndexOf(t *testing.T) {
	// Test with integers
	slice := []int{10, 20, 30, 40, 50}
	tst.AssertDeepEqual(t, generic.IndexOf(slice, 30), 2)
	tst.AssertDeepEqual(t, generic.IndexOf(slice, 60), -1)

	// Test with strings
	stringSlice := []string{"a", "b", "c", "d"}
	tst.AssertDeepEqual(t, generic.IndexOf(stringSlice, "c"), 2)
}

func TestUnique(t *testing.T) {
	// Test with duplicates
	input := []int{1, 2, 2, 3, 3, 3, 4}
	result := generic.Unique(input)
	expected := []int{1, 2, 3, 4}
	tst.AssertDeepEqual(t, result, expected)

	// Test with no duplicates
	noDupsInput := []int{1, 2, 3, 4}
	noDupsResult := generic.Unique(noDupsInput)
	tst.AssertDeepEqual(t, noDupsResult, noDupsInput)

	// Test with nil slice
	tst.AssertNil(t, generic.Unique[int](nil), "Unique with nil slice should return nil")

	// Test with strings
	stringInput := []string{"apple", "banana", "apple", "cherry", "banana"}
	stringResult := generic.Unique(stringInput)
	expectedStrings := []string{"apple", "banana", "cherry"}
	tst.AssertDeepEqual(t, stringResult, expectedStrings)
}

func TestReverse(t *testing.T) {
	// Test with integers
	input := []int{1, 2, 3, 4, 5}
	result := generic.Reverse(input)
	expected := []int{5, 4, 3, 2, 1}
	tst.AssertDeepEqual(t, result, expected)

	// Test with strings
	stringInput := []string{"a", "b", "c"}
	stringResult := generic.Reverse(stringInput)
	expectedStrings := []string{"c", "b", "a"}
	tst.AssertDeepEqual(t, stringResult, expectedStrings)

	// Test with nil slice
	tst.AssertNil(t, generic.Reverse[int](nil), "Reverse with nil slice should return nil")

	// Test with empty slice
	emptyResult := generic.Reverse([]int{})
	tst.AssertDeepEqual(t, len(emptyResult), 0)
}

func TestDeleteElement(t *testing.T) {
	// Test normal deletion
	input := []int{1, 2, 3, 4, 5}
	result := generic.DeleteElement(input, 2)
	expected := []int{1, 2, 4, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test deletion at beginning
	input2 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input2, 0)
	expected = []int{2, 3, 4, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test deletion at end
	input3 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input3, len(input3)-1)
	expected = []int{1, 2, 3, 4}
	tst.AssertDeepEqual(t, result, expected)

	// Test out of bounds index
	input4 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input4, 10)
	tst.AssertDeepEqual(t, result, input4)

	// Test negative index
	input5 := []int{1, 2, 3, 4, 5}
	result = generic.DeleteElement(input5, -1)
	tst.AssertDeepEqual(t, result, input5)
}

func TestInsertElement(t *testing.T) {
	// Test normal insertion
	input := []int{1, 2, 4, 5}
	result := generic.InsertElement(input, 2, 3)
	expected := []int{1, 2, 3, 4, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test insertion at beginning
	result = generic.InsertElement(input, 0, 0)
	expected = []int{0, 1, 2, 4, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test insertion out of bounds (should append)
	result = generic.InsertElement(input, 10, 6)
	expected = []int{1, 2, 4, 5, 6}
	tst.AssertDeepEqual(t, result, expected)
}

func TestDifference(t *testing.T) {
	// Test basic difference
	a := []int{1, 2, 3, 4, 5}
	b := []int{2, 4}
	result := generic.Difference(a, b)
	expected := []int{1, 3, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test with strings
	stringA := []string{"apple", "banana", "cherry", "date"}
	stringB := []string{"banana", "date"}
	stringResult := generic.Difference(stringA, stringB)
	expectedStrings := []string{"apple", "cherry"}
	tst.AssertDeepEqual(t, stringResult, expectedStrings)

	// Test with nil slice
	tst.AssertNil(t, generic.Difference[int](nil, []int{1, 2}), "Difference with nil slice should return nil")
}

func TestIntersection(t *testing.T) {
	// Test basic intersection
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	result := generic.Intersection(a, b)
	expected := []int{3, 4, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test with no common elements
	c := []int{8, 9, 10}
	result = generic.Intersection(a, c)
	tst.AssertDeepEqual(t, len(result), 0)

	// Test with duplicates
	d := []int{1, 1, 2, 2, 3}
	e := []int{1, 2, 2, 4}
	result = generic.Intersection(d, e)
	expected = []int{1, 2}
	tst.AssertDeepEqual(t, result, expected)
}

func TestUnion(t *testing.T) {
	// Test basic union
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	result := generic.Union(a, b)
	expected := []int{1, 2, 3, 4, 5}
	tst.AssertDeepEqual(t, result, expected)

	// Test with duplicates within slices
	c := []int{1, 1, 2}
	d := []int{2, 2, 3}
	result = generic.Union(c, d)
	expected = []int{1, 2, 3}
	tst.AssertDeepEqual(t, result, expected)
}

func TestChunk(t *testing.T) {
	// Test normal chunking
	input := []int{1, 2, 3, 4, 5, 6, 7}
	result := generic.Chunk(input, 3)
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	tst.AssertDeepEqual(t, result, expected)

	// Test with exact division
	result = generic.Chunk([]int{1, 2, 3, 4}, 2)
	expected = [][]int{{1, 2}, {3, 4}}
	tst.AssertDeepEqual(t, result, expected)

	// Test with chunk size larger than slice
	result = generic.Chunk([]int{1, 2}, 5)
	expected = [][]int{{1, 2}}
	tst.AssertDeepEqual(t, result, expected)

	// Test with zero or negative chunk size
	tst.AssertNil(t, generic.Chunk([]int{1, 2, 3}, 0), "Chunk with zero size should return nil")

	// Test with empty slice
	tst.AssertNil(t, generic.Chunk([]int{}, 2), "Chunk with empty slice should return nil")
}
