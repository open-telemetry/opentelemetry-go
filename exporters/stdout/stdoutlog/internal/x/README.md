# Experimental Features

This package documents the experimental features available for [go.opentelemetry.io/otel/exporters/stdout/stdoutlog].

## Self-Observability

The self-observability feature allows the stdout log exporter to emit metrics. When enabled, the exporter will record metrics for:

- Number of log records currently being processed (inflight)
- Total number of log records exported
- Duration of export operations

To enable this feature, set the `OTEL_GO_X_SELF_OBSERVABILITY` environment variable to `true`.

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.