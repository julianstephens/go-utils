# Logger Package

The `logger` package provides a unified structured logging interface for Go projects. It wraps logrus to offer log level control, custom formatting, and support for structured/contextual logging. The package is designed to be simple to import, idiomatic, and extensible.

## Features

- **Global and Instance-based Logging**: Simple global functions or custom logger instances
- **Structured Logging**: Contextual fields and structured data support
- **Dynamic Log Levels**: Thread-safe log level management at runtime
- **Custom Formatting**: JSON and text formatting with timestamps
- **Rotating File Output**: Automatic log rotation with configurable retention and compression
- **Concurrent Safe**: Full thread-safety with synchronized configuration changes

## Installation

```bash
go get github.com/julianstephens/go-utils/logger
```

## Quick Start

```go
package main

import "github.com/julianstephens/go-utils/logger"

func main() {
    // Configure
    logger.SetLogLevel("debug")
    
    // Basic logging
    logger.Infof("Server starting on port %d", 8080)
    
    // Structured logging
    logger.WithFields(map[string]interface{}{
        "user_id": "12345",
        "action":  "login",
    }).Info("User action recorded")
    
    // Set rotating file output (100MB files, 3 backups, 28 days retention)
    logger.SetFileOutput("logs/app.log")
    
    // Always close logger during shutdown to flush and close files
    defer logger.Close()
    
    // Custom instance
    customLog := logger.New()
    customLog.SetLogLevel("warn")
    customLog.Warnf("Warning: %s", "deprecated feature")
    defer customLog.Close()
}
```

## API Reference

### Global Functions

- `Infof/Debugf/Warnf/Errorf/Fatalf(format string, args ...interface{})` - Formatted logging at various levels
- `Info/Debug/Warn/Error/Fatal(args ...interface{})` - Unformatted logging
- `SetLogLevel(level string)` - Set log level ("debug", "info", "warn", "error", "fatal")
- `SetOutput(io.Writer)` - Change output destination (closes previous file if needed)
- `SetFormatter(logrus.Formatter)` - Set log formatter
- `SetFileOutput(filepath string)` - Set rotating file output with sensible defaults
- `SetFileOutputWithConfig(config FileRotationConfig)` - Set rotating file output with custom settings
- `Sync() error` - Flush pending logs to disk
- `Close() error` - Close underlying file (call during shutdown)
- `WithField(key string, value interface{}) *Logger` - Add single field
- `WithFields(fields map[string]interface{}) *Logger` - Add multiple fields
- `WithContext(ctx context.Context) *Logger` - Extract context values (trace ID, request ID, user ID) and add as fields
- `GetDefaultLogger() *Logger` - Access default logger instance

### Logger Instances

- `New() *Logger` - Create logger with default settings
- `NewWithOptions(output io.Writer, level logrus.Level, formatter logrus.Formatter) *Logger` - Create with custom options
- `SetOutput(io.Writer)` - Change output destination (closes previous file if needed)
- `SetFileOutput(filepath string) error` - Set rotating file output with sensible defaults
- `SetFileOutputWithConfig(config FileRotationConfig) error` - Set rotating file output with custom settings
- `Sync() error` - Flush pending logs to disk
- `Close() error` - Close underlying file (call during shutdown)
- `SafeLog()` - Recover from panics in logging operations (use with defer)
- `WithField(key string, value interface{}) *Logger` - Add single field
- `WithFields(fields map[string]interface{}) *Logger` - Add multiple fields
- `WithContext(ctx context.Context) *Logger` - Extract context values (trace ID, request ID, user ID) and add as fields
- Instance methods mirror global functions (Infof, SetLogLevel, etc.)

## Rotating File Output

