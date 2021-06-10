# "Passthrough" setup for OpenTelemetry

Some Go programs may wish to propagate context without recording spans. To do this in OpenTelemetry, simply install `TextMapPropagators`, but do not install a TracerProvider using the SDK. This works because the default TracerProvider implementation returns a "Non-Recording" span that keeps the context of the caller but does not record spans.

For example, when you initialize your global settings, the following will propagate context without recording spans:

```golang
// Setup Propagators only
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
```

But the following will propagate context _and_ create new, potentially recorded spans:

```golang
// Setup SDK
exp, _ := stdout.New(stdout.WithPrettyPrint())
tp = sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exp),
)
otel.SetTracerProvider(tp)
// Setup Propagators
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
```

## The Demo

The demo has the following call structure:

`Outer -> Passthrough -> Inner`

If all components had both an SDK and propagators registered, we would expect the trace to look like:

```
|-------outer---------|
 |-Passthrough recv-|
  |Passthrough send|
    |---inner---|
```

However, in this demo, only the outer and inner have TracerProvider backed by the SDK. All components have Propagators set. In this case, we expect to see:

```
|-------outer---------|
    |---inner---|
```
