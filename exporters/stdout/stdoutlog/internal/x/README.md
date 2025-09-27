# Experimental Features

The stdout log exporter contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added prior to stabilization so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Observability](#observability)

### Observability

The exporter provides observability features that allow you to monitor the exporter itself.

To opt-in, set the environment variable `OTEL_GO_X_SELF_OBSERVABILITY` to `true`.

When enabled, the exporter will record metrics for:

- Number of log records currently being processed (inflight)
- Total number of log records exported
- Duration of export operations

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
