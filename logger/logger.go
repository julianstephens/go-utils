package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileRotationConfig configures rotating file output for the logger.
type FileRotationConfig struct {
	Filename   string // Path to log file
	MaxSize    int    // Max file size in megabytes (default: 100)
	MaxBackups int    // Max backup files to retain (default: 3)
	MaxAge     int    // Max age in days before deletion (default: 28)
	Compress   bool   // Compress old logs (default: true)
}

// Logger wraps logrus to provide a unified logging interface for all julianstephens Go projects.
// It offers structured logging with configurable levels, custom formatting, and contextual logging support.
type Logger struct {
	entry *logrus.Entry
}

// New creates a new Logger instance with default configuration.
// By default, it logs to stdout with INFO level and JSON formatting.
func New() *Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetLevel(logrus.InfoLevel)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	return &Logger{
		entry: logrus.NewEntry(logrusLogger),
	}
}

// NewWithOptions creates a new Logger instance with custom configuration options.
func NewWithOptions(output io.Writer, level logrus.Level, formatter logrus.Formatter) *Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(output)
	logrusLogger.SetLevel(level)
	logrusLogger.SetFormatter(formatter)

	return &Logger{
		entry: logrus.NewEntry(logrusLogger),
	}
}

// SetLogLevel sets the logging level for the logger.
// Valid levels are: panic, fatal, error, warn, info, debug, trace
func (l *Logger) SetLogLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	l.entry.Logger.SetLevel(logLevel)
	return nil
}

// WithField adds a single field to the logger context and returns a new logger instance.
// This is useful for structured logging where you want to include contextual information.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		entry: l.entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the logger context and returns a new logger instance.
// This is useful for structured logging where you want to include multiple pieces of contextual information.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

// Debugf logs a message at debug level with printf-style formatting.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Infof logs a message at info level with printf-style formatting.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warnf logs a message at warning level with printf-style formatting.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Errorf logs a message at error level with printf-style formatting.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatalf logs a message at fatal level with printf-style formatting and then calls os.Exit(1).
// This should be used sparingly as it terminates the program.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Tracef logs a message at trace level with printf-style formatting.
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args...)
}

// Panicf logs a message at panic level with printf-style formatting and then panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

// Debug logs a message at debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Info logs a message at info level.
func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Warn logs a message at warning level.
func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Error logs a message at error level.
func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1).
// This should be used sparingly as it terminates the program.
func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Trace logs a message at trace level.
func (l *Logger) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

// Panic logs a message at panic level and then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

// GetLevel returns the current logging level.
func (l *Logger) GetLevel() string {
	return l.entry.Logger.GetLevel().String()
}

// SetFileOutput configures the logger to write to a rotating file with sensible defaults.
// The file will rotate when it reaches 100MB, keeping 3 backups for 28 days with compression enabled.
func (l *Logger) SetFileOutput(filepath string) error {
	config := FileRotationConfig{
		Filename:   filepath,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	return l.SetFileOutputWithConfig(config)
}

// SetFileOutputWithConfig configures rotating file output with custom settings.
func (l *Logger) SetFileOutputWithConfig(config FileRotationConfig) error {
	rotatingFile := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	l.entry.Logger.SetOutput(rotatingFile)
	return nil
}
