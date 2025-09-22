package dbutil

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ConnectionOptions holds configuration options for database connections.
type ConnectionOptions struct {
	// MaxOpenConns sets the maximum number of open connections to the database.
	MaxOpenConns int
	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	MaxIdleConns int
	// ConnMaxLifetime sets the maximum amount of time a connection may be reused.
	ConnMaxLifetime time.Duration
	// ConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	ConnMaxIdleTime time.Duration
	// PingTimeout sets the timeout for ping operations.
	PingTimeout time.Duration
	// RetryAttempts sets the number of retry attempts for failed operations.
	RetryAttempts int
	// RetryDelay sets the delay between retry attempts.
	RetryDelay time.Duration
}

// QueryOptions holds configuration options for query execution.
type QueryOptions struct {
	// Timeout sets the timeout for query execution.
	Timeout time.Duration
	// MaxRows limits the number of rows returned (0 means no limit).
	MaxRows int
	// FieldMapper is a function to map struct field names to database column names.
	FieldMapper func(string) string
}

// TransactionOptions holds configuration options for transactions.
type TransactionOptions struct {
	// Isolation sets the transaction isolation level.
	Isolation sql.IsolationLevel
	// ReadOnly sets whether the transaction is read-only.
	ReadOnly bool
	// Timeout sets the timeout for the entire transaction.
	Timeout time.Duration
}

// DefaultConnectionOptions returns sensible default connection options.
func DefaultConnectionOptions() *ConnectionOptions {
	return &ConnectionOptions{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		PingTimeout:     5 * time.Second,
		RetryAttempts:   3,
		RetryDelay:      time.Second,
	}
}

// DefaultQueryOptions returns sensible default query options.
func DefaultQueryOptions() *QueryOptions {
	return &QueryOptions{
		Timeout:     30 * time.Second,
		MaxRows:     0, // No limit
		FieldMapper: DefaultFieldMapper,
	}
}

// DefaultTransactionOptions returns sensible default transaction options.
func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Timeout:   time.Minute,
	}
}

// DefaultFieldMapper converts Go struct field names to database column names.
// It converts CamelCase to snake_case (e.g., "UserID" -> "user_id").
func DefaultFieldMapper(fieldName string) string {
	var result strings.Builder
	runes := []rune(fieldName)
	for i, r := range runes {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ConfigureDB configures a database connection with the provided options.
func ConfigureDB(db *sql.DB, opts *ConnectionOptions) error {
	if opts == nil {
		opts = DefaultConnectionOptions()
	}

	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetConnMaxLifetime(opts.ConnMaxLifetime)
	db.SetConnMaxIdleTime(opts.ConnMaxIdleTime)

	// Test the connection with ping
	ctx, cancel := context.WithTimeout(context.Background(), opts.PingTimeout)
	defer cancel()

	if err := PingWithRetry(ctx, db, opts.RetryAttempts, opts.RetryDelay); err != nil {
		return fmt.Errorf("dbutil: connection configuration failed: %w", err)
	}

	return nil
}

// PingWithContext pings the database with a context timeout.
func PingWithContext(ctx context.Context, db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("dbutil: ping failed: %w", err)
	}
	return nil
}

// PingWithRetry pings the database with retry logic.
func PingWithRetry(ctx context.Context, db *sql.DB, attempts int, delay time.Duration) error {
	var lastErr error
	for i := range attempts {
		if err := db.PingContext(ctx); err != nil {
			lastErr = err
			if i < attempts-1 {
				select {
				case <-ctx.Done():
					return fmt.Errorf("dbutil: ping context cancelled: %w", ctx.Err())
				case <-time.After(delay):
					continue
				}
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("dbutil: ping failed after %d attempts: %w", attempts, lastErr)
}

// QueryRow executes a query that is expected to return at most one row.
// It returns a *sql.Row which can be scanned into destination variables.
func QueryRow(ctx context.Context, db *sql.DB, query string, args ...any) *sql.Row {
	return db.QueryRowContext(ctx, query, args...)
}

// QueryRowTx is like QueryRow but uses a transaction.
func QueryRowTx(ctx context.Context, tx *sql.Tx, query string, args ...any) *sql.Row {
	return tx.QueryRowContext(ctx, query, args...)
}

// QueryRows executes a query and returns multiple rows.
// It's the caller's responsibility to close the returned *sql.Rows.
func QueryRows(ctx context.Context, db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("dbutil: query rows failed: %w", err)
	}
	return rows, nil
}

// QueryRowsTx is like QueryRows but uses a transaction.
func QueryRowsTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("dbutil: query rows (tx) failed: %w", err)
	}
	return rows, nil
}

// Exec executes a query without returning any rows.
// It returns the number of rows affected and any error encountered.
func Exec(ctx context.Context, db *sql.DB, query string, args ...any) (sql.Result, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("dbutil: exec failed: %w", err)
	}
	return result, nil
}

