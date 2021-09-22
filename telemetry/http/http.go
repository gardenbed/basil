// Package http is used for building observable HTTP servers and clients that automatically report logs, metrics, and traces.
package http

import (
	"fmt"
	"net/http"
	"regexp"
)

const (
	libraryName       = "github.com/basil/telemetry/http"
	requestUUIDHeader = "Request-UUID"
	clientNameHeader  = "Client-Name"
)

var (
	defaultIDRegexp = regexp.MustCompile("[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}")
)

// Options are optional configurations for middleware and clients.
type Options struct {
	IDRegexp *regexp.Regexp
}

func (opts Options) withDefaults() Options {
	if opts.IDRegexp == nil {
		opts.IDRegexp = defaultIDRegexp
	}

	return opts
}

// responseWriter extends the standard http.ResponseWriter.
type responseWriter struct {
	http.ResponseWriter
	StatusCode  int
	StatusClass string
}

func newResponseWriter(rw http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: rw,
	}
}

// WriteHeader overrides the http.WriteHeader.
func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)

	// Only capture the first value
	if r.StatusCode == 0 {
		r.StatusCode = statusCode
		r.StatusClass = fmt.Sprintf("%dxx", statusCode/100)
	}
}

// headerTextMapCarrier implements propagation.TextMapCarrier interface.
type headerTextMapCarrier struct {
	http.Header
}

func (c *headerTextMapCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Header))
	for k := range c.Header {
		keys = append(keys, k)
	}

	return keys
}
