package logger

import "time"

type Options struct {
	Environment bool // true para "production", false para "development"
}

type Field interface {
	Bool(key string, val bool) Field
	Int(key string, val int) Field
	Float64(key string, val float64) Field
	String(key string, val string) Field
	Time(key string, val time.Time) Field
	Int64(key string, val int64) Field
}

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	Field() Field
}
