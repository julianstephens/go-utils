package generic

// Contains returns true if the slice contains the specified value.
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// ContainsAll returns true if all elements in subset are present in mainSlice.
func ContainsAll[T comparable](mainSlice, subset []T) bool {
	if len(subset) == 0 {
		return true
	}

	mainMap := make(map[T]struct{})
	for _, item := range mainSlice {
		mainMap[item] = struct{}{}
	}

	for _, item := range subset {
		if _, found := mainMap[item]; !found {
			return false
		}
	}
	return true
}

// IndexOf returns the index of the first occurrence of value in the slice.
// Returns -1 if the value is not found.
func IndexOf[T comparable](slice []T, value T) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// Unique returns a new slice with duplicate elements removed.
// The order of elements is preserved, keeping the first occurrence of each element.
func Unique[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}
	seen := make(map[T]struct{})
	var result []T
	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Reverse returns a new slice with elements in reverse order.
func Reverse[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// DeleteElement removes an element from a slice at the specified index and returns the modified slice.
// Returns the original slice if the index is out of bounds.
func DeleteElement[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// InsertElement inserts an element at the specified index and returns the modified slice.
// If index is out of bounds, the element is appended to the end.
func InsertElement[T any](slice []T, index int, element T) []T {
	if index < 0 || index >= len(slice) {
		return append(slice, element)
	}
	slice = append(slice, Zero[T]())
	copy(slice[index+1:], slice[index:])
	slice[index] = element
	return slice
}

// Difference implements slice subtraction, returning elements in slice a that are not in slice b.
func Difference[T comparable](a, b []T) []T {
	if a == nil {
		return nil
	}
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// Intersection returns elements that exist in both slices.
func Intersection[T comparable](a, b []T) []T {
	if a == nil || b == nil {
		return nil
	}
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var intersection []T
	seen := make(map[T]struct{})
	for _, x := range a {
		if _, found := mb[x]; found {
			if _, alreadySeen := seen[x]; !alreadySeen {
				seen[x] = struct{}{}
				intersection = append(intersection, x)
			}
		}
	}
	return intersection
}

// Union returns elements that exist in either slice, with duplicates removed.
func Union[T comparable](a, b []T) []T {
	seen := make(map[T]struct{})
	var union []T
	
	for _, x := range a {
		if _, exists := seen[x]; !exists {
			seen[x] = struct{}{}
			union = append(union, x)
		}
	}
	
	for _, x := range b {
		if _, exists := seen[x]; !exists {
			seen[x] = struct{}{}
			union = append(union, x)
		}
	}
	
	return union
}

// Chunk divides a slice into chunks of the specified size.
// The last chunk may contain fewer elements if the slice length is not evenly divisible.
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 || len(slice) == 0 {
		return nil
	}
	
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}