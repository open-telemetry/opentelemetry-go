# OpenTelemetry/OpenTracing Bridge

## Getting started

`go get go.opentelemetry.io/otel/bridge/opentracing`

Assuming you have configured an OpenTelemetry TracerProvider, these will be the steps to follow to wire up the bridge:

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

## Interop from ot -> otel

In order to get ot spans properly into the otel context, so they can be propagated (both internally, and externally), you will need to explicitly use the `BridgeTracer` for creating your ot spans, rather than a bare ot `Tracer` instance.

When you have started an ot Span, make sure the otel knows about it like this:

```go
    contextWithOtSpan := span.ToContext(currentContext)
    contextWithOtelBridgeSpan := bridgeTracer.ContextWithSpanHook(contextWithOtSpan, span.span)
    // use the contextWithOtelBridgeSpan instance for calls that need to have the span propagated
```
