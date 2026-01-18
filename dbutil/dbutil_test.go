package dbutil_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/julianstephens/go-utils/dbutil"
	tst "github.com/julianstephens/go-utils/tests"
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
		tst.AssertTrue(t, result == test.expected, "DefaultFieldMapper result should match expected")
	}
}

func TestIsNoRowsError(t *testing.T) {
	tst.AssertTrue(t, dbutil.IsNoRowsError(sql.ErrNoRows), "IsNoRowsError should return true for sql.ErrNoRows")
	tst.AssertFalse(
		t,
		dbutil.IsNoRowsError(errors.New("other error")),
		"IsNoRowsError should return false for other errors",
	)
	tst.AssertFalse(t, dbutil.IsNoRowsError(nil), "IsNoRowsError should return false for nil")
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
		tst.AssertTrue(t, result == test.expected, "IsConnectionError result should match expected")
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
		tst.AssertTrue(t, result == test.expected, "IsContextError result should match expected")
	}
}

func TestDefaultOptions(t *testing.T) {
	t.Run("DefaultConnectionOptions", func(t *testing.T) {
		opts := dbutil.DefaultConnectionOptions()
		tst.AssertNotNil(t, opts, "DefaultConnectionOptions should not be nil")
		tst.AssertTrue(t, opts.MaxOpenConns > 0, "MaxOpenConns should be positive")
		tst.AssertTrue(t, opts.MaxIdleConns > 0, "MaxIdleConns should be positive")
	})

	t.Run("DefaultQueryOptions", func(t *testing.T) {
		opts := dbutil.DefaultQueryOptions()
		tst.AssertNotNil(t, opts, "DefaultQueryOptions should not be nil")
		tst.AssertTrue(t, opts.Timeout > 0, "Query Timeout should be positive")
		tst.AssertNotNil(t, opts.FieldMapper, "FieldMapper should not be nil")
	})

	t.Run("DefaultTransactionOptions", func(t *testing.T) {
		opts := dbutil.DefaultTransactionOptions()
		tst.AssertNotNil(t, opts, "DefaultTransactionOptions should not be nil")
		tst.AssertTrue(t, opts.Timeout > 0, "Transaction Timeout should be positive")
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
