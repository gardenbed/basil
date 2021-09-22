// Package grpc is used for building observable gRPC servers and clients that automatically report logs, metrics, and traces.
package grpc

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	libraryName    = "github.com/basil/telemetry/grpc"
	requestUUIDKey = "request-uuid"
	clientNameKey  = "client-name"
)

var (
	splitRegex = regexp.MustCompile(`/|\.`)
)

// Options are optional configurations for interceptors.
type Options struct {
	ExcludedMethods []string
}

func (opts Options) withDefaults() Options {
	return opts
}

// endpoint is a grpc endpoint.
type endpoint struct {
	Package string
	Service string
	Method  string
}

// fullMethod is in the form of /package.service/method
func parseEndpoint(fullMethod string) (endpoint, bool) {
	subs := splitRegex.Split(fullMethod, 4)
	if len(subs) != 4 {
		return endpoint{}, false
	}

	return endpoint{
		Package: subs[1],
		Service: subs[2],
		Method:  subs[3],
	}, true
}

// String implements the fmt.Stringer interface.
func (e endpoint) String() string {
	var s string
	if e.Package != "" && e.Service != "" && e.Method != "" {
		s = fmt.Sprintf("%s::%s::%s", e.Package, e.Service, e.Method)
	}
	return s
}

// serverStream extends the grpc.ServerStream.
type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	if s.ctx == nil {
		return s.ServerStream.Context()
	}
	return s.ctx
}

// ServerStreamWithContext returns a new grpc.ServerStream with a new context.
func ServerStreamWithContext(ctx context.Context, s grpc.ServerStream) grpc.ServerStream {
	if ss, ok := s.(*serverStream); ok {
		ss.ctx = ctx
		return ss
	}

	return &serverStream{
		ServerStream: s,
		ctx:          ctx,
	}
}

// metadataTextMapCarrier implements propagation.TextMapCarrier interface.
type metadataTextMapCarrier struct {
	md *metadata.MD
}

func (c *metadataTextMapCarrier) Get(key string) string {
	if vals := c.md.Get(key); len(vals) > 0 {
		return vals[0]
	}

	return ""
}

func (c *metadataTextMapCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataTextMapCarrier) Keys() []string {
	keys := make([]string, 0, len(*c.md))
	for k := range *c.md {
		keys = append(keys, k)
	}

	return keys
}
