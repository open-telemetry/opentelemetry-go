---
title: "Getting Started"
weight: 2
---

In this guide, you'll learn how to set up and get tracing telemetry from an HTTP
server using Go.

You can find more elaborate examples that use different OpenTelemetry libraries
in [example subfolders for each
library](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation)

## Installation

To begin, you'll want to install OpenTelemetry and the `net/http`
instrumentation package:

```
go get go.opentelemetry.io/otel \
  go.opentelemetry.io/otel/trace \
  go.opentelemetry.io/otel/sdk \
  go.opentelemetry.io/otel/exporters/otlp/otlptrace \
  go.opentelemetry.io/otel/exporters/stdout/stdouttrace \
  go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp \
  go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
```

## Create the sample HTTP Server

In a new project, create or edit your `main.go` file to be the following:

```go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
        "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

// Implement an HTTP Handler func to be instrumented later
func handleRollDice(w http.ResponseWriter, r *http.Request) {
	value := rand.Intn(6) + 1
	fmt.Fprintf(w, "%d", value)
}

func main() {
        http.HandleFunc("/rolldice", handleRollDice)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

When run, this will launch an HTTP server that does a "dice roll" whenever the
`/rolldice` route is accessed.

## Add HTTP Server instrumentation

[Instrumentation libraries]({{< relref "libraries" >}}) are used to create
instrumentation on your behalf. In this case, you can install OpenTelemetry and
the `net/http` instrumentation library so that calls to the server will start a
trace that contains data about the HTTP call.

Replace the contents of `main.go` with the following:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer      trace.Tracer
	serviceName string = "diceroller-service"
)

func newTraceProvider() *sdktrace.TracerProvider {
	exp, err :=
		stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)

	if err != nil {
		panic(err)
	}

	r, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String("v0.1.0"),
				attribute.String("environment", "getting-started"),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func handleRollDice(w http.ResponseWriter, r *http.Request) {
	// Create a child span called dice-roller that tracks only this function call
	_, span := tracer.Start(r.Context(), "dice-roller")
	defer span.End()

	value := rand.Intn(6) + 1
	fmt.Fprintf(w, "%d", value)
}

// Wrap the handleRollDice so that telemetry data
// can be automatically generated for it
func wrapHandler() {
	handler := http.HandlerFunc(handleRollDice)
	wrappedHandler := otelhttp.NewHandler(handler, "rolldice")
	http.Handle("/rolldice", wrappedHandler)
}

func main() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// Service name to be use by observability tool
			semconv.ServiceNameKey.String("roll-dice")))
	// Checking for errors
	if err != nil {
		fmt.Printf("Error adding %v to the tracer engine: %v", "applicationName", err)
	}

	collectorAddr := "localhost:4318"
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(collectorAddr),
	)

	// Checking for errors
	if err != nil {
		fmt.Printf("Error initializing the tracer exporter: %v", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Register context and baggage propagation.
	// Although not strictly necessary, for this sample,
	// it is required for distributed tracing.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tracer = tp.Tracer(serviceName)

	wrapHandler()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

This will initialize tracing and wrap the HTTP handler so that tracing data can
be generated for it automatically.

This code also initializes a `Tracer` and sets it, 

## Run the instrumented HTTP Server

When you run the app and access the `/rolldice` route, you'll see telemetry data
printed to the server process's standard out. The telemetry generated is a
single span that covers the lifetime of handling the request.

<details>
<summary>View example output</summary>

```json
{
        "Name": "rolldice",
        "SpanContext": {
                "TraceID": "d3bebe3482a3e6b2de1ad61c1c0aaae4",
                "SpanID": "4baf164509867e98",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "Parent": {
                "TraceID": "00000000000000000000000000000000",
                "SpanID": "0000000000000000",
                "TraceFlags": "00",
                "TraceState": "",
                "Remote": false
        },
        "SpanKind": 2,
        "StartTime": "0001-01-01T00:00:00Z",
        "EndTime": "0001-01-01T00:00:00Z",
        "Attributes": [
                {
                        "Key": "net.transport",
                        "Value": {
                                "Type": "STRING",
                                "Value": "ip_tcp"
                        }
                },
                {
                        "Key": "net.peer.ip",
                        "Value": {
                                "Type": "STRING",
                                "Value": "::1"
                        }
                },
                {
                        "Key": "net.peer.port",
                        "Value": {
                                "Type": "INT64",
                                "Value": 55121
                        }
                },
                {
                        "Key": "net.host.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "localhost"
                        }
                },
                {
                        "Key": "net.host.port",
                        "Value": {
                                "Type": "INT64",
                                "Value": 8080
                        }
                },
                {
                        "Key": "http.method",
                        "Value": {
                                "Type": "STRING",
                                "Value": "GET"
                        }
                },
                {
                        "Key": "http.target",
                        "Value": {
                                "Type": "STRING",
                                "Value": "/rolldice"
                        }
                },
                {
                        "Key": "http.server_name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "rolldice"
                        }
                },
                {
                        "Key": "http.user_agent",
                        "Value": {
                                "Type": "STRING",
                                "Value": "<browser user agent here>"
                        }
                },
                {
                        "Key": "http.scheme",
                        "Value": {
                                "Type": "STRING",
                                "Value": "http"
                        }
                },
                {
                        "Key": "http.host",
                        "Value": {
                                "Type": "STRING",
                                "Value": "localhost:8080"
                        }
                },
                {
                        "Key": "http.flavor",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.1"
                        }
                },
                {
                        "Key": "http.wrote_bytes",
                        "Value": {
                                "Type": "INT64",
                                "Value": 1
                        }
                },
                {
                        "Key": "http.status_code",
                        "Value": {
                                "Type": "INT64",
                                "Value": 200
                        }
                }
        ],
        "Events": null,
        "Links": null,
        "Status": {
                "Code": "Unset",
                "Description": ""
        },
        "DroppedAttributes": 0,
        "DroppedEvents": 0,
        "DroppedLinks": 0,
        "ChildSpanCount": 0,
        "Resource": [
                {
                        "Key": "environment",
                        "Value": {
                                "Type": "STRING",
                                "Value": "demo"
                        }
                },
                {
                        "Key": "service.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "diceroller-service"
                        }
                },
                {
                        "Key": "service.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "v0.1.0"
                        }
                },
                {
                        "Key": "telemetry.sdk.language",
                        "Value": {
                                "Type": "STRING",
                                "Value": "go"
                        }
                },
                {
                        "Key": "telemetry.sdk.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "opentelemetry"
                        }
                },
                {
                        "Key": "telemetry.sdk.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.7.0"
                        }
                }
        ],
        "InstrumentationLibrary": {
                "Name": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
                "Version": "semver:0.32.0",
                "SchemaURL": ""
        }
}
```

</details>

## Add manual instrumentation

Automatic instrumentation captures telemetry at the edges of your systems, such
as inbound and outbound HTTP requests, but it doesn’t capture what’s going on in
your application. For that you’ll need to write some [manual
instrumentation]({{< relref "manual" >}}). Here’s how you can easily link up
manual instrumentation with automatic instrumentation.

Add some top-level variables, namely, a `Tracer`:

```go
var (
	tracer      trace.Tracer
	serviceName string = "diceroller-service"
)
```

Next, modify the `handleRollDice` function as follows:

```go
func handleRollDice(w http.ResponseWriter, r *http.Request) {
	// Create a child span called dice-roller that tracks only this function call
	_, span := tracer.Start(r.Context(), "dice-roller")
	defer span.End()

	value := rand.Intn(6) + 1
	fmt.Fprintf(w, "%d", value)
}
```

The call to `tracer.Start` is how you can create a span that's connected to the
automatic instrumentation. Specifically, this will create a child span, where
its parent is the span representing the full lifetime of the request handling.

Finally, at the bottom of the `main` function, create the `Tracer`, after you
set up context propagation:

```go
tracer = tp.Tracer(serviceName)
```

The full sample code looks like this:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer      trace.Tracer
	serviceName string = "diceroller-service"
)

func newTraceProvider() *sdktrace.TracerProvider {
	exp, err :=
		stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)

	if err != nil {
		panic(err)
	}

	r, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String("v0.1.0"),
				attribute.String("environment", "getting-started"),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func handleRollDice(w http.ResponseWriter, r *http.Request) {
	// Create a child span called dice-roller that tracks only this function call
	_, span := tracer.Start(r.Context(), "dice-roller")
	defer span.End()

	value := rand.Intn(6) + 1
	fmt.Fprintf(w, "%d", value)
}

// Wrap the handleRollDice so that telemetry data
// can be automatically generated for it
func wrapHandler() {
	handler := http.HandlerFunc(handleRollDice)
	wrappedHandler := otelhttp.NewHandler(handler, "rolldice")
	http.Handle("/rolldice", wrappedHandler)
}

func main() {
	ctx := context.Background()

	tp := newTraceProvider()
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Register context and baggage propagation.
	// Although not strictly necessary, for this sample,
	// it is required for distributed tracing.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tracer = tp.Tracer(serviceName)

	wrapHandler()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

Now run the app again, and when you access the route you'll see similar output
as before, but this time with a span representing the manually-created one at
the top:

<details>
<summary>View example output</summary>

```json
{
        "Name": "dice-roller",
        "SpanContext": {
                "TraceID": "28ba593427c73784cfc271c6e48e7d2c",
                "SpanID": "1d437a88e928fde8",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "Parent": {
                "TraceID": "28ba593427c73784cfc271c6e48e7d2c",
                "SpanID": "fac18d9dbbc081ba",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "SpanKind": 1,
        "StartTime": "0001-01-01T00:00:00Z",
        "EndTime": "0001-01-01T00:00:00Z",
        "Attributes": null,
        "Events": null,
        "Links": null,
        "Status": {
                "Code": "Unset",
                "Description": ""
        },
        "DroppedAttributes": 0,
        "DroppedEvents": 0,
        "DroppedLinks": 0,
        "ChildSpanCount": 0,
        "Resource": [
                {
                        "Key": "environment",
                        "Value": {
                                "Type": "STRING",
                                "Value": "getting-started"
                        }
                },
                {
                        "Key": "service.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "diceroller-service"
                        }
                },
                {
                        "Key": "service.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "v0.1.0"
                        }
                },
                {
                        "Key": "telemetry.sdk.language",
                        "Value": {
                                "Type": "STRING",
                                "Value": "go"
                        }
                },
                {
                        "Key": "telemetry.sdk.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "opentelemetry"
                        }
                },
                {
                        "Key": "telemetry.sdk.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.7.0"
                        }
                }
        ],
        "InstrumentationLibrary": {
                "Name": "diceroller-service",
                "Version": "",
                "SchemaURL": ""
        }
}
```

</details>

You'll find that the `dice-roller` span has information about a parent span, and
the ID of that parent span matches the span created by automatic
instrumentation.

## Send traces to an OpenTelemetry Collector

The [OpenTelemetry Collector](/docs/collector/getting-started/) is a critical
component of most production deployments. Some examples of when it's beneficial
to use a collector:

* A single telemetry sink shared by multiple services, to reduce overhead of
  switching exporters
* Aggregating traces across multiple services, running on multiple hosts
* A central place to process traces prior to exporting them to a backend

Unless you have just a single service or are experimenting, you'll want to use a
collector in production deployments.

### Configure and run a local collector

First, write the following collector configuration code into `/tmp/`:

```yaml
# /tmp/otel-collector-config.yaml
receivers:
    otlp:
        protocols:
            grpc:
            http:
