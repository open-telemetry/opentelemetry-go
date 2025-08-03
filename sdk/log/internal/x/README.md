# Experimental Features

The log SDK contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These feature may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Self-Observability](#self-observability)

### Self-Observability

The log SDK can emit self-observability metrics to help monitor its performance and behavior.
To enable self-observability metrics set the `OTEL_GO_X_SELF_OBSERVABILITY` environment variable to the case-insensitive string of `"true"`.
All other values are ignored.

When enabled, the SDK will emit the following metrics using the global `MeterProvider`:

- `otel.sdk.processor.log.processed` - The number of log records for which the processing has finished, either successful or failed

#### Examples

Enable self-observability metrics.

```console
export OTEL_GO_X_SELF_OBSERVABILITY=true
```

Disable self-observability metrics.

```console
unset OTEL_GO_X_SELF_OBSERVABILITY
```

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
