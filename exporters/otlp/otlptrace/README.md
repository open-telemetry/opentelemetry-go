# OpenTelemetry-Go OTLP Span Exporter

OpenTelemetry Protocol (OTLP) Span Exporter.

To constructs a new Otlptrace Exporter and starts it:

```
exp, err := otlptrace.New(ctx, opts...)
```

## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace
```

## [`otlptrace`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace)

The `otlptrace` package provides an exporter implementing the OTel span exporter interface.
This exporter is configured using a client satisfying the `otlptrace.Client` interface.
This client handles the transformation of data into wire format and the transmission of that data to the collector.

## [`otlptracegrpc`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc)

The `otlptracegrpc` package implements a gRPC client to be used in the span exporter.

##  [`otlptracehttp`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp)

The `otlptracehttp` package implements a HTTP client to be used in the span exporter.
