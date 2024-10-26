package logger

import "time"

// Options defines the configuration options for the logger.
type Options struct {
	Environment bool // true for "production", false for "development"
}

// Field represents a log field with various data types.
type Field interface {
	Bool(key string, val bool) Field       // Adds a boolean field to the log entry.
	Int(key string, val int) Field         // Adds an integer field to the log entry.
	Float64(key string, val float64) Field // Adds a float64 field to the log entry.
	String(key string, val string) Field   // Adds a string field to the log entry.
	Time(key string, val time.Time) Field  // Adds a time field to the log entry.
	Int64(key string, val int64) Field     // Adds an int64 field to the log entry.
	Error(key string, val error) Field     // **Adicionado**: Adds an error field to the log entry.
	Uint64(key string, val uint64) Field
	Uint8(key string, val uint8) Field
}

// Logger provides methods for logging messages at different levels.
type Logger interface {
	Info(msg string, fields ...Field)  // Logs an informational message.
	Error(msg string, fields ...Field) // Logs an error message.
	Debug(msg string, fields ...Field) // Logs a debug message.
	Warn(msg string, fields ...Field)  // Logs a warning message.
	Fatal(msg string, fields ...Field) // Logs a fatal message and terminates the application.

	Field() Field // Creates a new log field for use in logging.
}
