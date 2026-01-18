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
    
    // Custom instance
    customLog := logger.New()
    customLog.SetLogLevel("warn")
    customLog.Warnf("Warning: %s", "deprecated feature")
}
```

## API Reference

### Global Functions

- `Infof/Debugf/Warnf/Errorf/Fatalf(format string, args ...interface{})` - Formatted logging at various levels
- `Info/Debug/Warn/Error/Fatal(args ...interface{})` - Unformatted logging
- `SetLogLevel(level string)` - Set log level ("debug", "info", "warn", "error", "fatal")
- `SetOutput(io.Writer)` - Change output destination
- `SetFormatter(logrus.Formatter)` - Set log formatter
- `SetFileOutput(filepath string)` - Set rotating file output with sensible defaults
- `SetFileOutputWithConfig(config FileRotationConfig)` - Set rotating file output with custom settings
- `WithField(key string, value interface{}) *Logger` - Add single field
- `WithFields(fields map[string]interface{}) *Logger` - Add multiple fields
- `GetDefaultLogger() *Logger` - Access default logger instance

### Logger Instances

- `New() *Logger` - Create logger with default settings
- `NewWithOptions(output io.Writer, level logrus.Level, formatter logrus.Formatter) *Logger` - Create with custom options
- `SetFileOutput(filepath string)` - Set rotating file output with sensible defaults
- `SetFileOutputWithConfig(config FileRotationConfig)` - Set rotating file output with custom settings
- Instance methods mirror global functions (Infof, SetLogLevel, WithField, etc.)

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

Use `SetFileOutputWithConfig()` for advanced control over rotation parameters.

```go
logger.SetFileOutputWithConfig(logger.FileRotationConfig{
    Filename:   "logs/app.log",
    MaxSize:    500,      // 500MB files
    MaxBackups: 5,        // Keep 5 old log files
    MaxAge:     90,       // Keep logs for 90 days
    Compress:   true,     // Compress old logs
})

// For custom instances
customLog := logger.New()
customLog.SetFileOutput("logs/custom.log")
```

### Rotation Behavior

- Files rotate when they reach the configured size limit
- Old logs are named with timestamps (e.g., `app.log.2025-01-15.01`)
- Logs older than `MaxAge` days are automatically deleted
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