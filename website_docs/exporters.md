---
title: Exporters
aliases: [/docs/instrumentation/go/exporting_data]
weight: 4
---

In order to visualize and analyze your [traces](/docs/concepts/signals/traces/)
and metrics, you will need to export them to a backend.

## OTLP Exporter

OpenTelemetry Protocol (OTLP) export is available in the
`go.opentelemetry.io/otel/exporters/otlp/otlptrace` and
`go.opentelemetry.io/otel/exporters/otlp/otlpmetric` packages.

Please find more documentation on
[GitHub](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/otlp)

## Jaeger Exporter

Jaeger export is available in the `go.opentelemetry.io/otel/exporters/jaeger`
package.

Please find more documentation on
[GitHub](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/jaeger)

## Prometheus Exporter

Prometheus export is available in the
`go.opentelemetry.io/otel/exporters/prometheus` package.

Please find more documentation on
[GitHub](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/prometheus)
