// Package health is used for implementing health checks for services.
package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	defaultTimeout = 5 * time.Second
)

var (
	logger   Logger
	checkers []Checker
	timeout  time.Duration
)

func init() {
	logger = new(voidLogger)
	checkers = make([]Checker, 0)
	timeout = defaultTimeout
}

// Logger is a simple interface for logging.
type Logger interface {
	Errorf(string, ...interface{})
}

type voidLogger struct{}

func (l *voidLogger) Errorf(string, ...interface{}) {}

// Checker is the interface for checking the health of a component.
type Checker interface {
	fmt.Stringer
	HealthCheck(context.Context) error
}

// SetLogger enables logging.
func SetLogger(l Logger) {
	logger = l
}

// RegisterChecker registers new health checkers.
func RegisterChecker(c ...Checker) {
	checkers = append(checkers, c...)
}

// SetTimeout sets the timeout for health checkers.
// If the timeout reaches, the health check is considered failed.
// The default timeout is 5 seconds.
func SetTimeout(d time.Duration) {
	timeout = d
}

// HandlerFunc returns an http handler function for checking the health of registered checkers.
func HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		group, ctx := errgroup.WithContext(ctx)

		for _, checker := range checkers {
			checker := checker // https://golang.org/doc/faq#closures_and_goroutines
			group.Go(func() error {
				err := checker.HealthCheck(ctx)
				if err != nil {
					logger.Errorf("Error on health checking: %s: %s", checker, err)
				}
				return err
			})
		}

		if err := group.Wait(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
