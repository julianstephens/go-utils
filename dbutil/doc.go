/*
Package dbutil provides utility functions and helpers for interacting with
databases in Go projects.

This package offers enhanced database/sql functionality with features like:
- Connection management with retry logic and context support
- Query execution helpers with automatic scanning and error context
- Transaction management utilities with rollback handling
- Context-aware operations for cancellation and timeout support
- Convenience generic helpers for scanning rows into slices/structs

Core types and functions
------------------------

- ConfigureDB(db *sql.DB, opts *ConnectionOptions) error
  - Helper to apply sensible defaults and connection pooling settings.

- PingWithContext(ctx context.Context, db *sql.DB, timeout time.Duration) error
  - Health check that respects context and timeout.

- QueryRow(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row
  - Thin wrapper providing consistent logging and error context.

- QueryRows[T any](ctx context.Context, db *sql.DB, query string, args ...interface{}) ([]T, error)
  - Generic helper to run queries and scan results into a slice of T.

- WithTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error
  - Run a function within a transaction and automatically rollback on error.

Examples
--------

Basic Query Row scan:

	var user User
	err := dbutil.QueryRowScan(ctx, db, &user, "SELECT id, name, email FROM users WHERE id = $1", 1)
	if err != nil {
			return err
	}

Query multiple rows into a slice using the generic helper:

	users, err := dbutil.QueryRows[User](ctx, db, "SELECT id, name, email FROM users WHERE active = $1", true)
	if err != nil {
			return err
	}

Safe transaction usage:

	err := dbutil.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "alice"); err != nil {
					return err
			}
			return nil
	})

The package is intentionally small and focused on making common database
operations less error-prone and easier to read. See the package tests for
additional usage patterns and edge cases.
*/
package dbutil
