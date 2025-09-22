package generic_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/generic"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestMap(t *testing.T) {
	// Test with integers to strings
	input := []int{1, 2, 3, 4, 5}
	result := generic.Map(input, func(x int) string {
		return strconv.Itoa(x)
	})
	expected := []string{"1", "2", "3", "4", "5"}
	tst.AssertDeepEqual(t, result, expected)

	// Test with strings to lengths
	stringInput := []string{"hello", "world", "go"}
	lengthResult := generic.Map(stringInput, func(s string) int {
		return len(s)
	})
	expectedLengths := []int{5, 5, 2}
	tst.AssertDeepEqual(t, lengthResult, expectedLengths)

	// Test with nil slice
	var nilSlice []int
	nilResult := generic.Map(nilSlice, func(x int) string { return strconv.Itoa(x) })
	tst.AssertNil(t, nilResult, "Map with nil slice should return nil")

	// Test with empty slice
	emptySlice := []int{}
	emptyResult := generic.Map(emptySlice, func(x int) string { return strconv.Itoa(x) })
	tst.AssertDeepEqual(t, len(emptyResult), 0)
}

func TestFilter(t *testing.T) {
	// Test with integers - filter even numbers
	input := []int{1, 2, 3, 4, 5, 6}
	result := generic.Filter(input, func(x int) bool { return x%2 == 0 })
	expected := []int{2, 4, 6}
	tst.AssertDeepEqual(t, result, expected)

	// Test with strings - filter by length
	stringInput := []string{"a", "hello", "go", "world"}
	stringResult := generic.Filter(stringInput, func(s string) bool { return len(s) > 2 })
	expectedStrings := []string{"hello", "world"}
	tst.AssertDeepEqual(t, stringResult, expectedStrings)

	// Test with nil slice
	var nilSlice []int
	nilResult := generic.Filter(nilSlice, func(x int) bool { return true })
	tst.AssertNil(t, nilResult, "Filter with nil slice should return nil")

	// Test with no matches
	noMatchResult := generic.Filter([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 })
	tst.AssertDeepEqual(t, len(noMatchResult), 0)
}

func TestReduce(t *testing.T) {
	// Test sum
	input := []int{1, 2, 3, 4, 5}
	sum := generic.Reduce(input, 0, func(acc, x int) int { return acc + x })
	tst.AssertDeepEqual(t, sum, 15)

	// Test product
	product := generic.Reduce(input, 1, func(acc, x int) int { return acc * x })
	tst.AssertDeepEqual(t, product, 120)

	// Test string concatenation
	words := []string{"hello", "world", "go"}
	sentence := generic.Reduce(words, "", func(acc, word string) string {
		if acc == "" {
			return word
		}
		return acc + " " + word
	})
	expected := "hello world go"
	tst.AssertDeepEqual(t, sentence, expected)

	// Test with empty slice
	emptySum := generic.Reduce([]int{}, 10, func(acc, x int) int { return acc + x })
	tst.AssertDeepEqual(t, emptySum, 10)
}

func TestFind(t *testing.T) {
	// Test finding existing element
	input := []int{1, 2, 3, 4, 5}
	result, found := generic.Find(input, func(x int) bool { return x > 3 })
	tst.AssertTrue(t, found, "Find should have found an element")
	tst.AssertDeepEqual(t, result, 4)

	// Test not finding element
	result, found = generic.Find(input, func(x int) bool { return x > 10 })
	tst.AssertFalse(t, found, "Find should not have found an element")
	tst.AssertDeepEqual(t, result, 0)

	// Test with strings
	words := []string{"apple", "banana", "cherry"}
	wordResult, wordFound := generic.Find(words, func(s string) bool { return strings.HasPrefix(s, "b") })
	tst.AssertTrue(t, wordFound, "Find should have found a word starting with 'b'")
	tst.AssertDeepEqual(t, wordResult, "banana")
}

func TestAny(t *testing.T) {
	// Test with some matching elements
	input := []int{1, 2, 3, 4, 5}
	result := generic.Any(input, func(x int) bool { return x > 3 })
	tst.AssertTrue(t, result, "Any should return true when some elements match")

	// Test with no matching elements
	result = generic.Any(input, func(x int) bool { return x > 10 })
	tst.AssertFalse(t, result, "Any should return false when no elements match")

	// Test with empty slice
	result = generic.Any([]int{}, func(x int) bool { return true })
	tst.AssertFalse(t, result, "Any should return false for empty slice")
}

func TestAll(t *testing.T) {
	// Test with all matching elements
	input := []int{2, 4, 6, 8}
	result := generic.All(input, func(x int) bool { return x%2 == 0 })
	tst.AssertTrue(t, result, "All should return true when all elements match")

	// Test with some non-matching elements
	input = []int{2, 3, 4, 6}
	result = generic.All(input, func(x int) bool { return x%2 == 0 })
	tst.AssertFalse(t, result, "All should return false when not all elements match")

	// Test with empty slice
	result = generic.All([]int{}, func(x int) bool { return false })
	tst.AssertTrue(t, result, "All should return true for empty slice")
}

func TestForEach(t *testing.T) {
	// Test side effect
	input := []int{1, 2, 3, 4, 5}
	sum := 0
	generic.ForEach(input, func(x int) {
		sum += x
	})
	tst.AssertDeepEqual(t, sum, 15)

	// Test with strings
	words := []string{"hello", "world"}
	var result []string
	generic.ForEach(words, func(s string) {
		result = append(result, strings.ToUpper(s))
	})
	expected := []string{"HELLO", "WORLD"}
	tst.AssertDeepEqual(t, result, expected)
}
