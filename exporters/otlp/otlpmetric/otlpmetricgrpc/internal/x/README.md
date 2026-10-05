# Experimental Features

The OTLP gRPC metric exporter contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go OTLP exporters prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Self-Observability](#self-observability)
- [Attribute Cache](#attribute-cache)

### Self-Observability

The OTLP gRPC metric exporter can emit self-observability metrics to track its own operation.

This experimental feature can be enabled by setting the `OTEL_GO_X_OBSERVABILITY` environment variable.
The value must be the case-insensitive string of `"true"` to enable the feature.
All other values are ignored.

When enabled, the exporter will emit the following metrics using the global MeterProvider:

- `otel.sdk.exporter.metric_data_point.exported`: Counter tracking successfully exported data points
- `otel.sdk.exporter.metric_data_point.inflight`: UpDownCounter tracking data points currently being exported  
- `otel.sdk.exporter.operation.duration`: Histogram tracking export operation duration in seconds

All metrics include attributes identifying the exporter component and destination server:

- `otel.component.type`: Type of component (e.g., "otlp_grpc_metric_exporter")
- `otel.component.name`: Unique component instance name (e.g., "otlp_grpc_metric_exporter/0")
- `server.address`: Server hostname or address
- `server.port`: Server port number

#### Examples

Enable self-observability metrics.

```console
export OTEL_GO_X_OBSERVABILITY=true
```

Disable self-observability metrics.

```console
unset OTEL_GO_X_OBSERVABILITY
```

### Attribute Cache

The OTLP gRPC metric exporter can cache the OTLP form of data point attribute sets across exports, so a set that repeats every export is converted only once.

This experimental feature is disabled by default.
To enable it, set the `OTEL_GO_X_OTLP_ATTRIBUTE_CACHE_SIZE` environment variable to the maximum number of attribute sets to cache.
The value must be a positive integer.
A value of `0` disables the cache, and all other values are ignored.

When the cache is full, new attribute sets are converted without being cached.
At most every 10 minutes, a full cache is cleared so it can adapt to attribute sets that change over time.
For the best results, set the value above the number of distinct attribute sets the process exports.
A smaller cache still speeds up the attribute sets it holds.
Memory use grows with the number of cached attribute sets and their size.

#### Examples

Cache up to 10000 attribute sets.

```console
export OTEL_GO_X_OTLP_ATTRIBUTE_CACHE_SIZE=10000
```

Disable the cache.

```console
unset OTEL_GO_X_OTLP_ATTRIBUTE_CACHE_SIZE
```

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
