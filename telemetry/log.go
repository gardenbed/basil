package telemetry

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is the logging level.
type Level int

// Logging level
const (
	LevelNone Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

// Logger is a levelled structured logger.
// It is concurrently safe to be used by multiple goroutines.
type Logger interface {
	Level() Level
	SetLevel(level string)
	With(kv ...interface{}) Logger
	Debug(message string, kv ...interface{})
	Debugf(format string, args ...interface{})
	Info(message string, kv ...interface{})
	Infof(format string, args ...interface{})
	Warn(message string, kv ...interface{})
	Warnf(format string, args ...interface{})
	Error(message string, kv ...interface{})
	Errorf(format string, args ...interface{})
	Close() error
}

// voidLogger implements a no-op logger.
type voidLogger struct{}

func (l *voidLogger) Level() Level                              { return LevelNone }
func (l *voidLogger) SetLevel(level string)                     {}
func (l *voidLogger) With(kv ...interface{}) Logger             { return l }
func (l *voidLogger) Debug(message string, kv ...interface{})   {}
func (l *voidLogger) Debugf(format string, args ...interface{}) {}
func (l *voidLogger) Info(message string, kv ...interface{})    {}
func (l *voidLogger) Infof(format string, args ...interface{})  {}
func (l *voidLogger) Warn(message string, kv ...interface{})    {}
func (l *voidLogger) Warnf(format string, args ...interface{})  {}
func (l *voidLogger) Error(message string, kv ...interface{})   {}
func (l *voidLogger) Errorf(format string, args ...interface{}) {}
func (l *voidLogger) Close() error                              { return nil }

type zapLogger struct {
	config *zap.Config
	logger *zap.SugaredLogger
}

func (l *zapLogger) Level() Level {
	switch l.config.Level.Level() {
	case zapcore.DebugLevel:
		return LevelDebug
	case zapcore.InfoLevel:
		return LevelInfo
	case zapcore.WarnLevel:
		return LevelWarn
	case zapcore.ErrorLevel:
		return LevelError
	default:
		return LevelNone
	}
}

func (l *zapLogger) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		l.config.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		l.config.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		l.config.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		l.config.Level.SetLevel(zapcore.ErrorLevel)
	case "none":
		fallthrough
	default:
		l.config.Level.SetLevel(zapcore.Level(99))
	}
}

func (l *zapLogger) With(kv ...interface{}) Logger {
	return &zapLogger{
		config: l.config,
		logger: l.logger.With(kv...),
	}
}

func (l *zapLogger) Debug(message string, kv ...interface{}) {
	l.logger.Debugw(message, kv...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *zapLogger) Info(message string, kv ...interface{}) {
	l.logger.Infow(message, kv...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *zapLogger) Warn(message string, kv ...interface{}) {
	l.logger.Warnw(message, kv...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *zapLogger) Error(message string, kv ...interface{}) {
	l.logger.Errorw(message, kv...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *zapLogger) Close() error {
	return l.logger.Sync()
}
