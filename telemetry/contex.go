package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const (
	uuidContextKey   = contextKey("UUID")
	loggerContextKey = contextKey("Logger")
	meterContextKey  = contextKey("Meter")
	tracerContextKey = contextKey("Tracer")
)

// ContextWithUUID creates a new context with a uuid.
func ContextWithUUID(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, uuidContextKey, uuid)
}

// UUIDFromContext retrieves a uuid from a context.
func UUIDFromContext(ctx context.Context) (string, bool) {
	uuid, ok := ctx.Value(uuidContextKey).(string)
	return uuid, ok
}

// ContextWithLogger returns a new context that holds a reference to a logger.
func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// LoggerFromContext returns a logger set on a context.
// If no logger found on the context, the singleton logger will be returned!
func LoggerFromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return logger
	}

	// Return the singleton logger as the default
	return singleton.Logger()
}

// ContextWithMeter returns a new context that holds a reference to a meter.
func ContextWithMeter(ctx context.Context, meter metric.Meter) context.Context {
	return context.WithValue(ctx, meterContextKey, meter)
}

// MeterFromContext returns a meter set on a context.
// If no meter found on the context, the singleton meter will be returned!
func MeterFromContext(ctx context.Context) metric.Meter {
	if meter, ok := ctx.Value(meterContextKey).(metric.Meter); ok {
		return meter
	}

	// Return the singleton meter as the default
	return singleton.Meter()
}

// ContextWithTracer returns a new context that holds a reference to a tracer.
func ContextWithTracer(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerContextKey, tracer)
}

// TracerFromContext returns a tracer set on a context.
// If no tracer found on the context, the singleton tracer will be returned!
func TracerFromContext(ctx context.Context) trace.Tracer {
	if tracer, ok := ctx.Value(tracerContextKey).(trace.Tracer); ok {
		return tracer
	}

	// Return the singleton tracer as the default
	return singleton.Tracer()
}
