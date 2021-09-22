package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/trace"

	"github.com/gardenbed/basil/telemetry"
)

// Client-side instruments for measurements.
type clientInstruments struct {
	total   metric.Int64Counter
	active  metric.Int64UpDownCounter
	latency metric.Int64Histogram
}

func newClientInstruments(meter metric.Meter) *clientInstruments {
	mm := metric.Must(meter)

	return &clientInstruments{
		total: mm.NewInt64Counter(
			"outgoing_http_requests_total",
			metric.WithUnit(unit.Dimensionless),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The total number of outgoing http requests (client-side)"),
		),
		active: mm.NewInt64UpDownCounter(
			"outgoing_http_requests_active",
			metric.WithUnit(unit.Dimensionless),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The number of in-flight outgoing http requests (client-side)"),
		),
		latency: mm.NewInt64Histogram(
			"outgoing_http_requests_latency",
			metric.WithUnit(unit.Milliseconds),
			metric.WithInstrumentationName(libraryName),
			metric.WithDescription("The duration of outgoing http requests in milliseconds (client-side)"),
		),
	}
}

// Client is a drop-in replacement for the standard http.Client.
// It is an observable http client with logging, metrics, and tracing.
type Client struct {
	client      *http.Client
	probe       telemetry.Probe
	opts        Options
	instruments *clientInstruments
}

// NewClient creates a new observable http client.
func NewClient(client *http.Client, probe telemetry.Probe, opts Options) *Client {
	opts = opts.withDefaults()
	instruments := newClientInstruments(probe.Meter())

	return &Client{
		client:      client,
		probe:       probe,
		opts:        opts,
		instruments: instruments,
	}
}

// CloseIdleConnections is the observable counterpart of standard http Client.CloseIdleConnections.
func (c *Client) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

// Get is the observable counterpart of standard http Client.Get.
// Using this method, request context (UUID and trace) will be auto-generated.
// If you have a context for the request, consider using the Do method.
func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

// Head is the observable counterpart of standard http Client.Head.
// Using this method, request context (UUID and trace) will be auto-generated.
// If you have a context for the request, consider using the Do method.
func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

// Post is the observable counterpart of standard http Client.Post.
// Using this method, request context (UUID and trace) will be auto-generated.
// If you have a context for the request, consider using the Do method.
func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return c.client.Do(req)
}

// PostForm is the observable counterpart of standard http Client.PostForm.
// Using this method, request context (UUID and trace) will be auto-generated.
// If you have a context for the request, consider using the Do method.
func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	contentType := "application/x-www-form-urlencoded"
	body := strings.NewReader(data.Encode())
	return c.Post(url, contentType, body)
}

// Do is the observable counterpart of standard http Client.Do.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	ctx := req.Context()
	kind := "client"
	method := req.Method
	url := req.URL.Path
	route := c.opts.IDRegexp.ReplaceAllString(url, ":id")

	meter := c.probe.Meter()
	tracer := c.probe.Tracer()

	// Increase the number of in-flight requests
	c.instruments.active.Add(ctx, 1,
		attribute.String("method", method),
		attribute.String("route", route),
	)

	// Make sure we decrease the number of in-flight requests
	defer c.instruments.active.Add(ctx, -1,
		attribute.String("method", method),
		attribute.String("route", route),
	)

	// Make sure the request has a UUID
	requestUUID, ok := telemetry.UUIDFromContext(ctx)
	if !ok || requestUUID == "" {
		requestUUID = uuid.New().String()
	}

	// Propagate request metadata by adding them to outgoing http request headers
	req.Header.Set(requestUUIDHeader, requestUUID)
	req.Header.Set(clientNameHeader, c.probe.Name())

	// Start a new span
	ctx, span := tracer.Start(ctx,
		"http-client-request",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	// Inject the context and the span context into the http headers
	otel.GetTextMapPropagator().Inject(ctx, &headerTextMapCarrier{
		Header: req.Header,
	})

	// Make the http call
	span.AddEvent("making http call")
	resp, err := c.client.Do(req)

	duration := time.Since(startTime).Milliseconds()

	var statusCode int
	var statusClass string

	if err == nil {
		statusCode = resp.StatusCode
		statusClass = fmt.Sprintf("%dxx", statusCode/100)
	}

	// Report metrics
	meter.RecordBatch(ctx,
		[]attribute.KeyValue{
			attribute.String("method", method),
			attribute.String("route", route),
			attribute.Int("status_code", statusCode),
			attribute.String("status_class", statusClass),
		},
		c.instruments.total.Measurement(1),
		c.instruments.latency.Measurement(duration),
	)

	// Report logs
	logger := c.probe.Logger()
	message := fmt.Sprintf("%s %s %d %dms", method, url, statusCode, duration)
	fields := []interface{}{
		"req.uuid", requestUUID,
		"req.kind", kind,
		"req.method", method,
		"req.url", url,
		"req.route", route,
		"resp.statusCode", statusCode,
		"resp.statusClass", statusClass,
		"resp.duration", duration,
		"traceId", span.SpanContext().TraceID().String(),
		"spanId", span.SpanContext().SpanID().String(),
	}

	// Determine the log level based on the result
	switch {
	case statusCode >= 500:
		logger.Error(message, fields...)
	case statusCode >= 400:
		logger.Warn(message, fields...)
	case statusCode >= 100:
		fallthrough
	default:
		logger.Info(message, fields...)
	}

	// Report the span
	span.SetAttributes(
		attribute.String("method", method),
		attribute.String("url", url),
		attribute.String("route", route),
		attribute.Int("status_code", statusCode),
	)

	return resp, err
}
