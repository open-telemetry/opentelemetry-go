---
title: "Getting Started"
weight: 2
---

Welcome to the OpenTelemetry for Go getting started guide! This guide will walk you through the basic steps in installing, instrumenting with, configuring, and exporting data from OpenTelemetry. Before you get started, be sure to have Go 1.15 or newer installed.

This guide will walk you the common situation where you already have an application that uses a library and want to add observability. The application will use asks users what Fibonacci number they would like generated and returns to them the computed value. To start, make a new directory named `fib` and add the following to a new file named `fibonacci.go` in that directory.

```go
package main

// Fibonacci returns the n-th fibonacci number.
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}
```

Now that you have your core logic added, you can build your application around it. Add a new `app.go` file with the following application logic.

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
)

// App is an Fibonacci computation application.
type App struct {
	r io.Reader
	l *log.Logger
}

// NewApp returns a new App.
func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		n, err := a.Poll()
		if err != nil {
			return err
		}

		a.Write(n)
	}
}

// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	return n, err
}

// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	f, err := Fibonacci(n)
	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}
```

With your application fully composed, you need a `main()` function to actually run the application. In a new `main.go` file add the following run logic.

```go
package main

import (
	"context"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "", 0)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
```

With the code complete it is time to run the application. Before you can do that you need to initialize this directory as a Go module. In your terminal, run the command `go mod init fib` in the `fib` directory. This will create a `go.mod` file, which is used by Go to manage imports. Now you should be able to run the application!

```sh
$ go run .
What Fibonacci number would you like to know:
42
Fibonacci(42) = 267914296
What Fibonacci number would you like to know:
^C
goodbye
```

# Trace Instrumentation

OpenTelemetry is split into two parts: an API to instrument code with, and SDKs that implement the API. To start integrating OpenTelemetry into any project, the API is used to define how telemetry is generated. To generate tracing telemetry you will use the OpenTelemetry Trace API from the `go.opentelemetry.io/otel/trace` package.

First, to install the necessary prerequisites for the Trace API, install the appropriate packages. Run the following command in your working directory.

```sh
go get go.opentelemetry.io/otel@v1.0.0-RC1 \
       go.opentelemetry.io/otel/trace@v1.0.0-RC1
```

Now you can instrument the application! First add the needed imports to your `app.go` file.

```go
import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)
```

With the imports handled you are almost ready to add tracing instrumentation to the application! First you need to consider how you will identify the telemetry you create as coming from the instrumentation library you will build. OpenTelemetry does this by naming `Tracer`s with the instrumentation library name. The first thing to add to `app.go` is a constant with the package name.

```go
// name is the Tracer name used to identify this instrumentation library.
const name = "fib"
```

Now that you have that out of the way, you can create traces from the appropriately named `Tracer` in the application. But first, what is a trace? And, how exactly should you build them for you application?

To back up a bit, a trace is a type of telemetry that represents work being done by a service. In a distributed system, a trace can be thought of as a 'stack trace', showing the work being done by each service as well as the upstream and downstream calls that its making to other services.

Each part of the work that a service performs is represented in the trace with a span. Those spans are not just an unordered collection, but are defined in a hierarchical relationship with each other. If that doesn't make sense now, don't worry. You will have a better understanding after we instrument the code, so let's get started.

Start by instrumenting the `Run` method.

```go
// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		var span trace.Span
		ctx, span = otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(ctx)
		if err != nil {
			span.End()
			return err
		}

		a.Write(ctx, n)
		span.End()
	}
}
```

The above code creates a trace every iteration of the for loop using a `Tracer` from the global `TracerProvider`. You will learn more about `TracerProvider`s and handle the other side of setting up a global `TracerProvider` when you install an SDK in a later section. For now, as an instrumentation author, all you need to worry about is that you are using an appropriately named `Tracer` from a `TracerProvider`.

Next, instrument the `Poll` method.

```go
// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d", &n)

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	// This is going to be a high cardinality attribute because the user can
	// pass an unbounded number of different values. This type of attribute
	// should be avoided in instrumentation for high performance systems.
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, err
}
```

Similar to the `Run` method instrumentation, this adds a span to the method to track the computation done there. However, it also adds an attribute to annotate the span. This annotation is something you can add when you think a user of your application will want to see the state or details about the run environment when looking at telemetry.

Finally, instrument the `Write` method.

```go
// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span = otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(n)
	}(ctx)
	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}
```

This method is instrumented with two spans. One to track the `Write` method itself, and another to track the call to the core logic with the `Fibonacci` function.

Now that you have instrumented code it should be clearer how spans are related with a hierarchy. In OpenTelemetry Go the span hierarchy is defined explicitly with a `context.Context`. These contexts can contain references to spans. When a span is created a context is passed and reference to the created span is stored in a new context also returned with the span. If that returned context is used when creating another span, the original span will become that span's parent. This hierarchy gives traces structure and can help identify how a system works. Based on what you instrumented above and this understanding of span hierarchy you should expect a trace for each execution of the run loop to look like this.

```
Run
├── Poll
└── Write
    └── Fibonacci
```

A `Run` span will be a parent to both a `Poll` and `Write` span, and the `Write` span will be a parent to a `Fibonacci` span.

Now how do you actually see the produced spans? To do this you will need to configure and install an SDK.

# SDK Installation

OpenTelemetry is designed to be modular in its implementation of the OpenTelemetry API. The OpenTelemetry Go project offers an SDK package, `go.opentelemetry.io/otel/sdk`, that implements this API and adheres to the OpenTelemetry specification. To start using this SDK you will first need to create an exporter, but before anything can happen we need to install some things. Run the following in the `fib` directory to install the trace STDOUT exporter and the SDK.

```sh
$ go get go.opentelemetry.io/otel/sdk@v1.0.0-RC1 \
         go.opentelemetry.io/otel/exporters/stdout/stdouttrace@v1.0.0-RC1
