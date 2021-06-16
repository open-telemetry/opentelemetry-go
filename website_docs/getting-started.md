---
title: "Getting Started"
weight: 2
---

Welcome to the OpenTelemetry for Go getting started guide! This guide will walk you the basic steps in installing, configuring, and exporting data from OpenTelemetry.

# Installation

OpenTelemetry packages for Go are available in the `go.opentelemetry.io/otel` namespace. You will need to add references to them in the `import` statement. We suggest using Go 1.15 or newer, for module support.

To get started with this guide, create a new directory and add a new file named `main.go` to it. In your terminal, run the command `go mod init main` in the same directory. This will create a `go.mod` file, which is used by Go to manage imports.

# Initialization and Configuration

To install the necessary prerequisites for OpenTelemetry, you'll want to run the following command in the directory with your `go.mod`:

`go get go.opentelemetry.io/otel@v0.20.0 go.opentelemetry.io/otel/sdk@v0.20.0 go.opentelemetry.io/otel/exporters/stdout@v0.20.0`

In your `main.go` file, you'll need to import several packages:

```go
package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)
```

These packages contain the basic requirements for OpenTelemetry Go - the API itself, the metrics and tracing SDK, and context propagation. The exact libraries and packages that you'll use in an application will vary depending on what features you need - for example, if you're writing a library that will be used by others, you don't need to require the SDK packages and will rely solely on the API. In general, you should configure the SDK in your code as close to program initialization as possible in order to capture telemetry at the earliest time it's available.

## Creating a Console Exporter

