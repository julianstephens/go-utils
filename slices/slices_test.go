package slices_test

import (
	"reflect"
	"testing"

	"github.com/julianstephens/go-utils/slices"
)

func TestIf(t *testing.T) {
	// Test with integers
	result := slices.If(true, 10, 20)
	if result != 10 {
		t.Errorf("If(true, 10, 20) = %v; expected 10", result)
	}

	result = slices.If(false, 10, 20)
	if result != 20 {
		t.Errorf("If(false, 10, 20) = %v; expected 20", result)
	}

	// Test with strings
	strResult := slices.If(true, "yes", "no")
	if strResult != "yes" {
		t.Errorf("If(true, 'yes', 'no') = %v; expected 'yes'", strResult)
	}

	strResult = slices.If(false, "yes", "no")
	if strResult != "no" {
		t.Errorf("If(false, 'yes', 'no') = %v; expected 'no'", strResult)
	}

	// Test with custom types
	type Person struct {
		Name string
	}
	person1 := Person{Name: "Alice"}
	person2 := Person{Name: "Bob"}

	personResult := slices.If(true, person1, person2)
	if personResult != person1 {
		t.Errorf("If with custom type failed: expected %v, got %v", person1, personResult)
	}

	// Test with nil values
	var nilPtr *string
	nonNilPtr := &[]string{"test"}[0]
	ptrResult := slices.If(false, nonNilPtr, nilPtr)
	if ptrResult != nilPtr {
		t.Errorf("If with nil pointer failed: expected nil, got %v", ptrResult)
	}

	// Test with slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	sliceResult := slices.If(true, slice1, slice2)
	if !reflect.DeepEqual(sliceResult, slice1) {
		t.Errorf("If with slices failed: expected %v, got %v", slice1, sliceResult)
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
			result := slices.DeleteElement(tt.slice, tt.index)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeleteElement(%v, %d) = %v; expected %v", tt.slice, tt.index, result, tt.expected)
			}
		})
	}
}

func TestDeleteElementString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	result := slices.DeleteElement(slice, 1)
	expected := []string{"apple", "cherry"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteElement with strings failed: expected %v, got %v", expected, result)
	}
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

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteElement with custom type failed: expected %v, got %v", expected, result)
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
			if result != tt.expected {
				t.Errorf("ContainsAll(%v, %v) = %v; expected %v", tt.mainSlice, tt.subset, result, tt.expected)
			}
		})
	}
}

func TestContainsAllString(t *testing.T) {
	mainSlice := []string{"apple", "banana", "cherry", "date"}
	subset := []string{"apple", "cherry"}

	result := slices.ContainsAll(mainSlice, subset)
	if !result {
		t.Errorf("ContainsAll with strings failed: expected true, got %v", result)
	}

	subset = []string{"apple", "grape"}
	result = slices.ContainsAll(mainSlice, subset)
	if result {
		t.Errorf("ContainsAll with strings failed: expected false, got %v", result)
	}
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
	if !result {
		t.Errorf("ContainsAll with custom type failed: expected true, got %v", result)
	}

	invalidSubset := []Color{
		{Name: "Red", Hex: "#FF0000"},
		{Name: "Yellow", Hex: "#FFFF00"},
	}

	result = slices.ContainsAll(colors, invalidSubset)
	if result {
		t.Errorf("ContainsAll with custom type failed: expected false, got %v", result)
	}
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
