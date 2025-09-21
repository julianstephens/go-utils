package dbutil

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

// QuerySlice executes a query and scans all rows into a slice of structs.
// dest should be a pointer to a slice of structs with appropriate db tags.
func QuerySlice(ctx context.Context, db *sql.DB, dest interface{}, query string, args ...interface{}) error {
	return QuerySliceWithOptions(ctx, db, dest, query, DefaultQueryOptions(), args...)
}

// QuerySliceTx is like QuerySlice but uses a transaction.
func QuerySliceTx(ctx context.Context, tx *sql.Tx, dest interface{}, query string, args ...interface{}) error {
	return QuerySliceWithOptionsTx(ctx, tx, dest, query, DefaultQueryOptions(), args...)
}

// QuerySliceWithOptions executes a query and scans all rows into a slice of structs with options.
func QuerySliceWithOptions(ctx context.Context, db *sql.DB, dest interface{}, query string, opts *QueryOptions, args ...interface{}) error {
	return querySliceImpl(ctx, func() (*sql.Rows, error) {
		return db.QueryContext(ctx, query, args...)
	}, dest, opts)
}

// QuerySliceWithOptionsTx is like QuerySliceWithOptions but uses a transaction.
func QuerySliceWithOptionsTx(ctx context.Context, tx *sql.Tx, dest interface{}, query string, opts *QueryOptions, args ...interface{}) error {
	return querySliceImpl(ctx, func() (*sql.Rows, error) {
		return tx.QueryContext(ctx, query, args...)
	}, dest, opts)
}

// querySliceImpl is the common implementation for QuerySlice functions.
func querySliceImpl(ctx context.Context, queryFn func() (*sql.Rows, error), dest interface{}, opts *QueryOptions) error {
	if opts == nil {
		opts = DefaultQueryOptions()
	}

	// Validate dest parameter
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return fmt.Errorf("dbutil: dest must be a non-nil pointer")
	}

	destElem := destValue.Elem()
	if destElem.Kind() != reflect.Slice {
		return fmt.Errorf("dbutil: dest must be a pointer to a slice")
	}

	sliceType := destElem.Type()
	elementType := sliceType.Elem()
	
	// Handle pointer to struct
	isPtr := false
	if elementType.Kind() == reflect.Ptr {
		isPtr = true
		elementType = elementType.Elem()
	}
	
	if elementType.Kind() != reflect.Struct {
		return fmt.Errorf("dbutil: dest must be a pointer to a slice of structs or pointers to structs")
	}

	// Get struct fields
	fields, err := getStructFields(elementType)
	if err != nil {
		return fmt.Errorf("dbutil: failed to analyze struct: %w", err)
	}

	// Execute query (timeout handling is done at caller level through ctx)
	rows, err := queryFn()
	if err != nil {
		return fmt.Errorf("dbutil: query slice failed: %w", err)
	}
	defer rows.Close()

	// Create new slice to hold results
	result := reflect.MakeSlice(sliceType, 0, 0)
	rowCount := 0

	for rows.Next() {
		// Check row limit
		if opts.MaxRows > 0 && rowCount >= opts.MaxRows {
			break
		}

		// Create new struct instance
		var elemValue reflect.Value
		if isPtr {
			elemValue = reflect.New(elementType)
		} else {
			elemValue = reflect.New(elementType).Elem()
		}

		// Create scan destinations
		scanDests := make([]interface{}, len(fields))
		for i, field := range fields {
			var fieldValue reflect.Value
			if isPtr {
				fieldValue = elemValue.Elem().FieldByName(field.Name)
			} else {
				fieldValue = elemValue.FieldByName(field.Name)
			}
			
			if !fieldValue.CanAddr() {
				return fmt.Errorf("dbutil: field %s cannot be addressed", field.Name)
			}
			scanDests[i] = fieldValue.Addr().Interface()
		}

		// Scan row into struct
		if err := rows.Scan(scanDests...); err != nil {
			return fmt.Errorf("dbutil: scan row failed: %w", err)
		}

		// Append to result slice
		if isPtr {
			result = reflect.Append(result, elemValue)
		} else {
			result = reflect.Append(result, elemValue)
		}
		rowCount++
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return fmt.Errorf("dbutil: row iteration failed: %w", err)
	}

	// Set the result
	destElem.Set(result)
	return nil
}

// QueryMap executes a query and returns the first row as a map[string]interface{}.
func QueryMap(ctx context.Context, db *sql.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	return queryMapImpl(ctx, func() (*sql.Rows, error) {
		return db.QueryContext(ctx, query, args...)
	})
}

