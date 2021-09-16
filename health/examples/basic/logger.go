package main

import "log"

// Logger is a simple logger implementing health.Logger interface.
type Logger struct{}

// Infof logs a message.
func (l *Logger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Errorf logs a message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
