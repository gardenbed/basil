package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gardenbed/basil/telemetry"
	httptelemetry "github.com/gardenbed/basil/telemetry/http"
)

const port = ":9001"

func main() {
	// Create a new probe
	probe := telemetry.NewProbe(
		telemetry.WithLogger("info"),
		telemetry.WithPrometheus(),
		telemetry.WithOpenTelemetry(false, true, "", nil),
		telemetry.WithMetadata("client", "0.1.0", map[string]string{
			"environment": "testing",
		}),
	)

	defer func() {
		if err := probe.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{},
	}

	client := httptelemetry.NewClient(c, probe, httptelemetry.Options{})

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:9000/users/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	probe.Logger().Info("received response",
		"content", string(bytes),
	)

	http.Handle("/metrics", probe)
	probe.Logger().Infof("starting http server on %s ...", port)
	panic(http.ListenAndServe(port, nil))
}