The SDK requires an exporter to be created. Exporters are packages that allow telemetry data to be emitted somewhere - either to the console (which is what we're doing here), or to a remote system or collector for further analysis and/or enrichment. OpenTelemetry supports a variety of exporters through its ecosystem including popular open source tools like Jaeger, Zipkin, and Prometheus.

To initialize the console exporter, add the following code to the file your `main.go` file -

```go
func main() {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace export pipeline: %v", err)
	}
```

This creates a new console exporter with basic options - `WithPrettyPrint` formats the text nicely when its printed, so that it's easier for humans to read.

## Creating a Tracer Provider

A trace is a type of telemetry that represents work being done by a service. In a distributed system, a trace can be thought of as a 'stack trace', showing the work being done by each service as well as the upstream and downstream calls that its making to other services.

OpenTelemetry requires a trace provider to be initialized in order to generate traces. A trace provider can have multiple span processors, which are components that allow for span data to be modified or exported after it's created.

To create a trace provider, add the following code to your `main.go` file -

```go
	ctx := context.Background()
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()
```

This block of code will create a new batch span processor, a type of span processor that batches up multiple spans over a period of time, that writes to the exporter we created in the previous step. You can see examples of other uses for span processors in [this file](https://github.com/open-telemetry/opentelemetry-go/blob/v0.16.0/sdk/trace/span_processor_example_test.go). We also created an instance of a Go context. It will be used later to store some important data.

## Creating a Meter Provider

A metric is a captured measurement about the execution of a computer program at run time. Examples of metrics can be "count the number of requests completed", "count the number of active requests", "capture a queue length" or "capture the number of cache misses".

OpenTelemetry requires a meter provider to be initialized in order to create instruments that will generate metrics. The way metrics are exported depends on the used system. For example, prometheus uses a pull model, while OTLP uses a push model. In this document we use an stdout exporter which uses the latter. Thus we need to create a push controller that will periodically push the collected metrics to the exporter.

To create a meter provider, add the following code to your `main.go` file -

```go
	metricExporter, err := stdoutmetric.New(
		stdoutmetric.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to initialize stdoutmetric export pipeline: %v", err)
	}

	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(5*time.Second),
	)

	err = pusher.Start(ctx)
	if err != nil {
		log.Fatalf("failed to initialize metric controller: %v", err)
	}

	// Handle this error in a sensible manner where possible
	defer func() { _ = pusher.Stop(ctx) }()
```

Again we create an exporter, this time using the `stdoutmetric` exporter package. Then we create a controller that uses a basic processor to aggregate and process metrics that are then sent to the exporter. The basic processor here uses a simple aggregator selector that decides what kind of an aggregator to use to aggregate measurements from a specific instrument. The processor also uses the exporter to learn how to prepare the aggregated measurements for the exporter to consume. The controller will periodically push aggregated measurements to the exporter.

## Setting Global Options

When using OpenTelemetry, it's a good practice to set a global tracer provider and a global meter provider. Doing so will make it easier for libraries and other dependencies that use the OpenTelemetry API to easily discover the SDK, and emit telemetry data. In addition, you'll want to configure context propagation options. Context propagation allows for OpenTelemetry to share values across multiple services - this includes trace identifiers, which ensure that all spans for a single request are part of the same trace, as well as baggage, which are arbitrary key/value pairs that you can use to pass observability data between services (for example, sharing a customer ID from one service to the next).

Setting up global options uses the `otel` package - add these options to your `main.go` file as shown -

```go
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(pusher.MeterProvider())
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
```

It's important to note that if you do not set a propagator, the default is to use the `NoOp` option, which means that context will not be shared between multiple services. To avoid that, we set up a composite propagator that consist of a baggage propagator and trace context propagator. That way, both trace information (trace IDs, span IDs, etc) and baggage will be propagated.

## Creating metric instruments

The next step is to create metric instruments that will capture measurements. There are two kinds of instruments: synchronous and asynchronous. Synchronous instruments capture measurements by explicitly calling the capture either by the application or by an instrumented library. Depending on the semantics of the measurements, we can say that synchronous instruments record or add measurements. Asynchronous instruments provide a callback that captures measurements. The callback is periodically called by meter in the background. We can say that asynchronous instrument performs observations.

Each measurement can be associated with attributes that can later be used by visualisation software to categorize and filter measurements. In case of synchronous instruments the attributes can be passed at the moment of capturing a measurement or can be passed when binding the instrument. Such a bound instrument can be later used to capture measurements without passing the attributes. In case of asynchronous instruments, the attributes are passed each time an observation is made explicitly in the callback.

To set up some metric instruments, add the following code to your `main.go` file -

```go
	fooKey := attribute.Key("ex.com/foo")
	barKey := attribute.Key("ex.com/bar")
	lemonsKey := attribute.Key("ex.com/lemons")
	anotherKey := attribute.Key("ex.com/another")

	commonAttributes := []attribute.KeyValue{lemonsKey.Int(10), attribute.String("A", "1"), attribute.String("B", "2"), attribute.String("C", "3")}

	meter := otel.Meter("ex.com/basic")

	observerCallback := func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, commonAttributes...)
	}
	_ = metric.Must(meter).NewFloat64ValueObserver("ex.com.one", observerCallback,
		metric.WithDescription("A ValueObserver set to 1.0"),
	)

	valueRecorder := metric.Must(meter).NewFloat64ValueRecorder("ex.com.two")

	boundRecorder := valueRecorder.Bind(commonAttributes...)
	defer boundRecorder.Unbind()
```

In this block we first create some keys and attributes that we will later use when capturing the measurements. Then we ask a global meter provider to give us a named meter instance ("ex.com/basic"). This acts as a way to namespace our instruments and make them distinct from other instruments in this process or another. Then we use the meter to create two instruments - an asynchronous value observer and a synchronous value recorder.

# Quick Start

Let's put the concepts we've just covered together, and create a trace and some measurements in a single process. In our main function, after the initialization code, add the following:

```go
	tracer := otel.Tracer("ex.com/basic")
	ctx = baggage.ContextWithValues(ctx,
		fooKey.String("foo1"),
		barKey.String("bar1"),
	)

	func(ctx context.Context) {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "operation")
		defer span.End()

		span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("bogons", 100)))
		span.SetAttributes(anotherKey.String("yes"))

		meter.RecordBatch(
			// Note: call-site variables added as context Entries:
			baggage.ContextWithValues(ctx, anotherKey.String("xyz")),
			commonAttributes,

			valueRecorder.Measurement(2.0),
		)

		func(ctx context.Context) {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "Sub operation...")
			defer span.End()

			span.SetAttributes(lemonsKey.String("five"))
			span.AddEvent("Sub span event")
			boundRecorder.Record(ctx, 1.3)
		}(ctx)
	}(ctx)
}
```

In this snippet, we're doing a few things. First, we're asking the global trace provider for an instance of a tracer, which is the object that manages spans for our service. We provide a name (`"ex.com/basic"`) too, which acts in the same way as a name we gave to our meter instance. Here we can also see the use of the Go context - it contains baggage items that are propagated to other places in our code and to other processes. Which means that baggage items should be used within limits as baggage may be sent over the network. The other use of the Go context is to store a reference to a span, so it can be propagated between function calls and processes.

Inside our function, we're creating a new span by calling `tracer.Start` with the context we just created, and a name. Passing the context will set our span as 'active' in it, which is used in our inner function to make a new child span. The name is important - every span needs a name, and these names are the primary method of indicating what a span represents.  Calling `defer span.End()` ensures that our span will complete once this function has finished its work. Spans can have attributes and events, which are metadata and log statements that help you interpret traces after-the-fact. Finally, in this code snippet we can see an example of creating a new function and propagating the span to it inside our code. When you run this program, you'll see that the 'Sub operation...' span has been created as a child of the 'operation' span.

We also record some measurements. Recording measurements with asynchronous instruments is controlled by SDK and the controller we use, so we do not need to do anything else after creating the instrument and passing the callback to it. For synchronous instruments there are two ways of recording measurements - either through the instrument, bounded or not (in our case it's a value recorder, so we use the `Record` function), or by making a batched measurement (with `meter.RecordBatch`). Batched measurements allow you to use multiple instruments to create measurement and record them once.
