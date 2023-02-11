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
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gardenbed/basil/telemetry"
)

// Server-side instruments for measurements.
type serverInstruments struct {
	panic   instrument.Int64Counter
	total   instrument.Int64Counter
	active  instrument.Int64UpDownCounter
	latency instrument.Int64Histogram
}

func newServerInstruments(m metric.Meter) *serverInstruments {
	panic, _ := m.Int64Counter(
		"incoming_grpc_requests_panic",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("The total number of panics happened in grpc handlers (server-side)"),
	)

	total, _ := m.Int64Counter(
		"incoming_grpc_requests_total",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("The total number of incoming grpc requests (server-side)"),
	)

	active, _ := m.Int64UpDownCounter(
		"incoming_grpc_requests_active",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("The number of in-flight incoming grpc requests (server-side)"),
	)

	latency, _ := m.Int64Histogram(
		"incoming_grpc_requests_latency",
		instrument.WithUnit(unit.Milliseconds),
		instrument.WithDescription("The duration of incoming grpc requests in milliseconds (server-side)"),
	)

	return &serverInstruments{
		panic:   panic,
		total:   total,
		active:  active,
		latency: latency,
	}
}

// ServerInterceptor creates interceptors with logging, metrics, and tracing for grpc servers.
type ServerInterceptor struct {
	probe       telemetry.Probe
	opts        Options
	instruments *serverInstruments
}

// NewServerInterceptor creates a new observable server interceptor.
func NewServerInterceptor(probe telemetry.Probe, opts Options) *ServerInterceptor {
	opts = opts.withDefaults()
	instruments := newServerInstruments(probe.Meter())

	return &ServerInterceptor{
		probe:       probe,
		opts:        opts,
		instruments: instruments,
	}
}

// ServerOptions return grpc server options for unary and stream interceptors.
// This can be used for making gRPC method handlers observable via logging, metrics, tracing, etc.
// It also observes and recovers panics that happened inside the method handlers.
func (i *ServerInterceptor) ServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(i.unaryInterceptor),
		grpc.StreamInterceptor(i.streamInterceptor),
	}
}

func (i *ServerInterceptor) callUnaryHandler(ctx context.Context, handler grpc.UnaryHandler, req interface{}) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
			i.probe.Logger().Errorf("Panic recovered: %v", r)
			i.instruments.panic.Add(ctx, 1)
		}
	}()

	resp, err = handler(ctx, req)

	return resp, err
}

func (i *ServerInterceptor) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	kind := "server"
	stream := false

	meter := i.probe.Meter()
	tracer := i.probe.Tracer()

	// Get the package, service, and method name for the request
	e, ok := parseEndpoint(info.FullMethod)
	if !ok {
		return i.callUnaryHandler(ctx, handler, req)
	}

	// Check excluded methods
	for _, m := range i.opts.ExcludedMethods {
		if e.Method == m {
			return i.callUnaryHandler(ctx, handler, req)
		}
	}

	packageAttr := attribute.String("package", e.Package)
	serviceAttr := attribute.String("service", e.Service)
	methodAttr := attribute.String("method", e.Method)
	streamAttr := attribute.Bool("stream", stream)

	// Handle the number of in-flight requests
	i.instruments.active.Add(ctx, 1, packageAttr, serviceAttr, methodAttr, streamAttr)
	defer i.instruments.active.Add(ctx, -1, packageAttr, serviceAttr, methodAttr, streamAttr)

	// Get grpc request metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	// Make sure the request has a UUID
	var requestUUID string
	if vals := md.Get(requestUUIDKey); len(vals) > 0 {
		requestUUID = vals[0]
	}
	if requestUUID == "" {
		requestUUID = uuid.New().String()
		md.Set(requestUUIDKey, requestUUID)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	// Get the name of client for the request if any
	var clientName string
	if vals := md.Get(clientNameKey); len(vals) > 0 {
		clientName = vals[0]
	}

	// Propagate request metadata by adding them to outgoing grpc response metadata
	header := metadata.New(map[string]string{
		requestUUIDKey: requestUUID,
		clientNameKey:  clientName,
	})
	_ = grpc.SendHeader(ctx, header)

	// Extract context from the grpc metadata
	ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataTextMapCarrier{md: &md})

	// Start a new span
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("%s (server unary)", e.Method),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	// Create a contextualized logger
	contextFields := []interface{}{
		"req.uuid", requestUUID,
		"req.kind", kind,
		"req.package", e.Package,
		"req.service", e.Service,
		"req.method", e.Method,
		"req.stream", stream,
		"traceId", span.SpanContext().TraceID().String(),
		"spanId", span.SpanContext().SpanID().String(),
	}
	if clientName != "" {
		contextFields = append(contextFields, "client.name", clientName)
	}
	logger := i.probe.Logger().With(contextFields...)

	// Augment the request context
	ctx = telemetry.ContextWithUUID(ctx, requestUUID)
	ctx = telemetry.ContextWithLogger(ctx, logger)
	ctx = telemetry.ContextWithMeter(ctx, meter)
	ctx = telemetry.ContextWithTracer(ctx, tracer)

	// Call gRPC method handler
	span.AddEvent("calling grpc method handler")
	res, err := i.callUnaryHandler(ctx, handler, req)

	duration := time.Since(startTime).Milliseconds()
	success := err == nil

	// Report metrics
	successAttr := attribute.Bool("success", success)
	i.instruments.total.Add(ctx, 1, packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)
	i.instruments.latency.Record(ctx, duration, packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)

	// Report logs
	message := fmt.Sprintf("%s %s %dms", kind, e, duration)
	fields := []interface{}{
		"resp.success", success,
		"resp.duration", duration,
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

	return res, err
}

