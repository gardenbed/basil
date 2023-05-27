package telemetry

import (
	"context"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	multierror "github.com/hashicorp/go-multierror"
	prom "github.com/prometheus/client_golang/prometheus"
	promcollector "github.com/prometheus/client_golang/prometheus/collectors"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	jaegerexporter "go.opentelemetry.io/otel/exporters/jaeger"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Probe encompasses a logger, meter, and tracer.
type Probe interface {
	http.Handler
	Name() string
	Logger() Logger
	Meter() metric.Meter
	Tracer() trace.Tracer
	Close(context.Context) error
}

type (
	closeFunc func(context.Context) error
	probe     struct {
		name        string
		logger      Logger
		meter       metric.Meter
		tracer      trace.Tracer
		promHandler http.Handler
		closeFuncs  []closeFunc
	}
)

func (p *probe) Name() string {
	return p.name
}

func (p *probe) Logger() Logger {
	return p.logger
}

func (p *probe) Meter() metric.Meter {
	return p.meter
}

func (p *probe) Tracer() trace.Tracer {
	return p.tracer
}

func (p *probe) Close(ctx context.Context) error {
	var err error
	for _, close := range p.closeFuncs {
		if e := close(ctx); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return err
}

func (p *probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.promHandler != nil {
		p.promHandler.ServeHTTP(w, r)
	}
}

// NewVoidProbe creates a new no-op probe.
func NewVoidProbe() Probe {
	return &probe{
		logger: new(voidLogger),
		meter:  noop.NewMeterProvider().Meter(""),
		tracer: trace.NewNoopTracerProvider().Tracer(""),
	}
}

// NewProbe creates a new probe.
func NewProbe(opts ...Option) Probe {
	o := optionsFromEnv()
	for _, opt := range opts {
		opt(&o)
	}

	p := &probe{
		name: o.name,
	}

	if o.loggerEnabled {
		var close closeFunc
		p.logger, close = createLogger(o)
		p.closeFuncs = append(p.closeFuncs, close)
	}

	if o.prometheusEnabled {
		p.meter, p.promHandler = createPrometheus(o)
	}

	if o.jaegerEnabled {
		var close closeFunc
		p.tracer, close = createJaeger(o)
		p.closeFuncs = append(p.closeFuncs, close)
	}

	if o.opentelemetryEnabled {
		var close closeFunc
		p.meter, p.tracer, close = createOpenTelemetry(o)
		p.closeFuncs = append(p.closeFuncs, close)
	}

	// Create void logger, meter, and/or tracer if they are not created

	if p.logger == nil {
		p.logger = new(voidLogger)
	}

	if p.meter == nil {
		p.meter = noop.NewMeterProvider().Meter("")
	}

	if p.tracer == nil {
		p.tracer = trace.NewNoopTracerProvider().Tracer("")
	}

	return p
}

func createResource(o options) *resource.Resource {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(o.name),
		semconv.ServiceVersionKey.String(o.version),
	}

	for k, v := range o.tags {
		attrs = append(attrs, attribute.String(k, v))
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	)

	return resource
}

func createLogger(o options) (Logger, closeFunc) {
	config := &zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Sampling:    nil,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		InitialFields:    make(map[string]interface{}),
	}

	if o.name != "" {
		config.InitialFields["service.name"] = o.name
	}

	if o.version != "" {
		config.InitialFields["service.version"] = o.version
	}

	for k, v := range o.tags {
		config.InitialFields[k] = v
	}

	switch strings.ToLower(o.loggerLevel) {
	case "debug":
		config.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		config.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		config.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		config.Level.SetLevel(zapcore.ErrorLevel)
	case "none":
		fallthrough
	default:
		config.Level.SetLevel(zapcore.Level(99))
	}

	l, _ := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(0),
	)

	close := func(context.Context) error {
		return l.Sync()
	}

	return &zapLogger{
		config: config,
		logger: l.Sugar(),
	}, close
}

