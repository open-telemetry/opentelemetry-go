# Experimental Features

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/trace/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace/x)

The Trace SDK contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the OpenTelemetry Go Trace SDK prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [OnEndingSpanProcessor](#onendingspanprocessor)
- [ProbabilitySampler](#probabilitysampler)

### OnEndingSpanProcessor

`OnEndingSpanProcessor` is an optional extension interface for
[`sdktrace.SpanProcessor`](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#SpanProcessor)
implementations. When a registered processor also implements this interface,
its `OnEnding` method is called during [`Span.End`](https://pkg.go.dev/go.opentelemetry.io/otel/trace#Span),
after the end timestamp is set but while the span is still mutable. This
allows processors to inspect and modify a span before it is exported — for
example, to implement tail-based sampling or span enrichment.

All `OnEnding` callbacks run in registration order, before any `OnEnd`
callback is invoked. The `OnEnding` method must not block or retain the span
past the call.

#### Usage

```go
import (
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/sdk/trace/x"
)

type myProcessor struct{}

func (p *myProcessor) OnStart(ctx context.Context, s sdktrace.ReadWriteSpan) {}
func (p *myProcessor) OnEnd(s sdktrace.ReadOnlySpan)                          {}
func (p *myProcessor) Shutdown(ctx context.Context) error                     { return nil }
func (p *myProcessor) ForceFlush(ctx context.Context) error                   { return nil }

// OnEnding implements x.OnEndingSpanProcessor.
func (p *myProcessor) OnEnding(s sdktrace.ReadWriteSpan) {
    // Modify the span before it becomes read-only.
    s.SetAttributes(attribute.Bool("sampled-by-tail", true))
}

// Compile-time assertion that myProcessor satisfies x.OnEndingSpanProcessor.
var _ x.OnEndingSpanProcessor = (*myProcessor)(nil)

tp := sdktrace.NewTracerProvider()
tp.RegisterSpanProcessor(&myProcessor{})
```

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

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.
