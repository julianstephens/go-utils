# Database Utilities Package

The `dbutil` package provides database utility functions and helpers for safe database interactions with connection management, query execution, transaction handling, and context support. It's designed to work with Go's standard `database/sql` package while adding convenience and safety features.

## Features

- **Connection Management**: Configuration and connection pooling
- **Query Execution**: Safe query execution with struct scanning
- **Transaction Management**: Automatic transaction handling with rollback
- **Context Support**: Cancellation and timeout support
- **Error Handling**: Enhanced error detection and classification
- **Struct Scanning**: Automatic scanning into structs
- **Field Mapping**: Customizable struct-to-column mapping

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
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    opts := dbutil.DefaultConnectionOptions()
    opts.MaxOpenConns = 25
    opts.MaxIdleConns = 5
    opts.ConnMaxLifetime = time.Hour
    
    if err := dbutil.ConfigureDB(db, opts); err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    if err := dbutil.PingWithContext(ctx, db, 5*time.Second); err != nil {
        log.Fatal("Database ping failed:", err)
    }
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
    db := setupDatabase()
    ctx := context.Background()

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

    var users []User
    err = dbutil.QuerySlice(ctx, db, &users,
        "SELECT id, name, email, created_at FROM users WHERE active = $1", true)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, u := range users {
        log.Printf("%s (%s)", u.Name, u.Email)
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

    err := dbutil.WithTransaction(ctx, db, func(tx *sql.Tx) error {
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

        _, err = dbutil.ExecTx(ctx, tx,
            "INSERT INTO audit_log (action, user_id) VALUES ($1, $2)",
            "user_created", userID)
        return err
    })

    if err != nil {
        log.Printf("Transaction failed: %v", err)
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

    txOpts := dbutil.TransactionOptions{
        Isolation: sql.LevelReadCommitted,
        ReadOnly:  false,
        Timeout:   30 * time.Second,
    }

    err := dbutil.WithTransactionOptions(ctx, db, txOpts, func(tx *sql.Tx) error {
        var count int
        err := dbutil.QueryRowScanTx(ctx, tx, &count,
            "SELECT COUNT(*) FROM users WHERE active = $1", true)
        if err != nil {
            return err
        }

        if count < 100 {
            _, err = dbutil.ExecTx(ctx, tx,
                "INSERT INTO users (name, email) VALUES ($1, $2)",
                "New User", "newuser@example.com")
        }
        return err
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

    exists, err := dbutil.Exists(ctx, db, 
        "SELECT 1 FROM users WHERE email = $1", "john@example.com")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("User exists: %t", exists)

    count, err := dbutil.Count(ctx, db, 
        "SELECT COUNT(*) FROM users WHERE active = $1", true)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Active users: %d", count)
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
```

### Field Mapping

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/julianstephens/go-utils/dbutil"
)

type UserProfile struct {
    UserID    int64  `db:"user_id"`
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
}

func main() {
    db := setupDatabase()
    ctx := context.Background()

    queryOpts := dbutil.QueryOptions{
        Timeout: 30 * time.Second,
        MaxRows: 100,
    }

    var profiles []UserProfile
    err := dbutil.QuerySliceWithOptions(ctx, db, &profiles, queryOpts,
        "SELECT user_id, first_name, last_name FROM user_profiles")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Options

### ConnectionOptions
- `MaxOpenConns` - Maximum open connections
- `MaxIdleConns` - Maximum idle connections  
- `ConnMaxLifetime` - Connection lifetime
- `ConnMaxIdleTime` - Connection idle time
- `PingTimeout` - Ping timeout
- `RetryAttempts` - Retry attempts
- `RetryDelay` - Retry delay

### QueryOptions
- `Timeout` - Query timeout
- `MaxRows` - Maximum rows (0 = no limit)
- `FieldMapper` - Field name mapper function

### TransactionOptions
- `Isolation` - Transaction isolation level
- `ReadOnly` - Read-only transaction flag
- `Timeout` - Transaction timeout

## API Reference

### Connection Management
- `ConfigureDB(db *sql.DB, opts ConnectionOptions) error` - Configure database
- `PingWithContext(ctx context.Context, db *sql.DB, timeout time.Duration) error` - Ping with timeout
- `PingWithRetry(ctx context.Context, db *sql.DB, attempts int, delay time.Duration) error` - Ping with retry

### Query Execution
- `QueryRowScan(ctx, db, dest, query, args...) error` - Query single row into struct
- `QuerySlice(ctx, db, dest, query, args...) error` - Query multiple rows into slice
- `QuerySliceWithOptions(ctx, db, dest, opts, query, args...) error` - Query with options
- `QueryMap(ctx, db, query, args...) (map[string]any, error)` - Query row into map
- `QueryMaps(ctx, db, query, args...) ([]map[string]any, error)` - Query rows into maps
- `QueryRow(ctx, db, query, args...) *sql.Row` - Raw single row
- `QueryRows(ctx, db, query, args...) (*sql.Rows, error)` - Raw multiple rows
- `Exec(ctx, db, query, args...) (sql.Result, error)` - Execute query

### Transaction Management
- `WithTransaction(ctx, db, fn) error` - Execute in transaction
- `WithTransactionOptions(ctx, db, opts, fn) error` - Execute with options
- `QueryRowScanTx(ctx, tx, dest, query, args...) error` - Query single row in tx
- `QuerySliceTx(ctx, tx, dest, query, args...) error` - Query slice in tx
- `QueryMapTx(ctx, tx, query, args...) (map[string]any, error)` - Query map in tx
- `QueryMapsTx(ctx, tx, query, args...) ([]map[string]any, error)` - Query maps in tx
- `QueryRowTx(ctx, tx, query, args...) *sql.Row` - Raw row in tx
- `QueryRowsTx(ctx, tx, query, args...) (*sql.Rows, error)` - Raw rows in tx
- `ExecTx(ctx, tx, query, args...) (sql.Result, error)` - Execute in tx

### Utility Functions
- `Exists(ctx, db, query, args...) (bool, error)` - Check if record exists
- `ExistsTx(ctx, tx, query, args...) (bool, error)` - Check existence in tx
- `Count(ctx, db, query, args...) (int64, error)` - Count records
- `CountTx(ctx, tx, query, args...) (int64, error)` - Count in tx

### Error Detection
- `IsNoRowsError(err) bool` - Check for sql.ErrNoRows
- `IsConnectionError(err) bool` - Check for connection errors
- `IsContextError(err) bool` - Check for context timeout/cancel

### Field Mapping
- `DefaultFieldMapper(fieldName) string` - CamelCase to snake_case mapper

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
2. **Use transactions** for operations requiring atomicity
3. **Handle specific error types** using provided detection functions
4. **Configure connection pooling** appropriately for your workload
5. **Use struct tags** to map Go fields to database columns
6. **Set appropriate timeouts** for different operation types
7. **Consider read-only transactions** for complex read operations

## Database Driver Compatibility

Works with any database driver implementing Go's `database/sql` interface:
- PostgreSQL (`github.com/lib/pq`)
- MySQL (`github.com/go-sql-driver/mysql`)
- SQLite (`github.com/mattn/go-sqlite3`)
- SQL Server (`github.com/denisenkom/go-mssqldb`)

## Integration

Works well with other go-utils packages:
- **logger**: Log database operations
- **config**: Manage database configuration