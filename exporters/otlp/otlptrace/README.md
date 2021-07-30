# OpenTelemetry-Go OTLP Span Exporter

[![Go Reference](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/otlp/otlptrace.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace)

OpenTelemetry Protocol (OTLP) Span Exporter.

## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace
```

## Example

To constructs a new OTLP trace Exporter, you can follow this
[`OTLP Trace Exporter Example`](https://github.com/open-telemetry/opentelemetry-go/blob/main/exporters/otlp/otlptrace/example_test.go).

Also,if you are not familiar with how to export trace and metric data from the OpenTelemetry-Go SDK to
the OpenTelemetry Collector, you can reference this
[`OpenTelemetry Collector Traces Example`](https://github.com/open-telemetry/opentelemetry-go/tree/main/example/otel-collector).


## [`otlptrace`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace)

The `otlptrace` package provides an exporter implementing the OTel span exporter interface.
This exporter is configured using a client satisfying the `otlptrace.Client` interface.
This client handles the transformation of data into wire format and the transmission of that data to the collector.

## [`otlptracegrpc`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc)

The `otlptracegrpc` package implements a gRPC client to be used in the span exporter.

##  [`otlptracehttp`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp)

The `otlptracehttp` package implements a HTTP client to be used in the span exporter.
