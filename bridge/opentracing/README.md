# OpenTelemetry/OpenTracing Bridge

## Getting started

`go get go.opentelemetry.io/otel/bridge/opentracing`

Assuming you have configured an OpenTelemetry `TracerProvider`, these will be the steps to follow to wire up the bridge:

```go
import (
	"go.opentelemetry.io/otel"
	otelBridge "go.opentelemetry.io/otel/bridge/opentracing"
)

func main() {
	/* Create tracerProvider and configure OpenTelemetry ... */
	
	otelTracer := tracerProvider.Tracer("tracer_name")
	// Use the bridgeTracer as your OpenTracing tracer.
	bridgeTracer, wrapperTracerProvider := otelBridge.NewTracerPair(otelTracer)
	// Set the wrapperTracerProvider as the global OpenTelemetry
	// TracerProvider so instrumentation will use it by default.
	otel.SetTracerProvider(wrapperTracerProvider)

	/* ... */
}
```

## Interop from trace context from OpenTracing to OpenTelemetry

In order to get OpenTracing spans properly into the OpenTelemetry context, so they can be propagated (both internally, and externally), you will need to explicitly use the `BridgeTracer` for creating your OpenTracing spans, rather than a bare OpenTracing `Tracer` instance.

When you have started an OpenTracing Span, make sure the OpenTelemetry knows about it like this:

```go
	ctxWithOTSpan := opentracing.ContextWithSpan(ctx, otSpan)
	ctxWithOTAndOTelSpan := bridgeTracer.ContextWithSpanHook(ctxWithOTSpan, otSpan)
	// Propagate the otSpan to both OpenTracing and OpenTelemetry
	// instrumentation by using the ctxWithOTAndOTelSpan context.
```

## Extended Functionality

The bridge functionality can be extended beyond the OpenTracing API.

Any [`trace.SpanContext`](https://pkg.go.dev/go.opentelemetry.io/otel/trace#SpanContext) method can be accessed as following:

```go
type spanContextProvider interface {
	IsSampled() bool
	TraceID() trace.TraceID
	SpanID() trace.SpanID
	TraceFlags() trace.TraceFlags
	... // any other available method can be added here to access it
}

var sc opentracing.SpanContext = ...
if s, ok := sc.(spanContextProvider); ok {
	// Use TraceID by s.TraceID()
	// Use SpanID by s.SpanID()
	// Use TraceFlags by s.TraceFlags()
	...
}
```

## Migrating from OpenTracing to OpenTelemetry

If your codebase (or libraries you depend on) are still instrumented with OpenTracing, you can migrate incrementally by:
1) installing OpenTelemetry SDK,
2) wiring the OpenTracing bridge (shim),
3) progressively replacing OpenTracing instrumentation with OpenTelemetry instrumentation.

For a general, vendor-neutral migration approach, see the OpenTelemetry migration guide:
- [Migrating from OpenTracing](https://opentelemetry.io/docs/migration/opentracing/) :contentReference[oaicite:4]{index=4}

### Minimal Go example (OpenTracing API -> OpenTelemetry SDK)

This repository provides an OpenTracing bridge that forwards OpenTracing API calls to the OpenTelemetry SDK.
In Go, you can create a pair of tracers using `NewTracerPair()`:
- `BridgeTracer` implements the OpenTracing API (use it as the OpenTracing global tracer)
- `WrapperTracer` implements the OpenTelemetry API and cooperates with the bridge

```go
// Pseudo-code (exporter setup omitted for brevity)

import (
"github.com/opentracing/opentracing-go"
"go.opentelemetry.io/otel"
"go.opentelemetry.io/otel/sdk/trace"
otelBridge "go.opentelemetry.io/otel/bridge/opentracing"
)

func main() {
// 1) Configure OpenTelemetry SDK TracerProvider (exporter, resource, sampler, etc.)
tp := trace.NewTracerProvider(/* ... */)

// 2) Create an OpenTelemetry tracer you want to use under the hood
otelTracer := tp.Tracer("example-service")

// 3) Create a tracer pair for bridging
bridgeTracer, wrapperTracerProvider := otelBridge.NewTracerPair(otelTracer) // (*BridgeTracer, *WrapperTracerProvider) :contentReference[oaicite:2]{index=2}

// 4) Register globals:
// - OpenTracing global tracer for existing OpenTracing instrumentation
// - OpenTelemetry TracerProvider so new OTel instrumentation cooperates with OpenTracing context
opentracing.SetGlobalTracer(bridgeTracer)
otel.SetTracerProvider(wrapperTracerProvider)

/* ... */
}
```
