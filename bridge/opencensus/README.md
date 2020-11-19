# OpenCensus Bridge

The OpenCensus Bridge helps facilitate the migration of an application from OpenCensus to OpenTelemetry.

## The Problem: Mixing OpenCensus and OpenTelemetry libraries

In a perfect world, one would simply migrate their entire go application --including custom instrumentation, libraries, and exporters-- from OpenCensus to OpenTelemetry all at once.  In the real world, dependency constraints, third-party ownership of libraries, or other reasons may require mixing OpenCensus and OpenTelemetry libraries in a single application.

However, if you create the following spans in a go application:

```golang
ctx, ocSpan := opencensus.StartSpan(context.Background(), "OuterSpan")
defer ocSpan.End()
ctx, otSpan := opentelemetryTracer.Start(ctx, "MiddleSpan")
defer otSpan.End()
ctx, ocSpan := opencensus.StartSpan(ctx, "InnerSpan")
defer ocSpan.End()
```

OpenCensus reports (to OpenCensus exporters):

```
[--------OuterSpan------------]
    [----InnerSpan------]
```

OpenTelemetry reports (to OpenTelemetry exporters):

```
   [-----MiddleSpan--------]
```

Instead, I would prefer (to a single set of exporters):

```
[--------OuterSpan------------]
   [-----MiddleSpan--------]
    [----InnerSpan------]
```

### The bridge solution

The bridge implements the OpenCensus trace API using OpenTelemetry.  This would cause, for example, a span recorded with OpenCensus' `StartSpan()` method to be equivalent to recording a span using OpenTelemetry's `tracer.Start()` method.  Funneling all tracing API calls to OpenTelemetry APIs results in the desired unified span hierarchy.

### User Journey

1. Instantiate OpenTelemetry SDK and Exporters
2. Override OpenCensus' DefaultTracer with the bridge
3. Migrate libraries from OpenCensus to OpenTelemetry
4. Remove OpenCensus Exporters

To override OpenCensus' DefaultTracer with the bridge:
```golang
import (
    octrace "go.opencensus.io/trace"
    "go.opentelemetry.io/otel/bridge/opencensus"
    "go.opentelemetry.io/otel"
)

tracer := otel.GetTracerProvider().Tracer("bridge")
octrace.DefaultTracer = opencensus.NewTracer(tracer)
```

Be sure to set the `Tracer` name to your instrumentation package name instead of `"bridge"`.

### Incompatibilities

OpenCensus and OpenTelemetry APIs are not entirely compatible.  If the bridge finds any incompatibilities, it will log them.  Incompatibilities include:

* Custom OpenCensus Samplers specified during StartSpan are ignored.
* Links cannot be added to OpenCensus spans.
* OpenTelemetry Debug or Deferred trace flags are dropped after an OpenCensus span is created.
