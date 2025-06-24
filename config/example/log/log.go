// Package log provides a toy logger for the examples.
package log

import (
	"log"
	"strings"
)

const (
	// None level
	None uint = iota
	// Error level
	Error
	// Warn level
	Warn
	// Info level
	Info
	// Debug level
	Debug
)

// Logger is a simple logger with four different levels: Debug, Info, Warn, Error
type Logger struct {
	level uint
}

// SetLevel changes the verbosity level
func (l *Logger) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "none":
		l.level = None
	case "error":
		l.level = Error
	case "warn":
		l.level = Warn
	case "info":
		l.level = Info
	case "debug":
		l.level = Debug
	}
}

// Debug logs in debug level.
func (l *Logger) Debug(v ...any) {
	if l.level >= Debug {
		log.Print(v...)
	}
}

// Debugf logs in debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level >= Debug {
		log.Printf(format, args...)
	}
}

// Info logs in info level.
func (l *Logger) Info(v ...any) {
	if l.level >= Info {
		log.Print(v...)
	}
}

// Infof logs in info level.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level >= Info {
		log.Printf(format, args...)
	}
}

// Warn logs in warn level.
func (l *Logger) Warn(v ...any) {
	if l.level >= Warn {
		log.Print(v...)
	}
}

// Warnf logs in warn level.
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level >= Warn {
		log.Printf(format, args...)
	}
}

// Error logs in error level.
func (l *Logger) Error(v ...any) {
	if l.level >= Error {
		log.Print(v...)
	}
}

// Errorf logs in error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level >= Error {
		log.Printf(format, args...)
	}
}
