package telemetry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsFromEnv(t *testing.T) {
	type keyval struct {
		name  string
		value string
	}

	tests := []struct {
		name            string
		envars          []keyval
		expectedOptions options
	}{
		{
			name:   "Defaults",
			envars: []keyval{},
			expectedOptions: options{
				loggerLevel:                   "info",
				opentelemetryCollectorAddress: "localhost:55680",
				tags:                          map[string]string{},
			},
		},
		{
			name: "All",
			envars: []keyval{
				{"PROBE_NAME", "my-service"},
				{"PROBE_VERSION", "0.1.0"},
				{"PROBE_TAG_ENVIRONMENT", "testing"},
				{"PROBE_LOGGER_ENABLED", "true"},
				{"PROBE_LOGGER_LEVEL", "warn"},
				{"PROBE_PROMETHEUS_ENABLED", "true"},
				{"PROBE_JAEGER_ENABLED", "true"},
				{"PROBE_JAEGER_AGENT_HOST", "localhost"},
				{"PROBE_JAEGER_AGENT_PORT", "6832"},
				{"PROBE_JAEGER_COLLECTOR_ENDPOINT", "http://localhost:14268/api/traces"},
				{"PROBE_JAEGER_COLLECTOR_USERNAME", "username"},
				{"PROBE_JAEGER_COLLECTOR_PASSWORD", "password"},
				{"PROBE_OPENTELEMETRY_ENABLED", "true"},
				{"PROBE_OPENTELEMETRY_COLLECTOR_ADDRESS", "localhost:55680"},
			},
			expectedOptions: options{
				name:    "my-service",
				version: "0.1.0",
				tags: map[string]string{
					"environment": "testing",
				},
				loggerEnabled:                     true,
				loggerLevel:                       "warn",
				prometheusEnabled:                 true,
				jaegerEnabled:                     true,
				jaegerAgentHost:                   "localhost",
				jaegerAgentPort:                   "6832",
				jaegerCollectorEndpoint:           "http://localhost:14268/api/traces",
				jaegerCollectorUsername:           "username",
				jaegerCollectorPassword:           "password",
				opentelemetryEnabled:              true,
				opentelemetryCollectorAddress:     "localhost:55680",
				opentelemetryCollectorCredentials: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			for _, envar := range tc.envars {
				if err := os.Setenv(envar.name, envar.value); err != nil {
					t.Fatalf("Failed to set environment variable %s: %s", envar.name, err)
				}
				defer os.Unsetenv(envar.name)
			}

			options := optionsFromEnv()

			assert.Equal(t, tc.expectedOptions, options)
		})
	}
}

func TestOption(t *testing.T) {
	tests := []struct {
		name            string
		options         *options
		option          Option
		expectedOptions *options
	}{
		{
			name:    "WithMetadata",
			options: &options{},
			option: WithMetadata("my-service", "0.1.0", map[string]string{
				"environment": "testing",
			}),
			expectedOptions: &options{
				name:    "my-service",
				version: "0.1.0",
				tags: map[string]string{
					"environment": "testing",
				},
			},
		},
		{
			name:    "WithLogger",
			options: &options{},
			option:  WithLogger("warn"),
			expectedOptions: &options{
				loggerEnabled: true,
				loggerLevel:   "warn",
			},
		},
		{
			name:    "WithPrometheus",
			options: &options{},
			option:  WithPrometheus(),
			expectedOptions: &options{
				prometheusEnabled: true,
			},
		},
		{
			name:    "WithJaeger",
			options: &options{},
			option:  WithJaeger("localhost", "6832", "http://localhost:14268/api/traces", "username", "password"),
			expectedOptions: &options{
				jaegerEnabled:           true,
				jaegerAgentHost:         "localhost",
				jaegerAgentPort:         "6832",
				jaegerCollectorEndpoint: "http://localhost:14268/api/traces",
				jaegerCollectorUsername: "username",
				jaegerCollectorPassword: "password",
			},
		},
		{
			name:    "WithJaeger_Defaults",
			options: &options{},
			option:  WithJaeger("", "", "", "", ""),
			expectedOptions: &options{
				jaegerEnabled: true,
			},
		},
		{
			name:    "WithOpenTelemetry",
			options: &options{},
			option:  WithOpenTelemetry("localhost:55680", nil),
			expectedOptions: &options{
				opentelemetryEnabled:              true,
				opentelemetryCollectorAddress:     "localhost:55680",
				opentelemetryCollectorCredentials: nil,
			},
		},
		{
			name:    "WithOpenTelemetry_Defaults",
			options: &options{},
			option:  WithOpenTelemetry("", nil),
			expectedOptions: &options{
				opentelemetryEnabled:          true,
				opentelemetryCollectorAddress: "localhost:55680",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(T *testing.T) {
			tc.option(tc.options)

			assert.Equal(t, tc.expectedOptions, tc.options)
		})
	}
}