// ExecTx is like Exec but uses a transaction.
func ExecTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("dbutil: exec (tx) failed: %w", err)
	}
	return result, nil
}

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func WithTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	return WithTransactionOptions(ctx, db, DefaultTransactionOptions(), fn)
}

// WithTransactionOptions is like WithTransaction but allows specifying transaction options.
func WithTransactionOptions(ctx context.Context, db *sql.DB, opts *TransactionOptions, fn func(*sql.Tx) error) error {
	if opts == nil {
		opts = DefaultTransactionOptions()
	}

	// Create transaction context with timeout if specified
	txCtx := ctx
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		txCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Begin transaction with options
	txOpts := &sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	}

	tx, err := db.BeginTx(txCtx, txOpts)
	if err != nil {
		return fmt.Errorf("dbutil: begin transaction failed: %w", err)
	}

	// Ensure rollback on panic or error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		}
	}()

	// Execute function
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("dbutil: transaction failed: %w (rollback error: %v)", err, rbErr)
		}
		return fmt.Errorf("dbutil: transaction failed: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("dbutil: commit transaction failed: %w", err)
	}

	return nil
}

// QueryRowScan executes a query that returns a single row and scans the result into dest.
// dest should be a pointer to a struct with appropriate db tags.
func QueryRowScan(ctx context.Context, db *sql.DB, dest any, query string, args ...any) error {
	return queryRowScanImpl(ctx, func() *sql.Row {
		return db.QueryRowContext(ctx, query, args...)
	}, dest)
}

// QueryRowScanTx is like QueryRowScan but uses a transaction.
func QueryRowScanTx(ctx context.Context, tx *sql.Tx, dest any, query string, args ...any) error {
	return queryRowScanImpl(ctx, func() *sql.Row {
		return tx.QueryRowContext(ctx, query, args...)
	}, dest)
}

// queryRowScanImpl is the common implementation for QueryRowScan functions.
func queryRowScanImpl(_ context.Context, queryFn func() *sql.Row, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Pointer || destValue.IsNil() {
		return fmt.Errorf("dbutil: dest must be a non-nil pointer")
	}

	destElem := destValue.Elem()
	if destElem.Kind() != reflect.Struct {
		return fmt.Errorf("dbutil: dest must be a pointer to a struct")
	}

	// Get struct fields and their db tags
	fields, err := getStructFields(destElem.Type())
	if err != nil {
		return fmt.Errorf("dbutil: failed to analyze struct: %w", err)
	}

	// Create slice of pointers to scan into
	scanDests := make([]any, len(fields))
	for i, field := range fields {
		fieldValue := destElem.FieldByName(field.Name)
		if !fieldValue.CanAddr() {
			return fmt.Errorf("dbutil: field %s cannot be addressed", field.Name)
		}
		scanDests[i] = fieldValue.Addr().Interface()
	}

	// Execute query and scan
	row := queryFn()
	if err := row.Scan(scanDests...); err != nil {
		return fmt.Errorf("dbutil: query row scan failed: %w", err)
	}

	return nil
}

// structField represents a struct field with its database column information.
type structField struct {
	Name   string
	Column string
	Index  int
}

// getStructFields extracts struct fields with db tags.
func getStructFields(t reflect.Type) ([]structField, error) {
	var fields []structField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get db tag or use field name
		column := field.Tag.Get("db")
		if column == "" {
			column = DefaultFieldMapper(field.Name)
		} else if column == "-" {
			continue // Skip fields marked with db:"-"
		}

		fields = append(fields, structField{
			Name:   field.Name,
			Column: column,
			Index:  i,
		})
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("no scannable fields found")
	}

	return fields, nil
}

// IsNoRowsError checks if an error is a sql.ErrNoRows error.
func IsNoRowsError(err error) bool {
	return err == sql.ErrNoRows
}

// IsConnectionError checks if an error is a database connection error.
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}

	// Check for driver.ErrBadConn
	if err == driver.ErrBadConn {
		return true
	}

	// Check for common connection error patterns
	errStr := strings.ToLower(err.Error())
	connectionErrors := []string{
		"connection refused",
		"connection reset",
		"connection lost",
		"bad connection",
		"network is unreachable",
		"no such host",
		"timeout",
	}

	for _, connErr := range connectionErrors {
		if strings.Contains(errStr, connErr) {
			return true
		}
	}

	return false
}

// IsContextError checks if an error is a context cancellation or timeout error.
func IsContextError(err error) bool {
	if err == nil {
		return false
	}
	return err == context.Canceled || err == context.DeadlineExceeded
}
