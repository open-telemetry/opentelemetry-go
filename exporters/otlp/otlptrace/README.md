# OpenTelemetry-Go OTLP Span Exporter

OpenTelemetry Protocol (OTLP) Span Exporter.

To constructs a new Otlptrace Exporter ,you can follow this 
 [`example`](https://github.com/open-telemetry/opentelemetry-go/blob/main/exporters/otlp/otlptrace/example_test.go).

## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace
```

## [`otlptrace`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace)

`otlptrace` package implements a span exporter that uses a `otlptrace.Client` interface.

`otlptrace.Client` manages connections to the collector, handles the transformation of data into wire format,
and the transmission of that data to the collector.

##example

The `otlptracegrpc` package implements a gRPC span exporter.

##  [`otlptracehttp`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp)

The `otlptracehttp` package implements a HTTP span exporter.
