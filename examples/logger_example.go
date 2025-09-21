package main

import (
	"os"

	"github.com/julianstephens/go-utils/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	// Example 1: Using global logger functions (simplest approach)
	logger.SetLogLevel("debug")
	logger.Infof("Server starting on port %d", 8080)
	logger.Debugf("Debug information: %s", "configuration loaded")
	logger.Warnf("Warning: %s", "deprecated feature in use")

	// Example 2: Using structured logging with fields
	logger.WithField("user_id", "12345").Info("User logged in")
	logger.WithFields(map[string]interface{}{
		"user_id":    "12345",
		"action":     "login",
		"ip_address": "192.168.1.100",
	}).Info("User action recorded")

	// Example 3: Creating a custom logger instance
	customLogger := logger.New()
	customLogger.SetLogLevel("warn")
	customLogger.Infof("This won't appear (below warn level)")
	customLogger.Warnf("This will appear: %s", "warning message")

	// Example 4: Logger with custom configuration
	textLogger := logger.NewWithOptions(
		os.Stdout,
		logrus.DebugLevel,
		&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		},
	)
	textLogger.Infof("Using text formatter: %s", "colored output")

	// Example 5: Contextual logging for request processing
	requestLogger := logger.WithFields(map[string]interface{}{
		"request_id": "req-abc123",
		"method":     "POST",
		"path":       "/api/users",
	})

	requestLogger.Info("Processing request")
	requestLogger.Debugf("Request body size: %d bytes", 256)
	requestLogger.Info("Request completed successfully")

	// Example 6: Error handling with context
	if err := processOrder("order-456"); err != nil {
		logger.WithFields(map[string]interface{}{
			"order_id": "order-456",
			"error":    err.Error(),
		}).Error("Failed to process order")
	}
}

func processOrder(orderID string) error {
	// Simulate processing
	logger.WithField("order_id", orderID).Info("Starting order processing")

	// Simulate an error
	return nil // In real code, this might return an actual error
}
