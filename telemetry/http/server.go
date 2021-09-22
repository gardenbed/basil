package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/trace"

	"github.com/gardenbed/basil/telemetry"
)

// Server-side instruments for measurements.
type serverInstruments struct {
	panic   metric.Int64Counter
	total   metric.Int64Counter
	active  metric.Int64UpDownCounter
	latency metric.Int64Histogram
}

func newServerInstruments(meter metric.Meter) *serverInstruments {
	mm := metric.Must(meter)

	return &serverInstruments{
		panic: mm.NewInt64Counter(
			"incoming_http_requests_panic",
			metric.WithUnit(unit.Dimensionless),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The total number of panics happened in http handlers (server-side)"),
		),
		total: mm.NewInt64Counter(
			"incoming_http_requests_total",
			metric.WithUnit(unit.Dimensionless),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The total number of incoming http requests (server-side)"),
		),
		active: mm.NewInt64UpDownCounter(
			"incoming_http_requests_active",
			metric.WithUnit(unit.Dimensionless),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The number of in-flight incoming http requests (server-side)"),
		),
		latency: mm.NewInt64Histogram(
			"incoming_http_requests_latency",
			metric.WithUnit(unit.Milliseconds),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The duration of incoming http requests in milliseconds (server-side)"),
		),
	}
}

// Middleware creates observable http handlers with logging, metrics, and tracing.
type Middleware struct {
	probe       telemetry.Probe
	opts        Options
	instruments *serverInstruments
}

// NewMiddleware creates a new observable http middleware.
func NewMiddleware(probe telemetry.Probe, opts Options) *Middleware {
	opts = opts.withDefaults()
	instruments := newServerInstruments(probe.Meter())

	return &Middleware{
		probe:       probe,
		opts:        opts,
		instruments: instruments,
	}
}

func (m *Middleware) callHandlerFunc(ctx context.Context, handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			m.probe.Logger().Errorf("Panic recovered: %v", r)
			m.instruments.panic.Add(ctx, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	handler(w, r)
}

// Wrap wraps an existing http handler function and returns a new observable handler function.
// This can be used for making http handlers observable via logging, metrics, tracing, etc.
// It also observes and recovers panics that happened inside the inner http handler.
func (m *Middleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := r.Context()
		kind := "server"
		method := r.Method
		url := r.URL.Path
		route := m.opts.IDRegexp.ReplaceAllString(url, ":id")

		meter := m.probe.Meter()
		tracer := m.probe.Tracer()

		// Increase the number of in-flight requests
		m.instruments.active.Add(ctx, 1,
			attribute.String("method", method),
			attribute.String("route", route),
		)

		// Make sure we decrease the number of in-flight requests
		defer m.instruments.active.Add(ctx, -1,
			attribute.String("method", method),
			attribute.String("route", route),
		)

		// Make sure the request has a UUID
		requestUUID := r.Header.Get(requestUUIDHeader)
		if requestUUID == "" {
			requestUUID = uuid.New().String()
			r.Header.Set(requestUUIDHeader, requestUUID)
		}

		// Get the name of client for the request if any
		clientName := r.Header.Get(clientNameHeader)

		// Propagate request metadata by adding them to outgoing http response headers
		w.Header().Set(requestUUIDHeader, requestUUID)
		w.Header().Set(clientNameHeader, clientName)

		// Extract context from the http headers
		ctx = otel.GetTextMapPropagator().Extract(ctx, &headerTextMapCarrier{
			Header: r.Header,
		})

		// Start a new span
		ctx, span := tracer.Start(ctx,
			"http-server-request",
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Create a contextualized logger
		contextFields := []interface{}{
			"req.uuid", requestUUID,
			"req.kind", kind,
			"req.method", method,
			"req.url", url,
			"req.route", route,
			"traceId", span.SpanContext().TraceID().String(),
			"spanId", span.SpanContext().SpanID().String(),
		}
		if clientName != "" {
			contextFields = append(contextFields, "client.name", clientName)
		}
		logger := m.probe.Logger().With(contextFields...)

		// Augment the request context
		ctx = telemetry.ContextWithUUID(ctx, requestUUID)
		ctx = telemetry.ContextWithLogger(ctx, logger)
		ctx = telemetry.ContextWithMeter(ctx, meter)
		ctx = telemetry.ContextWithTracer(ctx, tracer)
		req := r.WithContext(ctx)

		// Create a wrapped response writer, so we can know about the response
		rw := newResponseWriter(w)

		// Call http handler
		span.AddEvent("calling http handler")
		m.callHandlerFunc(ctx, next, rw, req)

		duration := time.Since(startTime).Milliseconds()
		statusCode := rw.StatusCode
		statusClass := rw.StatusClass

		// Report metrics
		meter.RecordBatch(ctx,
			[]attribute.KeyValue{
				attribute.String("method", method),
				attribute.String("route", route),
				attribute.Int("status_code", statusCode),
				attribute.String("status_class", statusClass),
			},
			m.instruments.total.Measurement(1),
			m.instruments.latency.Measurement(duration),
		)

		// Report logs
		message := fmt.Sprintf("%s %s %d %dms", method, url, statusCode, duration)
		fields := []interface{}{
			"resp.statusCode", statusCode,
			"resp.statusClass", statusClass,
			"resp.duration", duration,
		}

		// Determine the log level based on the result
		switch {
		case statusCode >= 500:
			logger.Error(message, fields...)
		case statusCode >= 400:
			logger.Warn(message, fields...)
		case statusCode >= 100:
			fallthrough
		default:
			logger.Info(message, fields...)
		}

		// Report the span
		span.SetAttributes(
			attribute.String("method", method),
			attribute.String("url", url),
			attribute.String("route", route),
			attribute.Int("status_code", statusCode),
		)
	}
}
