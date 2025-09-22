# Logger Package

The `logger` package provides a unified structured logging interface for Go projects. It wraps logrus to offer log level control, custom formatting, and support for structured/contextual logging. The package is designed to be simple to import, idiomatic, and extensible.

## Features

- **Global and Instance-based Logging**: Use global functions for simplicity or create custom logger instances
- **Structured Logging**: Support for contextual fields and structured data
- **Log Level Control**: Dynamic log level management
- **Custom Formatting**: JSON and text formatting options with timestamps
- **Thread-Safe**: Safe for concurrent use across multiple goroutines
- **Contextual Logging**: Add fields and context to log entries

## Installation

```bash
go get github.com/julianstephens/go-utils/logger
```

## Usage

### Global Logger (Simple Approach)

```go
package main

import "github.com/julianstephens/go-utils/logger"

func main() {
    // Set log level
    logger.SetLogLevel("debug")
    
    // Basic logging
    logger.Infof("Server starting on port %d", 8080)
    logger.Debugf("Debug information: %s", "configuration loaded")
    logger.Warnf("Warning: %s", "deprecated feature in use")
    logger.Errorf("Error occurred: %v", err)
    
    // Structured logging with fields
    logger.WithField("user_id", "12345").Info("User logged in")
    
    logger.WithFields(map[string]interface{}{
        "user_id":    "12345",
        "action":     "login",
        "ip_address": "192.168.1.100",
    }).Info("User action recorded")
}
```

### Custom Logger Instance

```go
package main

import (
    "os"
    "github.com/julianstephens/go-utils/logger"
    "github.com/sirupsen/logrus"
)

func main() {
    // Create a custom logger instance
    customLogger := logger.New()
    customLogger.SetLogLevel("warn")
    
    customLogger.Infof("This won't appear (below warn level)")
    customLogger.Warnf("This will appear: %s", "warning message")
    
    // Logger with custom configuration and text formatting
    textLogger := logger.NewWithOptions(
        os.Stdout,
        logrus.DebugLevel,
        &logrus.TextFormatter{
            FullTimestamp: true,
            ForceColors:   true,
        },
    )
    textLogger.Infof("Using text formatter: %s", "colored output")
}
```

### Contextual Logging

```go
package main

import "github.com/julianstephens/go-utils/logger"

func main() {
    // Create a contextual logger for request processing
    requestLogger := logger.WithFields(map[string]interface{}{
        "request_id": "req-abc123",
        "method":     "POST",
        "path":       "/api/users",
    })

    requestLogger.Info("Processing request")
    requestLogger.Debugf("Request body size: %d bytes", 256)
    requestLogger.Info("Request completed successfully")
}
```

### Error Handling with Context

```go
package main

import (
    "errors"
    "github.com/julianstephens/go-utils/logger"
)

func processOrder(orderID string) error {
    logger.WithField("order_id", orderID).Info("Starting order processing")

    // Simulate some processing...
    if orderID == "invalid" {
        return errors.New("invalid order ID")
    }

    logger.WithField("order_id", orderID).Info("Order processing completed")
    return nil
}

func main() {
    orderID := "order-456"
    
    if err := processOrder(orderID); err != nil {
        logger.WithFields(map[string]interface{}{
            "order_id": orderID,
            "error":    err.Error(),
        }).Error("Failed to process order")
    }
}
```

### Application Logging Setup

```go
package main

import (
    "os"
    "github.com/julianstephens/go-utils/logger"
    "github.com/sirupsen/logrus"
)

func main() {
    // Configure logging based on environment
    if os.Getenv("DEBUG") == "true" {
        logger.SetLogLevel("debug")
    } else {
        logger.SetLogLevel("info")
    }

    // Application logger with consistent fields
    appLogger := logger.WithFields(map[string]interface{}{
        "service": "user-service",
        "version": "v1.2.3",
    })

    appLogger.Info("Application starting")

    // Simulate application work
    appLogger.Debug("Loading configuration")
    appLogger.Info("Database connection established")
    appLogger.Warn("Cache miss for user data")

    appLogger.Info("Application ready to serve requests")
}
```

### HTTP Middleware Integration

```go
package main

import (
    "net/http"
    "time"
    "github.com/julianstephens/go-utils/logger"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Create request logger
        requestLogger := logger.WithFields(map[string]interface{}{
            "method":     r.Method,
            "path":       r.URL.Path,
            "remote_ip":  r.RemoteAddr,
            "user_agent": r.UserAgent(),
        })

        requestLogger.Info("Request started")

        // Call next handler
        next.ServeHTTP(w, r)

        // Log completion
        requestLogger.WithField("duration", time.Since(start)).Info("Request completed")
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    logger.WithField("handler", "hello").Debug("Processing hello request")
    w.Write([]byte("Hello, World!"))
}

func main() {
    logger.SetLogLevel("debug")

    mux := http.NewServeMux()
    mux.HandleFunc("/hello", helloHandler)

    // Wrap with logging middleware
    loggedMux := loggingMiddleware(mux)

    logger.Info("Starting HTTP server on :8080")
    http.ListenAndServe(":8080", loggedMux)
}
```

