package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// defaultLogger is the package-level logger instance used by the global functions.
var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

// SetOutput sets the output destination for the default logger.
func SetOutput(output io.Writer) {
	defaultLogger.entry.Logger.SetOutput(output)
}

// SetFormatter sets the formatter for the default logger.
func SetFormatter(formatter logrus.Formatter) {
	defaultLogger.entry.Logger.SetFormatter(formatter)
}

// SetLogLevel sets the logging level for the default logger.
// Valid levels are: panic, fatal, error, warn, info, debug, trace
func SetLogLevel(level string) error {
	return defaultLogger.SetLogLevel(level)
}

// WithField adds a single field to the default logger context and returns a new logger instance.
// This is useful for structured logging where you want to include contextual information.
func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields adds multiple fields to the default logger context and returns a new logger instance.
// This is useful for structured logging where you want to include multiple pieces of contextual information.
func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}

// Debugf logs a message at debug level with printf-style formatting using the default logger.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof logs a message at info level with printf-style formatting using the default logger.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warnf logs a message at warning level with printf-style formatting using the default logger.
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Errorf logs a message at error level with printf-style formatting using the default logger.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatalf logs a message at fatal level with printf-style formatting using the default logger and then calls os.Exit(1).
// This should be used sparingly as it terminates the program.
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Tracef logs a message at trace level with printf-style formatting using the default logger.
func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args...)
}

// Panicf logs a message at panic level with printf-style formatting using the default logger and then panics.
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}
// Debug logs a message at debug level using the default logger.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info logs a message at info level using the default logger.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn logs a message at warning level using the default logger.
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error logs a message at error level using the default logger.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Fatal logs a message at fatal level using the default logger and then calls os.Exit(1).
// This should be used sparingly as it terminates the program.
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Trace logs a message at trace level using the default logger.
func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

// Tracef logs a message at trace level with printf-style formatting using the default logger.
func Tracef(format string, args ...interface{}) {
	defaultLogger.Tracef(format, args...)
}

// Panic logs a message at panic level using the default logger and then panics.
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Panicf logs a message at panic level with printf-style formatting using the default logger and then panics.
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}
// GetDefaultLogger returns the default logger instance for advanced usage.
func GetDefaultLogger() *Logger {
	return defaultLogger
}