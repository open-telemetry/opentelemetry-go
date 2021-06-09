# OpenTelemetry Collector Go Exporter

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/otlp)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp)

This exporter exports OpenTelemetry spans and metrics to the OpenTelemetry Collector.

## Installation and Setup

The exporter can be installed using standard `go` functionality.

```bash
go get -u go.opentelemetry.io/otel/exporters/otlp
```

A new exporter can be created using the `New` function.
