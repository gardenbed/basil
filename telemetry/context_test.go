package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
)

func TestContextWithUUID(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		uuid string
	}{
		{
			name: "OK",
			ctx:  context.Background(),
			uuid: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := ContextWithUUID(tc.ctx, tc.uuid)

			assert.Equal(t, tc.uuid, ctx.Value(uuidContextKey))
		})
	}
}

func TestUUIDFromContext(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		expectedUUID string
		expectedOK   bool
	}{
		{
			name:         "WithoutUUID",
			ctx:          context.Background(),
			expectedUUID: "",
			expectedOK:   false,
		},
		{
			name:         "WithUUID",
			ctx:          context.WithValue(context.Background(), uuidContextKey, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			expectedUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			expectedOK:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uuid, ok := UUIDFromContext(tc.ctx)

			assert.Equal(t, tc.expectedUUID, uuid)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestContextWithLogger(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		logger Logger
	}{
		{
			name:   "OK",
			ctx:    context.Background(),
			logger: new(voidLogger),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := ContextWithLogger(tc.ctx, tc.logger)

			assert.Equal(t, tc.logger, ctx.Value(loggerContextKey))
		})
	}
}

func TestLoggerFromContext(t *testing.T) {
	logger := new(zapLogger)

	tests := []struct {
		name           string
		ctx            context.Context
		expectedLogger Logger
	}{
		{
			name:           "SingletonLogger",
			ctx:            context.Background(),
			expectedLogger: new(voidLogger),
		},
		{
			name:           "CustomLogger",
			ctx:            context.WithValue(context.Background(), loggerContextKey, logger),
			expectedLogger: logger,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := LoggerFromContext(tc.ctx)

			assert.Equal(t, tc.expectedLogger, logger)
		})
	}
}

func TestContextWithMeter(t *testing.T) {

	tests := []struct {
		name  string
		ctx   context.Context
		meter metric.Meter
	}{
		{
			name:  "OK",
			ctx:   context.Background(),
			meter: global.Meter(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := ContextWithMeter(tc.ctx, tc.meter)

			assert.Equal(t, tc.meter, ctx.Value(meterContextKey))
		})
	}
}

func TestMeterFromContext(t *testing.T) {
	meter := global.Meter("")

	tests := []struct {
		name          string
		ctx           context.Context
		expectedMeter metric.Meter
	}{
		{
			name:          "SingletonMeter",
			ctx:           context.Background(),
			expectedMeter: metric.NewNoopMeterProvider().Meter(""),
		},
		{
			name:          "CustomMeter",
			ctx:           context.WithValue(context.Background(), meterContextKey, meter),
			expectedMeter: meter,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meter := MeterFromContext(tc.ctx)

			assert.Equal(t, tc.expectedMeter, meter)
		})
	}
}

func TestContextWithTracer(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		tracer trace.Tracer
	}{
		{
			name:   "OK",
			ctx:    context.Background(),
			tracer: otel.Tracer(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := ContextWithTracer(tc.ctx, tc.tracer)

			assert.Equal(t, tc.tracer, ctx.Value(tracerContextKey))
		})
	}
}

func TestTracerFromContext(t *testing.T) {
	tracer := otel.Tracer("")

	tests := []struct {
		name           string
		ctx            context.Context
		expectedTracer trace.Tracer
	}{
		{
			name:           "SingletonTracer",
			ctx:            context.Background(),
			expectedTracer: trace.NewNoopTracerProvider().Tracer(""),
		},
		{
			name:           "CustomTracer",
			ctx:            context.WithValue(context.Background(), tracerContextKey, tracer),
			expectedTracer: tracer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer := TracerFromContext(tc.ctx)

			assert.Equal(t, tc.expectedTracer, tracer)
		})
	}
}
