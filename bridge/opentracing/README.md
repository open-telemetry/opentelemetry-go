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

While the bridge does not expose functionality that is not implemented by OpenTelemetry, it does expose some that is part of OpenTelemetry API and not OpenTracing API.

**`SpanContext.IsSampled`**

Proxies underlying `trace.IsSampled` method (see [documentation](https://pkg.go.dev/go.opentelemetry.io/otel/trace#SpanContext.IsSampled)). In order to use it, you have to cast it:

```go
type samplable interface {
	IsSampled() bool
}

var sc opentracing.SpanContext = ...
if sc.(samplable).IsSampled() {
	// Span is expected to be sampled.
} else {
	// Span will be discarded.
}
```
