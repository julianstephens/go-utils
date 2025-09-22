package dbutil_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/julianstephens/go-utils/dbutil"
)

// Test structs
type User struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

type UserWithDefaults struct {
	ID    int64  `db:"id"`
	Name  string // No db tag, should use default field mapper
	Email string `db:"email"`
}

func TestDefaultFieldMapper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ID", "i_d"},                            // Each letter is uppercase, so gets separated
		{"UserID", "user_i_d"},                   // Same here - ID becomes i_d
		{"FirstName", "first_name"},              // This works as expected
		{"XMLHttpRequest", "x_m_l_http_request"}, // Each uppercase letter gets separated
		{"simple", "simple"},                     // No uppercase, no change
		{"CamelCase", "camel_case"},              // Standard camel case
	}

	for _, test := range tests {
		result := dbutil.DefaultFieldMapper(test.input)
		if result != test.expected {
			t.Errorf("DefaultFieldMapper(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestIsNoRowsError(t *testing.T) {
	if !dbutil.IsNoRowsError(sql.ErrNoRows) {
		t.Error("IsNoRowsError should return true for sql.ErrNoRows")
	}

	if dbutil.IsNoRowsError(errors.New("other error")) {
		t.Error("IsNoRowsError should return false for other errors")
	}

	if dbutil.IsNoRowsError(nil) {
		t.Error("IsNoRowsError should return false for nil")
	}
}

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{driver.ErrBadConn, true},
		{errors.New("connection refused"), true},
		{errors.New("Connection Reset"), true},
		{errors.New("timeout occurred"), true},
		{errors.New("no such host"), true},
		{errors.New("some other error"), false},
	}

	for _, test := range tests {
		result := dbutil.IsConnectionError(test.err)
		if result != test.expected {
			t.Errorf("IsConnectionError(%v) = %v, expected %v", test.err, result, test.expected)
		}
	}
}

func TestIsContextError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{context.Canceled, true},
		{context.DeadlineExceeded, true},
		{errors.New("other error"), false},
	}

	for _, test := range tests {
		result := dbutil.IsContextError(test.err)
		if result != test.expected {
			t.Errorf("IsContextError(%v) = %v, expected %v", test.err, result, test.expected)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	t.Run("DefaultConnectionOptions", func(t *testing.T) {
		opts := dbutil.DefaultConnectionOptions()
		if opts == nil {
			t.Error("DefaultConnectionOptions returned nil")
		}
		if opts.MaxOpenConns <= 0 {
			t.Error("DefaultConnectionOptions should have positive MaxOpenConns")
		}
		if opts.MaxIdleConns <= 0 {
			t.Error("DefaultConnectionOptions should have positive MaxIdleConns")
		}
	})

	t.Run("DefaultQueryOptions", func(t *testing.T) {
		opts := dbutil.DefaultQueryOptions()
		if opts == nil {
			t.Error("DefaultQueryOptions returned nil")
		}
		if opts.Timeout <= 0 {
			t.Error("DefaultQueryOptions should have positive Timeout")
		}
		if opts.FieldMapper == nil {
			t.Error("DefaultQueryOptions should have FieldMapper")
		}
	})

	t.Run("DefaultTransactionOptions", func(t *testing.T) {
		opts := dbutil.DefaultTransactionOptions()
		if opts == nil {
			t.Error("DefaultTransactionOptions returned nil")
		}
		if opts.Timeout <= 0 {
			t.Error("DefaultTransactionOptions should have positive Timeout")
		}
	})
}

// Benchmark tests
func BenchmarkDefaultFieldMapper(b *testing.B) {
	fieldName := "VeryLongFieldNameWithManyCamelCaseWords"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dbutil.DefaultFieldMapper(fieldName)
	}
}
