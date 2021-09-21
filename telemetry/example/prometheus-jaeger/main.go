package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gardenbed/basil/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric"
)

type instruments struct {
	reqCounter  metric.Int64Counter
	reqDuration metric.Float64Histogram
}

func newInstruments(meter metric.Meter) *instruments {
	mm := metric.Must(meter)

	return &instruments{
		reqCounter:  mm.NewInt64Counter("requests_total", metric.WithDescription("the total number of requests")),
		reqDuration: mm.NewFloat64Histogram("request_duration_seconds", metric.WithDescription("the duration of requests in seconds")),
	}
}

type server struct {
	probe       telemetry.Probe
	instruments *instruments
}

func (s *server) Handle(ctx context.Context) {
	// Tracing
	ctx, span := s.probe.Tracer().Start(ctx, "handle-request")
	defer span.End()

	start := time.Now()
	s.fetch(ctx)
	s.respond(ctx)
	duration := time.Since(start)

	labels := []attribute.KeyValue{
		attribute.String("method", "GET"),
		attribute.String("endpoint", "/user"),
		attribute.Int("statusCode", 200),
	}

	// Metrics
	s.probe.Meter().RecordBatch(ctx, labels,
		s.instruments.reqCounter.Measurement(1),
		s.instruments.reqDuration.Measurement(duration.Seconds()),
	)

	// Logging
	s.probe.Logger().Info("request handled successfully.",
		"method", "GET",
		"endpoint", "/user",
		"statusCode", 200,
	)
}

func (s *server) fetch(ctx context.Context) {
	_, span := s.probe.Tracer().Start(ctx, "read-database")
	defer span.End()

	time.Sleep(50 * time.Millisecond)
}

func (s *server) respond(ctx context.Context) {
	_, span := s.probe.Tracer().Start(ctx, "send-response")
	defer span.End()

	time.Sleep(10 * time.Millisecond)
}

func main() {
	ctx := context.Background()

	// Creating a new probe and set it as the singleton
	p := telemetry.NewProbe(
		telemetry.WithLogger("info"),
		telemetry.WithPrometheus(),
		telemetry.WithJaeger("", "", "", "", ""),
		telemetry.WithMetadata("my-service", "0.1.0", map[string]string{
			"environment": "example",
		}),
	)
	defer p.Close(ctx)

	telemetry.Set(p)

	srv := &server{
		probe:       p,
		instruments: newInstruments(p.Meter()),
	}

	// Create a context
	m, _ := baggage.NewMember("tenancy", "testing")
	b, _ := baggage.New(m)
	ctx = baggage.ContextWithBaggage(ctx, b)

	srv.Handle(ctx)

	// Serving metrics endpoint
	http.Handle("/metrics", p)
	_ = http.ListenAndServe(":8080", nil)
}
