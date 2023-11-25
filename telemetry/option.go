package telemetry

import (
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/credentials"
)

type (
	options struct {
		// Metadata
		name    string
		version string
		tags    map[string]string

		// Logger
		logger

		// Prometheus
		prometheus

		// OpenTelemetry
		opentelemetry
	}

	logger struct {
		enabled bool
		level   string
	}

	prometheus struct {
		enabled bool
	}

	opentelemetry struct {
		meterEnabled         bool
		tracerEnabled        bool
		collectorAddress     string
		collectorCredentials credentials.TransportCredentials
	}
)

func optionsFromEnv() options {
	o := options{}

	// Metadata
	o.name = os.Getenv("PROBE_NAME")
	o.version = os.Getenv("PROBE_VERSION")
	o.tags = map[string]string{}

	for _, env := range os.Environ() {
		pair := strings.Split(env, "=")
		if strings.HasPrefix(pair[0], "PROBE_TAG_") {
			tag := strings.TrimPrefix(pair[0], "PROBE_TAG_")
			tag = strings.ToLower(tag)
			o.tags[tag] = pair[1]
		}
	}

	// Logger
	o.logger.enabled, _ = strconv.ParseBool(os.Getenv("PROBE_LOGGER_ENABLED"))
	o.logger.level = os.Getenv("PROBE_LOGGER_LEVEL")
	if o.logger.level == "" {
		o.logger.level = "info"
	}

	// Prometheus
	o.prometheus.enabled, _ = strconv.ParseBool(os.Getenv("PROBE_PROMETHEUS_ENABLED"))

	// OpenTelemetry
	o.opentelemetry.meterEnabled, _ = strconv.ParseBool(os.Getenv("PROBE_OPENTELEMETRY_METER_ENABLED"))
	o.opentelemetry.tracerEnabled, _ = strconv.ParseBool(os.Getenv("PROBE_OPENTELEMETRY_TRACER_ENABLED"))
	o.opentelemetry.collectorAddress = os.Getenv("PROBE_OPENTELEMETRY_COLLECTOR_ADDRESS")
	if o.opentelemetry.collectorAddress == "" {
		o.opentelemetry.collectorAddress = "localhost:55680"
	}

	return o
}

// Option is used for configuring a probe.
type Option func(*options)

// WithMetadata is the option for specifying and reporting metadata.
// All arguments are optional.
func WithMetadata(name, version string, tags map[string]string) Option {
	return func(o *options) {
		o.name = name
		o.version = version
		o.tags = tags
	}
}

// WithLogger is the option for enabling the logger.
func WithLogger(level string) Option {
	return func(o *options) {
		o.logger.enabled = true
		o.logger.level = level
	}
}

// WithPrometheus is the option for enabling Prometheus.
func WithPrometheus() Option {
	return func(o *options) {
		o.prometheus.enabled = true
	}
}

// WithOpenTelemetry is the option for enabling OpenTelemetry Collector.
// collectorCredentials is optional. If not specified, the connection will be insecure.
// The default collector address is localhost:55680.
func WithOpenTelemetry(meterEnabled, tracerEnabled bool, collectorAddress string, collectorCredentials credentials.TransportCredentials) Option {
	if collectorAddress == "" {
		collectorAddress = "localhost:55680"
	}

	return func(o *options) {
		o.opentelemetry.meterEnabled = meterEnabled
		o.opentelemetry.tracerEnabled = tracerEnabled
		o.opentelemetry.collectorAddress = collectorAddress
		o.opentelemetry.collectorCredentials = collectorCredentials
	}
}
