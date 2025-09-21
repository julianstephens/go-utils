package generic_test

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/julianstephens/go-utils/generic"
)

func TestKeys(t *testing.T) {
	// Test with string keys
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := generic.Keys(m)
	sort.Strings(keys) // Sort for consistent comparison
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Keys failed: expected %v, got %v", expected, keys)
	}

	// Test with int keys
	intMap := map[int]string{1: "one", 2: "two", 3: "three"}
	intKeys := generic.Keys(intMap)
	sort.Ints(intKeys)
	expectedInts := []int{1, 2, 3}
	if !reflect.DeepEqual(intKeys, expectedInts) {
		t.Errorf("Keys with int keys failed: expected %v, got %v", expectedInts, intKeys)
	}

	// Test with nil map
	if result := generic.Keys[string, int](nil); result != nil {
		t.Errorf("Keys with nil map should return nil, got %v", result)
	}

	// Test with empty map
	emptyMap := make(map[string]int)
	emptyKeys := generic.Keys(emptyMap)
	if len(emptyKeys) != 0 {
		t.Errorf("Keys with empty map should return empty slice, got %v", emptyKeys)
	}
}

func TestValues(t *testing.T) {
	// Test with string values
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := generic.Values(m)
	sort.Ints(values) // Sort for consistent comparison
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("Values failed: expected %v, got %v", expected, values)
	}

	// Test with string values
	stringMap := map[int]string{1: "one", 2: "two", 3: "three"}
	stringValues := generic.Values(stringMap)
	sort.Strings(stringValues)
	expectedStrings := []string{"one", "three", "two"}
	if !reflect.DeepEqual(stringValues, expectedStrings) {
		t.Errorf("Values with string values failed: expected %v, got %v", expectedStrings, stringValues)
	}

	// Test with nil map
	if result := generic.Values[string, int](nil); result != nil {
		t.Errorf("Values with nil map should return nil, got %v", result)
	}
}

func TestHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test existing key
	if !generic.HasKey(m, "b") {
		t.Error("HasKey should return true for existing key")
	}

	// Test non-existing key
	if generic.HasKey(m, "d") {
		t.Error("HasKey should return false for non-existing key")
	}

	// Test with nil map
	if generic.HasKey[string, int](nil, "a") {
		t.Error("HasKey should return false for nil map")
	}
}

func TestHasValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test existing value
	if !generic.HasValue(m, 2) {
		t.Error("HasValue should return true for existing value")
	}

	// Test non-existing value
	if generic.HasValue(m, 4) {
		t.Error("HasValue should return false for non-existing value")
	}

	// Test with nil map
	if generic.HasValue[string, int](nil, 1) {
		t.Error("HasValue should return false for nil map")
	}
}

func TestMapToSlice(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test converting to key-value pairs
	result := generic.MapToSlice(m, func(k string, v int) string {
		return k + ":" + strconv.Itoa(v)
	})
	sort.Strings(result) // Sort for consistent comparison
	expected := []string{"a:1", "b:2", "c:3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapToSlice failed: expected %v, got %v", expected, result)
	}

	// Test with nil map
	if result := generic.MapToSlice[string, int, string](nil, func(k string, v int) string { return k }); result != nil {
		t.Errorf("MapToSlice with nil map should return nil, got %v", result)
	}
}

func TestSliceToMap(t *testing.T) {
	type Person struct {
		ID   int
		Name string
	}

	people := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	// Test converting slice to map
	result := generic.SliceToMap(people,
		func(p Person) int { return p.ID },
		func(p Person) string { return p.Name })

	expected := map[int]string{1: "Alice", 2: "Bob", 3: "Charlie"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceToMap failed: expected %v, got %v", expected, result)
	}

	// Test with nil slice
	if result := generic.SliceToMap[Person, int, string](nil,
		func(p Person) int { return p.ID },
		func(p Person) string { return p.Name }); result != nil {
		t.Errorf("SliceToMap with nil slice should return nil, got %v", result)
	}

	// Test with duplicate keys (last one wins)
	duplicates := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 1, Name: "Alice2"},
	}
	dupResult := generic.SliceToMap(duplicates,
		func(p Person) int { return p.ID },
		func(p Person) string { return p.Name })
	if dupResult[1] != "Alice2" {
		t.Errorf("SliceToMap with duplicates should keep last value, got %v", dupResult[1])
	}
}

func TestSliceToMapBy(t *testing.T) {
	type Person struct {
		ID   int
		Name string
	}

	people := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	result := generic.SliceToMapBy(people, func(p Person) int { return p.ID })
	expected := map[int]Person{1: {ID: 1, Name: "Alice"}, 2: {ID: 2, Name: "Bob"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceToMapBy failed: expected %v, got %v", expected, result)
	}
}

func TestFilterMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Test filtering by value
	result := generic.FilterMap(m, func(k string, v int) bool { return v%2 == 0 })
	expected := map[string]int{"b": 2, "d": 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterMap failed: expected %v, got %v", expected, result)
	}

	// Test filtering by key
	result = generic.FilterMap(m, func(k string, v int) bool { return k >= "c" })
	expected = map[string]int{"c": 3, "d": 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterMap by key failed: expected %v, got %v", expected, result)
	}

	// Test with nil map
	if result := generic.FilterMap[string, int](nil, func(k string, v int) bool { return true }); result != nil {
		t.Errorf("FilterMap with nil map should return nil, got %v", result)
	}

	// Test with no matches
	noMatch := generic.FilterMap(m, func(k string, v int) bool { return v > 10 })
	if len(noMatch) != 0 {
		t.Errorf("FilterMap with no matches should return empty map, got %v", noMatch)
	}
}

func TestMapMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test transforming keys and values
	result := generic.MapMap(m, func(k string, v int) (int, string) {
		return v, k + strconv.Itoa(v)
	})
	expected := map[int]string{1: "a1", 2: "b2", 3: "c3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapMap failed: expected %v, got %v", expected, result)
	}

	// Test with nil map
	if result := generic.MapMap[string, int, int, string](nil, func(k string, v int) (int, string) {
		return v, k
	}); result != nil {
		t.Errorf("MapMap with nil map should return nil, got %v", result)
	}
}

func TestMergeMap(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	m3 := map[string]int{"c": 5, "d": 6}

	result := generic.MergeMap(m1, m2, m3)
	expected := map[string]int{"a": 1, "b": 3, "c": 5, "d": 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MergeMap failed: expected %v, got %v", expected, result)
	}

	// Test with no maps
	emptyResult := generic.MergeMap[string, int]()
	if len(emptyResult) != 0 {
		t.Errorf("MergeMap with no maps should return empty map, got %v", emptyResult)
	}

	// Test with single map
	singleResult := generic.MergeMap(m1)
	if !reflect.DeepEqual(singleResult, m1) {
		t.Errorf("MergeMap with single map should return copy of that map, got %v", singleResult)
	}
}

func TestCopyMap(t *testing.T) {
	original := map[string]int{"a": 1, "b": 2, "c": 3}
	copy := generic.CopyMap(original)

	// Test that copy has same content
	if !reflect.DeepEqual(copy, original) {
		t.Errorf("CopyMap failed: expected %v, got %v", original, copy)
	}

	// Test that they are different maps (not same reference)
	copy["d"] = 4
	if len(original) == len(copy) {
		t.Error("CopyMap should create a separate map instance")
	}

	// Test with nil map
	if result := generic.CopyMap[string, int](nil); result != nil {
		t.Errorf("CopyMap with nil map should return nil, got %v", result)
	}
}