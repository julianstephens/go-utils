package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/logger"
	tst "github.com/julianstephens/go-utils/tests"
	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	log := logger.New()
	tst.AssertNotNil(t, log, "New() should return a non-nil logger")

	// Test default level is info
	tst.AssertDeepEqual(t, log.GetLevel(), "info")
}

func TestNewWithOptions(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.DebugLevel, &logrus.TextFormatter{})

	tst.AssertNotNil(t, log, "NewWithOptions() should return a non-nil logger")
	tst.AssertDeepEqual(t, log.GetLevel(), "debug")
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
				tst.AssertNotNil(t, err, "Expected error for invalid level")
				return
			}

			tst.AssertNoError(t, err)
			tst.AssertDeepEqual(t, log.GetLevel(), tt.expected)
		})
	}
}

func TestLoggingMethods(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.DebugLevel, &logrus.JSONFormatter{})

	tests := []struct {
		name    string
		logFunc func()
		level   string
		message string
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
			tst.AssertTrue(t, output != "", "Expected log output, got empty string")

			// Parse JSON log entry
			var logEntry map[string]interface{}
			tst.AssertNoError(t, json.Unmarshal([]byte(output), &logEntry))

			// Check level
			tst.AssertDeepEqual(t, logEntry["level"], tt.level)

			// Check message
			tst.AssertDeepEqual(t, logEntry["msg"], tt.message)

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
}

func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

	contextLogger := log.WithField("user_id", "12345")
	contextLogger.Info("test message")

	output := buf.String()
	tst.AssertTrue(t, output != "", "Expected log output, got empty string")

	// Parse JSON log entry
	var logEntry map[string]interface{}
	tst.AssertNoError(t, json.Unmarshal([]byte(output), &logEntry))

	// Check field is present
	tst.AssertDeepEqual(t, logEntry["user_id"], "12345")

	// Check message
	tst.AssertDeepEqual(t, logEntry["msg"], "test message")
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
	tst.AssertTrue(t, output != "", "Expected log output, got empty string")

	// Parse JSON log entry
	var logEntry map[string]interface{}
	tst.AssertNoError(t, json.Unmarshal([]byte(output), &logEntry))

	// Check all fields are present
	for key, expectedValue := range fields {
		if logEntry[key] != expectedValue {
			t.Errorf("Expected field '%s' to be '%v', got '%v'", key, expectedValue, logEntry[key])
		}
	}

	// Check message
	tst.AssertDeepEqual(t, logEntry["msg"], "user action")
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
		t.Errorf(
			"First entry should be warning level with 'warn message', got level='%s' msg='%s'",
			firstEntry["level"],
			firstEntry["msg"],
		)
	}

	// Check second line is error
	var secondEntry map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &secondEntry); err != nil {
		t.Fatalf("Failed to parse second JSON log entry: %v", err)
	}
	if secondEntry["level"] != "error" || secondEntry["msg"] != "error message" {
		t.Errorf(
			"Second entry should be error level with 'error message', got level='%s' msg='%s'",
			secondEntry["level"],
			secondEntry["msg"],
		)
	}
}

func TestClose(t *testing.T) {
	t.Run("close logger without file output", func(t *testing.T) {
		log := logger.New()
		err := log.Close()
		tst.AssertNil(t, err, "Close() should return nil when no file output is configured")
	})

	t.Run("close logger with file output", func(t *testing.T) {
		tmpFile := t.TempDir() + "/test.log"
		log := logger.New()
		err := log.SetFileOutput(tmpFile)
		tst.AssertNil(t, err, "SetFileOutput should not return an error")

		err = log.Close()
		tst.AssertNil(t, err, "Close() should successfully close the file output")
	})

	t.Run("close logger with custom config", func(t *testing.T) {
		tmpFile := t.TempDir() + "/test.log"
		log := logger.New()
		maxBackups := 2
		maxAge := 7

		config := logger.FileRotationConfig{
			Filename:   tmpFile,
			MaxSize:    50,
			MaxBackups: &maxBackups,
			MaxAge:     &maxAge,
			Compress:   true,
		}

		err := log.SetFileOutputWithConfig(config)
		tst.AssertNil(t, err, "SetFileOutputWithConfig should not return an error")

		err = log.Close()
		tst.AssertNil(t, err, "Close() should successfully close the file output")
	})
}

