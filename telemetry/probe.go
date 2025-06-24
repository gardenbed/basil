package telemetry

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	multierror "github.com/hashicorp/go-multierror"
	prom "github.com/prometheus/client_golang/prometheus"
	promcollector "github.com/prometheus/client_golang/prometheus/collectors"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
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
		meter:  metricnoop.NewMeterProvider().Meter(""),
		tracer: tracenoop.NewTracerProvider().Tracer(""),
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

	if o.logger.enabled {
		var close closeFunc
		p.logger, close = createLogger(o)
		p.closeFuncs = append(p.closeFuncs, close)
	}

	if o.prometheus.enabled {
		p.meter, p.promHandler = createPrometheus(o)
	}

	if o.opentelemetry.meterEnabled {
		var close closeFunc
		p.meter, close = createOpenTelemetryMeter(o)
		p.closeFuncs = append(p.closeFuncs, close)
	}

	if o.opentelemetry.tracerEnabled {
		var close closeFunc
		p.tracer, close = createOpenTelemetryTracer(o)
		p.closeFuncs = append(p.closeFuncs, close)
	}

	// Create void logger, meter, and/or tracer if they are not created

	if p.logger == nil {
		p.logger = new(voidLogger)
	}

	if p.meter == nil {
		p.meter = metricnoop.NewMeterProvider().Meter("")
	}

	if p.tracer == nil {
		p.tracer = tracenoop.NewTracerProvider().Tracer("")
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

	switch strings.ToLower(o.logger.level) {
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

	l := zap.Must(config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(0),
	))

	close := func(context.Context) error {
		err := l.Sync()

		// This is a workaround for this issue: https://github.com/uber-go/zap/issues/880
		if err != nil && strings.Contains(err.Error(), "sync /dev/stdout") {
			err = nil
		}

		return err
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

func createOpenTelemetryMeter(o options) (metric.Meter, closeFunc) {
	ctx := context.Background()
	resource := createResource(o)

	// ====================> Meter Provider <====================

	metricEndpointOpt := otlpmetricgrpc.WithEndpoint(o.opentelemetry.collectorAddress)

	var metricAuthOpt otlpmetricgrpc.Option
	if o.opentelemetry.collectorCredentials == nil {
		metricAuthOpt = otlpmetricgrpc.WithInsecure()
	} else {
		metricAuthOpt = otlpmetricgrpc.WithTLSCredentials(o.opentelemetry.collectorCredentials)
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

	// ====================> Set Globals <====================

	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter(o.name)
	close := meterProvider.Shutdown

	return meter, close
}

func createOpenTelemetryTracer(o options) (trace.Tracer, closeFunc) {
	ctx := context.Background()
	resource := createResource(o)

	// ====================> Trace Provider <====================

	traceEndpointOpt := otlptracegrpc.WithEndpoint(o.opentelemetry.collectorAddress)

	var traceAuthOpt otlptracegrpc.Option
	if o.opentelemetry.collectorCredentials == nil {
		traceAuthOpt = otlptracegrpc.WithInsecure()
	} else {
		traceAuthOpt = otlptracegrpc.WithTLSCredentials(o.opentelemetry.collectorCredentials)
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

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer := traceProvider.Tracer(o.name)
	close := traceProvider.Shutdown

	return tracer, close
}
