# Experimental Features

The Metric SDK contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go Metric SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These feature may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Cardinality Limit](#cardinality-limit)

### Cardinality Limit

The cardinality limit is the hard limit on the number of metric datapoints that can be
collected for a single instrument in a single collect cycle. By default, a limit of 2000
is applied.

To override the default limit, set the `OTEL_GO_X_CARDINALITY_LIMIT` environment variable
to the integer limit to use. Setting this to a zero or negative value means no limit is
applied. Invalid values are ignored, and the default is used instead.

This is provided for backward compatibility; new code should use
[WithCardinalityLimit](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#WithCardinalityLimit) instead.

#### Examples

Set the cardinality limit to 5000.

```console
export OTEL_GO_X_CARDINALITY_LIMIT=5000
```

Disable the cardinality limit.

```console
export OTEL_GO_X_CARDINALITY_LIMIT=-1
```

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
