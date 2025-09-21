/*
Package logger provides a unified structured logging interface for all julianstephens Go projects.

This package wraps logrus to offer log level control, custom formatting, and support for
structured/contextual logging. It is designed to be simple to import, idiomatic, and
extensible for future needs such as external log sinks.

Basic Usage:

The logger package provides both instance-based and global logging functions:

	package main

	import (
		"github.com/julianstephens/go-utils/logger"
	)

	func main() {
		// Using global functions (simple approach)
		logger.SetLogLevel("debug")
		logger.Infof("Server starting on port %d", 8080)
		logger.Debugf("Debug information: %s", "some details")
		logger.Errorf("An error occurred: %v", err)

		// Using structured logging with fields
		logger.WithField("user_id", "12345").Info("User logged in")
		logger.WithFields(map[string]interface{}{
			"user_id": "12345",
			"action":  "login",
		}).Info("User action recorded")
	}

Instance-based Usage:

For more control, you can create your own logger instances:

	// Create a new logger with default settings
	log := logger.New()
	log.SetLogLevel("warn")
	log.Infof("Application initialized")

	// Create a logger with custom configuration
	customLogger := logger.NewWithOptions(
		os.Stderr,
		logrus.DebugLevel,
		&logrus.TextFormatter{FullTimestamp: true},
	)

Log Levels:

The logger supports the following log levels (from highest to lowest priority):
  - panic - Logs and then calls panic()
  - fatal - Logs and then calls os.Exit(1)
  - error - Error conditions
  - warn  - Warning conditions
  - info  - Informational messages (default)
  - debug - Debug-level messages
  - trace - Very detailed debug information

Structured Logging:

The logger supports structured logging through the WithField and WithFields methods:

	// Add a single field
	contextLogger := logger.WithField("request_id", "req-123")
	contextLogger.Info("Processing request")

	// Add multiple fields
	contextLogger = logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"user_id":    "user-456",
		"action":     "create_order",
	})
	contextLogger.Info("Order created successfully")

Formatting:

By default, the logger uses JSON formatting with ISO 8601 timestamps. You can customize
the formatter when creating a logger instance:

	import "github.com/sirupsen/logrus"

	// Use text formatting instead of JSON
	textLogger := logger.NewWithOptions(
		os.Stdout,
		logrus.InfoLevel,
		&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		},
	)

Integration:

This logger is designed to replace existing logging implementations in other julianstephens
repositories. It provides the same API surface (Infof, Debugf, Warnf, Errorf, Fatalf,
SetLogLevel) while adding structured logging capabilities.

For HTTP middleware integration, the logger can be easily adapted to work with existing
middleware that expects a *log.Logger:

	// Create a standard library compatible logger wrapper if needed
	stdLogger := log.New(os.Stdout, "", 0)

Thread Safety:

The logger is thread-safe and can be safely used across multiple goroutines.
*/
package logger
