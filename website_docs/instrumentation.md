---
title: "Instrumentation"
weight: 3
---

Instrumentation is the process of adding observability code to your application. There are two general types of instrumentation - automatic, and manual - and you should be familiar with both in order to effectively instrument your software.

# Creating Spans

Spans are created by tracers, which can be acquired from a Tracer Provider.

```go
ctx := context.Background()
tracer := provider.Tracer("example/main")
var span trace.Span
ctx, span = tracer.Start(ctx, "helloWorld")
defer span.End()
```

In Go, the `context` package is used to store the active span. When you start a span, you'll get a handle on not only the span that's created, but the modified context that contains it. When starting a new span using this context, a parent-child relationship will automatically be created between the two spans, as seen here:

```go
func parentFunction() {
  ctx := context.Background()
  var parentSpan trace.Span
  ctx, parentSpan = tracer.Start(ctx, "parent")
  defer parentSpan.End()
  // call our child function
  childFunction(ctx)
  // do more work, when this function ends, parentSpan will complete.
}

func childFunction(ctx context.Context) {
  var childSpan trace.Span
  ctx, childSpan = tracer.Start(ctx, "child")
  defer childSpan.End()
  // do work here, when this function returns, childSpan will complete.
}
```

Once a span has completed, it is immutable and can no longer be modified.

## Attributes

Attributes are keys and values that are applied as metadata to your spans and are useful for aggregating, filtering, and grouping traces. Attributes can be added at span creation, or at any other time during the lifecycle of a span before it has completed.

```go
// setting attributes at creation...
ctx, span = tracer.Start(ctx, "attributesAtCreation", trace.WithAttributes(label.String("hello", "world")))
// ... and after creation
span.SetAttributes(label.Bool("isTrue", true), label.String("stringAttr", "hi!"))
```

Attribute keys can be precomputed, as well -

```go
var myKey = label.Key("myCoolAttribute")
span.SetAttributes(myKey.String("a value"))
```

### Semantic Attributes

Semantic Attributes are attributes that are defined by the OpenTelemetry Specification in order to provide a shared set of attribute keys across multiple languages, frameworks, and runtimes for common concepts like HTTP methods, status codes, user agents, and more. These attributes are available in the `go.opentelemetry.io/otel/semconv` package.

Tracing semantic conventions can be found [in this document](https://github.com/open-telemetry/opentelemetry-specification/tree/main/specification/trace/semantic_conventions)

## Events

An event is a human-readable message on a span that represents "something happening" during it's lifetime. For example, imagine a function that requires exclusive access to a resource that is under a mutex. An event could be created at two points - once, when we try to gain access to the resource, and another when we acquire the mutex.

```go
span.AddEvent("Acquiring lock")
mutex.Lock()
span.AddEvent("Got lock, doing work...")
// do stuff
span.AddEvent("Unlocking")
mutex.Unlock()
```

A useful characteristic of events is that their timestamps are displayed as offsets from the beginning of the span, allowing you to easily see how much time elapsed between them.

Events can also have attributes of their own -

```go
span.AddEvent("Cancelled wait due to external signal", trace.WithAttributes(label.Int("pid", 4328), label.String("signal", "SIGHUP")))
```

# Creating Metrics

The metrics API is currently unstable, documentation TBA.

# Propagators and Context

Traces can extend beyond a single process. This requires _context propagation_, a mechanism where identifiers for a trace are sent to remote processes.

In order to propagate trace context over the wire, a propagator must be registered with the OpenTelemetry SDK.

```go
import (
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/propagation"
)
...
otel.SetTextMapPropagator(propagation.TraceContext{})
```

> OpenTelemetry also supports the B3 header format, for compatibility with existing tracing systems (`go.opentelemetry.io/contrib/propagators/b3`) that do not support the W3C TraceContext standard.

After configuring context propagation, you'll most likely want to use automatic instrumentation to handle the behind-the-scenes work of actually managing serializing the context.

# Automatic Instrumentation

Automatic instrumentation, broadly, refers to instrumentation code that you didn't write. OpenTelemetry for Go supports this process through wrappers and helper functions around many popular frameworks and libraries. You can find a current list [here](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation), as well as at the [registry](/registry).
