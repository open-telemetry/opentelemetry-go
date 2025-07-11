# Semantic Convention Changes

The `go.opentelemetry.io/otel/semconv/v1.33.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.32.0` with the following exceptions.

## Metric instrument type fixes

### `goconv.MemoryUse`

The underlying metric instrument type for `goconv.MemoryUse` has been corrected be an `Int64ObservableUpDownCounter` instead of an `Int64ObservableCounter`.
This change aligns with the semantic conventions for memory usage metrics.

### `systemconv.MemoryUsage`

The underlying metric instrument type for `systemconv.MemoryUsage` has been corrected be an `Int64ObservableUpDownCounter` instead of an `Int64ObservableGauge`.
This change aligns with the semantic conventions for memory usage metrics.

## Dropped deprecations

The following declarations have been deprecated in the [OpenTelemetry Semantic Conventions].
Refer to the respective documentation in that repository for deprecation instructions for each type.

- `FeatureFlagEvaluationErrorMessage`
- `FeatureFlagEvaluationErrorMessageKey`

[OpenTelemetry Semantic Conventions]: https://github.com/open-telemetry/semantic-conventions