func (i *ServerInterceptor) callStreamHandler(ctx context.Context, handler grpc.StreamHandler, srv interface{}, stream grpc.ServerStream) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
			i.probe.Logger().Errorf("Panic recovered: %v", r)
			i.instruments.panic.Add(ctx, 1)
		}
	}()

	err = handler(srv, stream)

	return err
}

func (i *ServerInterceptor) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	ctx := ss.Context()
	kind := "server"
	stream := true

	meter := i.probe.Meter()
	tracer := i.probe.Tracer()

	// Get the package, service, and method name for the request
	e, ok := parseEndpoint(info.FullMethod)
	if !ok {
		return i.callStreamHandler(ctx, handler, srv, ss)
	}

	// Check excluded methods
	for _, m := range i.opts.ExcludedMethods {
		if e.Method == m {
			return i.callStreamHandler(ctx, handler, srv, ss)
		}
	}

	packageAttr := attribute.String("package", e.Package)
	serviceAttr := attribute.String("service", e.Service)
	methodAttr := attribute.String("method", e.Method)
	streamAttr := attribute.Bool("stream", stream)

	// Handle the number of in-flight requests
	i.instruments.active.Add(ctx, 1, packageAttr, serviceAttr, methodAttr, streamAttr)
	defer i.instruments.active.Add(ctx, -1, packageAttr, serviceAttr, methodAttr, streamAttr)

	// Get grpc request metadata (an incoming grpc request context is guaranteed to have metadata)
	md, _ := metadata.FromIncomingContext(ctx)
	md = md.Copy()

	// Make sure the request has a UUID
	var requestUUID string
	if vals := md.Get(requestUUIDKey); len(vals) > 0 {
		requestUUID = vals[0]
	}
	if requestUUID == "" {
		requestUUID = uuid.New().String()
		md.Set(requestUUIDKey, requestUUID)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	// Get the name of client for the request if any
	var clientName string
	if vals := md.Get(clientNameKey); len(vals) > 0 {
		clientName = vals[0]
	}

	// Propagate request metadata by adding them to outgoing grpc response metadata
	header := metadata.New(map[string]string{
		requestUUIDKey: requestUUID,
		clientNameKey:  clientName,
	})
	_ = ss.SendHeader(header)

	// Extract context from the grpc metadata
	ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataTextMapCarrier{md: &md})

	// Start a new span
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("%s (server stream)", e.Method),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	// Create a contextualized logger
	contextFields := []interface{}{
		"req.uuid", requestUUID,
		"req.kind", kind,
		"req.package", e.Package,
		"req.service", e.Service,
		"req.method", e.Method,
		"req.stream", stream,
		"traceId", span.SpanContext().TraceID().String(),
		"spanId", span.SpanContext().SpanID().String(),
	}
	if clientName != "" {
		contextFields = append(contextFields, "client.name", clientName)
	}
	logger := i.probe.Logger().With(contextFields...)

	// Augment the request context
	ctx = telemetry.ContextWithUUID(ctx, requestUUID)
	ctx = telemetry.ContextWithLogger(ctx, logger)
	ctx = telemetry.ContextWithMeter(ctx, meter)
	ctx = telemetry.ContextWithTracer(ctx, tracer)
	ss = ServerStreamWithContext(ctx, ss)

	// Call gRPC method handler
	span.AddEvent("calling grpc method handler")
	err := i.callStreamHandler(ctx, handler, srv, ss)

	duration := time.Since(startTime).Milliseconds()
	success := err == nil

	// Report metrics
	successAttr := attribute.Bool("success", success)
	i.instruments.total.Add(ctx, 1, packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)
	i.instruments.latency.Record(ctx, duration, packageAttr, serviceAttr, methodAttr, streamAttr, successAttr)

	// Report logs
	message := fmt.Sprintf("%s %s %dms", kind, e, duration)
	fields := []interface{}{
		"resp.success", success,
		"resp.duration", duration,
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
