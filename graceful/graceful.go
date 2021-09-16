// Package graceful provides graceful start, graceful retry, and graceful stop!
// It can be used for:
//   - Gracefully starting servers (http, grpc, etc.) and clients (external services, databases, message queues, etc.).
//   - Gracefully retrying lost connections to external services, databases, message queues, etc.
//   - Gracefully stopping servers (http, grpc, etc.) and clients (external services, databases, message queues, etc.).
package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	successCode        = 0
	errorCode          = 1
	defaultMaxRetry    = 5
	defaultGracePeriod = 30 * time.Second
)

var (
	logger      Logger
	clients     []Client
	servers     []Server
	maxRetry    int
	gracePeriod time.Duration
)

func init() {
	logger = new(voidLogger)
	clients = make([]Client, 0)
	servers = make([]Server, 0)
	maxRetry = defaultMaxRetry
	gracePeriod = defaultGracePeriod
}

// Logger is a simple interface for logging.
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

type voidLogger struct{}

func (l *voidLogger) Debugf(string, ...interface{}) {}
func (l *voidLogger) Infof(string, ...interface{})  {}
func (l *voidLogger) Warnf(string, ...interface{})  {}
func (l *voidLogger) Errorf(string, ...interface{}) {}

// Client is the generic interface for a client (external service, database, message queue, etc.).
type Client interface {
	fmt.Stringer
	Connect() error
	Disconnect(context.Context) error
}

// Server is the generic interface for a server (http, grpc, etc.).
type Server interface {
	fmt.Stringer
	ListenAndServe() error
	Shutdown(context.Context) error
}

// SetLogger enables logging.
func SetLogger(l Logger) {
	logger = l
}

// RegisterClient registers new clients (external services, databases, message queues, etc.).
func RegisterClient(c ...Client) {
	clients = append(clients, c...)
}

// RegisterServer registers new servers (http, grpc, etc.).
func RegisterServer(s ...Server) {
	servers = append(servers, s...)
}

// SetMaxRetry sets the maximum number of retries.
// When a client connection is lost, the library retries to establish the connection.
// The default maximum retry is 5.
func SetMaxRetry(n int) {
	maxRetry = n
}

// SetGracePeriod sets the timeout for gracefully stopping servers and clients.
// If the timeout reaches, the process will be terminated ungracefully.
// The default grace period is 30 seconds.
func SetGracePeriod(d time.Duration) {
	gracePeriod = d
}

func pow2(exponent int) int {
	r := 1
	for n := 0; n < exponent; n++ {
		r *= 2
	}
	return r
}

// connectWithRetry tries to run a connect function successfully with exponential backoff.
//   Example: with maxRetry = 5 --> 0s, 1s, 2s, 4s, 8s, 16s
func connectWithRetry(c Client, retries int) error {
	err := c.Connect()
	if err == nil {
		logger.Infof("Client connected successfully: %s", c)
		return nil
	}

	if retries == 0 {
		logger.Errorf("Failed to connect client: %s: %s", c, err)
		return err
	}

	logger.Debugf("Retrying connecting client: %s", c)
	factor := pow2(maxRetry - retries)
	backoff := time.Duration(factor) * time.Second
	time.Sleep(backoff)
	return connectWithRetry(c, retries-1)
}

func terminateGracefully(termServers bool) error {
	group := new(errgroup.Group)
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	// Gracefully shutting down servers
	if termServers {
		for _, s := range servers {
			s := s // https://golang.org/doc/faq#closures_and_goroutines
			group.Go(func() error {
				if err := s.Shutdown(ctx); err != nil {
					logger.Errorf("Failed to shutdown server: %s: %s", s, err)
					return err
				}
				logger.Infof("Server gracefully shutdown: %s", s)
				return nil
			})
		}
	}

	// Gracefully disconnecting clients
	for _, c := range clients {
		c := c // https://golang.org/doc/faq#closures_and_goroutines
		group.Go(func() error {
			if err := c.Disconnect(ctx); err != nil {
				logger.Errorf("Failed to disconnect client: %s: %s", c, err)
				return err
			}
			logger.Infof("Client gracefully disconnected: %s", c)
			return nil
		})
	}

	// Wait for all servers and clients to gracefully terminate
	// Wait blocks until all functions have returned, then returns the first non-nil error (if any)
	if err := group.Wait(); err != nil {
		return err
	}

	logger.Infof("Gracefully Terminated!")
	return nil
}

// StartAndWait behaves as follows:
//
//   1. It tries to connect all clients each in a new goroutine.
//     - If a client fails to connect, it will be automatically retried for a limited number of times with exponential backoff.
//     - If at least one client fails to connect (after retries), a graceful termination will be initiated.
//   2. Once all clients are connected successfully, all servers start listening each in a new goroutine.
//     - If any server errors, a graceful termination will be initiated.
//   3. Then, this method blocks the current goroutine until one of the following conditions happen:
//     - If any of SIGHUP, SIGINT, SIGQUIT, SIGTERM signals is sent, a graceful termination will be initiated.
//     - If any of the above signals is sent for the second time before the graceful termination is completed, the process will exit immediately with an error code.
//
func StartAndWait() int {
	// First, starting clients asynchronously
	clientGroup := new(errgroup.Group)
	for _, c := range clients {
		c := c // https://golang.org/doc/faq#closures_and_goroutines
		clientGroup.Go(func() error {
			return connectWithRetry(c, maxRetry)
		})
	}

	// Wait for all clients to successfully connect
	// If at least one client fails to connect, we gracefully terminate all clients
	if err := clientGroup.Wait(); err != nil {
		logger.Errorf("Initiating graceful termination due to client error: %s", err)
		_ = terminateGracefully(false)
		return errorCode
	}

	// Starting servers asynchronously
	errCh := make(chan error)
	for _, s := range servers {
		go func(s Server) {
			logger.Infof("Server starting to listen: %s", s)
			if err := s.ListenAndServe(); err != nil {
				logger.Errorf("Error occurred for server: %s: %s", s, err)
				errCh <- err
			}
		}(s)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,  // 1: Hangup signal (controlling terminal is closed)
		syscall.SIGINT,  // 2: Terminal interrupt signal (Ctrl+C)
		syscall.SIGQUIT, // 3: Terminal quit signal (Ctrl+\)
		syscall.SIGTERM, // 5: Termination signal (sent by Kubernetes for graceful termination)
	)

	// Waiting for a server error or a termination signal
	select {
	case err := <-errCh:
		logger.Errorf("Initiating graceful termination due to server error: %s", err)
		_ = terminateGracefully(true)
		return errorCode

	case sig := <-sigCh:
		logger.Infof("Initiating graceful termination: signal received: %s", sig)

		go func() {
			sig := <-sigCh
			logger.Errorf("Terminating immediately: second signal received: %s", sig)
			os.Exit(1)
		}()

		if err := terminateGracefully(true); err != nil {
			return errorCode
		}
		return successCode
	}
}
