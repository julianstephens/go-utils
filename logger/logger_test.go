package logger_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/logger"
	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	log := logger.New()
	if log == nil {
		t.Fatal("New() should return a non-nil logger")
	}

	// Test default level is info
	if log.GetLevel() != "info" {
		t.Errorf("Expected default level to be 'info', got '%s'", log.GetLevel())
	}
}

func TestNewWithOptions(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.DebugLevel, &logrus.TextFormatter{})

	if log == nil {
		t.Fatal("NewWithOptions() should return a non-nil logger")
	}

	if log.GetLevel() != "debug" {
		t.Errorf("Expected level to be 'debug', got '%s'", log.GetLevel())
	}
}

func TestSetLogLevel(t *testing.T) {
	log := logger.New()

	tests := []struct {
		level    string
		expected string
		hasError bool
	}{
		{"debug", "debug", false},
		{"info", "info", false},
		{"warn", "warning", false},
		{"error", "error", false},
		{"fatal", "fatal", false},
		{"panic", "panic", false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			err := log.SetLogLevel(tt.level)
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for invalid level '%s'", tt.level)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error setting level '%s': %v", tt.level, err)
				return
			}

			if log.GetLevel() != tt.expected {
				t.Errorf("Expected level '%s', got '%s'", tt.expected, log.GetLevel())
			}
		})
	}
}

func TestLoggingMethods(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.DebugLevel, &logrus.JSONFormatter{})

	tests := []struct {
		name     string
		logFunc  func()
		level    string
		message  string
	}{
		{
			name:    "Debug",
			logFunc: func() { log.Debug("debug message") },
			level:   "debug",
			message: "debug message",
		},
		{
			name:    "Info",
			logFunc: func() { log.Info("info message") },
			level:   "info",
			message: "info message",
		},
		{
			name:    "Warn",
			logFunc: func() { log.Warn("warn message") },
			level:   "warning",
			message: "warn message",
		},
		{
			name:    "Error",
			logFunc: func() { log.Error("error message") },
			level:   "error",
			message: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if output == "" {
				t.Fatal("Expected log output, got empty string")
			}

			// Parse JSON log entry
			var logEntry map[string]interface{}
			if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
				t.Fatalf("Failed to parse JSON log output: %v", err)
			}

			// Check level
			if logEntry["level"] != tt.level {
				t.Errorf("Expected level '%s', got '%s'", tt.level, logEntry["level"])
			}

			// Check message
			if logEntry["msg"] != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, logEntry["msg"])
			}

			// Check timestamp exists
			if _, exists := logEntry["time"]; !exists {
				t.Error("Expected timestamp field in log output")
			}
		})
	}
}

func TestFormattedLoggingMethods(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.DebugLevel, &logrus.JSONFormatter{})

	tests := []struct {
		name     string
		logFunc  func()
		level    string
		contains string
	}{
		{
			name:     "Debugf",
			logFunc:  func() { log.Debugf("debug %s %d", "test", 123) },
			level:    "debug",
			contains: "debug test 123",
		},
		{
			name:     "Infof",
			logFunc:  func() { log.Infof("info %s %d", "test", 456) },
			level:    "info",
			contains: "info test 456",
		},
		{
			name:     "Warnf",
			logFunc:  func() { log.Warnf("warn %s %d", "test", 789) },
			level:    "warning",
			contains: "warn test 789",
		},
		{
			name:     "Errorf",
			logFunc:  func() { log.Errorf("error %s %d", "test", 101112) },
			level:    "error",
			contains: "error test 101112",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if output == "" {
				t.Fatal("Expected log output, got empty string")
			}

			// Parse JSON log entry
			var logEntry map[string]interface{}
			if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
				t.Fatalf("Failed to parse JSON log output: %v", err)
			}

			// Check level
			if logEntry["level"] != tt.level {
				t.Errorf("Expected level '%s', got '%s'", tt.level, logEntry["level"])
			}

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
}

func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

	contextLogger := log.WithField("user_id", "12345")
	contextLogger.Info("test message")

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
	if logEntry["user_id"] != "12345" {
		t.Errorf("Expected user_id field to be '12345', got '%v'", logEntry["user_id"])
	}

	// Check message
	if logEntry["msg"] != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", logEntry["msg"])
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

	fields := map[string]interface{}{
		"user_id":    "12345",
		"request_id": "req-abcdef",
		"action":     "login",
	}

	contextLogger := log.WithFields(fields)
	contextLogger.Info("user action")

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
		if logEntry[key] != expectedValue {
			t.Errorf("Expected field '%s' to be '%v', got '%v'", key, expectedValue, logEntry[key])
		}
	}

	// Check message
	if logEntry["msg"] != "user action" {
		t.Errorf("Expected message 'user action', got '%s'", logEntry["msg"])
	}
}

func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.WarnLevel, &logrus.JSONFormatter{})

	// These should not appear in output (below warn level)
	log.Debug("debug message")
	log.Info("info message")

	// These should appear in output (warn level and above)
	log.Warn("warn message")
	log.Error("error message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should only have 2 lines (warn and error)
	if len(lines) != 2 {
		t.Errorf("Expected 2 log lines, got %d: %v", len(lines), lines)
	}

	// Check first line is warn
	var firstEntry map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &firstEntry); err != nil {
		t.Fatalf("Failed to parse first JSON log entry: %v", err)
	}
	if firstEntry["level"] != "warning" || firstEntry["msg"] != "warn message" {
		t.Errorf("First entry should be warning level with 'warn message', got level='%s' msg='%s'", firstEntry["level"], firstEntry["msg"])
	}

	// Check second line is error
	var secondEntry map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &secondEntry); err != nil {
		t.Fatalf("Failed to parse second JSON log entry: %v", err)
	}
	if secondEntry["level"] != "error" || secondEntry["msg"] != "error message" {
		t.Errorf("Second entry should be error level with 'error message', got level='%s' msg='%s'", secondEntry["level"], secondEntry["msg"])
	}
}