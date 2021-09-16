package telemetry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"google.golang.org/grpc/credentials"

	"github.com/stretchr/testify/assert"
)

func TestProbe_Name(t *testing.T) {
	tests := []struct {
		name  string
		probe *probe
	}{
		{
			name: "OK",
			probe: &probe{
				name: "my-service",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.probe.name, tc.probe.Name())
		})
	}
}

func TestProbe_Logger(t *testing.T) {
	tests := []struct {
		name  string
		probe *probe
	}{
		{
			name: "OK",
			probe: &probe{
				logger: &voidLogger{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.probe.logger, tc.probe.Logger())
		})
	}
}

func TestProbe_Meter(t *testing.T) {
	tests := []struct {
		name  string
		probe *probe
	}{
		{
			name: "OK",
			probe: &probe{
				meter: global.Meter(""),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.probe.meter, tc.probe.Meter())
		})
	}
}

func TestProbe_Tracer(t *testing.T) {
	tests := []struct {
		name  string
		probe *probe
	}{
		{
			name: "OK",
			probe: &probe{
				tracer: otel.Tracer(""),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.probe.tracer, tc.probe.Tracer())
		})
	}
}

func TestProbe_Close(t *testing.T) {
	tests := []struct {
		name          string
		probe         *probe
		ctx           context.Context
		expectedError string
	}{
		{
			name: "Success",
			probe: &probe{
				closeFuncs: []closeFunc{
					func(context.Context) error {
						return nil
					},
				},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name: "Fail",
			probe: &probe{
				closeFuncs: []closeFunc{
					func(context.Context) error {
						return errors.New("error on closing")
					},
				},
			},
			ctx:           context.Background(),
			expectedError: "error on closing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.probe.Close(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}

func TestProbe_ServeHTTP(t *testing.T) {
	tests := []struct {
		name               string
		probe              *probe
		req                *http.Request
		expectedStatusCode int
	}{
		{
			name: "OK",
			probe: &probe{
				promHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			req:                httptest.NewRequest("GET", "/metrics", nil),
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			tc.probe.ServeHTTP(resp, tc.req)

			statusCode := resp.Result().StatusCode
			assert.Equal(t, tc.expectedStatusCode, statusCode)
		})
	}
}

func TestNewVoidProbe(t *testing.T) {
	probe := NewVoidProbe()
	assert.NotNil(t, probe)
}

func TestNewProbe(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "NoOption",
			opts: []Option{},
		},
		{
			name: "PrometheusAndJaeger",
			opts: []Option{
				WithMetadata("my-service", "0.1.0", map[string]string{
					"environment": "testing",
				}),
				WithLogger("warn"),
				WithPrometheus(),
				WithJaeger("localhost", "6832", "", "", ""),
			},
		},
		{
			name: "OpenTelemetry",
			opts: []Option{
				WithMetadata("my-service", "0.1.0", map[string]string{
					"environment": "testing",
				}),
				WithLogger("warn"),
				WithOpenTelemetry("localhost:55680", nil),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(T *testing.T) {
			probe := NewProbe(tc.opts...)
			defer probe.Close(context.Background())

			assert.NotNil(t, probe)
			assert.NotNil(t, probe.Logger())
			assert.NotNil(t, probe.Meter())
			assert.NotNil(t, probe.Tracer())
		})
	}
}

func TestCreateLogger(t *testing.T) {
	tests := []struct {
		name          string
		options       options
		expectedLevel Level
	}{
		{
			name: "Production",
			options: options{
				name:    "my-service",
				version: "0.1.0",
				tags: map[string]string{
					"environment": "testing",
				},
				loggerLevel: "warn",
			},
			expectedLevel: LevelWarn,
		},
		{
			name: "LogLevelDebug",
			options: options{
				name:        "my-service",
				loggerLevel: "debug",
			},
			expectedLevel: LevelDebug,
		},
		{
			name: "LogLevelInfo",
			options: options{
				name:        "my-service",
				loggerLevel: "info",
			},
			expectedLevel: LevelInfo,
		},
		{
			name: "LogLevelWarn",
			options: options{
				name:        "my-service",
				loggerLevel: "warn",
			},
			expectedLevel: LevelWarn,
		},
		{
			name: "LogLevelError",
			options: options{
				name:        "my-service",
				loggerLevel: "error",
			},
			expectedLevel: LevelError,
		},
		{
			name: "LogLevelNone",
			options: options{
				name:        "my-service",
				loggerLevel: "none",
			},
			expectedLevel: LevelNone,
		},
		{
			name: "InvalidLogLevel",
			options: options{
				name:          "my-service",
				loggerEnabled: true,
				loggerLevel:   "invalid",
			},
			expectedLevel: LevelNone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(T *testing.T) {
			logger, close := createLogger(tc.options)
			defer close(context.Background())

			assert.NotNil(t, logger)
			assert.NotNil(t, close)
			assert.Equal(t, tc.expectedLevel, logger.Level())
		})
	}
}

func TestCreatePrometheus(t *testing.T) {
	tests := []struct {
		name    string
		options options
	}{
		{
			name: "Production",
			options: options{
				name:              "my-service",
				prometheusEnabled: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meter, handler := createPrometheus(tc.options)

			assert.NotNil(t, meter)
			assert.NotNil(t, handler)
		})
	}
}

func TestCreateJaeger(t *testing.T) {
	tests := []struct {
		name    string
		options options
	}{
		{
			name: "WithAgent",
			options: options{
				name: "my-service",
				tags: map[string]string{
					"environment": "testing",
				},
				jaegerEnabled:   true,
				jaegerAgentHost: "localhost",
				jaegerAgentPort: "6832",
			},
		},
		{
			name: "WithCollector",
			options: options{
				name: "my-service",
				tags: map[string]string{
					"environment": "testing",
				},
				jaegerEnabled:           true,
				jaegerCollectorEndpoint: "http://localhost:14268/api/traces",
				jaegerCollectorUsername: "username",
				jaegerCollectorPassword: "password",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer, close := createJaeger(tc.options)
			defer close(context.Background())

			assert.NotNil(t, tracer)
			assert.NotNil(t, close)
		})
	}
}

func TestCreateOpenTelemetry(t *testing.T) {
	tests := []struct {
		name    string
		options options
	}{
		{
			name: "Insecure",
			options: options{
				name:                          "my-service",
				opentelemetryEnabled:          true,
				opentelemetryCollectorAddress: "localhost:55680",
			},
		},
		{
			name: "Secure",
			options: options{
				name:                              "my-service",
				opentelemetryEnabled:              true,
				opentelemetryCollectorAddress:     "localhost:55680",
				opentelemetryCollectorCredentials: credentials.NewTLS(nil),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meter, tracer, close := createOpenTelemetry(tc.options)
			defer close(context.Background())

			assert.NotNil(t, meter)
			assert.NotNil(t, tracer)
			assert.NotNil(t, close)
		})
	}
}
