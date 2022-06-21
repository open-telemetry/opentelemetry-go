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

Create a new go module in a fresh directory:

```shell
mkdir otel-getting-started
cd otel-getting-started
go mod init main
```

Next, install OpenTelemetry and the `net/http` instrumentation package:

```
go get go.opentelemetry.io/otel \
  go.opentelemetry.io/otel/trace \
  go.opentelemetry.io/otel/sdk \
  go.opentelemetry.io/otel/exporters/otlp/otlptrace \
  go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp \
  go.opentelemetry.io/otel/exporters/stdout/stdouttrace \
  go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
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

When run, this will launch an HTTP server that does a "dice roll" and writes its
value to the response whenever the `/rolldice` path is accessed. For example,
accessing <http://localhost:8080/rolldice> in your browser will show the result
of a single "dice roll".

## Initialize tracing and add HTTP server instrumentation

[Instrumentation libraries](libraries.md) are used to create instrumentation on
your behalf. In this case, you can install OpenTelemetry and the `net/http`
instrumentation library so that calls to the server will start a trace that
contains data about the HTTP call.

To use the instrumentation library for `net/http`, you'll need to initialize
tracing and wrap the previously-defined HTTP handler.

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
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// Initialize a Tracerprovider, which is necessary to generate traces
// and export them to the console.
func newTracerProvider() *sdktrace.TracerProvider {
	exporter, err :=
		stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)

	if err != nil {
		panic(err)
	}

        // This includes the following resources:
        //
        // - sdk.language, sdk.version
        // - service.name, service.version, environment
        //
        // Including these resources is a good practice because it is commonly
        // used by various tracing backends to let you more accurately 
        // analyze your telemetry data.
	resource, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("diceroller-service"),
				semconv.ServiceVersionKey.String("v0.1.0"),
				attribute.String("environment", "demo"),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
}

// Same handler as before
func handleRollDice(w http.ResponseWriter, r *http.Request) {
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

	tp := newTracerProvider()
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

	wrapHandler()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

There's several things going on here:

* Creating an exporter that writes telemetry data to standard out
* Ensuring the right `Resource` attributes are initialized
* Initializing a `TracerProvider` and handling its shutdown
* Initializing the `net/http` instrumentation library and wrapping the HTTP
  handler

A `Resource` is useful metadata included in every trace. A `TracerProvider` is
used to create traces by letting you create a `Tracer`. You'll create your own
`Tracer` later in this article.

For details, see [Initializing Tracing](manual.md#initiallizing-a-new-tracer).

It may seem like quite a bit of code, but the good news is that very little
needs to change when you add more instrumentation later.

## Run the instrumented HTTP Server

When you run the app and access the `/rolldice` path, you'll see the same "dice
roll" values as before in the method you accessed it (browser, `curl`, etc.) and
telemetry data printed to the server process's standard out. The telemetry
generated is a single span that covers the lifetime of handling the request.

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
your application. For that, you’ll need to write some [manual
instrumentation](manual.md). Here’s how you can easily link up manual
instrumentation with automatic instrumentation.

Add a `Tracer` at the top level. Note that a `serviceName` and `serviceVersion`
are also added, since it is generally a good ideal to centralize constants like
this.

```go
import (
        //...
        "go.opentelemetry.io/otel"
        //...
)
//...
var (
	tracer         trace.Tracer
	serviceName    string = "diceroller-service"
	serviceVersion string = "0.1.0"
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

Creating a `Tracer` from a `TracerProvider` is necessary to let you create trace
spans manually.

The full code sample, with additional required imports, looks like this:

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
	tracer         trace.Tracer
	serviceName    string = "diceroller-service"
	serviceVersion string = "0.1.0"
)

func newTraceProvider() *sdktrace.TracerProvider {
	exporter, err :=
		stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)

	if err != nil {
		panic(err)
	}

        // This includes the following resources:
        //
        // - sdk.language, sdk.version
        // - service.name, service.version, environment
        //
        // Including these resources is a good practice because it is commonly
        // used by various tracing backends to let you more accurately 
        // analyze your telemetry data.
	resource, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String(serviceVersion),
				attribute.String("environment", "getting-started"),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
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

