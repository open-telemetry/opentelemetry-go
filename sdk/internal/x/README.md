# Experimental Features

The SDK contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These feature may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Resource](#resource)
- [Observability](#observability)
- [Per-Series Start Timestamps](#per-series-start-timestamps)

### Resource

[OpenTelemetry resource semantic conventions] include many attribute definitions that are defined as experimental.
To have experimental semantic conventions be added by [resource detectors] set the `OTEL_GO_X_RESOURCE` environment variable.
The value set must be the case-insensitive string of `"true"` to enable the feature.
All other values are ignored.

<!-- TODO: document what attributes are added by which detector -->

[OpenTelemetry resource semantic conventions]: https://opentelemetry.io/docs/specs/semconv/resource/
[resource detectors]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/resource#Detector

#### Examples

Enable experimental resource semantic conventions.

```console
export OTEL_GO_X_RESOURCE=true
```

Disable experimental resource semantic conventions.

```console
unset OTEL_GO_X_RESOURCE
```

### Observability

The SDK provides an observability feature that allows you to monitor the SDK itself.

To opt-in, set the `OTEL_GO_X_OBSERVABILITY` environment variable to `true`.
The `OTEL_GO_X_SELF_OBSERVABILITY` environment variable is also recognized as an alias.

When enabled, the SDK components (e.g. the trace and metric SDKs) create metrics
about their own operation using the global `MeterProvider`.
Please see the [Semantic conventions for OpenTelemetry SDK metrics] documentation
for more details on these metrics.

[Semantic conventions for OpenTelemetry SDK metrics]: https://github.com/open-telemetry/semantic-conventions/blob/v1.37.0/docs/otel/sdk-metrics.md

#### Examples

Enable SDK self-observability metrics.

```console
export OTEL_GO_X_OBSERVABILITY=true
```

Disable SDK self-observability metrics.

```console
unset OTEL_GO_X_OBSERVABILITY
```

### Per-Series Start Timestamps

[OpenTelemetry metrics data model] defines a per-series `StartTimeUnixNano`
that reflects when each individual attribute set was first observed, rather
than a single start time shared by every data point of an instrument.

To have the metric SDK report per-series start timestamps, set the
`OTEL_GO_X_PER_SERIES_START_TIMESTAMPS` environment variable.
The value set must be the case-insensitive string value of `"true"` to enable
the feature. All other values are ignored.

[OpenTelemetry metrics data model]: https://opentelemetry.io/docs/specs/otel/metrics/data-model/

#### Examples

Enable per-series start timestamps.

```console
export OTEL_GO_X_PER_SERIES_START_TIMESTAMPS=true
```

Disable per-series start timestamps.

```console
unset OTEL_GO_X_PER_SERIES_START_TIMESTAMPS
```

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