func TestRotatingFileOutput(t *testing.T) {
	t.Run("SetFileOutput creates and writes to log file", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "app.log")

		log := logger.New()
		_ = log.SetLogLevel("info")

		err := log.SetFileOutput(logFile)
		tst.AssertNil(t, err, "SetFileOutput should not return an error")
		defer func() { _ = log.Close() }()

		// Write logs
		log.Info("test message 1")
		log.Infof("test message %d", 2)

		// Verify file was created
		_, err = os.Stat(logFile)
		tst.AssertNil(t, err, "Log file should be created")

		// Verify logs were written to file
		content, err := os.ReadFile(logFile)
		tst.AssertNil(t, err, "Should be able to read log file")
		tst.AssertNotNil(t, content, "Log file should contain content")

		contentStr := string(content)
		if !strings.Contains(contentStr, "test message 1") {
			t.Errorf("Log file should contain 'test message 1', got: %s", contentStr)
		}
	})

	t.Run("SetFileOutputWithConfig uses custom settings", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "custom.log")
		maxBackups := 2
		maxAge := 7

		config := logger.FileRotationConfig{
			Filename:   logFile,
			MaxSize:    100,
			MaxBackups: &maxBackups,
			MaxAge:     &maxAge,
			Compress:   false,
		}

		log := logger.New()
		_ = log.SetLogLevel("info")

		err := log.SetFileOutputWithConfig(config)
		tst.AssertNil(t, err, "SetFileOutputWithConfig should not return an error")
		defer func() { _ = log.Close() }()

		// Write logs
		log.Info("custom config test")

		// Verify file was created
		_, err = os.Stat(logFile)
		tst.AssertNil(t, err, "Log file should be created with custom config")

		// Verify logs were written
		content, err := os.ReadFile(logFile)
		tst.AssertNil(t, err, "Should be able to read log file")

		if !strings.Contains(string(content), "custom config test") {
			t.Error("Log file should contain custom config test message")
		}
	})

	t.Run("SetFileOutput with optional MaxBackups only", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "maxbackups.log")
		maxBackups := 3

		config := logger.FileRotationConfig{
			Filename:   logFile,
			MaxSize:    100,
			MaxBackups: &maxBackups,
			MaxAge:     nil, // No age limit
			Compress:   false,
		}

		log := logger.New()
		err := log.SetFileOutputWithConfig(config)
		tst.AssertNil(t, err, "SetFileOutputWithConfig with MaxBackups only should not error")
		defer func() { _ = log.Close() }()

		log.Info("maxbackups test")
		_, err = os.Stat(logFile)
		tst.AssertNil(t, err, "Log file should be created")
	})

	t.Run("SetFileOutput with optional MaxAge only", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "maxage.log")
		maxAge := 14

		config := logger.FileRotationConfig{
			Filename:   logFile,
			MaxSize:    100,
			MaxBackups: nil, // No backup limit
			MaxAge:     &maxAge,
			Compress:   false,
		}

		log := logger.New()
		err := log.SetFileOutputWithConfig(config)
		tst.AssertNil(t, err, "SetFileOutputWithConfig with MaxAge only should not error")
		defer func() { _ = log.Close() }()

		log.Info("maxage test")
		_, err = os.Stat(logFile)
		tst.AssertNil(t, err, "Log file should be created")
	})

	t.Run("SetFileOutput with neither MaxBackups nor MaxAge", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "no_constraints.log")

		config := logger.FileRotationConfig{
			Filename:   logFile,
			MaxSize:    100,
			MaxBackups: nil, // No backup limit
			MaxAge:     nil, // No age limit
			Compress:   false,
		}

		log := logger.New()
		err := log.SetFileOutputWithConfig(config)
		tst.AssertNil(t, err, "SetFileOutputWithConfig without constraints should not error")
		defer func() { _ = log.Close() }()

		log.Info("no constraints test")
		_, err = os.Stat(logFile)
		tst.AssertNil(t, err, "Log file should be created")
	})

	t.Run("multiple logs are written to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "multi.log")

		log := logger.New()
		_ = log.SetLogLevel("debug")

		err := log.SetFileOutput(logFile)
		tst.AssertNil(t, err, "SetFileOutput should not return an error")
		defer func() { _ = log.Close() }()

		// Write multiple logs
		log.Debug("debug message")
		log.Info("info message")
		log.Warn("warn message")
		log.Error("error message")

		// Read and verify all messages are in file
		content, err := os.ReadFile(logFile)
		tst.AssertNil(t, err, "Should be able to read log file")

		contentStr := string(content)
		messages := []string{"debug message", "info message", "warn message", "error message"}
		for _, msg := range messages {
			if !strings.Contains(contentStr, msg) {
				t.Errorf("Log file should contain '%s'", msg)
			}
		}
	})

	t.Run("JSON format is used for file output", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "json.log")

		log := logger.New()
		_ = log.SetLogLevel("info")

		err := log.SetFileOutput(logFile)
		tst.AssertNil(t, err, "SetFileOutput should not return an error")
		defer func() { _ = log.Close() }()

		log.WithFields(map[string]interface{}{
			"user_id": "123",
			"action":  "login",
		}).Info("user logged in")

		// Read and parse JSON
		content, err := os.ReadFile(logFile)
		tst.AssertNil(t, err, "Should be able to read log file")

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(lines) == 0 {
			t.Fatal("Log file should contain at least one line")
		}

		var entry map[string]interface{}
		err = json.Unmarshal([]byte(lines[0]), &entry)
		tst.AssertNil(t, err, "First log line should be valid JSON")

		if entry["msg"] != "user logged in" {
			t.Errorf("Expected msg 'user logged in', got '%v'", entry["msg"])
		}
		if entry["user_id"] != "123" {
			t.Errorf("Expected user_id '123', got '%v'", entry["user_id"])
		}
		if entry["action"] != "login" {
			t.Errorf("Expected action 'login', got '%v'", entry["action"])
		}
	})

	t.Run("multiple SetFileOutput calls replace the output", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile1 := filepath.Join(tmpDir, "first.log")
		logFile2 := filepath.Join(tmpDir, "second.log")

		log := logger.New()
		_ = log.SetLogLevel("info")

		// Set first file
		err := log.SetFileOutput(logFile1)
		tst.AssertNil(t, err, "First SetFileOutput should not error")

		log.Info("message to first file")

		// Set second file (should replace output)
		err = log.SetFileOutput(logFile2)
		tst.AssertNil(t, err, "Second SetFileOutput should not error")
		defer func() { _ = log.Close() }()

		log.Info("message to second file")

		// Verify first file has first message
		content1, err := os.ReadFile(logFile1)
		tst.AssertNil(t, err, "Should be able to read first log file")
		if !strings.Contains(string(content1), "message to first file") {
			t.Error("First log file should contain 'message to first file'")
		}

		// Verify second file has second message
		content2, err := os.ReadFile(logFile2)
		tst.AssertNil(t, err, "Should be able to read second log file")
		if !strings.Contains(string(content2), "message to second file") {
			t.Error("Second log file should contain 'message to second file'")
		}
	})
}
func TestSetOutput(t *testing.T) {
	t.Run("SetOutput changes output destination", func(t *testing.T) {
		var buf1 bytes.Buffer
		var buf2 bytes.Buffer

		log := logger.New()
		_ = log.SetLogLevel("info")

		// Set first output
		log.SetOutput(&buf1)
		log.Info("message to buffer 1")

		// Verify message in first buffer
		if !strings.Contains(buf1.String(), "message to buffer 1") {
			t.Error("Message should be in first buffer")
		}

		// Set second output
		log.SetOutput(&buf2)
		log.Info("message to buffer 2")

		// Verify message in second buffer
		if !strings.Contains(buf2.String(), "message to buffer 2") {
			t.Error("Message should be in second buffer")
		}

		// Verify message not in first buffer anymore
		if strings.Contains(buf1.String(), "message to buffer 2") {
			t.Error("Second message should not be in first buffer")
		}
	})

	t.Run("SetOutput closes previous file output", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")

		log := logger.New()
		_ = log.SetLogLevel("info")

		// Set file output
		_ = log.SetFileOutput(logFile)
		log.Info("message 1")

		var buf bytes.Buffer
		// Set buffer output (should close previous file)
		log.SetOutput(&buf)
		log.Info("message 2")

		// Verify both messages work
		if !strings.Contains(buf.String(), "message 2") {
			t.Error("Message should be in buffer output")
		}

		_ = log.Close()
	})
}

