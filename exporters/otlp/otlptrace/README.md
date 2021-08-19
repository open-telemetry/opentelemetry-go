# OpenTelemetry-Go OTLP Span Exporter

[![Go Reference](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/otlp/otlptrace.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace)

[OpenTelemetry Protocol Exporter](https://github.com/open-telemetry/opentelemetry-specification/blob/v1.5.0/specification/protocol/exporter.md) implementation.

## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace
```

## Examples

- [example_test.go](example_test.go)
- [../../../example/otel-collector](../../../example/otel-collector)


## [`otlptrace`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace)

The `otlptrace` package provides an exporter implementing the OTel span exporter interface.
This exporter is configured using a client satisfying the `otlptrace.Client` interface.
This client handles the transformation of data into wire format and the transmission of that data to the collector.

## [`otlptracegrpc`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc)

The `otlptracegrpc` package implements a gRPC client to be used in the span exporter.

##  [`otlptracehttp`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp)

The `otlptracehttp` package implements a HTTP client to be used in the span exporter.
