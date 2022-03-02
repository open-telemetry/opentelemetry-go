---
title: "Getting Started"
weight: 2
---

Welcome to the OpenTelemetry for Go getting started guide! This guide will walk you through the basic steps in installing, instrumenting with, configuring, and exporting data from OpenTelemetry. Before you get started, be sure to have Go 1.16 or newer installed.

Understand how a system is functioning when it is failing or having issues is critical to resolving those issues. One strategy to understand this is with tracing. This guide shows how the OpenTelemetry Go project can be used to trace an example application. You will start with an application that computes Fibonacci numbers for users, and from there you will add instrumentation to produce tracing telemetry with OpenTelemetry Go.

For reference, a complete example of the code you will build can be found [here](https://github.com/open-telemetry/opentelemetry-go/tree/main/example/fib).

To start building the application, make a new directory named `fib` to house our Fibonacci project. Next, add the following to a new file named `fib.go` in that directory.

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

With your core logic added, you can now build your application around it. Add a new `app.go` file with the following application logic.

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
)

// App is a Fibonacci computation application.
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
		n, err := a.Poll(ctx)
		if err != nil {
			return err
		}

		a.Write(ctx, n)
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
	"os/signal"
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

With the code complete it is almost time to run the application. Before you can do that you need to initialize this directory as a Go module. From your terminal, run the command `go mod init fib` in the `fib` directory. This will create a `go.mod` file, which is used by Go to manage imports. Now you should be able to run the application!

```sh
$ go run .
What Fibonacci number would you like to know:
42
Fibonacci(42) = 267914296
What Fibonacci number would you like to know:
^C
goodbye
```

The application can be exited with CTRL+C. You should see a similar output as above, if not make sure to go back and fix any errors.

# Trace Instrumentation

OpenTelemetry is split into two parts: an API to instrument code with, and SDKs that implement the API. To start integrating OpenTelemetry into any project, the API is used to define how telemetry is generated. To generate tracing telemetry in your application you will use the OpenTelemetry Trace API from the [`go.opentelemetry.io/otel/trace`] package.

First, you need to install the necessary packages for the Trace API. Run the following command in your working directory.

```sh
go get go.opentelemetry.io/otel \
       go.opentelemetry.io/otel/trace
```

Now that the packages installed you can start updating your application with imports you will use in the `app.go` file.

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

With the imports added, you can start instrumenting.

The OpenTelemetry Tracing API provides a [`Tracer`] to create traces. These [`Tracer`]s are designed to be associated with one instrumentation library. That way telemetry they produce can be understood to come from that part of a code base. To uniquely identify your application to the [`Tracer`] you will use create a constant with the package name in `app.go`.

```go
// name is the Tracer name used to identify this instrumentation library.
const name = "fib"
```

Using the full-qualified package name, something that should be unique for Go packages, is the standard way to identify a [`Tracer`]. If your example package name differs, be sure to update the name you use here to match.

Everything should be in place now to start tracing your application. But first, what is a trace? And, how exactly should you build them for you application?

To back up a bit, a trace is a type of telemetry that represents work being done by a service. A trace is a record of the connection(s) between participants processing a transaction, often through client/server requests processing and other forms of communication.

Each part of the work that a service performs is represented in the trace by a span. Those spans are not just an unordered collection. Like the call stack of our application, those spans are defined with relationships to one another. The "root" span is the only span without a parent, it represents how a service request is started. All other spans have a parent relationship to another span in the same trace.

If this last part about span relationships doesn't make complete sense now, don't worry. The most important takeaway is that each part of your code, which does some work, should be represented as a span. You will have a better understanding of these span relationships after you instrument your code, so let's get started.

Start by instrumenting the `Run` method.

```go
// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		// Each execution of the run loop, we should get a new "root" span and context.
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}

		a.Write(newCtx, n)
		span.End()
	}
}
```

