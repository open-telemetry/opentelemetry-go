# Metric SDK

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric)

See the [package setup example] for configuring a meter provider and the
[view example] for customizing metric streams.
Applications choose a [manual reader] or [periodic reader] for collecting
and exporting measurements, rather than importing metricdata directly.

[package setup example]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#example-package
[view example]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#example-View
[manual reader]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#NewManualReader
[periodic reader]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#NewPeriodicReader
