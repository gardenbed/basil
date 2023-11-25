package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gardenbed/basil/telemetry"
	httptelemetry "github.com/gardenbed/basil/telemetry/http"
)

const port = ":9000"

func main() {
	// Create a new probe and set it as the singleton
	probe := telemetry.NewProbe(
		telemetry.WithLogger("info"),
		telemetry.WithPrometheus(),
		telemetry.WithOpenTelemetry(false, true, "", nil),
		telemetry.WithMetadata("server", "0.1.0", map[string]string{
			"environment": "testing",
		}),
	)
	defer probe.Close(context.Background())

	mid := httptelemetry.NewMiddleware(probe, httptelemetry.Options{})

	handler := mid.Wrap(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(50 * time.Millisecond)

		ctx := req.Context()
		logger := telemetry.LoggerFromContext(ctx)
		tracer := telemetry.TracerFromContext(ctx)

		_, span := tracer.Start(ctx, "database-read")
		defer span.End()

		time.Sleep(100 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "Hello, world!")

		logger.Debug("responded back!")
	})

	http.Handle("/users/", handler)
	http.Handle("/metrics", probe)
	probe.Logger().Info("starting http server on %s ...", port)
	panic(http.ListenAndServe(port, nil))
}
