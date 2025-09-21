package generic_test

import (
	"strconv"
	"testing"

	"github.com/julianstephens/go-utils/generic"
)

// Benchmark tests to ensure performance is reasonable

func BenchmarkMap(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Map(input, func(x int) string {
			return strconv.Itoa(x)
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Filter(input, func(x int) bool {
			return x%2 == 0
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Reduce(input, 0, func(acc, x int) int {
			return acc + x
		})
	}
}

func BenchmarkContains(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Contains(input, 500)
	}
}

func BenchmarkUnique(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i % 100 // Create duplicates
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Unique(input)
	}
}

func BenchmarkDifference(b *testing.B) {
	a := make([]int, 1000)
	b_slice := make([]int, 200)
	for i := range a {
		a[i] = i
	}
	for i := range b_slice {
		b_slice[i] = i * 2
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Difference(a, b_slice)
	}
}

func BenchmarkKeys(b *testing.B) {
	m := make(map[int]string, 1000)
	for i := 0; i < 1000; i++ {
		m[i] = strconv.Itoa(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Keys(m)
	}
}

func BenchmarkValues(b *testing.B) {
	m := make(map[int]string, 1000)
	for i := 0; i < 1000; i++ {
		m[i] = strconv.Itoa(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Values(m)
	}
}

func BenchmarkSliceToMap(b *testing.B) {
	type Item struct {
		ID   int
		Name string
	}

	items := make([]Item, 1000)
	for i := range items {
		items[i] = Item{ID: i, Name: strconv.Itoa(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.SliceToMap(items,
			func(item Item) int { return item.ID },
			func(item Item) string { return item.Name })
	}
}

func BenchmarkIf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = generic.If(i%2 == 0, "even", "odd")
	}
}

func BenchmarkPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = generic.Ptr(i)
	}
}

func BenchmarkDeref(b *testing.B) {
	value := 42
	ptr := &value

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generic.Deref(ptr)
	}
}