func createPrometheus(o options) (metric.Meter, http.Handler) {
	resource := createResource(o)

	exporter, err := promexporter.New()
	if err != nil {
		panic(err)
	}

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithReader(exporter),
		metricsdk.WithResource(resource),
	)

	otel.SetMeterProvider(provider)
	meter := provider.Meter(o.name)

	// Create a new Prometheus registry
	registry := prom.NewRegistry()
	registry.MustRegister(promcollector.NewGoCollector())
	registry.MustRegister(promcollector.NewProcessCollector(
		promcollector.ProcessCollectorOpts{
			Namespace: strings.ReplaceAll(o.name, "-", "_"),
		},
	))

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return meter, handler
}

func createJaeger(o options) (trace.Tracer, closeFunc) {
	resource := createResource(o)

	var endpointOpt jaegerexporter.EndpointOption
	switch {
	case o.jaegerAgentHost != "" || o.jaegerAgentPort != "":
		endpointOpt = jaegerexporter.WithAgentEndpoint(
			jaegerexporter.WithAgentHost(o.jaegerAgentHost),
			jaegerexporter.WithAgentPort(o.jaegerAgentPort),
			jaegerexporter.WithAttemptReconnectingInterval(5*time.Second),
		)
	case o.jaegerCollectorEndpoint != "" || (o.jaegerCollectorUsername != "" && o.jaegerCollectorPassword != ""):
		endpointOpt = jaegerexporter.WithCollectorEndpoint(
			jaegerexporter.WithEndpoint(o.jaegerCollectorEndpoint),
			jaegerexporter.WithUsername(o.jaegerCollectorUsername),
			jaegerexporter.WithPassword(o.jaegerCollectorPassword),
		)
	}

	exporter, err := jaegerexporter.New(endpointOpt)
	if err != nil {
		panic(err)
	}

	// TODO: Use a smarter sampler
	sampler := tracesdk.AlwaysSample()

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithResource(resource),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(sampler),
	)

	otel.SetTracerProvider(provider)
	tracer := provider.Tracer(o.name)

	return tracer, provider.Shutdown
}

func createOpenTelemetry(o options) (metric.Meter, trace.Tracer, closeFunc) {
	ctx := context.Background()
	resource := createResource(o)

	// ====================> Meter Provider <====================

	metricEndpointOpt := otlpmetricgrpc.WithEndpoint(o.opentelemetryCollectorAddress)

	var metricAuthOpt otlpmetricgrpc.Option
	if o.opentelemetryCollectorCredentials == nil {
		metricAuthOpt = otlpmetricgrpc.WithInsecure()
	} else {
		metricAuthOpt = otlpmetricgrpc.WithTLSCredentials(o.opentelemetryCollectorCredentials)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, metricEndpointOpt, metricAuthOpt)
	if err != nil {
		panic(err)
	}

	meterProvider := metricsdk.NewMeterProvider(
		metricsdk.WithReader(
			metricsdk.NewPeriodicReader(metricExporter),
		),
		metricsdk.WithResource(resource),
	)

	// ====================> Trace Provider <====================

	traceEndpointOpt := otlptracegrpc.WithEndpoint(o.opentelemetryCollectorAddress)

	var traceAuthOpt otlptracegrpc.Option
	if o.opentelemetryCollectorCredentials == nil {
		traceAuthOpt = otlptracegrpc.WithInsecure()
	} else {
		traceAuthOpt = otlptracegrpc.WithTLSCredentials(o.opentelemetryCollectorCredentials)
	}

	traceExporter, err := otlptracegrpc.New(ctx, traceEndpointOpt, traceAuthOpt)
	if err != nil {
		panic(err)
	}

	bsp := tracesdk.NewBatchSpanProcessor(traceExporter)

	// TODO: Use a smarter sampler
	sampler := tracesdk.AlwaysSample()

	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithResource(resource),
		tracesdk.WithSpanProcessor(bsp),
		tracesdk.WithSampler(sampler),
	)

	// ====================> Set Globals <====================

	otel.SetMeterProvider(meterProvider)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	meter := meterProvider.Meter(o.name)
	tracer := traceProvider.Tracer(o.name)

	close := func(ctx context.Context) error {
		g := new(errgroup.Group)
		g.Go(func() error {
			return meterProvider.Shutdown(ctx)
		})
		g.Go(func() error {
			return traceProvider.Shutdown(ctx)
		})
		return g.Wait()
	}

	return meter, tracer, close
}
