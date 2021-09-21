[![Go Doc][godoc-image]][godoc-url]

# Telemetry

This package can be used for building observable applications in Go.
It aims to unify the three pillars of observability in one single package that is *easy-to-use* and *hard-to-misuse*.

This package leverages the [OpenTelemetry](https://opentelemetry.io) API.
OpenTelemetry is a great initiative that has brought all different standards and APIs for observability under one umbrella.
However, due to the requirements for interoperability with existing systems, OpenTelemetry is highly abstract and complex.
Many packages, configurations, and options make the developer experience not so pleasant.
Furthermore, due to the changing nature of this project, OpenTelemetry specification changes often so does the Go library for OpenTelemetry.
Currently, the tracing API is *stable*, the metric API is in *alpha* stage, and the logging API is *frozen*.

IMHO, this is not how a single unified observability API should be.
Hopefully, many of these issues will go away once all APIs reache to v1.0.0.
This package intends to provide a very minimal and yet practical API for observability
by hiding the complexity of configuring and using OpenTelemetry API.

A probe encompasses a logger, meter, and tracer.
It offers a unified developer experience for building observable applications.

## The Three Pillars of Observability

### Logging

Logs are used for auditing purposes (sometimes for debugging with limited capabilities).
When looking at logs, you need to know what to look for ahead of time (known unknowns vs. unknown unknowns).
Since log data can have any arbitrary shape and size, they cannot be used for real-time computational purposes.
Logs are hard to track across different and distributed processes. Logs are also very expensive at scale.

### Metrics

Metrics are regular time-series data with low and fixed cardinality.
They are aggregated by time. Metrics are used for **real-time** monitoring purposes.
Using metrics we can implement **SLIs** (service-level indicators), **SLOs** (service-level objectives), and automated alerts.
Metrics are very good at taking the distribution of data into account.
Metrics cannot be used with high-cardinality data.

### Tracing

Traces are used for debugging and tracking requests across different processes and services.
They can be used for identifying performance bottlenecks.
Due to their very data-heavy nature, traces in real-world applications need to be sampled.
Insights extracted from traces cannot be aggregated since they are sampled.
In other words, information captured by one trace does not tell anything about how
this trace is compared against other traces, and what is the distribution of data.

## Quick Start

You can find basic examples [here](./example).

## Options

Most options can be set through environment variables.
This lets SRE people change how the observability pipeline is configured without making any code changes.
Options set explicity in the code will override those set by environment variables.

| Environment Variable | Description |
|----------------------|-------------|
| `PROBE_NAME` | The name of service or application. |
| `PROBE_VERSION` | The version of service or application. |
| `PROBE_TAG_*` | Each variable prefixed with `PROBE_TAG_` represents a tag for the service or application. |
| `PROBE_LOGGER_ENABLED` | Whether or not to create a logger (boolean). |
| `PROBE_LOGGER_LEVEL` | The verbosity level for the logger (`debug`, `info`, `warn`, `error`, or `none`). |
| `PROBE_PROMETHEUS_ENABLED` | Whether or not to configure and create a Prometheus meter (boolean). |
| `PROBE_JAEGER_ENABLED` | Whether or not to configure and create a Jaeger tracer (boolean). |
| `PROBE_JAEGER_AGENT_HOST` | The Jaeger agent host (i.e. `localhost`). |
| `PROBE_JAEGER_AGENT_PORT` | The Jaeger agent port (i.e. `6832`). |
| `PROBE_JAEGER_COLLECTOR_ENDPOINT` | The full URL to the Jaeger HTTP Thrift collector (i.e. `http://localhost:14268/api/traces`). |
| `PROBE_JAEGER_COLLECTOR_USERNAME` | The username for Jaeger collector endpoint if basic auth is required. |
| `PROBE_JAEGER_COLLECTOR_PASSWORD` | The password for Jaeger collector endpoint if basic auth is required. |
| `PROBE_OPENTELEMETRY_ENABLED` | Whether or not to configure and create an OpenTelemetry Collector meter and tracer (boolean). |
| `PROBE_OPENTELEMETRY_COLLECTOR_ADDRESS` | The address to OpenTelemetry collector (i.e. `localhost:55680`). |

## Documentation

  - **Logging**
    - [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap)
  - **Metrics**
    - [Metrics API](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/api.md)
    - [go.opentelemetry.io/otel/metric](https://pkg.go.dev/go.opentelemetry.io/otel/metric)
  - **Tracing**
    - [Tracing API](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md)
    - [go.opentelemetry.io/otel/trace](https://pkg.go.dev/go.opentelemetry.io/otel/trace)
  - **OpenTelemetry**
    - [Collector Configuration](https://opentelemetry.io/docs/collector/configuration)
    - [Collector Architecture](hhttps://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/design.md)


[godoc-url]: https://pkg.go.dev/github.com/gardenbed/basil/telemetry
[godoc-image]: https://pkg.go.dev/badge/github.com/gardenbed/basil/telemetry
