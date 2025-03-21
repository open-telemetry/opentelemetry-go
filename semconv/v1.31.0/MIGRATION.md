# Semantic Convention Changes

The `go.opentelemetry.io/otel/semconv/v1.31.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.30.0` with the following exceptions.

## Dropped deprecations

The following declarations have been deprecated in the [OpenTelemetry Semantic Conventions].
Refer to the respective documentation in that repository for deprecation instructions for each type.

- `CodeFilepathKey`
- `CodeNamespaceKey`
- `CodeNamespace`
- `GenAIOpenaiRequestResponseFormatJSONObject`
- `GenAIOpenaiRequestResponseFormatJSONSchema`
- `GenAIOpenaiRequestResponseFormatKey`
- `GenAIOpenaiRequestResponseFormatText`

[OpenTelemetry Semantic Conventions]: https://github.com/open-telemetry/semantic-conventions
