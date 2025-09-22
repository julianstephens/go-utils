package logger_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/logger"
	tst "github.com/julianstephens/go-utils/tests"
	"github.com/sirupsen/logrus"
)

func TestGlobalSetLogLevel(t *testing.T) {
	// Test setting valid levels
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}

	for _, level := range validLevels {
		t.Run(level, func(t *testing.T) {
			err := logger.SetLogLevel(level)
			tst.AssertNoError(t, err)
		})
	}

	// Test setting invalid level
	err := logger.SetLogLevel("invalid")
	tst.AssertNotNil(t, err, "Expected error when setting invalid log level")
}

func TestGlobalLoggingMethods(t *testing.T) {
	// Redirect global logger output to a buffer for testing
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLogLevel("debug")

	tests := []struct {
		name    string
		logFunc func()
		level   string
		message string
	}{
		{
			name:    "Global Debug",
			logFunc: func() { logger.Debug("global debug message") },
			level:   "debug",
			message: "global debug message",
		},
		{
			name:    "Global Info",
			logFunc: func() { logger.Info("global info message") },
			level:   "info",
			message: "global info message",
		},
		{
			name:    "Global Warn",
			logFunc: func() { logger.Warn("global warn message") },
			level:   "warning",
			message: "global warn message",
		},
		{
			name:    "Global Error",
			logFunc: func() { logger.Error("global error message") },
			level:   "error",
			message: "global error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			tst.AssertTrue(t, output != "", "Expected log output, got empty string")

			// Parse JSON log entry
			var logEntry map[string]interface{}
			tst.AssertNoError(t, json.Unmarshal([]byte(output), &logEntry))

			// Check level
			tst.AssertDeepEqual(t, logEntry["level"], tt.level)

			// Check message
			tst.AssertDeepEqual(t, logEntry["msg"], tt.message)
		})
	}

	// Reset to stdout to avoid affecting other tests
	logger.SetOutput(os.Stdout)
}

func TestGlobalFormattedLoggingMethods(t *testing.T) {
	// Redirect global logger output to a buffer for testing
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLogLevel("debug")

	tests := []struct {
		name     string
		logFunc  func()
		level    string
		contains string
	}{
		{
			name:     "Global Debugf",
			logFunc:  func() { logger.Debugf("global debug %s %d", "test", 123) },
			level:    "debug",
			contains: "global debug test 123",
		},
		{
			name:     "Global Infof",
			logFunc:  func() { logger.Infof("global info %s %d", "test", 456) },
			level:    "info",
			contains: "global info test 456",
		},
		{
			name:     "Global Warnf",
			logFunc:  func() { logger.Warnf("global warn %s %d", "test", 789) },
			level:    "warning",
			contains: "global warn test 789",
		},
		{
			name:     "Global Errorf",
			logFunc:  func() { logger.Errorf("global error %s %d", "test", 101112) },
			level:    "error",
			contains: "global error test 101112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			tst.AssertTrue(t, output != "", "Expected log output, got empty string")

			// Parse JSON log entry
			var logEntry map[string]interface{}
			tst.AssertNoError(t, json.Unmarshal([]byte(output), &logEntry))

			// Check level
			tst.AssertDeepEqual(t, logEntry["level"], tt.level)

			// Check message contains expected text
			msg, ok := logEntry["msg"].(string)
			if !ok {
				t.Fatal("Expected message to be a string")
			}
			if !strings.Contains(msg, tt.contains) {
				t.Errorf("Expected message to contain '%s', got '%s'", tt.contains, msg)
			}
		})
	}

	// Reset to stdout to avoid affecting other tests
	logger.SetOutput(os.Stdout)
}

func TestGlobalWithField(t *testing.T) {
	// Redirect global logger output to a buffer for testing
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLogLevel("info")

	contextLogger := logger.WithField("service", "test-service")
	contextLogger.Info("service started")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	// Parse JSON log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log output: %v", err)
	}

	// Check field is present
	if logEntry["service"] != "test-service" {
		t.Errorf("Expected service field to be 'test-service', got '%v'", logEntry["service"])
	}

	// Check message
	if logEntry["msg"] != "service started" {
		t.Errorf("Expected message 'service started', got '%s'", logEntry["msg"])
	}

	// Reset to stdout to avoid affecting other tests
	logger.SetOutput(os.Stdout)
}

func TestGlobalWithFields(t *testing.T) {
	// Redirect global logger output to a buffer for testing
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLogLevel("info")

	fields := map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
		"env":     "test",
		"port":    8080,
	}

	contextLogger := logger.WithFields(fields)
	contextLogger.Info("application configuration")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	// Parse JSON log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log output: %v", err)
	}

	// Check all fields are present
	for key, expectedValue := range fields {
		// Handle numeric conversion for port
		if key == "port" {
			if portVal, ok := logEntry[key].(float64); ok {
				if int(portVal) != expectedValue.(int) {
					t.Errorf("Expected field '%s' to be %v, got %v", key, expectedValue, int(portVal))
				}
			} else {
				t.Errorf("Expected field '%s' to be numeric, got %T", key, logEntry[key])
			}
		} else {
			if logEntry[key] != expectedValue {
				t.Errorf("Expected field '%s' to be '%v', got '%v'", key, expectedValue, logEntry[key])
			}
		}
	}

	// Check message
	if logEntry["msg"] != "application configuration" {
		t.Errorf("Expected message 'application configuration', got '%s'", logEntry["msg"])
	}

	// Reset to stdout to avoid affecting other tests
	logger.SetOutput(os.Stdout)
}

func TestGetDefaultLogger(t *testing.T) {
	defaultLogger := logger.GetDefaultLogger()
	if defaultLogger == nil {
		t.Fatal("GetDefaultLogger() should return a non-nil logger")
	}

	// Test that we can use the returned logger
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})

	defaultLogger.Info("test message from default logger")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output from default logger, got empty string")
	}

	// Parse JSON log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log output: %v", err)
	}

	if logEntry["msg"] != "test message from default logger" {
		t.Errorf("Expected message 'test message from default logger', got '%s'", logEntry["msg"])
	}

	// Reset to stdout to avoid affecting other tests
	logger.SetOutput(os.Stdout)
}
