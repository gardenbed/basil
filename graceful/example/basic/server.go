package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// Server is a simple http server implementing graceful.Client interface.
type Server struct {
	name string
	*http.Server
}

// NewServer creates a new server.
func NewServer(name string, port uint16) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello, World!")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	// This is required for gracefully shutting down connections that have undergone ALPN protocol upgrade or that have been hijacked.
	//   - The BaseContext will be provided to all incoming requests.
	//   - When a server shutdown is initiated, the base context will be cancelled.
	//   - The handler for upgraded or hijacked connection (websocket, etc.) should handle context cancellation.
	server.RegisterOnShutdown(cancel)

	return &Server{
		name:   name,
		Server: server,
	}
}

func (s *Server) String() string {
	return s.name
}
