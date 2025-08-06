package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance
type Logger struct {
	level  LogLevel
	format string
	logger *log.Logger
}

// Config holds logger configuration
type Config struct {
	Level  string // debug, info, warn, error, fatal
	Format string // text, json
}

// New creates a new logger instance
func New(config Config) *Logger {
	level := parseLogLevel(config.Level)
	format := config.Format
	if format == "" {
		format = "text"
	}

	logger := log.New(os.Stdout, "", 0)

	return &Logger{
		level:  level,
		format: format,
		logger: logger,
	}
}

// Default returns a default logger
func Default() *Logger {
	return New(Config{
		Level:  "info",
		Format: "text",
	})
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, message, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, message, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, message, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, message, fields...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, fields ...interface{}) {
	l.log(FATAL, message, fields...)
	os.Exit(1)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, fmt.Sprintf(format, args...))
	}
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, fmt.Sprintf(format, args...))
	}
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, fmt.Sprintf(format, args...))
	}
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, fmt.Sprintf(format, args...))
	}
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FATAL, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// In a more advanced implementation, this would create a new logger
	// with the fields attached for structured logging
	return l
}

// log handles the actual logging
func (l *Logger) log(level LogLevel, message string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	caller := getCaller()

	var logMessage string
	if l.format == "json" {
		logMessage = l.formatJSON(timestamp, level, message, caller, fields...)
	} else {
		logMessage = l.formatText(timestamp, level, message, caller, fields...)
	}

	l.logger.Print(logMessage)
}

// formatText formats log message as text
func (l *Logger) formatText(timestamp string, level LogLevel, message, caller string, fields ...interface{}) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", timestamp))
	parts = append(parts, fmt.Sprintf("[%s]", level.String()))
	parts = append(parts, fmt.Sprintf("[%s]", caller))
	parts = append(parts, message)

	if len(fields) > 0 {
		fieldsStr := fmt.Sprintf("%v", fields)
		parts = append(parts, fieldsStr)
	}

	return strings.Join(parts, " ")
}

// formatJSON formats log message as JSON
func (l *Logger) formatJSON(timestamp string, level LogLevel, message, caller string, fields ...interface{}) string {
	// Simple JSON formatting - in production, use a proper JSON library
	fieldsStr := ""
	if len(fields) > 0 {
		fieldsStr = fmt.Sprintf(",\"fields\":%v", fields)
	}

	return fmt.Sprintf(`{"timestamp":"%s","level":"%s","caller":"%s","message":"%s"%s}`,
		timestamp, level.String(), caller, message, fieldsStr)
}

// getCaller returns the calling function information
func getCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}

	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	filename := parts[len(parts)-1]

	return fmt.Sprintf("%s:%d", filename, line)
}

// parseLogLevel parses string log level to LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// Global logger instance
var globalLogger = Default()

// Global logging functions
func Debug(message string, fields ...interface{}) {
	globalLogger.Debug(message, fields...)
}

func Info(message string, fields ...interface{}) {
	globalLogger.Info(message, fields...)
}

func Warn(message string, fields ...interface{}) {
	globalLogger.Warn(message, fields...)
}

func Error(message string, fields ...interface{}) {
	globalLogger.Error(message, fields...)
}

func Fatal(message string, fields ...interface{}) {
	globalLogger.Fatal(message, fields...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// SetGlobalLogger sets the global logger
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}