func TestWithContext(t *testing.T) {
	t.Run("WithContext extracts trace-id", func(t *testing.T) {
		var buf bytes.Buffer
		log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

		//nolint:staticcheck
		ctx := context.WithValue(context.Background(), "trace-id", "trace-123")
		log.WithContext(ctx).Info("test message")

		var entry map[string]interface{}
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		_ = json.Unmarshal([]byte(lines[0]), &entry)

		if entry["trace_id"] != "trace-123" {
			t.Errorf("Expected trace_id 'trace-123', got '%v'", entry["trace_id"])
		}
	})

	t.Run("WithContext extracts request-id", func(t *testing.T) {
		var buf bytes.Buffer
		log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

		//nolint:staticcheck
		ctx := context.WithValue(context.Background(), "request-id", "req-456")
		log.WithContext(ctx).Info("test message")

		var entry map[string]interface{}
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		_ = json.Unmarshal([]byte(lines[0]), &entry)

		if entry["request_id"] != "req-456" {
			t.Errorf("Expected request_id 'req-456', got '%v'", entry["request_id"])
		}
	})

	t.Run("WithContext extracts user-id", func(t *testing.T) {
		var buf bytes.Buffer
		log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

		//nolint:staticcheck
		ctx := context.WithValue(context.Background(), "user-id", "user-789")
		log.WithContext(ctx).Info("test message")

		var entry map[string]interface{}
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		_ = json.Unmarshal([]byte(lines[0]), &entry)

		if entry["user_id"] != "user-789" {
			t.Errorf("Expected user_id 'user-789', got '%v'", entry["user_id"])
		}
	})

	t.Run("WithContext handles alternate key formats", func(t *testing.T) {
		var buf bytes.Buffer
		log := logger.NewWithOptions(&buf, logrus.InfoLevel, &logrus.JSONFormatter{})

		// Test traceId (camelCase)
		//nolint:staticcheck
		ctx := context.WithValue(context.Background(), "traceId", "trace-alt")
		//nolint:staticcheck
		ctx = context.WithValue(ctx, "requestId", "req-alt")
		//nolint:staticcheck
		ctx = context.WithValue(ctx, "userId", "user-alt")

		log.WithContext(ctx).Info("test message")

		var entry map[string]interface{}
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		_ = json.Unmarshal([]byte(lines[0]), &entry)

		if entry["trace_id"] != "trace-alt" {
			t.Errorf("Expected trace_id 'trace-alt', got '%v'", entry["trace_id"])
		}
		if entry["request_id"] != "req-alt" {
			t.Errorf("Expected request_id 'req-alt', got '%v'", entry["request_id"])
		}
		if entry["user_id"] != "user-alt" {
			t.Errorf("Expected user_id 'user-alt', got '%v'", entry["user_id"])
		}
	})

	t.Run("WithContext with nil context returns logger unchanged", func(t *testing.T) {
		log := logger.New()
		result := log.WithContext(context.TODO())
		tst.AssertNotNil(t, result, "WithContext(context.TODO()) should return a logger")
	})

	t.Run("WithContext with empty context returns logger unchanged", func(t *testing.T) {
		log := logger.New()
		ctx := context.Background()
		result := log.WithContext(ctx)
		tst.AssertNotNil(t, result, "WithContext with empty context should return a logger")
	})
}

func TestSync(t *testing.T) {
	t.Run("Sync succeeds on stdout output", func(t *testing.T) {
		log := logger.New()
		err := log.Sync()
		tst.AssertNil(t, err, "Sync should succeed")
	})

	t.Run("Sync succeeds on file output", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "sync_test.log")

		log := logger.New()
		_ = log.SetFileOutput(logFile)
		defer func() { _ = log.Close() }()

		log.Info("test message")
		err := log.Sync()
		tst.AssertNil(t, err, "Sync should succeed on file output")

		// Verify file was written
		_, err = os.Stat(logFile)
		tst.AssertNil(t, err, "Log file should exist")
	})

	t.Run("Close calls Sync before closing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "close_sync_test.log")

		log := logger.New()
		_ = log.SetFileOutput(logFile)

		log.Info("test message before close")
		err := log.Close()
		tst.AssertNil(t, err, "Close should succeed")

		// Verify file was written (Sync was called)
		content, err := os.ReadFile(logFile)
		tst.AssertNil(t, err, "Should be able to read log file")
		if !strings.Contains(string(content), "test message before close") {
			t.Error("Log file should contain message from before close")
		}
	})
}
