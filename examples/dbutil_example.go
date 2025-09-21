package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/julianstephens/go-utils/dbutil"
	
	// Example with postgres driver (commented out to avoid dependency)
	// _ "github.com/lib/pq"
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

func runDBUtilExample() {
	fmt.Println("=== DBUtil Package Example ===\n")

	// This is a demonstration of how to use the dbutil package
	// In a real application, you would have an actual database connection
	fmt.Println("1. Database Connection Configuration:")
	
	// Configure database connection options
	opts := &dbutil.ConnectionOptions{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		PingTimeout:     5 * time.Second,
		RetryAttempts:   3,
		RetryDelay:      time.Second,
	}
	
	fmt.Printf("  Max Open Connections: %d\n", opts.MaxOpenConns)
	fmt.Printf("  Max Idle Connections: %d\n", opts.MaxIdleConns)
	fmt.Printf("  Connection Max Lifetime: %v\n", opts.ConnMaxLifetime)
	fmt.Printf("  Ping Timeout: %v\n", opts.PingTimeout)
	fmt.Printf("  Retry Attempts: %d\n", opts.RetryAttempts)

	fmt.Println("\n2. Query Options:")
	
	queryOpts := dbutil.DefaultQueryOptions()
	fmt.Printf("  Default Timeout: %v\n", queryOpts.Timeout)
	fmt.Printf("  Max Rows (0 = unlimited): %d\n", queryOpts.MaxRows)
	
	// Demonstrate field mapping
	fmt.Println("\n3. Field Mapping Examples:")
	fields := []string{"ID", "UserID", "FirstName", "CreatedAt", "XMLHttpRequest"}
	for _, field := range fields {
		mapped := dbutil.DefaultFieldMapper(field)
		fmt.Printf("  %s -> %s\n", field, mapped)
	}

	fmt.Println("\n4. Transaction Options:")
	
	txOpts := dbutil.DefaultTransactionOptions()
	fmt.Printf("  Isolation Level: %v\n", txOpts.Isolation)
	fmt.Printf("  Read Only: %t\n", txOpts.ReadOnly)
	fmt.Printf("  Timeout: %v\n", txOpts.Timeout)

	fmt.Println("\n5. Error Detection Examples:")
	
	// Test error detection functions
	testErrors := []error{
		sql.ErrNoRows,
		context.Canceled,
		context.DeadlineExceeded,
		fmt.Errorf("connection refused"),
		fmt.Errorf("timeout occurred"),
		fmt.Errorf("some other error"),
	}
	
	for _, err := range testErrors {
		fmt.Printf("  Error: %v\n", err)
		fmt.Printf("    IsNoRowsError: %t\n", dbutil.IsNoRowsError(err))
		fmt.Printf("    IsContextError: %t\n", dbutil.IsContextError(err))
		fmt.Printf("    IsConnectionError: %t\n", dbutil.IsConnectionError(err))
		fmt.Println()
	}

	// Example usage patterns (commented out since we don't have a real database)
	fmt.Println("6. Example Usage Patterns:")
	fmt.Println(`
  // Connect to database
  db, err := sql.Open("postgres", "postgres://user:pass@localhost/db?sslmode=disable")
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()
  
  // Configure connection
  err = dbutil.ConfigureDB(db, opts)
  if err != nil {
      log.Fatal(err)
  }
  
  ctx := context.Background()
  
  // Query single row into struct
  var user User
  err = dbutil.QueryRowScan(ctx, db, &user, 
      "SELECT id, name, email, created_at FROM users WHERE id = $1", 1)
  if err != nil {
      log.Printf("Query failed: %v", err)
  }
  
  // Query multiple rows into slice
  var users []User
  err = dbutil.QuerySlice(ctx, db, &users,
      "SELECT id, name, email, created_at FROM users WHERE active = $1", true)
  if err != nil {
      log.Printf("Query failed: %v", err)
  }
  
  // Execute within transaction
  err = dbutil.WithTransaction(ctx, db, func(tx *sql.Tx) error {
      _, err := dbutil.ExecTx(ctx, tx, 
          "INSERT INTO users (name, email) VALUES ($1, $2)", 
          "John Doe", "john@example.com")
      if err != nil {
          return err
      }
      
      _, err = dbutil.ExecTx(ctx, tx,
          "INSERT INTO audit_log (action, user_name) VALUES ($1, $2)",
          "user_created", "John Doe")
      return err
  })
  if err != nil {
      log.Printf("Transaction failed: %v", err)
  }
  
  // Check if record exists
  exists, err := dbutil.Exists(ctx, db, 
      "SELECT 1 FROM users WHERE email = $1", "john@example.com")
  if err != nil {
      log.Printf("Exists check failed: %v", err)
  }
  fmt.Printf("User exists: %t\n", exists)
  
  // Count records
  count, err := dbutil.Count(ctx, db, "SELECT COUNT(*) FROM users WHERE active = $1", true)
  if err != nil {
      log.Printf("Count failed: %v", err)
  }
  fmt.Printf("Active users: %d\n", count)`)

	fmt.Println("\nThe dbutil package provides safe, idiomatic database operations with:")
	fmt.Println("  • Enhanced error context and handling")
	fmt.Println("  • Connection management utilities")
	fmt.Println("  • Context-aware operations")
	fmt.Println("  • Transaction management with automatic rollback")
	fmt.Println("  • Struct scanning helpers")
	fmt.Println("  • Query execution utilities")
}

func main() {
	runDBUtilExample()
}