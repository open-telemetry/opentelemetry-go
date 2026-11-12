# Experimental Features

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/trace/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace/x)

The Trace SDK contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go Trace SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [ProbabilitySampler](#probabilitysampler)
- [Context fallback](#context-fallback)

### ProbabilitySampler

`ProbabilitySampler` is a threshold-based sampler that conforms to the [OpenTelemetry specification's ProbabilitySampler](https://opentelemetry.io/docs/specs/otel/trace/sdk/#probabilitysampler).

It uses the least significant 56 bits of the trace ID (per [W3C Trace Context Level 2 Random Trace ID Flag](https://www.w3.org/TR/trace-context-2/#random-trace-id-flag)) for deterministic sampling decisions and propagates the sampling threshold via the `th` sub-key in the W3C `ot` tracestate vendor key.

#### Usage

```go
import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/x"
)

tp := sdktrace.NewTracerProvider(
	sdktrace.WithSampler(
		sdktrace.ParentBased(x.ProbabilitySampler(0.5)),
	),
)
```

### Context fallback

`WithContextFallback` supplies incoming trace context and baggage when a span starts without a parent.
It is an experimental, opt-in policy for a single `TracerProvider`.
It can connect a command's work to the process that launched it, including instrumentation that starts with `context.Background()`.
Providers created without this option keep their existing behavior.

#### Process environment usage

Supply the application's configured propagator and an environment-variable carrier from [`go.opentelemetry.io/contrib/propagators/envcar`](https://pkg.go.dev/go.opentelemetry.io/contrib/propagators/envcar):

```go
import (
	"context"

	"go.opentelemetry.io/contrib/propagators/envcar"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/x"
)

propagator := propagation.NewCompositeTextMapPropagator(
	propagation.TraceContext{},
	propagation.Baggage{},
)
tp := sdktrace.NewTracerProvider(
	x.WithContextFallback(propagator, &envcar.Carrier{}),
)
defer tp.Shutdown(context.Background())

ctx, span := tp.Tracer("example.com/worker").Start(context.Background(), "job")
defer span.End()
// Pass ctx to the work belonging to this job.
_ = ctx
```

The option extracts once when `WithContextFallback` is called and retains only the extracted span context and baggage.
Later changes to the carrier or process environment do not change that snapshot.
The supplied propagator determines which fields are read; the SDK does not assume a propagation format or read the global propagator.
For the composite propagator above, `envcar` reads `TRACEPARENT`, `TRACESTATE`, and `BAGGAGE`, using the specification's [environment-variable key normalization](https://opentelemetry.io/docs/specs/otel/context/env-carriers/#key-name-normalization).
The SDK does not launch child processes or inject their environments; the carrier also supports [explicit child-process injection](https://pkg.go.dev/go.opentelemetry.io/contrib/propagators/envcar#example-Carrier-ChildProcess).

#### Modes and parent selection

Use `WithContextFallbackMode` to select a mode:

| Mode | Incoming trace context | Incoming baggage |
| --- | --- | --- |
| `ContextFallbackParent` (default) | Used as the parent of an eligible span. | Merged into its returned context. |
| `ContextFallbackLink` | Appended as a link to an eligible span. | Merged into its returned context. |
| `ContextFallbackNone` | Not extracted or applied. | Not extracted or applied. |

An invalid mode disables fallback and reports an error through the OpenTelemetry error handler.
In parent or link mode, a nil propagator or carrier also disables fallback and reports an error.
When multiple fallback options configure one provider, the last one takes precedence.

```go
tp := sdktrace.NewTracerProvider(
	x.WithContextFallback(propagator, carrier,
		x.WithContextFallbackMode(x.ContextFallbackLink),
	),
)
```

Fallback applies only when the input context contains no span context and `trace.WithNewRoot()` was not requested.
Existing parents, including remote and non-recording parents, take precedence.
Partially populated span contexts are also preserved.
`trace.WithNewRoot()` skips both the incoming trace context and incoming baggage for that start; baggage already in the supplied context is retained.
The SDK cannot distinguish an intentionally empty span context from one that was never set.

Parent selection and link insertion happen before trace ID generation and sampling.
Parent mode preserves the incoming remote flag, trace flags, and tracestate when presenting the parent to the sampler.
Link mode preserves existing links, appends the incoming link if its trace and span IDs are not already present, and exposes that link to the sampler.
Recording remains subject to the configured sampler and span limits.
No process span is created.

#### Choose the scope deliberately

This policy applies to **every eligible span started by the provider**, throughout its lifetime.
For example, if a long-running HTTP server uses this provider, a later request without trace headers can inherit the process startup context even when that request is unrelated to the process invocation.
The SDK has no indication that those two operations are unrelated.
Ending the first job span does not consume or clear the fallback.

Use a provider with fallback only for work that belongs to the same invocation, or use `trace.WithNewRoot()` for independent operations.
For services that handle unrelated work, ordinary scoped extraction often expresses the boundary more precisely:

```go
invocationCtx := propagator.Extract(context.Background(), &envcar.Carrier{})
tp := sdktrace.NewTracerProvider() // No provider-wide fallback.
defer tp.Shutdown(context.Background())

ctx, span := tp.Tracer("example.com/worker").Start(invocationCtx, "job")
defer span.End()
// Pass ctx only to work belonging to this invocation.
_ = ctx
```

In that version, a later operation started with a separate `context.Background()` does not inherit the invocation context.
See the [compiled examples](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace/x#WithContextFallback) for both approaches.

#### Baggage and context lifetime

At an eligible start, incoming baggage is merged with the input context's baggage; existing application values win on key conflicts.
Incoming members that would exceed the W3C baggage limits are skipped, preserving application members.
Valid baggage can be used even if no valid trace context was extracted.
The returned context preserves the caller's cancellation, deadline, and other values.
Pass that returned context to subsequent work so that normal Go context propagation carries its span and baggage.
Changes or removals made to baggage in a descendant context are preserved when starting child spans; the startup baggage is not reapplied to those descendants.
An independent eligible root receives the original snapshot again.

This feature explores SDK-managed receiving-context behavior discussed in [specification issue #5398](https://github.com/open-telemetry/opentelemetry-specification/issues/5398) and [specification PR #5421](https://github.com/open-telemetry/opentelemetry-specification/pull/5421).
Parent mode experiments with a change to the current [Trace API parent-selection rule](https://opentelemetry.io/docs/specs/otel/trace/api/#determining-the-parent-span-from-a-context), which makes a span a root when its supplied context contains no span.
The specification does not currently define this persistent fallback policy; its scope and precedence need agreement before stabilization.

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