Now run the app again, and when you access the path you'll see similar output as
before, but this time with a span representing the manually-created one at the
top:

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

First, write the following collector configuration code into the `/tmp/`
directory on your local machine:

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

Then in a separate terminal window, run the docker command to acquire and run
the collector based on this configuration:

```
docker run -p 4318:4318 \
    -v /tmp/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
    otel/opentelemetry-collector:latest \
    --config=/etc/otel-collector-config.yaml
```

You will now have an OpenTelemetry Collector instance running locally.

### Modify the code to export spans via OTLP

The next step is to modify the code to send spans to the Collector via OTLP
instead of the console. Change the code to create an OTLP exporter, replacing
the console exporter from before with an OTLP HTTP exporter that talks to the
local endpoint the collector is listening on.

```go
import (
        //...
	"context"
        //..
)
exporter, err :=
	otlptracehttp.New(ctx,
		// WithInsecure lets us use http instead of https.
		// This is just for local development.
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(collectorAddr),
	)
```

The full code sample, with additional required imports, looks like this:

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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer         trace.Tracer
	serviceName    string = "diceroller-service"
	serviceVersion string = "0.1.0"
	collectorAddr  string = "localhost:4318" // HTTP endpoint for collector
)

func newTraceProvider(ctx context.Context) *sdktrace.TracerProvider {
	exporter, err :=
		otlptracehttp.New(ctx,
			// WithInsecure lets us use http instead of https.
			// This is just for local development.
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(collectorAddr),
		)

	if err != nil {
		panic(err)
	}

        // This includes the following resources:
        //
        // - sdk.language, sdk.version
        // - service.name, service.version, environment
        //
        // Including these resources is a good practice because it is commonly
        // used by various tracing backends to let you more accurately 
        // analyze your telemetry data.
	resource, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String(serviceVersion),
				attribute.String("environment", "getting-started"),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
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

When you run this server and access the `/rolldice` path, the collector process
will emit a log showing the two spans created from the request:

<details>
<summary>View example output</summary>

```
2022-05-10T16:20:09.715Z        INFO    loggingexporter/logging_exporter.go:41   TracesExporter  {"#spans": 2}
2022-05-10T16:20:09.716Z        DEBUG   loggingexporter/logging_exporter.go:51   ResourceSpans #0
Resource labels:
     -> environment: STRING(getting-started)
     -> service.name: STRING(diceroller-service)
     -> service.version: STRING(v0.1.0)
     -> telemetry.sdk.language: STRING(go)
     -> telemetry.sdk.name: STRING(opentelemetry)
     -> telemetry.sdk.version: STRING(1.7.0)
InstrumentationLibrarySpans #0
InstrumentationLibrary diceroller-service 
Span #0
    Trace ID       : 14f135f6f7414dc3d9272ba51dfc940a
    Parent ID      : d8995bd024c68216
    ID             : 0a587f78f3991936
    Name           : dice-roller
    Kind           : SPAN_KIND_INTERNAL
    Start time     : 2022-05-10 16:20:05.471544872 +0000 UTC
    End time       : 2022-05-10 16:20:05.471550922 +0000 UTC
    Status code    : STATUS_CODE_UNSET
    Status message : 
InstrumentationLibrarySpans #1
InstrumentationLibrary go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp semver:0.32.0
Span #0
    Trace ID       : 14f135f6f7414dc3d9272ba51dfc940a
    Parent ID      : 
    ID             : d8995bd024c68216
    Name           : rolldice
    Kind           : SPAN_KIND_SERVER
    Start time     : 2022-05-10 16:20:05.471522982 +0000 UTC
    End time       : 2022-05-10 16:20:05.471559732 +0000 UTC
    Status code    : STATUS_CODE_UNSET
    Status message : 
Attributes:
     -> net.transport: STRING(ip_tcp)
     -> net.peer.ip: STRING(10.20.150.7)
     -> net.peer.port: INT(36424)
     -> net.host.name: STRING(8080-cartermp-otelsamples-o7vjrp16ull.ws-us44.gitpod.io)
     -> http.method: STRING(GET)
     -> http.target: STRING(/rolldice?vscodeBrowserReqId=1652199605336)
     -> http.server_name: STRING(rolldice)
     -> http.client_ip: STRING(24.22.216.124)
     -> http.user_agent: STRING(Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36)
     -> http.scheme: STRING(http)
     -> http.host: STRING(8080-cartermp-otelsamples-o7vjrp16ull.ws-us44.gitpod.io)
     -> http.flavor: STRING(1.1)
     -> http.wrote_bytes: INT(1)
     -> http.status_code: INT(200)

2022-05-10T17:11:04.807Z        INFO    loggingexporter/logging_exporter.go:41
TracesExporter  {"#spans": 2} 2022-05-10T17:11:04.807Z        DEBUG
loggingexporter/logging_exporter.go:51   ResourceSpans #0 Resource labels: ->
     environment: STRING(getting-started) -> service.name:
     STRING(diceroller-service) -> service.version: STRING(0.1.0) ->
     telemetry.sdk.language: STRING(go) -> telemetry.sdk.name:
     STRING(opentelemetry) -> telemetry.sdk.version: STRING(1.7.0)
     InstrumentationLibrarySpans #0 InstrumentationLibrary diceroller-service
     Span #0 Trace ID       : 13d8990c0d82fb6a7fe6307bc94acc28 Parent ID      :
2e82b9ec5b4853ff ID             : 1e2cd857b1f1f7fd Name           : dice-roller
Kind           : SPAN_KIND_INTERNAL Start time     : 2022-05-10
17:11:02.092278673 +0000 UTC End time       : 2022-05-10 17:11:02.092283423
    +0000 UTC Status code    : STATUS_CODE_UNSET Status message :
    InstrumentationLibrarySpans #1 InstrumentationLibrary
    go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp semver:0.32.0
    Span #0 Trace ID       : 13d8990c0d82fb6a7fe6307bc94acc28 Parent ID      :
    ID             : 2e82b9ec5b4853ff Name           : rolldice Kind           :
    SPAN_KIND_SERVER Start time     : 2022-05-10 17:11:02.092262193 +0000 UTC
    End time       : 2022-05-10 17:11:02.092290013 +0000 UTC Status code    :
    STATUS_CODE_UNSET Status message : Attributes: -> net.transport:
    STRING(ip_tcp) -> net.peer.ip: STRING(10.20.150.7) -> net.peer.port:
INT(36426) -> net.host.name:
STRING(8080-cartermp-otelsamples-o7vjrp16ull.ws-us44.gitpod.io) -> http.method:
STRING(GET) -> http.target: STRING(/rolldice) -> http.server_name:
    STRING(rolldice) -> http.client_ip: STRING(24.22.216.124) ->
    http.user_agent: STRING(Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)
    AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36)
    -> http.scheme: STRING(http) -> http.host:
    STRING(8080-cartermp-otelsamples-o7vjrp16ull.ws-us44.gitpod.io) ->
    http.flavor: STRING(1.1) -> http.wrote_bytes: INT(1) -> http.status_code:
    INT(200)
```

</details>

## Next steps

There are several options available for automatic instrumentation and Go. See
[Using instrumentation libraries](libraries.md) to learn about them and how to
configure them.

There’s a lot more to manual instrumentation than just creating a child span. To
learn details about initializing manual instrumentation and many more parts of
the OpenTelemetry API you can use, see [Manual Instrumentation](manual.md).

Finally, there are several options for exporting your telemetry data with
OpenTelemetry. To learn how to export your data to a preferred backend, see
[Processing and Exporting Data](exporting_data.md).