```

Now add the needed imports to `main.go`.

```go
import (
	"context"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)
```

## Creating a Console Exporter

The SDK connects telemetry from the OpenTelemetry API to exporters. Exporters are packages that allow telemetry data to be emitted somewhere - either to the console (which is what we're doing here), or to a remote system or collector for further analysis and/or enrichment. OpenTelemetry supports a variety of exporters through its ecosystem including popular open source tools like Jaeger, Zipkin, and Prometheus.

To initialize the console exporter, add the following function to the `main.go` file:

```go
// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
```

This creates a new console exporter with basic options. You will use this function later when you configure the SDK to send telemetry data to it, but first you need to make sure that data is identifiable.

## Creating a Resource

Telemetry data can be crucial to solving issues with a service. The catch is, you need a way to identify what service, or even what service instance, that data is coming from. OpenTelemetry uses a `Resource` to represent the entity producing telemetry. Add the following function to the `main.go` file to create an appropriate `Resource` for the application.

```go
// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("fib"),
		semconv.ServiceVersionKey.String("v0.1.0"),
		attribute.String("environment", "demo"),
	)
}
```

Any information you would like to associate with all telemetry data the SDK handles can be added to the returned `Resource`. This is done by registering the `Resource` with the `TracerProvider`. Something you can now create!

## Installing a Tracer Provider

You have your application instrumented to produce telemetry data and you have an exporter to send that data to the console, but how are they connected? This is where the `TracerProvider` is used. It is a centralized point where instrumentation will get a `Tracer` from and configures these delegated `Tracer`s where and how to send the data they produce.

The pipelines that receive and ultimately transmit data to exporters are called `SpanProcessor`s. A `TracerProvider` can be configured to have multiple span processors, but for this example you will configure one that sends to data to your exporter in batches. Update your `main` function in `main.go` with the following.

```go
func main() {
	l := log.New(os.Stdout, "", 0)

	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

    /* … */
}
```

There's a fair amount going on here. First you are creating an console exporter that will export to a file. With the exporter ready, it can be registered in a new `TracerProvider`. This is done with a `BatchSpanProcessor` when it is passed to the `trace.WithBatcher` option. Batching data is a good practice and will help not overload systems downstream. Finally, with the `TracerProvider` created, you are deferring a function to flush and stop it, and registering it as the global OpenTelemetry `TracerProvider`.

Do you remember in the previous instrumentation section when we used the global `TracerProvider` to get a `Tracer`? This last step, registering the `TracerProvider` globally, is what will connect that instrumentation's `Tracer` with this `TracerProvider`. This pattern, using a global `TracerProvider`, is convenient, but not always appropriate. `TracerProvider`s can be explicitly passed to instrumentation or inferred from a context that contains a span. For this simple example using a global provider makes sense, but for more complex or distributed codebases these other ways of passing `TracerProvider`s may make more sense.

# Putting It All Together

You should have a working application that produces trace telemetry data! Give it a try.

```sh
$ go run .
What Fibonacci number would you like to know:
42
Fibonacci(42) = 267914296
What Fibonacci number would you like to know:
^C
goodbye
```

A new file named `traces.txt` should be created in your working directory. All the traces created from running your application should be in there!

# (Bonus) Errors

At this point you have a working application that is producing tracing telemetry data. Unfortunately, it was discovered there is an error in the core functionality of the `fib` module.

```sh
$ go run .
What Fibonacci number would you like to know:
100
Fibonacci(100) = 3736710778780434371
# …
```

But the 100-th Fibonacci number is `354224848179261915075`, not `3736710778780434371`! This application is only meant as a demo, but it shouldn't return wrong values. Update the `Fibonacci` function to return an error instead of computing incorrect values.

```go
// Fibonacci returns the n-th fibonacci number. An error is returned if the
// fibonacci number cannot be represented as a uint64.
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	if n > 93 {
		return 0, fmt.Errorf("unsupported fibonacci number %d: too large", n)
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}
```

Great, you have fixed the code, but it would be ideal to include errors returned to a user in the telemetry data. Luckily, spans can be configured to communicate this information. Update the `Write` method in `app.go` with the following code.

```go
// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span = otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		f, err := Fibonacci(n)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return f, err
	}(ctx)
    /* … */
}
```

The `Poll` method can be updated as well with a similar fix. If the user gave bad data before the application would fail but the telemetry data did not reflect this failure. Now you can add the following updates to the `Poll` method and this error will be captured in the data.

```go
// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d", &n)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
    /* … */
}
```

All that is left is updating imports for the `app.go` file to include the `go.opentelemetry.io/otel/codes` package.

```go
import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)
```

With these fixes in place and the instrumentation updated, re-trigger the bug.

```sh
$ go run .
What Fibonacci number would you like to know:
100
Fibonacci(100): unsupported fibonacci number 100: too large
What Fibonacci number would you like to know:
^C
goodbye
```

Excellent! The application no longer returns wrong values, and looking at the telemetry data in the `traces.txt` file you should see the error captured as an event.

```
"Events": [
	{
		"Name": "exception",
		"Attributes": [
			{
				"Key": "exception.type",
				"Value": {
					"Type": "STRING",
					"Value": "*errors.errorString"
				}
			},
			{
				"Key": "exception.message",
				"Value": {
					"Type": "STRING",
					"Value": "unsupported fibonacci number 100: too large"
				}
			}
		],
        ...
    }
]
```
