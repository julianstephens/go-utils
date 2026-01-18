package generic_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/julianstephens/go-utils/generic"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestKeys(t *testing.T) {
	// Test with string keys
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := generic.Keys(m)
	sort.Strings(keys) // Sort for consistent comparison
	expected := []string{"a", "b", "c"}
	tst.AssertDeepEqual(t, keys, expected)

	// Test with int keys
	intMap := map[int]string{1: "one", 2: "two", 3: "three"}
	intKeys := generic.Keys(intMap)
	sort.Ints(intKeys)
	expectedInts := []int{1, 2, 3}
	tst.AssertDeepEqual(t, intKeys, expectedInts)

	// Test with nil map
	tst.AssertNil(t, generic.Keys[string, int](nil), "Keys with nil map should return nil")

	// Test with empty map
	emptyMap := make(map[string]int)
	emptyKeys := generic.Keys(emptyMap)
	tst.AssertDeepEqual(t, len(emptyKeys), 0)
}

func TestValues(t *testing.T) {
	// Test with string values
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := generic.Values(m)
	sort.Ints(values) // Sort for consistent comparison
	expected := []int{1, 2, 3}
	tst.AssertDeepEqual(t, values, expected)

	// Test with string values
	stringMap := map[int]string{1: "one", 2: "two", 3: "three"}
	stringValues := generic.Values(stringMap)
	sort.Strings(stringValues)
	expectedStrings := []string{"one", "three", "two"}
	tst.AssertDeepEqual(t, stringValues, expectedStrings)

	// Test with nil map
	tst.AssertNil(t, generic.Values[string, int](nil), "Values with nil map should return nil")
}

func TestHasKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test existing key
	tst.AssertTrue(t, generic.HasKey(m, "b"), "HasKey should return true for existing key")

	// Test non-existing key
	tst.AssertFalse(t, generic.HasKey(m, "d"), "HasKey should return false for non-existing key")

	// Test with nil map
	tst.AssertFalse(t, generic.HasKey[string, int](nil, "a"), "HasKey should return false for nil map")
}

func TestHasValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test existing value
	tst.AssertTrue(t, generic.HasValue(m, 2), "HasValue should return true for existing value")

	// Test non-existing value
	tst.AssertFalse(t, generic.HasValue(m, 4), "HasValue should return false for non-existing value")

	// Test with nil map
	tst.AssertFalse(t, generic.HasValue[string, int](nil, 1), "HasValue should return false for nil map")
}

func TestMapToSlice(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test converting to key-value pairs
	result := generic.MapToSlice(m, func(k string, v int) string {
		return k + ":" + strconv.Itoa(v)
	})
	sort.Strings(result) // Sort for consistent comparison
	expected := []string{"a:1", "b:2", "c:3"}
	tst.AssertDeepEqual(t, result, expected)

	// Test with nil map
	tst.AssertNil(
		t,
		generic.MapToSlice[string, int, string](nil, func(k string, v int) string { return k }),
		"MapToSlice with nil map should return nil",
	)
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
	tst.AssertDeepEqual(t, result, expected)

	// Test with nil slice
	tst.AssertNil(t, generic.SliceToMap[Person, int, string](nil,
		func(p Person) int { return p.ID },
		func(p Person) string { return p.Name }), "SliceToMap with nil slice should return nil")

	// Test with duplicate keys (last one wins)
	duplicates := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 1, Name: "Alice2"},
	}
	dupResult := generic.SliceToMap(duplicates,
		func(p Person) int { return p.ID },
		func(p Person) string { return p.Name })
	tst.AssertDeepEqual(t, dupResult[1], "Alice2")
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
	tst.AssertDeepEqual(t, result, expected)
}

func TestFilterMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Test filtering by value
	result := generic.FilterMap(m, func(k string, v int) bool { return v%2 == 0 })
	expected := map[string]int{"b": 2, "d": 4}
	tst.AssertDeepEqual(t, result, expected)

	// Test filtering by key
	result = generic.FilterMap(m, func(k string, v int) bool { return k >= "c" })
	expected = map[string]int{"c": 3, "d": 4}
	tst.AssertDeepEqual(t, result, expected)

	// Test with nil map
	tst.AssertNil(
		t,
		generic.FilterMap[string, int](nil, func(k string, v int) bool { return true }),
		"FilterMap with nil map should return nil",
	)

	// Test with no matches
	noMatch := generic.FilterMap(m, func(k string, v int) bool { return v > 10 })
	tst.AssertDeepEqual(t, len(noMatch), 0)
}

func TestMapMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Test transforming keys and values
	result := generic.MapMap(m, func(k string, v int) (int, string) {
		return v, k + strconv.Itoa(v)
	})
	expected := map[int]string{1: "a1", 2: "b2", 3: "c3"}
	tst.AssertDeepEqual(t, result, expected)

	// Test with nil map
	tst.AssertNil(t, generic.MapMap[string, int, int, string](nil, func(k string, v int) (int, string) {
		return v, k
	}), "MapMap with nil map should return nil")
}

func TestMergeMap(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	m3 := map[string]int{"c": 5, "d": 6}

	result := generic.MergeMap(m1, m2, m3)
	expected := map[string]int{"a": 1, "b": 3, "c": 5, "d": 6}
	tst.AssertDeepEqual(t, result, expected)

	// Test with no maps
	emptyResult := generic.MergeMap[string, int]()
	tst.AssertDeepEqual(t, len(emptyResult), 0)

	// Test with single map
	singleResult := generic.MergeMap(m1)
	tst.AssertDeepEqual(t, singleResult, m1)
}

func TestCopyMap(t *testing.T) {
	original := map[string]int{"a": 1, "b": 2, "c": 3}
	copy := generic.CopyMap(original)

	// Test that copy has same content
	tst.AssertDeepEqual(t, copy, original)

	// Test that they are different maps (not same reference)
	copy["d"] = 4
	if len(original) == len(copy) {
		t.Error("CopyMap should create a separate map instance")
	}

	// Test with nil map
	tst.AssertNil(t, generic.CopyMap[string, int](nil), "CopyMap with nil map should return nil")
}
