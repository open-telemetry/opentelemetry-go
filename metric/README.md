# Metric API

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/metric)](https://pkg.go.dev/go.opentelemetry.io/otel/metric)

The Metric API defines the telemetry recording interfaces for instrumentation,
while applications configure the [metric SDK] to collect and export
measurements.
See the [synchronous meter example] and [asynchronous meter example] for
instrumentation patterns.

[metric SDK]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric
[synchronous meter example]: https://pkg.go.dev/go.opentelemetry.io/otel/metric#example-Meter-synchronous
[asynchronous meter example]: https://pkg.go.dev/go.opentelemetry.io/otel/metric#example-Meter-asynchronous-single
