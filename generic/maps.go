package generic

import (
	"maps"
)

// Keys returns a slice containing all keys from the map.
// The order of keys is not guaranteed.
func Keys[K comparable, V any](m map[K]V) []K {
	if m == nil {
		return nil
	}
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice containing all values from the map.
// The order of values is not guaranteed.
func Values[K comparable, V any](m map[K]V) []V {
	if m == nil {
		return nil
	}
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// HasKey returns true if the map contains the specified key.
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

// HasValue returns true if the map contains the specified value.
// This operation is O(n) as it needs to check all values.
func HasValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

// MapToSlice converts a map to a slice using a transformation function.
// The function f receives the key and value and returns the transformed element.
func MapToSlice[K comparable, V any, T any](m map[K]V, f func(K, V) T) []T {
	if m == nil {
		return nil
	}
	result := make([]T, 0, len(m))
	for k, v := range m {
		result = append(result, f(k, v))
	}
	return result
}

// SliceToMap converts a slice to a map using key and value extraction functions.
// If multiple elements produce the same key, the last one will be kept.
func SliceToMap[T any, K comparable, V any](slice []T, keyFunc func(T) K, valueFunc func(T) V) map[K]V {
	if slice == nil {
		return nil
	}
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		key := keyFunc(item)
		value := valueFunc(item)
		result[key] = value
	}
	return result
}

// SliceToMapBy converts a slice to a map using the provided key function.
// The values in the map are the original slice elements.
func SliceToMapBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	return SliceToMap(slice, keyFunc, func(t T) T { return t })
}

// FilterMap returns a new map containing only the key-value pairs that satisfy the predicate.
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	if m == nil {
		return nil
	}
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapMap applies a transformation function to each key-value pair and returns a new map.
func MapMap[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2)) map[K2]V2 {
	if m == nil {
		return nil
	}
	result := make(map[K2]V2, len(m))
	for k, v := range m {
		newKey, newValue := f(k, v)
		result[newKey] = newValue
	}
	return result
}

// MergeMap merges multiple maps into a new map.
// If the same key exists in multiple maps, the value from the last map takes precedence.
func MergeMap[K comparable, V any](mergeMaps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range mergeMaps {
		maps.Copy(result, m)
	}
	return result
}

// CopyMap creates a shallow copy of the map.
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	result := make(map[K]V, len(m))
	maps.Copy(result, m)
	return result
}
