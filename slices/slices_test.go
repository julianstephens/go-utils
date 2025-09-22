package slices_test

import (
	"testing"

	"github.com/julianstephens/go-utils/slices"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestIf(t *testing.T) {
	// Test with integers
	result := slices.If(true, 10, 20)
	tst.AssertTrue(t, result == 10, "If(true, 10, 20) should return 10")

	result = slices.If(false, 10, 20)
	tst.AssertTrue(t, result == 20, "If(false, 10, 20) should return 20")

	// Test with strings
	strResult := slices.If(true, "yes", "no")
	tst.AssertTrue(t, strResult == "yes", "If with strings should return yes when true")

	strResult = slices.If(false, "yes", "no")
	tst.AssertTrue(t, strResult == "no", "If with strings should return no when false")

	// Test with custom types
	type Person struct {
		Name string
	}
	person1 := Person{Name: "Alice"}
	person2 := Person{Name: "Bob"}

	personResult := slices.If(true, person1, person2)
	tst.AssertDeepEqual(t, personResult, person1)

	// Test with nil values
	var nilPtr *string
	nonNilPtr := &[]string{"test"}[0]
	ptrResult := slices.If(false, nonNilPtr, nilPtr)
	tst.AssertNil(t, ptrResult, "If with nil pointer should return nil")

	// Test with slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	sliceResult := slices.If(true, slice1, slice2)
	tst.AssertDeepEqual(t, sliceResult, slice1)
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
		{
			name:     "duplicates in a",
			a:        []string{"a", "b", "a", "c"},
			b:        []string{"b"},
			expected: []string{"a", "a", "c"},
		},
		{
			name:     "duplicates in b",
			a:        []string{"a", "b", "c"},
			b:        []string{"b", "b", "d"},
			expected: []string{"a", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slices.Difference(tt.a, tt.b)
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty (could be nil vs empty), treat as equal
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
			result := slices.DeleteElement(tt.slice, tt.index)
			tst.AssertDeepEqual(t, result, tt.expected)
		})
	}
}

func TestDeleteElementString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	result := slices.DeleteElement(slice, 1)
	expected := []string{"apple", "cherry"}

	tst.AssertDeepEqual(t, result, expected)
}

func TestDeleteElementCustomType(t *testing.T) {
	type Product struct {
		ID   int
		Name string
	}

	products := []Product{
		{ID: 1, Name: "Laptop"},
		{ID: 2, Name: "Mouse"},
		{ID: 3, Name: "Keyboard"},
	}

	result := slices.DeleteElement(products, 1)
	expected := []Product{
		{ID: 1, Name: "Laptop"},
		{ID: 3, Name: "Keyboard"},
	}

	tst.AssertDeepEqual(t, result, expected)
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
		{
			name:      "both slices empty",
			mainSlice: []int{},
			subset:    []int{},
			expected:  true,
		},
		{
			name:      "duplicate elements in subset",
			mainSlice: []int{1, 2, 3, 4},
			subset:    []int{2, 2, 3},
			expected:  true,
		},
		{
			name:      "single element match",
			mainSlice: []int{1},
			subset:    []int{1},
			expected:  true,
		},
		{
			name:      "single element no match",
			mainSlice: []int{1},
			subset:    []int{2},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slices.ContainsAll(tt.mainSlice, tt.subset)
			tst.AssertTrue(t, result == tt.expected, "ContainsAll result should match expected")
		})
	}
}

func TestContainsAllString(t *testing.T) {
	mainSlice := []string{"apple", "banana", "cherry", "date"}
	subset := []string{"apple", "cherry"}

	result := slices.ContainsAll(mainSlice, subset)
	tst.AssertTrue(t, result, "ContainsAll with strings should return true")

	subset = []string{"apple", "grape"}
	result = slices.ContainsAll(mainSlice, subset)
	tst.AssertFalse(t, result, "ContainsAll with strings should return false for missing element")
}

func TestContainsAllCustomType(t *testing.T) {
	type Color struct {
		Name string
		Hex  string
	}

	colors := []Color{
		{Name: "Red", Hex: "#FF0000"},
		{Name: "Green", Hex: "#00FF00"},
		{Name: "Blue", Hex: "#0000FF"},
	}

	subset := []Color{
		{Name: "Red", Hex: "#FF0000"},
		{Name: "Blue", Hex: "#0000FF"},
	}

	result := slices.ContainsAll(colors, subset)
	tst.AssertTrue(t, result, "ContainsAll with custom type should return true")

	invalidSubset := []Color{
		{Name: "Red", Hex: "#FF0000"},
		{Name: "Yellow", Hex: "#FFFF00"},
	}

	result = slices.ContainsAll(colors, invalidSubset)
	tst.AssertFalse(t, result, "ContainsAll with custom type should return false for missing element")
}

// Benchmark tests to ensure performance is reasonable
func BenchmarkIf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = slices.If(i%2 == 0, "even", "odd")
	}
}

func BenchmarkDifference(b *testing.B) {
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	subset := []string{"b", "d", "f", "h"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slices.Difference(a, subset)
	}
}

func BenchmarkDeleteElement(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slices.DeleteElement(slice, 5)
	}
}

func BenchmarkContainsAll(b *testing.B) {
	main := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	subset := []int{2, 4, 6, 8}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = slices.ContainsAll(main, subset)
	}
}
