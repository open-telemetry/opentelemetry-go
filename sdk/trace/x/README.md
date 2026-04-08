# Experimental Features

The Trace SDK contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go Trace SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [TraceIDRatioBased Sampler](#traceidratiobased-sampler)

### TraceIDRatioBased Sampler

`TraceIDRatioBased` is a threshold-based sampler that conforms to the [OpenTelemetry specification's TraceIdRatioBased sampler](https://opentelemetry.io/docs/specs/otel/trace/sdk/#traceidratiobased).

It uses the least significant 56 bits of the trace ID (per [W3C Trace Context Level 2 Random Trace ID Flag](https://www.w3.org/TR/trace-context-2/#random-trace-id-flag)) for deterministic sampling decisions and propagates the sampling threshold via the `th` sub-key in the W3C `ot` tracestate vendor key.

#### Usage

```go
import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/x"
)

tp := sdktrace.NewTracerProvider(
	sdktrace.WithSampler(
		sdktrace.ParentBased(x.TraceIDRatioBased(0.5)),
	),
)
```

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
