# OpenTelemetry/OpenTracing Bridge

## Deprecation

The `go.opentelemetry.io/otel/bridge/opentracing` module and its packages are deprecated as of v1.48.0.
v1.48.0 is the final planned release of this module.
OpenTracing compatibility requirements have been [deprecated in the OpenTelemetry specification](https://opentelemetry.io/blog/2026/deprecating-opentracing-compatibility/).
Use native OpenTelemetry APIs and SDKs instead; see [Migrating from OpenTracing](https://opentelemetry.io/docs/compatibility/migration/opentracing/).

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
