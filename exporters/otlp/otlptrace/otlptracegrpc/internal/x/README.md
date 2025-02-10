# Experimental Features

The OTLP trace gRPC exporter contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These feature may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [SDK Self-Observability](#sdk-self-observability)

### SDK Self-Observability

To enable experimental metric and trace instrumentation in SDKs, set the `OTEL_GO_X_SELF_OBSERVABILITY` environment variable.
If enabled, this instrumentation uses the global `TracerProvider` and `MeterProvider`.
The value set must be the case-insensitive string of `"true"` to enable the feature.
All other values are ignored.

#### Examples

Enable experimental sdk self observability.

```console
export OTEL_GO_X_SELF_OBSERVABILITY=true
```

Disable experimental sdk self observability.

```console
unset OTEL_GO_X_SELF_OBSERVABILITY
```

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
