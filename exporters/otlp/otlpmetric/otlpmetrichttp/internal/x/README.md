# Experimental Features

The `otlpmetrichttp` exporter contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the `otlpmetrichttp` exporter prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Observability](#observability)
- [Attribute Cache](#attribute-cache)

### Observability

The `otlpmetrichttp` exporter can be configured to provide observability about itself using OpenTelemetry metrics.

To opt-in, set the environment variable `OTEL_GO_X_OBSERVABILITY` to `true`.

When enabled, the exporter will create the following metrics using the global `MeterProvider`:

- `otel.sdk.exporter.metric_data_point.inflight`
- `otel.sdk.exporter.metric_data_point.exported`
- `otel.sdk.exporter.operation.duration`

Please see the [Semantic conventions for OpenTelemetry SDK metrics] documentation for more details on these metrics.

[Semantic conventions for OpenTelemetry SDK metrics]: https://github.com/open-telemetry/semantic-conventions/blob/v1.43.0/docs/otel/sdk-metrics.md

### Attribute Cache

The `otlpmetrichttp` exporter can cache the OTLP form of data point attribute sets across exports, so a set that repeats every export is converted only once.

To opt-in, set the environment variable `OTEL_GO_X_OTLP_ATTRIBUTE_CACHE_SIZE` to the maximum number of attribute sets to cache, as a positive integer.
The cache is disabled when the variable is unset or `0`, and all other values are ignored.

When the cache is full, new attribute sets are converted without being cached.
At most every 10 minutes, a full cache is cleared so it can adapt to attribute sets that change over time.
For the best results, set the value above the number of distinct attribute sets the process exports.
A smaller cache still speeds up the attribute sets it holds.
Memory use grows with the number of cached attribute sets and their size.

Cached attributes are shared across exports, so code in the exporter must not modify them.

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
