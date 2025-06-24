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
				logger: logger{
					level: "info",
				},
				opentelemetry: opentelemetry{
					collectorAddress: "localhost:55680",
				},
				tags: map[string]string{},
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
				{"PROBE_OPENTELEMETRY_METER_ENABLED", "true"},
				{"PROBE_OPENTELEMETRY_TRACER_ENABLED", "true"},
				{"PROBE_OPENTELEMETRY_COLLECTOR_ADDRESS", "localhost:55680"},
			},
			expectedOptions: options{
				name:    "my-service",
				version: "0.1.0",
				tags: map[string]string{
					"environment": "testing",
				},
				logger: logger{
					enabled: true,
					level:   "warn",
				},
				prometheus: prometheus{
					enabled: true,
				},
				opentelemetry: opentelemetry{
					meterEnabled:         true,
					tracerEnabled:        true,
					collectorAddress:     "localhost:55680",
					collectorCredentials: nil,
				},
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

				defer func() {
					assert.NoError(t, os.Unsetenv(envar.name))
				}()
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
				logger: logger{
					enabled: true,
					level:   "warn",
				},
			},
		},
		{
			name:    "WithPrometheus",
			options: &options{},
			option:  WithPrometheus(),
			expectedOptions: &options{
				prometheus: prometheus{
					enabled: true,
				},
			},
		},
		{
			name:    "WithOpenTelemetry",
			options: &options{},
			option:  WithOpenTelemetry(true, true, "localhost:55680", nil),
			expectedOptions: &options{
				opentelemetry: opentelemetry{
					meterEnabled:         true,
					tracerEnabled:        true,
					collectorAddress:     "localhost:55680",
					collectorCredentials: nil,
				},
			},
		},
		{
			name:    "WithOpenTelemetry_Defaults",
			options: &options{},
			option:  WithOpenTelemetry(true, true, "", nil),
			expectedOptions: &options{
				opentelemetry: opentelemetry{
					meterEnabled:     true,
					tracerEnabled:    true,
					collectorAddress: "localhost:55680",
				},
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