The above code creates a span for every iteration of the for loop. The span is created using a [`Tracer`] from the [global `TracerProvider`](https://pkg.go.dev/go.opentelemetry.io/otel#GetTracerProvider). You will learn more about [`TracerProvider`]s and handle the other side of setting up a global [`TracerProvider`] when you install an SDK in a later section. For now, as an instrumentation author, all you need to worry about is that you are using an appropriately named [`Tracer`] from a [`TracerProvider`] when you write `otel.Tracer(name)`.

Next, instrument the `Poll` method.

```go
// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, err
}
```

Similar to the `Run` method instrumentation, this adds a span to the method to track the computation performed. However, it also adds an attribute to annotate the span. This annotation is something you can add when you think a user of your application will want to see the state or details about the run environment when looking at telemetry.

Finally, instrument the `Write` method.

```go
// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
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

This method is instrumented with two spans. One to track the `Write` method itself, and another to track the call to the core logic with the `Fibonacci` function. Do you see how context is passed through the spans? Do you see how this also defines the relationship between spans?

In OpenTelemetry Go the span relationships are defined explicitly with a `context.Context`. When a span is created a context is returned alongside the span. That context will contain a reference to the created span. If that context is used when creating another span the two spans will be related. The original span will become the new span's parent, and as a corollary, the new span is said to be a child of the original. This hierarchy gives traces structure, structure that helps show a computation path through a system. Based on what you instrumented above and this understanding of span relationships you should expect a trace for each execution of the run loop to look like this.

```
Run
├── Poll
└── Write
    └── Fibonacci
```

A `Run` span will be a parent to both a `Poll` and `Write` span, and the `Write` span will be a parent to a `Fibonacci` span.

Now how do you actually see the produced spans? To do this you will need to configure and install an SDK.

# SDK Installation

OpenTelemetry is designed to be modular in its implementation of the OpenTelemetry API. The OpenTelemetry Go project offers an SDK package, [`go.opentelemetry.io/otel/sdk`], that implements this API and adheres to the OpenTelemetry specification. To start using this SDK you will first need to create an exporter, but before anything can happen we need to install some packages. Run the following in the `fib` directory to install the trace STDOUT exporter and the SDK.

```sh
$ go get go.opentelemetry.io/otel/sdk \
         go.opentelemetry.io/otel/exporters/stdout/stdouttrace
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

The SDK connects telemetry from the OpenTelemetry API to exporters. Exporters are packages that allow telemetry data to be emitted somewhere - either to the console (which is what we're doing here), or to a remote system or collector for further analysis and/or enrichment. OpenTelemetry supports a variety of exporters through its ecosystem including popular open source tools like [Jaeger](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/jaeger), [Zipkin](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/zipkin), and [Prometheus](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus).

To initialize the console exporter, add the following function to the `main.go` file:

```go
// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}
```

This creates a new console exporter with basic options. You will use this function later when you configure the SDK to send telemetry data to it, but first you need to make sure that data is identifiable.

## Creating a Resource

Telemetry data can be crucial to solving issues with a service. The catch is, you need a way to identify what service, or even what service instance, that data is coming from. OpenTelemetry uses a [`Resource`] to represent the entity producing telemetry. Add the following function to the `main.go` file to create an appropriate [`Resource`] for the application.

```go
// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
```

Any information you would like to associate with all telemetry data the SDK handles can be added to the returned [`Resource`]. This is done by registering the [`Resource`] with the [`TracerProvider`]. Something you can now create!

## Installing a Tracer Provider

You have your application instrumented to produce telemetry data and you have an exporter to send that data to the console, but how are they connected? This is where the [`TracerProvider`] is used. It is a centralized point where instrumentation will get a [`Tracer`] from and funnels the telemetry data from these [`Tracer`]s to export pipelines.

The pipelines that receive and ultimately transmit data to exporters are called [`SpanProcessor`]s. A [`TracerProvider`] can be configured to have multiple span processors, but for this example you will only need to configure only one. Update your `main` function in `main.go` with the following.

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

There's a fair amount going on here. First you are creating a console exporter that will export to a file. You are then registering the exporter with a new [`TracerProvider`]. This is done with a [`BatchSpanProcessor`] when it is passed to the [`trace.WithBatcher`] option. Batching data is a good practice and will help not overload systems downstream. Finally, with the [`TracerProvider`] created, you are deferring a function to flush and stop it, and registering it as the global OpenTelemetry [`TracerProvider`].

Do you remember in the previous instrumentation section when we used the global [`TracerProvider`] to get a [`Tracer`]? This last step, registering the [`TracerProvider`] globally, is what will connect that instrumentation's [`Tracer`] with this [`TracerProvider`]. This pattern, using a global [`TracerProvider`], is convenient, but not always appropriate. [`TracerProvider`]s can be explicitly passed to instrumentation or inferred from a context that contains a span. For this simple example using a global provider makes sense, but for more complex or distributed codebases these other ways of passing [`TracerProvider`]s may make more sense.

# Putting It All Together

You should now have a working application that produces trace telemetry data! Give it a try.

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

At this point you have a working application and it is producing tracing telemetry data. Unfortunately, it was discovered that there is an error in the core functionality of the `fib` module.

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
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
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

With this change any error returned from the `Fibonacci` function will mark that span as an error and record an event describing the error.

This is a great start, but it is not the only error returned in from the application. If a user makes a request for a non unsigned integer value the application will fail. Update the `Poll` method with a similar fix to capture this error in the telemetry data.

```go
// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, nil
}
```

All that is left is updating imports for the `app.go` file to include the [`go.opentelemetry.io/otel/codes`] package.

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

# What's Next

This guide has walked you through adding tracing instrumentation to an
application and using a console exporter to send telemetry data to a file. There
are many other topics to cover in OpenTelemetry, but you should be ready to
start adding OpenTelemetry Go to your projects at this point. Go instrument your
code!

For more information about instrumenting your code and things you can do with
spans, refer to the [Instrumenting]({{< relref "manual" >}})
documentation. Likewise, advanced topics about processing and exporting
telemetry data can be found in the [Processing and Exporting Data]({{< relref
"exporting_data" >}}) documentation.

[`go.opentelemetry.io/otel/trace`]: https://pkg.go.dev/go.opentelemetry.io/otel/trace
[`go.opentelemetry.io/otel/sdk`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk
[`go.opentelemetry.io/otel/codes`]: https://pkg.go.dev/go.opentelemetry.io/otel/codes
[`Tracer`]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#Tracer
[`TracerProvider`]: https://pkg.go.dev/go.opentelemetry.io/otel/trace#TracerProvider
[`Resource`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/resource#Resource
[`SpanProcessor`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#SpanProcessor
[`BatchSpanProcessor`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#NewBatchSpanProcessor
[`trace.WithBatcher`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#WithBatcher
