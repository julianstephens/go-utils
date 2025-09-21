package generic_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/generic"
)

func TestMap(t *testing.T) {
	// Test with integers to strings
	input := []int{1, 2, 3, 4, 5}
	result := generic.Map(input, func(x int) string {
		return strconv.Itoa(x)
	})
	expected := []string{"1", "2", "3", "4", "5"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map failed: expected %v, got %v", expected, result)
	}

	// Test with strings to lengths
	stringInput := []string{"hello", "world", "go"}
	lengthResult := generic.Map(stringInput, func(s string) int {
		return len(s)
	})
	expectedLengths := []int{5, 5, 2}
	if !reflect.DeepEqual(lengthResult, expectedLengths) {
		t.Errorf("Map with strings failed: expected %v, got %v", expectedLengths, lengthResult)
	}

	// Test with nil slice
	var nilSlice []int
	nilResult := generic.Map(nilSlice, func(x int) string { return strconv.Itoa(x) })
	if nilResult != nil {
		t.Errorf("Map with nil slice should return nil, got %v", nilResult)
	}

	// Test with empty slice
	emptySlice := []int{}
	emptyResult := generic.Map(emptySlice, func(x int) string { return strconv.Itoa(x) })
	if len(emptyResult) != 0 {
		t.Errorf("Map with empty slice should return empty slice, got %v", emptyResult)
	}
}

func TestFilter(t *testing.T) {
	// Test with integers - filter even numbers
	input := []int{1, 2, 3, 4, 5, 6}
	result := generic.Filter(input, func(x int) bool { return x%2 == 0 })
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter failed: expected %v, got %v", expected, result)
	}

	// Test with strings - filter by length
	stringInput := []string{"a", "hello", "go", "world"}
	stringResult := generic.Filter(stringInput, func(s string) bool { return len(s) > 2 })
	expectedStrings := []string{"hello", "world"}
	if !reflect.DeepEqual(stringResult, expectedStrings) {
		t.Errorf("Filter with strings failed: expected %v, got %v", expectedStrings, stringResult)
	}

	// Test with nil slice
	var nilSlice []int
	nilResult := generic.Filter(nilSlice, func(x int) bool { return true })
	if nilResult != nil {
		t.Errorf("Filter with nil slice should return nil, got %v", nilResult)
	}

	// Test with no matches
	noMatchResult := generic.Filter([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 })
	if len(noMatchResult) != 0 {
		t.Errorf("Filter with no matches should return empty slice, got %v", noMatchResult)
	}
}

func TestReduce(t *testing.T) {
	// Test sum
	input := []int{1, 2, 3, 4, 5}
	sum := generic.Reduce(input, 0, func(acc, x int) int { return acc + x })
	if sum != 15 {
		t.Errorf("Reduce sum failed: expected 15, got %v", sum)
	}

	// Test product
	product := generic.Reduce(input, 1, func(acc, x int) int { return acc * x })
	if product != 120 {
		t.Errorf("Reduce product failed: expected 120, got %v", product)
	}

	// Test string concatenation
	words := []string{"hello", "world", "go"}
	sentence := generic.Reduce(words, "", func(acc, word string) string {
		if acc == "" {
			return word
		}
		return acc + " " + word
	})
	expected := "hello world go"
	if sentence != expected {
		t.Errorf("Reduce string concatenation failed: expected %v, got %v", expected, sentence)
	}

	// Test with empty slice
	emptySum := generic.Reduce([]int{}, 10, func(acc, x int) int { return acc + x })
	if emptySum != 10 {
		t.Errorf("Reduce with empty slice should return initial value, got %v", emptySum)
	}
}

func TestFind(t *testing.T) {
	// Test finding existing element
	input := []int{1, 2, 3, 4, 5}
	result, found := generic.Find(input, func(x int) bool { return x > 3 })
	if !found {
		t.Error("Find should have found an element")
	}
	if result != 4 {
		t.Errorf("Find failed: expected 4, got %v", result)
	}

	// Test not finding element
	result, found = generic.Find(input, func(x int) bool { return x > 10 })
	if found {
		t.Error("Find should not have found an element")
	}
	if result != 0 {
		t.Errorf("Find should return zero value when not found, got %v", result)
	}

	// Test with strings
	words := []string{"apple", "banana", "cherry"}
	wordResult, wordFound := generic.Find(words, func(s string) bool { return strings.HasPrefix(s, "b") })
	if !wordFound {
		t.Error("Find should have found a word starting with 'b'")
	}
	if wordResult != "banana" {
		t.Errorf("Find failed: expected 'banana', got %v", wordResult)
	}
}

func TestAny(t *testing.T) {
	// Test with some matching elements
	input := []int{1, 2, 3, 4, 5}
	result := generic.Any(input, func(x int) bool { return x > 3 })
	if !result {
		t.Error("Any should return true when some elements match")
	}

	// Test with no matching elements
	result = generic.Any(input, func(x int) bool { return x > 10 })
	if result {
		t.Error("Any should return false when no elements match")
	}

	// Test with empty slice
	result = generic.Any([]int{}, func(x int) bool { return true })
	if result {
		t.Error("Any should return false for empty slice")
	}
}

func TestAll(t *testing.T) {
	// Test with all matching elements
	input := []int{2, 4, 6, 8}
	result := generic.All(input, func(x int) bool { return x%2 == 0 })
	if !result {
		t.Error("All should return true when all elements match")
	}

	// Test with some non-matching elements
	input = []int{2, 3, 4, 6}
	result = generic.All(input, func(x int) bool { return x%2 == 0 })
	if result {
		t.Error("All should return false when not all elements match")
	}

	// Test with empty slice
	result = generic.All([]int{}, func(x int) bool { return false })
	if !result {
		t.Error("All should return true for empty slice")
	}
}

func TestForEach(t *testing.T) {
	// Test side effect
	input := []int{1, 2, 3, 4, 5}
	sum := 0
	generic.ForEach(input, func(x int) {
		sum += x
	})
	if sum != 15 {
		t.Errorf("ForEach failed: expected sum to be 15, got %v", sum)
	}

	// Test with strings
	words := []string{"hello", "world"}
	var result []string
	generic.ForEach(words, func(s string) {
		result = append(result, strings.ToUpper(s))
	})
	expected := []string{"HELLO", "WORLD"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ForEach with strings failed: expected %v, got %v", expected, result)
	}
}