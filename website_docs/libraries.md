---
title: Using instrumentation libraries
weight: 3
linkTitle: Libraries
aliases: [/docs/instrumentation/go/using_instrumentation_libraries, /docs/instrumentation/go/automatic_instrumentation]
---

Go does not support truly automatic instrumentation like other languages today. Instead, you'll need to depend on [instrumentation libraries](/docs/reference/specification/glossary/#instrumentation-library) that generate telemetry data for a particular instrumented library. For example, the instrumentation library for `net/http` will automatically create spans that track inbound and outbound requests once you configure it in your code.

## Setup

Each instrumentation library is a package. In general, this means you need to `go get` the appropriate package:

```console
go get go.opentelemetry.io/contrib/instrumentation/{import-path}/otel{package-name}
```

And then configure it in your code based on what the library requires to be activated.

## Example with `net/http`

As an example, here's how you can set up automatic instrumentation for inbound HTTP requests for `net/http`:

First, get the `net/http` instrumentation library:

```console
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
```

Next, use the library to wrap an HTTP handler in your code:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Package-level tracer.
// This should be configured in your code setup instead of here.
var tracer = otel.Tracer("github.com/full/path/to/mypkg")

// sleepy mocks work that your application does.
func sleepy(ctx context.Context) {
	_, span := tracer.Start(ctx, "sleep")
	defer span.End()

	sleepTime := 1 * time.Second
	time.Sleep(sleepTime)
	span.SetAttributes(attribute.Int("sleep.duration", int(sleepTime)))
}

// httpHandler is an HTTP handler function that is going to be instrumented.
func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! I am instrumented automatically!")
	ctx := r.Context()
	sleepy(ctx)
}

func main() {
	// Wrap your httpHandler function.
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "hello-instrumented")
	http.Handle("/hello-instrumented", wrappedHandler)

	// And start the HTTP serve.
	log.Fatal(http.ListenAndServe(":3030", nil))
}
```

Assuming that you have a `Tracer` and [exporter]({{< relref "exporting_data" >}}) configured, this code will:

* Start an HTTP server on port `3030`
* Automatically generate a span for each inbound HTTP request to `/hello-instrumented`
* Create a child span of the automatically-generated one that tracks the work done in `sleepy`

Connecting manual instrumentation you write in your app with instrumentation generated from a library is essential to get good observability into your apps and services.

## Available packages

A full list of instrumentation libraries available can be found in the [OpenTelemetry registry](/registry/?language=go&component=instrumentation).

## Next steps

Instrumentation libraries can do things like generate telemetry data for inbound and outbound HTTP requests, but they don't instrument your actual application.

To get richer telemetry data, use [manual instrumentation]({{< relref "manual" >}}) to enrich your telemetry data from instrumentation libraries with instrumentation from your running application.
