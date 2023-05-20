package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gardenbed/basil/telemetry"
)

// Client-side instruments for measurements.
type clientInstruments struct {
	total   instrument.Int64Counter
	active  instrument.Int64UpDownCounter
	latency instrument.Int64Histogram
}

func newClientInstruments(m metric.Meter) *clientInstruments {
	total, _ := m.Int64Counter(
		"outgoing_grpc_requests_total",
		instrument.WithDescription("The total number of outgoing grpc requests (client-side)"),
	)

	active, _ := m.Int64UpDownCounter(
		"outgoing_grpc_requests_active",
		instrument.WithDescription("The number of in-flight outgoing grpc requests (client-side)"),
	)

	latency, _ := m.Int64Histogram(
		"outgoing_grpc_requests_latency",
		instrument.WithUnit("ms"),
		instrument.WithDescription("The duration of outgoing grpc requests in milliseconds (client-side)"),
	)

	return &clientInstruments{
		total:   total,
		active:  active,
		latency: latency,
	}
}

// ClientInterceptor is used for creating interceptors with logging, metrics, and tracing for grpc clients.
type ClientInterceptor struct {
	probe       telemetry.Probe
	opts        Options
	instruments *clientInstruments
}

// NewClientInterceptor creates a new observable client interceptor.
func NewClientInterceptor(probe telemetry.Probe, opts Options) *ClientInterceptor {
	opts = opts.withDefaults()
	instruments := newClientInstruments(probe.Meter())

	return &ClientInterceptor{
		probe:       probe,
		opts:        opts,
		instruments: instruments,
	}
}

// DialOptions return grpc dial options for unary and stream interceptors.
// This can be used for making gRPC method calls observable via logging, metrics, tracing, etc.
func (i *ClientInterceptor) DialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(i.unaryInterceptor),
		grpc.WithStreamInterceptor(i.streamInterceptor),
	}
}

func (i *ClientInterceptor) unaryInterceptor(ctx context.Context, fullMethod string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	startTime := time.Now()
	kind := "client"
	stream := false

	tracer := i.probe.Tracer()

	// Get the package, service, and method name for the request
	e, ok := parseEndpoint(fullMethod)
	if !ok {
		return invoker(ctx, fullMethod, req, res, cc, opts...)
	}

	// Check excluded methods
	for _, m := range i.opts.ExcludedMethods {
		if e.Method == m {
			return invoker(ctx, fullMethod, req, res, cc, opts...)
		}
	}

	packageAttr := attribute.String("package", e.Package)
	serviceAttr := attribute.String("service", e.Service)
	methodAttr := attribute.String("method", e.Method)
	streamAttr := attribute.Bool("stream", stream)

	// Handle the number of in-flight requests
	reqOpt := metric.WithAttributes(packageAttr, serviceAttr, methodAttr, streamAttr)
	i.instruments.active.Add(ctx, 1, reqOpt)
	defer i.instruments.active.Add(ctx, -1, reqOpt)

	// Make sure the request has a UUID
	requestUUID, ok := telemetry.UUIDFromContext(ctx)
	if !ok || requestUUID == "" {
		requestUUID = uuid.New().String()
	}

	// Get grpc request metadata
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	// Propagate request metadata by adding them to outgoing grpc request metadata
	md.Set(requestUUIDKey, requestUUID)
	md.Set(clientNameKey, i.probe.Name())
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Start a new span
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("%s (client unary)", e.Method),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	// Inject the context and the span context into the grpc metadata
	otel.GetTextMapPropagator().Inject(ctx, &metadataTextMapCarrier{md: &md})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Call gRPC method invoker
	span.AddEvent("invoking grpc method")
	err := invoker(ctx, fullMethod, req, res, cc, opts...)

	duration := time.Since(startTime).Milliseconds()
	success := err == nil

	// Report metrics
	successAttr := attribute.Bool("success", success)
	resOpt := metric.WithAttributes(packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)
	i.instruments.total.Add(ctx, 1, resOpt)
	i.instruments.latency.Record(ctx, duration, resOpt)

	// Report logs
	logger := i.probe.Logger()
	message := fmt.Sprintf("%s %s %dms", kind, e, duration)
	fields := []interface{}{
		"req.uuid", requestUUID,
		"req.kind", kind,
		"req.package", e.Package,
		"req.service", e.Service,
		"req.method", e.Method,
		"req.stream", stream,
		"resp.success", success,
		"resp.duration", duration,
		"traceId", span.SpanContext().TraceID().String(),
		"spanId", span.SpanContext().SpanID().String(),
	}
	if err != nil {
		fields = append(fields, "grpc.error", err.Error())
	}

	// Determine the log level based on the result
	if success {
		logger.Info(message, fields...)
	} else {
		logger.Error(message, fields...)
	}

	// Report the span
	span.SetAttributes(packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)
	if err != nil {
		code := codes.Code(status.Code(err))
		span.SetStatus(code, err.Error())
	}

	return err
}

