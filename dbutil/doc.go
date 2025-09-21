/*
Package dbutil provides utility functions and helpers for interacting with databases in Go projects.

This package offers enhanced database/sql functionality with features like:
- Connection management with retry logic and context support
- Query execution helpers with automatic scanning and error context
- Transaction management utilities with rollback handling
- Context-aware operations for cancellation and timeout support
- Comprehensive error handling with meaningful context

Basic Usage:

The dbutil package provides both simple helper functions and advanced configuration options:

	package main

	import (
		"context"
		"database/sql"
		"log"

		"github.com/julianstephens/go-utils/dbutil"
		_ "github.com/lib/pq" // postgres driver
	)

	type User struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	func main() {
		// Open database connection
		db, err := sql.Open("postgres", "postgres://user:pass@localhost/db?sslmode=disable")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		ctx := context.Background()

		// Simple query execution
		var user User
		err = dbutil.QueryRow(ctx, db, "SELECT id, name, email FROM users WHERE id = $1", 1).
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}

		// Query multiple rows
		users, err := dbutil.QueryRows[User](ctx, db, "SELECT id, name, email FROM users")
		if err != nil {
			log.Fatalf("Query rows failed: %v", err)
		}

		// Execute with transaction
		err = dbutil.WithTransaction(ctx, db, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John", "john@example.com")
			return err
		})
		if err != nil {
			log.Fatalf("Transaction failed: %v", err)
		}
	}

Connection Management:

The package provides utilities for managing database connections with health checks and retries:

	// Connection options
	opts := &dbutil.ConnectionOptions{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 30,
		PingTimeout:     time.Second * 5,
		RetryAttempts:   3,
		RetryDelay:      time.Second,
	}

	// Configure database connection
	err := dbutil.ConfigureDB(db, opts)
	if err != nil {
		log.Fatalf("Failed to configure database: %v", err)
	}

	// Health check with context
	err = dbutil.PingWithContext(ctx, db, time.Second*10)
	if err != nil {
		log.Printf("Database health check failed: %v", err)
	}

Query Execution:

Enhanced query execution with automatic scanning and error context:

	// Query single row with automatic scanning
	var user User
	err := dbutil.QueryRowScan(ctx, db, &user, "SELECT id, name, email FROM users WHERE id = $1", 1)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}

	// Query multiple rows with automatic scanning
	var users []User
	err = dbutil.QuerySlice(ctx, db, &users, "SELECT id, name, email FROM users WHERE active = $1", true)
	if err != nil {
		log.Printf("Query failed: %v", err)
	}

	// Execute query with custom options
	opts := &dbutil.QueryOptions{
		Timeout:     time.Second * 30,
		MaxRows:     1000,
		FieldMapper: dbutil.DefaultFieldMapper,
	}
	err = dbutil.QuerySliceWithOptions(ctx, db, &users, "SELECT * FROM users", opts)

Transaction Management:

Safe transaction handling with automatic rollback on errors:

	// Simple transaction wrapper
	err := dbutil.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Multiple operations in transaction
		if _, err := tx.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice"); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO audit_log (action) VALUES ($1)", "user_created"); err != nil {
			return err
		}
		return nil
	})

	// Transaction with options
	txOpts := &dbutil.TransactionOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Timeout:   time.Minute,
	}
	err = dbutil.WithTransactionOptions(ctx, db, txOpts, func(tx *sql.Tx) error {
		// Transaction logic here
		return nil
	})

Error Handling:

All functions provide enhanced error context to help with debugging:

	// Instead of: "no rows in result set"
	// You get: "dbutil: query row failed: no rows in result set"

	// Database connection errors include context:
	// "dbutil: connection failed after 3 attempts: connection refused"

The package is designed for reuse across Go projects and provides consistent and safe
database access patterns with comprehensive error handling and context support.
*/
package dbutil