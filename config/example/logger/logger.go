package logger

import "log"

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
	Level uint
}

// SetLevel changes the verbosity level
func (l *Logger) SetLevel(level string) {
	switch level {
	case "None":
		l.Level = None
	case "Error":
		l.Level = Error
	case "Warn":
		l.Level = Warn
	case "Info":
		l.Level = Info
	case "Debug":
		l.Level = Debug
	}
}

// Debugf logs in debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Level >= Debug {
		log.Printf(format, args...)
	}
}

// Infof logs in info level.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.Level >= Info {
		log.Printf(format, args...)
	}
}

// Warnf logs in warn level.
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.Level >= Warn {
		log.Printf(format, args...)
	}
}

// Errorf logs in error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.Level >= Error {
		log.Printf(format, args...)
	}
}
