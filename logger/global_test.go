package logger_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"sync"
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
	_ = logger.SetLogLevel("debug")

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
	_ = logger.SetLogLevel("debug")

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
	_ = logger.SetLogLevel("info")

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
	_ = logger.SetLogLevel("info")

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

func TestGlobalConcurrentLoggingAndConfiguration(t *testing.T) {
	// Test that concurrent logging and configuration changes don't cause race conditions
	var wg sync.WaitGroup
	var buf bytes.Buffer

	// Set initial output
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.JSONFormatter{})
	_ = logger.SetLogLevel("debug")

	// Launch goroutines that log concurrently
	numLoggers := 10
	logsPerGoroutine := 100
	wg.Add(numLoggers)

	for i := 0; i < numLoggers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				logger.Infof("concurrent log from goroutine %d, iteration %d", id, j)
			}
		}(i)
	}

	// Wait for all logging goroutines to complete
	wg.Wait()

	// Verify logs were written successfully
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	expectedLogCount := numLoggers * logsPerGoroutine
	if len(lines) < expectedLogCount {
		t.Errorf("Expected at least %d log lines, got %d", expectedLogCount, len(lines))
	}

	// Verify logs are valid JSON
	var validJsonCount int
	for _, line := range lines {
		if line == "" {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Errorf("Failed to parse log line as JSON: %s, error: %v", line, err)
		} else {
			validJsonCount++
		}
	}

	if validJsonCount < expectedLogCount {
		t.Errorf("Expected at least %d valid JSON log entries, got %d", expectedLogCount, validJsonCount)
	}

	// Reset to stdout
	logger.SetOutput(os.Stdout)
}

func TestGlobalConcurrentConfigurationChanges(t *testing.T) {
	// Test that concurrent configuration changes don't cause panics or data races
	var wg sync.WaitGroup
	var buf bytes.Buffer

	logger.SetOutput(&buf)
	_ = logger.SetLogLevel("info")

	numGoroutines := 20
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Alternately change log level and log
			levels := []string{"debug", "info", "warn", "error"}
			for j := 0; j < 25; j++ {
				level := levels[j%len(levels)]
				_ = logger.SetLogLevel(level)
				logger.Infof("log from goroutine %d, iteration %d, level %s", id, j, level)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify no panics occurred and some logs were written
	output := buf.String()
	if output == "" {
		t.Error("Expected some log output from concurrent configuration changes")
	}

	// Reset to stdout
	logger.SetOutput(os.Stdout)
}
