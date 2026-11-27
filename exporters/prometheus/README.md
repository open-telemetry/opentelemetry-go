# Prometheus Exporter

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/prometheus)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus)

Create a Prometheus exporter using [prometheus.New](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus#New)
and register it with the MeterProvider using [metric.WithReader](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#WithReader).
Metrics are then exposed through a Prometheus HTTP handler.

See the [Prometheus exporter guide](https://opentelemetry.io/docs/languages/go/exporters/#prometheus-experimental)
for an integration example.
