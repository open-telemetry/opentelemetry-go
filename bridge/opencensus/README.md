# OpenTelemetry/OpenCensus Bridge

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/bridge/opencensus)](https://pkg.go.dev/go.opentelemetry.io/otel/bridge/opencensus)

## Deprecation

The `go.opentelemetry.io/otel/bridge/opencensus` module and package are deprecated as of v1.48.0.
v1.48.0 is the final planned release of this module.
OpenCensus has been [sunset](https://opentelemetry.io/blog/2023/sunsetting-opencensus/), and its [compatibility requirements are deprecated in the OpenTelemetry specification](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/compatibility/opencensus.md).
Use native OpenTelemetry APIs and SDKs instead; see [Migrating from OpenCensus](https://opentelemetry.io/docs/specs/otel/compatibility/opencensus/#migration-path).
