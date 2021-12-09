---
title: Automatic Instrumentation
weight: 3
---

Automatic instrumentation in Go is done by pulling in dependencies for the component you'd like to have automatically instrumented. For example, you can use the automatic instrumentation package for `net/http` to automatically create spans that track inbound and outbound requests from your app or service.

## Setup

Each automatic instrumentation package is pulled in as a dependency for your app.

In general, this means you `go get` the appropriate package:

```console
go get go.opentelemetry.io/contrib/instrumentation/{import-path}/otel{package-name}
```

And you use it in your code directly (typically wrapping something).

## Example with `net/http`

As an example, here's how you can set up automatic instrumentation for inbound HTTP requests:

First, get the `net/http` automatic instrumentation package:

```console
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
```

Next, wrap an HTTP handler in your code:

```go

// Import the package
import (
  // ...
  "time"
  "net/http"
  "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
  // ...
)

// ...

// Have a function that does some work,
// which will connect with automatic instruemtnation later in the code
func sleepy(ctx context.Context) {
    _, span := tracer.Start(ctx, "sleep")
    defer span.End()

    sleepTime := 1 * time.Second
    time.Sleep(sleepTime)

    span.SetAttributes(attribute.Int("sleep.duration", sleepTime))
}

// Have some http handler function to instrument
func httpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World! I am instrumented autoamtically!")

    ctx := r.Context()

    // Do some work
    sleepy(ctx)
}

// Instantiate and wrap the httpHandler function that you just defined
handler := http.HandlerFunc(httpHandler)
wrappedHandler := otelhttp.NewHandler(handler, "hello-instrumented")
http.Handle("/hello-instrumented", wrappedHandler)

// And start the HTTP server
log.Fatal(http.ListenAndServe(":3030", nil))
```

Assuming that you have a `Tracer` and [exporter](exporting_data.md) configured, this code will:

* Start an HTTP server on port `3030`
* Generate a span for each inbound HTTP request to `/hello-instrumented`
* Create a child span of the automatically-generated one that tracks the work done in `sleepy`

Connecting manual instrumentation with automatic instrumentation is essential to get good observability into your apps and services.

## Available packages

A full list of packages that offer automatic instrumentation for particular libraries can be found in the [OpenTelementry registry](https://opentelemetry.io/registry/?language=go&component=instrumentation).

## Next steps

Automatic instrumentation can do things like generate telemtry data for inbound and outbound HTTP requests, but it doesn't cover your actual application data.

To get richer telemetry data, use [manual instrumentatiion](instrumentation.md) to enrich your traces with information about your running application.
