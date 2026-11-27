# SDK Trace

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/trace)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace)

To get started with the Trace SDK, construct a provider with
[`NewTracerProvider`], register a span exporter using [`WithBatcher`], and call
[`TracerProvider.Shutdown`] before your application exits to flush pending
spans and release resources.

[`NewTracerProvider`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#NewTracerProvider
[`WithBatcher`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#WithBatcher
[`TracerProvider.Shutdown`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProvider.Shutdown
