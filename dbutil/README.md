# Database Utilities Package

The `dbutil` package provides database utility functions and helpers for safe database interactions with connection management, query execution, transaction handling, and context support. It's designed to work with Go's standard `database/sql` package while adding convenience and safety features.

## Features

- **Connection Management**: Configuration helpers for database connections
- **Query Execution**: Safe query execution with struct scanning
- **Transaction Management**: Automatic transaction handling with rollback
- **Context Support**: All operations are context-aware for cancellation and timeouts
- **Error Handling**: Enhanced error detection and classification
- **Struct Scanning**: Automatic scanning of query results into structs
- **Field Mapping**: Customizable mapping between struct fields and database columns

## Installation

```bash
go get github.com/julianstephens/go-utils/dbutil
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/dbutil"
    _ "github.com/lib/pq" // postgres driver
)

func main() {
    // Connect to database
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/mydb?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Configure database connection
    opts := dbutil.DefaultConnectionOptions()
    opts.MaxOpenConns = 25
    opts.MaxIdleConns = 5
    opts.ConnMaxLifetime = time.Hour
    
    if err := dbutil.ConfigureDB(db, opts); err != nil {
        log.Fatal(err)
    }

    // Test connection
    ctx := context.Background()
    if err := dbutil.PingWithTimeout(ctx, db, 5*time.Second); err != nil {
        log.Fatal("Database ping failed:", err)
    }

    log.Println("Database connected successfully")
}
```

### Struct Scanning

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/dbutil"
)

type User struct {
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
}

func main() {
    db := setupDatabase() // your database setup
    ctx := context.Background()

    // Query single row into struct
    var user User
    err := dbutil.QueryRowScan(ctx, db, &user, 
        "SELECT id, name, email, created_at FROM users WHERE id = $1", 1)
    if err != nil {
        if dbutil.IsNoRowsError(err) {
            log.Println("User not found")
        } else {
            log.Printf("Query failed: %v", err)
        }
        return
    }

    log.Printf("User: %+v", user)

    // Query multiple rows into slice
    var users []User
    err = dbutil.QuerySlice(ctx, db, &users,
        "SELECT id, name, email, created_at FROM users WHERE active = $1", true)
    if err != nil {
        log.Printf("Query failed: %v", err)
        return
    }

    log.Printf("Found %d users", len(users))
    for _, u := range users {
        log.Printf("User: %s (%s)", u.Name, u.Email)
    }
}
```

### Transaction Management

```go
package main

import (
    "context"
    "database/sql"
    "log"
    
    "github.com/julianstephens/go-utils/dbutil"
)

func main() {
    db := setupDatabase()
    ctx := context.Background()

    // Execute within transaction with automatic rollback on error
    err := dbutil.WithTransaction(ctx, db, func(tx *sql.Tx) error {
        // Insert user
        result, err := dbutil.ExecTx(ctx, tx, 
            "INSERT INTO users (name, email) VALUES ($1, $2)", 
            "John Doe", "john@example.com")
        if err != nil {
            return err
        }

        userID, err := result.LastInsertId()
        if err != nil {
            return err
        }

        // Insert audit log
        _, err = dbutil.ExecTx(ctx, tx,
            "INSERT INTO audit_log (action, user_id) VALUES ($1, $2)",
            "user_created", userID)
        if err != nil {
            return err
        }

        log.Printf("User created with ID: %d", userID)
        return nil
    })

    if err != nil {
        log.Printf("Transaction failed: %v", err)
    } else {
        log.Println("Transaction completed successfully")
    }
}
```

### Advanced Transaction Options

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/dbutil"
)

func main() {
    db := setupDatabase()
    ctx := context.Background()

    // Transaction with custom options
    txOpts := dbutil.TransactionOptions{
        Isolation: sql.LevelReadCommitted,
        ReadOnly:  false,
        Timeout:   30 * time.Second,
    }

    err := dbutil.WithTransactionOptions(ctx, db, txOpts, func(tx *sql.Tx) error {
        // Perform read operations
        var count int
        err := dbutil.QueryRowScanTx(ctx, tx, &count,
            "SELECT COUNT(*) FROM users WHERE active = $1", true)
        if err != nil {
            return err
        }

        log.Printf("Active users: %d", count)

        // Perform write operations if needed
        if count < 100 {
            _, err = dbutil.ExecTx(ctx, tx,
                "INSERT INTO users (name, email) VALUES ($1, $2)",
                "New User", "newuser@example.com")
            return err
        }

        return nil
    })

    if err != nil {
        log.Printf("Transaction failed: %v", err)
    }
}
```