func (i *ClientInterceptor) streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	startTime := time.Now()
	kind := "client"
	stream := true

	tracer := i.probe.Tracer()

	// Get the package, service, and method name for the request
	e, ok := parseEndpoint(fullMethod)
	if !ok {
		return streamer(ctx, desc, cc, fullMethod, opts...)
	}

	// Check excluded methods
	for _, m := range i.opts.ExcludedMethods {
		if e.Method == m {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}
	}

	packageAttr := attribute.String("package", e.Package)
	serviceAttr := attribute.String("service", e.Service)
	methodAttr := attribute.String("method", e.Method)
	streamAttr := attribute.Bool("stream", stream)

	// Handle the number of in-flight requests
	reqOpt := metric.WithAttributes(packageAttr, serviceAttr, methodAttr, streamAttr)
	i.instruments.active.Add(ctx, 1, reqOpt)
	i.instruments.active.Add(ctx, -1, reqOpt)

	// Make sure the request has a UUID
	requestUUID, ok := telemetry.UUIDFromContext(ctx)
	if !ok || requestUUID == "" {
		requestUUID = uuid.New().String()
	}

	// Get grpc request metadata
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	// Propagate request metadata by adding them to outgoing grpc request metadata
	md.Set(requestUUIDKey, requestUUID)
	md.Set(clientNameKey, i.probe.Name())
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Start a new span
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("%s (client stream)", e.Method),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	// Inject the context and the span context into the grpc metadata
	otel.GetTextMapPropagator().Inject(ctx, &metadataTextMapCarrier{md: &md})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Call gRPC method streamer
	span.AddEvent("invoking grpc method")
	cs, err := streamer(ctx, desc, cc, fullMethod, opts...)

	duration := time.Since(startTime).Milliseconds()
	success := err == nil

	// Report metrics
	successAttr := attribute.Bool("success", success)
	resOpt := metric.WithAttributes(packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)
	i.instruments.total.Add(ctx, 1, resOpt)
	i.instruments.latency.Record(ctx, duration, resOpt)

	// Report logs
	logger := i.probe.Logger()
	message := fmt.Sprintf("%s %s %dms", kind, e, duration)
	fields := []interface{}{
		"req.uuid", requestUUID,
		"req.kind", kind,
		"req.package", e.Package,
		"req.service", e.Service,
		"req.method", e.Method,
		"req.stream", stream,
		"resp.success", success,
		"resp.duration", duration,
		"traceId", span.SpanContext().TraceID().String(),
		"spanId", span.SpanContext().SpanID().String(),
	}
	if err != nil {
		fields = append(fields, "grpc.error", err.Error())
	}

	// Determine the log level based on the result
	if success {
		logger.Info(message, fields...)
	} else {
		logger.Error(message, fields...)
	}

	// Report the span
	span.SetAttributes(packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)
	if err != nil {
		code := codes.Code(status.Code(err))
		span.SetStatus(code, err.Error())
	}

	return cs, err
}
