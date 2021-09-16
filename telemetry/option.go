package telemetry

import (
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/credentials"
)

type options struct {
	// Metadata
	name    string
	version string
	tags    map[string]string

	// Logger
	loggerEnabled bool
	loggerLevel   string

	// Prometheus
	prometheusEnabled bool

	// Jaeger
	jaegerEnabled           bool
	jaegerAgentHost         string
	jaegerAgentPort         string
	jaegerCollectorEndpoint string
	jaegerCollectorUsername string
	jaegerCollectorPassword string

	// OpenTelemetry
	opentelemetryEnabled              bool
	opentelemetryCollectorAddress     string
	opentelemetryCollectorCredentials credentials.TransportCredentials
}

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
	o.loggerEnabled, _ = strconv.ParseBool(os.Getenv("PROBE_LOGGER_ENABLED"))
	o.loggerLevel = os.Getenv("PROBE_LOGGER_LEVEL")
	if o.loggerLevel == "" {
		o.loggerLevel = "info"
	}

	// Prometheus
	o.prometheusEnabled, _ = strconv.ParseBool(os.Getenv("PROBE_PROMETHEUS_ENABLED"))

	// Jaeger
	o.jaegerEnabled, _ = strconv.ParseBool(os.Getenv("PROBE_JAEGER_ENABLED"))
	o.jaegerAgentHost = os.Getenv("PROBE_JAEGER_AGENT_HOST")
	o.jaegerAgentPort = os.Getenv("PROBE_JAEGER_AGENT_PORT")
	o.jaegerCollectorEndpoint = os.Getenv("PROBE_JAEGER_COLLECTOR_ENDPOINT")
	o.jaegerCollectorUsername = os.Getenv("PROBE_JAEGER_COLLECTOR_USERNAME")
	o.jaegerCollectorPassword = os.Getenv("PROBE_JAEGER_COLLECTOR_PASSWORD")

	// OpenTelemetry
	o.opentelemetryEnabled, _ = strconv.ParseBool(os.Getenv("PROBE_OPENTELEMETRY_ENABLED"))
	o.opentelemetryCollectorAddress = os.Getenv("PROBE_OPENTELEMETRY_COLLECTOR_ADDRESS")
	if o.opentelemetryCollectorAddress == "" {
		o.opentelemetryCollectorAddress = "localhost:55680"
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
		o.loggerEnabled = true
		o.loggerLevel = level
	}
}

// WithPrometheus is the option for enabling Prometheus.
func WithPrometheus() Option {
	return func(o *options) {
		o.prometheusEnabled = true
	}
}

// WithJaeger is the option for enabling Jaeger.
// Only one of agentEndpoint or collectorEndpoint is required.
// collectorUsername and collectorPassword are optional.
// The default agent endpoint is localhost:6832.
func WithJaeger(agentHost, agentPort, collectorEndpoint, collectorUsername, collectorPassword string) Option {
	return func(o *options) {
		o.jaegerEnabled = true
		o.jaegerAgentHost = agentHost
		o.jaegerAgentPort = agentPort
		o.jaegerCollectorEndpoint = collectorEndpoint
		o.jaegerCollectorUsername = collectorUsername
		o.jaegerCollectorPassword = collectorPassword
	}
}

// WithOpenTelemetry is the option for enabling OpenTelemetry Collector.
// collectorCredentials is optional. If not specified, the connection will be insecure.
// The default collector address is localhost:55680.
func WithOpenTelemetry(collectorAddress string, collectorCredentials credentials.TransportCredentials) Option {
	if collectorAddress == "" {
		collectorAddress = "localhost:55680"
	}

	return func(o *options) {
		o.opentelemetryEnabled = true
		o.opentelemetryCollectorAddress = collectorAddress
		o.opentelemetryCollectorCredentials = collectorCredentials
	}
}