### Utility Operations

```go
package main

import (
    "context"
    "database/sql"
    "log"
    
    "github.com/julianstephens/go-utils/dbutil"
)

func main() {
    db := setupDatabase()
    ctx := context.Background()

    // Check if record exists
    exists, err := dbutil.Exists(ctx, db, 
        "SELECT 1 FROM users WHERE email = $1", "john@example.com")
    if err != nil {
        log.Printf("Exists check failed: %v", err)
        return
    }
    log.Printf("User exists: %t", exists)

    // Count records
    count, err := dbutil.Count(ctx, db, 
        "SELECT COUNT(*) FROM users WHERE active = $1", true)
    if err != nil {
        log.Printf("Count failed: %v", err)
        return
    }
    log.Printf("Active users: %d", count)

    // Get single value
    var maxID int64
    err = dbutil.QueryRowScan(ctx, db, &maxID,
        "SELECT MAX(id) FROM users")
    if err != nil && !dbutil.IsNoRowsError(err) {
        log.Printf("Max ID query failed: %v", err)
        return
    }
    log.Printf("Max user ID: %d", maxID)
}
```

### Error Handling

```go
package main

import (
    "context"
    "database/sql"
    "errors"
    "log"
    
    "github.com/julianstephens/go-utils/dbutil"
)

func getUserByID(ctx context.Context, db *sql.DB, userID int64) (*User, error) {
    var user User
    err := dbutil.QueryRowScan(ctx, db, &user,
        "SELECT id, name, email, created_at FROM users WHERE id = $1", userID)
    
    if err != nil {
        if dbutil.IsNoRowsError(err) {
            return nil, errors.New("user not found")
        }
        if dbutil.IsContextError(err) {
            return nil, errors.New("request timeout or cancelled")
        }
        if dbutil.IsConnectionError(err) {
            return nil, errors.New("database connection failed")
        }
        return nil, err
    }

    return &user, nil
}

func main() {
    db := setupDatabase()
    ctx := context.Background()

    user, err := getUserByID(ctx, db, 123)
    if err != nil {
        log.Printf("Error getting user: %v", err)
        return
    }

    log.Printf("Found user: %+v", user)
}
```

### Field Mapping

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "strings"
    
    "github.com/julianstephens/go-utils/dbutil"
)

type UserProfile struct {
    UserID    int64  `db:"user_id"`
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
    XMLData   string `db:"xml_data"`
}

func main() {
    db := setupDatabase()
    ctx := context.Background()

    // Use custom field mapper for struct fields to database columns
    queryOpts := dbutil.QueryOptions{
        Timeout: 30 * time.Second,
        MaxRows: 100,
        FieldMapper: func(fieldName string) string {
            // Convert CamelCase to snake_case
            return dbutil.DefaultFieldMapper(fieldName)
        },
    }

    var profiles []UserProfile
    err := dbutil.QuerySliceWithOptions(ctx, db, &profiles, queryOpts,
        "SELECT user_id, first_name, last_name, xml_data FROM user_profiles")
    if err != nil {
        log.Printf("Query failed: %v", err)
        return
    }

    log.Printf("Found %d profiles", len(profiles))
}
```

## Configuration Options

### ConnectionOptions

```go
type ConnectionOptions struct {
    MaxOpenConns    int           // Maximum open connections
    MaxIdleConns    int           // Maximum idle connections  
    ConnMaxLifetime time.Duration // Connection lifetime
    ConnMaxIdleTime time.Duration // Connection idle time
    PingTimeout     time.Duration // Ping timeout
    RetryAttempts   int           // Retry attempts
    RetryDelay      time.Duration // Retry delay
}

