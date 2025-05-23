# Semantic Convention Changes

The `go.opentelemetry.io/otel/semconv/v1.33.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.32.0` with the following exceptions.

## Dropped deprecations

The following declarations have been deprecated in the [OpenTelemetry Semantic Conventions].
Refer to the respective documentation in that repository for deprecation instructions for each type.

- `FeatureFlagEvaluationErrorMessage`
- `FeatureFlagEvaluationErrorMessageKey`

[OpenTelemetry Semantic Conventions]: https://github.com/open-telemetry/semantic-conventions