exporters:
    logging:
        loglevel: debug
processors:
    batch:
service:
    pipelines:
        traces:
            receivers: [otlp]
            exporters: [logging]
            processors: [batch]
```

Then run the docker command to acquire and run the collector based on this
configuration:

```
docker run -p 4318:4318 \
    -v /tmp/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
    otel/opentelemetry-collector:latest \
    --config=/etc/otel-collector-config.yaml
```

You will now have an OpenTelemetry Collector instance running locally.

### Modify the code to export spans via OTLP

The next step is to modify the code to send spans to the Collector via OTLP
instead of the console.

To do this, install the OTLP exporter packages:

```
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace \
  go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp \
```

Next, change the code to create an OTLP exporter, replacing the console exporter
from before:

```go
exp, err := otlptrace.New(ctx, otlptracehttp.NewClient())
```

The full code sample looks like this:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer      trace.Tracer
	serviceName string = "diceroller-service"
)

func newTraceProvider(ctx context.Context) *sdktrace.TracerProvider {
	exp, err := otlptrace.New(ctx, otlptracehttp.NewClient())

	if err != nil {
		panic(err)
	}

	r, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String("v0.1.0"),
				attribute.String("environment", "getting-started"),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func handleRollDice(w http.ResponseWriter, r *http.Request) {
	// Create a child span called dice-roller that tracks only this function call
	_, span := tracer.Start(r.Context(), "dice-roller")
	defer span.End()

	value := rand.Intn(6) + 1
	fmt.Fprintf(w, "%d", value)
}

// Wrap the handleRollDice so that telemetry data
// can be automatically generated for it
func wrapHandler() {
	handler := http.HandlerFunc(handleRollDice)
	wrappedHandler := otelhttp.NewHandler(handler, "rolldice")
	http.Handle("/rolldice", wrappedHandler)
}

func main() {
	ctx := context.Background()

	tp := newTraceProvider(ctx)
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Register context and baggage propagation.
	// Although not strictly necessary, for this sample,
	// it is required for distributed tracing.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tracer = tp.Tracer(serviceName)

	wrapHandler()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

TODO

you should now start to see traces from the OTel collector process. Except
phillip can't seem to get that working because, you know, he is not good at
computers.

## Next steps

There are several options available for automatic instrumentation and Go. See
[Using instrumentation libraries]({{< relref "libraries" >}}) to learn about
them and how to configure them.

There’s a lot more to manual instrumentation than just creating a child span. To
learn details about initializing manual instrumentation and many more parts of
the OpenTelemetry API you can use, see [Manual Instrumentation]({{< relref
"manual" >}}).

Finally, there are several options for exporting your telemetry data with
OpenTelemetry. To learn how to export your data to a preferred backend, see
[Processing and Exporting Data]({{< relref "exporting_data" >}}).
