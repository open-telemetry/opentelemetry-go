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

### `SpanContext.IsSampled`

Return the underlying OpenTelemetry [`Span.IsSampled`](https://pkg.go.dev/go.opentelemetry.io/otel/trace#SpanContext.IsSampled) value by converting a `bridgeSpanContext`.

```go
type samplable interface {
	IsSampled() bool
}

var sc opentracing.SpanContext = ...
if s, ok := sc.(samplable); ok && s.IsSampled() {
	// Do something with sc knowing it is sampled.
}
```

### `SpanContext.TraceID`

Return the underlying OpenTelemetry `Span.TraceID` value by converting a `bridgeSpanContext`.

```go
type traceId interface {
	TraceID() bool
}

var sc opentracing.SpanContext = ...
if s, ok := sc.(traceId); ok {
	// Use TraceID by s.TraceID()
}
```

### `SpanContext.SpanID`

Return the underlying OpenTelemetry `Span.SpanID` value by converting a `bridgeSpanContext`.

```go
type spanId interface {
	SpanID() bool
}

var sc opentracing.SpanContext = ...
if s, ok := sc.(spanId); ok {
	// Use SpanID by s.SpanID()
}
```
