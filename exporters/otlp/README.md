# OpenTelemetry Collector Go Exporter

[![GoDoc](https://godoc.org/go.opentelemetry.io/otel?status.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp)


This exporter converts OpenTelemetry [SpanData](https://github.com/open-telemetry/opentelemetry-go/blob/6769330394f78192df01cb59299e9e0f2e5e977b/sdk/export/trace/trace.go#L49) 
to OpenTelemetry Protocol [Span](https://github.com/open-telemetry/opentelemetry-proto/blob/c20698d5bb483cf05de1a7c0e134b7c57e359674/opentelemetry/proto/trace/v1/trace.proto#L46)
and exports them to OpenTelemetry Collector.


## Installation

```bash
$ go get -u go.opentelemetry.io/otel/exporters/otlp
```