// QueryMapTx is like QueryMap but uses a transaction.
func QueryMapTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (map[string]interface{}, error) {
	return queryMapImpl(ctx, func() (*sql.Rows, error) {
		return tx.QueryContext(ctx, query, args...)
	})
}

// QueryMaps executes a query and returns all rows as []map[string]interface{}.
func QueryMaps(ctx context.Context, db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return queryMapsImpl(ctx, func() (*sql.Rows, error) {
		return db.QueryContext(ctx, query, args...)
	})
}

// QueryMapsTx is like QueryMaps but uses a transaction.
func QueryMapsTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return queryMapsImpl(ctx, func() (*sql.Rows, error) {
		return tx.QueryContext(ctx, query, args...)
	})
}

// queryMapImpl is the common implementation for QueryMap functions.
func queryMapImpl(ctx context.Context, queryFn func() (*sql.Rows, error)) (map[string]interface{}, error) {
	rows, err := queryFn()
	if err != nil {
		return nil, fmt.Errorf("dbutil: query map failed: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("dbutil: query map failed: %w", err)
		}
		return nil, sql.ErrNoRows
	}

	return scanRowToMap(rows)
}

// queryMapsImpl is the common implementation for QueryMaps functions.
func queryMapsImpl(ctx context.Context, queryFn func() (*sql.Rows, error)) ([]map[string]interface{}, error) {
	rows, err := queryFn()
	if err != nil {
		return nil, fmt.Errorf("dbutil: query maps failed: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		rowMap, err := scanRowToMap(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dbutil: query maps iteration failed: %w", err)
	}

	return result, nil
}

// scanRowToMap scans a single row into a map[string]interface{}.
func scanRowToMap(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("dbutil: get columns failed: %w", err)
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("dbutil: get column types failed: %w", err)
	}

	// Create interface{} slice for scanning
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Scan the row
	if err := rows.Scan(scanArgs...); err != nil {
		return nil, fmt.Errorf("dbutil: scan row to map failed: %w", err)
	}

	// Convert to map, handling NULL values
	result := make(map[string]interface{})
	for i, column := range columns {
		value := values[i]
		
		// Handle NULL values
		if value == nil {
			result[column] = nil
			continue
		}

		// Convert byte slices to strings for text columns using ScanType
		if b, ok := value.([]byte); ok {
			if columnTypes[i] != nil {
				if columnTypes[i].ScanType() == reflect.TypeOf("") {
					result[column] = string(b)
					continue
				}
			}
		}

		result[column] = value
	}

	return result, nil
}

// isTextColumn checks if a database type should be treated as text.
func isTextColumn(dbType string) bool {
	textTypes := []string{
		"VARCHAR", "TEXT", "CHAR", "NVARCHAR", "NTEXT", "NCHAR",
		"STRING", "CLOB", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT",
	}
	
	for _, textType := range textTypes {
		if dbType == textType {
			return true
		}
	}
	
	return false
}

// Exists checks if a query returns any rows.
func Exists(ctx context.Context, db *sql.DB, query string, args ...interface{}) (bool, error) {
	return existsImpl(ctx, func() (*sql.Rows, error) {
		return db.QueryContext(ctx, query, args...)
	})
}

// ExistsTx is like Exists but uses a transaction.
func ExistsTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (bool, error) {
	return existsImpl(ctx, func() (*sql.Rows, error) {
		return tx.QueryContext(ctx, query, args...)
	})
}

// existsImpl is the common implementation for Exists functions.
func existsImpl(ctx context.Context, queryFn func() (*sql.Rows, error)) (bool, error) {
	rows, err := queryFn()
	if err != nil {
		return false, fmt.Errorf("dbutil: exists query failed: %w", err)
	}
	defer rows.Close()

	exists := rows.Next()
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("dbutil: exists check failed: %w", err)
	}

	return exists, nil
}

// Count executes a COUNT query and returns the result.
func Count(ctx context.Context, db *sql.DB, query string, args ...interface{}) (int64, error) {
	return countImpl(ctx, func() *sql.Row {
		return db.QueryRowContext(ctx, query, args...)
	})
}

// CountTx is like Count but uses a transaction.
func CountTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (int64, error) {
	return countImpl(ctx, func() *sql.Row {
		return tx.QueryRowContext(ctx, query, args...)
	})
}

// countImpl is the common implementation for Count functions.
func countImpl(ctx context.Context, queryFn func() *sql.Row) (int64, error) {
	var count int64
	err := queryFn().Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("dbutil: count query failed: %w", err)
	}
	return count, nil
}