// Get default options
opts := dbutil.DefaultConnectionOptions()
```

### QueryOptions

```go
type QueryOptions struct {
    Timeout     time.Duration           // Query timeout
    MaxRows     int                     // Maximum rows (0 = no limit)
    FieldMapper func(string) string     // Field name mapper
}

// Get default options
opts := dbutil.DefaultQueryOptions()
```

### TransactionOptions

```go
type TransactionOptions struct {
    Isolation sql.IsolationLevel // Transaction isolation level
    ReadOnly  bool               // Read-only transaction
    Timeout   time.Duration      // Transaction timeout
}

// Get default options
opts := dbutil.DefaultTransactionOptions()
```

## API Reference

### Connection Management
- `ConfigureDB(db *sql.DB, opts ConnectionOptions) error` - Configure database connection
- `PingWithTimeout(ctx context.Context, db *sql.DB, timeout time.Duration) error` - Ping with timeout
- `DefaultConnectionOptions() ConnectionOptions` - Get default connection options

### Query Execution
- `QueryRowScan(ctx context.Context, db *sql.DB, dest interface{}, query string, args ...interface{}) error` - Query single row
- `QuerySlice(ctx context.Context, db *sql.DB, dest interface{}, query string, args ...interface{}) error` - Query multiple rows
- `QuerySliceWithOptions(ctx context.Context, db *sql.DB, dest interface{}, opts QueryOptions, query string, args ...interface{}) error` - Query with options
- `Exec(ctx context.Context, db *sql.DB, query string, args ...interface{}) (sql.Result, error)` - Execute query

### Transaction Management
- `WithTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error` - Execute in transaction
- `WithTransactionOptions(ctx context.Context, db *sql.DB, opts TransactionOptions, fn func(*sql.Tx) error) error` - Execute with options
- `QueryRowScanTx(ctx context.Context, tx *sql.Tx, dest interface{}, query string, args ...interface{}) error` - Query in transaction
- `ExecTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)` - Execute in transaction

### Utility Functions
- `Exists(ctx context.Context, db *sql.DB, query string, args ...interface{}) (bool, error)` - Check if record exists
- `Count(ctx context.Context, db *sql.DB, query string, args ...interface{}) (int64, error)` - Count records

### Error Detection
- `IsNoRowsError(err error) bool` - Check if error is sql.ErrNoRows
- `IsContextError(err error) bool` - Check if error is context-related
- `IsConnectionError(err error) bool` - Check if error is connection-related

### Field Mapping
- `DefaultFieldMapper(fieldName string) string` - Default CamelCase to snake_case mapper

## Supported Struct Tags

Use the `db` tag to specify database column names:

```go
type User struct {
    ID       int64  `db:"id"`
    UserName string `db:"user_name"`
    Email    string `db:"email_address"`
}
```

## Thread Safety

All functions in the dbutil package are thread-safe and can be called concurrently from multiple goroutines. The package properly handles the underlying database/sql thread safety guarantees.

## Best Practices

1. **Always use context** for cancellation and timeouts
2. **Use transactions** for operations that need atomicity
3. **Handle specific error types** using the provided error detection functions
4. **Configure connection pooling** appropriately for your workload
5. **Use struct tags** to clearly map between Go fields and database columns
6. **Set appropriate timeouts** for different types of operations
7. **Consider using read-only transactions** for complex read operations

## Database Driver Compatibility

The dbutil package works with any database driver that implements Go's `database/sql` interface, including:

- PostgreSQL (`github.com/lib/pq`)
- MySQL (`github.com/go-sql-driver/mysql`)
- SQLite (`github.com/mattn/go-sqlite3`)
- SQL Server (`github.com/denisenkom/go-mssqldb`)

## Integration

Works well with other go-utils packages:

```go
// Use with logger for database operation logging
logger.WithField("query", "SELECT * FROM users").Debug("Executing query")

// Use with config for database configuration
dbURL := cfg.Database.URL
db, err := sql.Open("postgres", dbURL)
```