# Experimental Features

This package documents experimental features for the `go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc` package.

## Observability

Observability metrics for the OTLP gRPC metric exporter can be enabled by setting the `OTEL_GO_X_OBSERVABILITY` environment variable to `true`.

When enabled, the exporter will emit the following metrics:

- `otel.sdk.exporter.metric_data_point.inflight`: The number of metric data points currently being exported
- `otel.sdk.exporter.metric_data_point.exported`: The number of metric data points successfully exported
- `otel.sdk.exporter.operation.duration`: The duration of export operations

These metrics follow the [OpenTelemetry Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/general/metrics/).

**Note**: This is an experimental feature and may change or be removed in future versions.