The logger supports automatic log rotation with configurable retention and compression using [lumberjack](https://github.com/natefinch/lumberjack).

### Simple Usage (Recommended)

Use `SetFileOutput()` with sensible defaults: 100MB file size, 3 backups, 28-day retention, compression enabled.

```go
package main

import "github.com/julianstephens/go-utils/logger"

func main() {
    // Set up rotating file logs in one line
    logger.SetFileOutput("logs/app.log")
    
    // All logging now goes to rotating files
    logger.Infof("Application started")
    logger.WithFields(map[string]interface{}{
        "version": "1.0.0",
    }).Info("Service ready")
}
```

### Custom Configuration

Use `SetFileOutputWithConfig()` for advanced control over rotation parameters. Both `MaxBackups` and `MaxAge` are optional - set to `nil` to disable that constraint.

```go
// Example 1: Only use file size and backup count (no age limit)
maxBackups := 10
logger.SetFileOutputWithConfig(logger.FileRotationConfig{
    Filename:   "logs/app.log",
    MaxSize:    500,        // 500MB files
    MaxBackups: &maxBackups, // Keep 10 old log files
    MaxAge:     nil,        // No age limit
    Compress:   true,
})

// Example 2: Only use file size and retention age (no backup limit)
maxAge := 90
logger.SetFileOutputWithConfig(logger.FileRotationConfig{
    Filename:   "logs/app.log",
    MaxSize:    500,      // 500MB files
    MaxBackups: nil,      // No backup limit
    MaxAge:     &maxAge,  // Keep logs for 90 days
    Compress:   true,
})

// Example 3: Use both backup count and age limit (full control)
maxBackups := 5
maxAge := 90
logger.SetFileOutputWithConfig(logger.FileRotationConfig{
    Filename:   "logs/app.log",
    MaxSize:    500,         // 500MB files
    MaxBackups: &maxBackups, // Keep 5 old log files
    MaxAge:     &maxAge,     // Keep logs for 90 days
    Compress:   true,
})

// For custom instances
customLog := logger.New()
customLog.SetFileOutput("logs/custom.log")
```

### Rotation Behavior

- Files rotate when they reach the configured size limit
- Old logs are named with timestamps (e.g., `app.log.2025-01-15.01`)
- Logs older than `MaxAge` days are automatically deleted (if `MaxAge` is set)
- Old backups beyond `MaxBackups` count are automatically deleted (if `MaxBackups` is set)
- Old logs are compressed to `.gz` format if `Compress` is enabled
- Original `app.log` always contains the most recent logs

## Thread Safety

The logger is fully thread-safe for concurrent use. Synchronization protects:
- **Logging operations**: Thread-safe through logrus's internal mutexes
- **Configuration changes**: Protected by `sync.RWMutex` on global logger config
- **Concurrent logging and configuration**: Safe to change settings during concurrent logging

No additional synchronization required from callers.

## Best Practices

1. Use structured logging with fields for better searchability
2. Set appropriate log levels (debug in dev, info/warn in prod)
3. Include context (request IDs, user IDs, operation names)
4. Use consistent field names across your application
5. Don't log sensitive information (passwords, tokens, PII)
6. Use global logger for simplicity or instances for separate components
7. Add timing information for performance monitoring
8. Always call Close() during graceful shutdown to flush logs
9. Handle SetFileOutput() and SetLogLevel() errors appropriately

## Graceful Shutdown

Proper logger shutdown ensures all logs are written to disk before the application exits. This is critical in production environments.

### Basic Shutdown Pattern

```go
func main() {
    logger.SetLogLevel("info")
    logger.SetFileOutput("logs/app.log")
    // IMPORTANT: Always close logger on shutdown
    defer logger.Close()
    
    // ... application code ...
}
```

### Complete Shutdown Pattern (Recommended)

```go
func main() {
    logger.SetLogLevel("info")
    if err := logger.SetFileOutput("logs/app.log"); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to set log file: %v\n", err)
        os.Exit(1)
    }
    
    // Ensure graceful shutdown in all exit scenarios
    defer logger.Close()
    defer logger.Sync()  // Extra flush before close
    
    // ... application code ...
}
```

### Shutdown with Multiple Loggers

```go
func main() {
    globalLog := logger.GetDefaultLogger()
    globalLog.SetFileOutput("logs/app.log")
    
    serviceLog := logger.New()
    serviceLog.SetFileOutput("logs/service.log")
    
    // Close all loggers in order
    defer globalLog.Close()
    defer serviceLog.Close()
    
    // ... application code ...
}
```

### Shutdown in HTTP Servers

```go
func main() {
    logger.SetFileOutput("logs/server.log")
    defer logger.Close()
    defer logger.Sync()
    
    server := &http.Server{Addr: ":8080"}
    
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt)
        <-sigint
        
        logger.Info("Shutting down server")
        logger.Sync()  // Flush logs before closing
        server.Close()
    }()
    
    logger.Info("Server starting")
    server.ListenAndServe()
    logger.Info("Server stopped")
}
```

### Panic Recovery in Logging

Use `SafeLog()` to prevent logging errors from crashing your application:

```go
func criticalOperation() {
    defer logger.SafeLog()  // Recover from any panics
    
    // Even if logging panics here, execution continues
    logger.WithFields(map[string]interface{}{
        "operation": "critical",
        "status":    "running",
    }).Info("Processing critical task")
}
```

## Context Integration

Extract request-scoped values from context automatically:

```go
// In HTTP middleware
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // WithContext automatically extracts:
        // - trace-id, traceId, or traceID
        // - request-id, requestId, or requestID
        // - user-id, userId, or userID
        logger.WithContext(r.Context()).Infof("%s %s", r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

// Logging with context values
func handleRequest(ctx context.Context) {
    logger.WithContext(ctx).Info("processing request")
    // Output includes: trace_id, request_id, user_id fields
}
```

## Initialization Patterns

### Simple Pattern (Most Applications)

```go
func main() {
    logger.SetLogLevel("info")
    logger.SetFileOutput("logs/app.log")
    defer logger.Close()
    
    logger.Info("Application started")
}
```

### Structured Pattern (Microservices)

```go
func init() {
    logger.SetLogLevel(os.Getenv("LOG_LEVEL"))
    if path := os.Getenv("LOG_FILE"); path != "" {
        config := logger.FileRotationConfig{
            Filename:   path,
            MaxSize:    parseInt(os.Getenv("LOG_SIZE"), 100),
            MaxBackups: parseIntPtr(os.Getenv("LOG_BACKUPS")),
            MaxAge:     parseIntPtr(os.Getenv("LOG_AGE")),
            Compress:   parseBool(os.Getenv("LOG_COMPRESS"), true),
        }
        _ = logger.SetFileOutputWithConfig(config)
    }
}

func main() {
    defer logger.Close()
    defer logger.Sync()
    
    logger.Info("Service started")
}
```

### Component Pattern (Multiple Loggers)

```go
func newService(name string) *Service {
    log := logger.New()
    _ = log.SetLogLevel("debug")
    _ = log.SetFileOutput(fmt.Sprintf("logs/%s.log", name))
    
    return &Service{log: log}
}

func (s *Service) Shutdown() error {
    return s.log.Close()
}
```

## Integration

Logger integrates well with other go-utils packages:

```go
// With config package
if cfg.App.Debug {
    logger.SetLogLevel("debug")
}

// With cliutil package
if cliutil.HasFlag(os.Args, "--verbose") {
    logger.SetLogLevel("debug")
}
```

## Error Handling

Logger methods handle errors gracefully:
- Invalid log levels default to "info"
- Nil values are handled safely
- Circular references are detected
- Malformed format strings are logged as-is