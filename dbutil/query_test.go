package dbutil_test

import (
	"testing"

	"github.com/julianstephens/go-utils/dbutil"
)

// Test query slice parameter validation
func TestQuerySliceValidation(t *testing.T) {
	t.Run("nil dest parameter", func(t *testing.T) {
		opts := dbutil.DefaultQueryOptions()
		err := dbutil.QuerySliceWithOptions(nil, nil, nil, "SELECT * FROM users", opts)
		if err == nil {
			t.Error("QuerySliceWithOptions should fail with nil dest")
		}
	})

	t.Run("non-pointer dest parameter", func(t *testing.T) {
		opts := dbutil.DefaultQueryOptions()
		var users []User
		err := dbutil.QuerySliceWithOptions(nil, nil, users, "SELECT * FROM users", opts)
		if err == nil {
			t.Error("QuerySliceWithOptions should fail with non-pointer dest")
		}
	})
}

// Test QueryRowScan parameter validation
func TestQueryRowScanValidation(t *testing.T) {
	t.Run("nil dest parameter", func(t *testing.T) {
		err := dbutil.QueryRowScan(nil, nil, nil, "SELECT * FROM users WHERE id = $1", 1)
		if err == nil {
			t.Error("QueryRowScan should fail with nil dest")
		}
	})

	t.Run("non-pointer dest parameter", func(t *testing.T) {
		var user User
		err := dbutil.QueryRowScan(nil, nil, user, "SELECT * FROM users WHERE id = $1", 1)
		if err == nil {
			t.Error("QueryRowScan should fail with non-pointer dest")
		}
	})

	t.Run("non-struct dest parameter", func(t *testing.T) {
		var str string
		err := dbutil.QueryRowScan(nil, nil, &str, "SELECT * FROM users WHERE id = $1", 1)
		if err == nil {
			t.Error("QueryRowScan should fail with non-struct dest")
		}
	})
}

// Benchmark tests for query functions
func BenchmarkQuerySliceOptions(b *testing.B) {
	opts := dbutil.DefaultQueryOptions()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Just test options creation performance
		_ = opts.FieldMapper("TestField")
	}
}