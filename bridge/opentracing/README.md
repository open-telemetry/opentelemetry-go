# OpenTelemetry/OpenTracing Bridge

### Getting started

`go get go.opentelemetry.io/otel/bridge/opentracing`

Assuming you have configured an OpenTelemetry TracerProvider, these will be the steps to follow to wire up the bridge:

```go
  import {
    "go.opentelemetry.io/otel"
    otelBridge "go.opentelemetry.io/otel/bridge/opentracing"
  }
  
  ...

  otelTracer := tracerProvider.Tracer("tracer_name")
  // use the bridgeTracer as your ot tracer, and the wrapperTracerProvider for otel instrumentation.
  bridgeTracer, wrapperTracerProvider := otelBridge.NewTracerPair(otelTracer)
  otel.SetTracerProvider(wrapperTracerProvider)
```

### Interop from ot -> otel

In order to get ot spans properly into the otel context, so they can be propagated (both internally, and externally), you will need to explicitly use the `BridgeTracer` for creating your ot spans, rather than a bare ot `Tracer` instance.

When you have started an ot Span, make sure the otel knows about it like this:

```go
    contextWithOtSpan := span.ToContext(currentContext)
    contextWithOtelBridgeSpan := bridgeTracer.ContextWithSpanHook(contextWithOtSpan, span.span)
    // use the contextWithOtelBridgeSpan instance for calls that need to have the span propagated
```