### Database Operation Logging

```go
package main

import (
    "database/sql"
    "time"
    "github.com/julianstephens/go-utils/logger"
)

func getUserFromDB(userID string) (*User, error) {
    queryLogger := logger.WithFields(map[string]interface{}{
        "operation": "select_user",
        "user_id":   userID,
    })

    queryLogger.Debug("Executing database query")
    start := time.Now()

    // Simulate database query
    time.Sleep(50 * time.Millisecond)

    queryLogger.WithField("duration", time.Since(start)).Debug("Database query completed")

    // Simulate user not found
    if userID == "notfound" {
        queryLogger.Info("User not found")
        return nil, sql.ErrNoRows
    }

    queryLogger.Info("User retrieved successfully")
    return &User{ID: userID, Name: "John Doe"}, nil
}

type User struct {
    ID   string
    Name string
}

func main() {
    logger.SetLogLevel("debug")

    userID := "user123"
    user, err := getUserFromDB(userID)
    if err != nil {
        logger.WithFields(map[string]interface{}{
            "user_id": userID,
            "error":   err.Error(),
        }).Error("Failed to get user from database")
        return
    }

    logger.WithFields(map[string]interface{}{
        "user_id":   user.ID,
        "user_name": user.Name,
    }).Info("User operation completed")
}
```

## API Reference

### Global Functions

#### Basic Logging
- `Debugf(format string, args ...interface{})` - Debug level logging with formatting
- `Infof(format string, args ...interface{})` - Info level logging with formatting  
- `Warnf(format string, args ...interface{})` - Warning level logging with formatting
- `Errorf(format string, args ...interface{})` - Error level logging with formatting
- `Fatalf(format string, args ...interface{})` - Fatal level logging with formatting (exits)

#### Configuration
- `SetLogLevel(level string)` - Set global log level ("debug", "info", "warn", "error", "fatal")

#### Structured Logging
- `WithField(key string, value interface{}) *logrus.Entry` - Add a single field
- `WithFields(fields map[string]interface{}) *logrus.Entry` - Add multiple fields

### Logger Instance Methods

#### Creation
- `New() *Logger` - Create new logger instance with default settings
- `NewWithOptions(output io.Writer, level logrus.Level, formatter logrus.Formatter) *Logger` - Create logger with custom options

#### Instance Methods
- `SetLogLevel(level string)` - Set log level for this instance
- `Debugf/Infof/Warnf/Errorf/Fatalf(format string, args ...interface{})` - Formatted logging
- `WithField(key string, value interface{}) *logrus.Entry` - Add field to this logger
- `WithFields(fields map[string]interface{}) *logrus.Entry` - Add fields to this logger

### Log Levels

Available log levels (in order of severity):
- `debug` - Detailed information for debugging
- `info` - General information about program execution
- `warn` - Warning messages for potentially harmful situations
- `error` - Error messages for error conditions
- `fatal` - Critical errors that cause program termination

### Supported Field Types

Logger fields support any JSON-serializable type:
- Basic types: `string`, `int`, `float64`, `bool`
- Complex types: `map[string]interface{}`, `[]interface{}`
- Custom types that implement `fmt.Stringer` or are JSON-serializable

## Formatting

### JSON Formatting (Default)
```json
{
  "level": "info",
  "msg": "User logged in",
  "time": "2023-10-15T14:30:45Z",
  "user_id": "12345",
  "action": "login"
}
```

### Text Formatting
```
INFO[2023-10-15T14:30:45Z] User logged in  action=login user_id=12345
```

## Thread Safety

The logger is thread-safe and can be safely used across multiple goroutines. Both global functions and logger instances handle concurrent access properly.

## Best Practices

1. **Use structured logging** with fields instead of string formatting for better searchability
2. **Set appropriate log levels** in different environments (debug in dev, info/warn in prod)
3. **Include context** in log messages (request IDs, user IDs, operation names)
4. **Use consistent field names** across your application
5. **Don't log sensitive information** (passwords, tokens, personal data)
6. **Use the global logger for simplicity** or instances for different components
7. **Add timing information** for performance monitoring

## Integration with Other Packages

The logger integrates well with other go-utils packages:

```go
// Use with config package
if cfg.App.Debug {
    logger.SetLogLevel("debug")
} else {
    logger.SetLogLevel(cfg.App.LogLevel)
}

// Use with cliutil package
if cliutil.HasFlag(os.Args, "--verbose") {
    logger.SetLogLevel("debug")
}
```

## Error Handling

Logger methods handle errors gracefully:
- Invalid log levels default to "info"
- Malformed format strings are logged as-is
- Nil field values are handled safely
- Circular references in complex objects are detected