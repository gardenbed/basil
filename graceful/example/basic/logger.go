package main

import "log"

// Logger is a simple logger implementing graceful.Logger interface.
type Logger struct{}

// Debugf logs a message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Infof logs a message.
func (l *Logger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Warnf logs a message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Errorf logs a message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
