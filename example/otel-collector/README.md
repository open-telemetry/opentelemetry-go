# OpenTelemetry Collector Traces Example

This example illustrates how to export traces from the OpenTelemetry-Go SDK to the OpenTelemetry Collector, and from there to a Jaeger instance.
The complete flow is:

`otel-collector-demo -> otel-collector -> Jaeger`

# Prerequisites

The demo assumes you already have both an instance of the OpenTelemetry Collector and Jaeger up and running.
For information about how to set these up, please follow the corresponding documentation:

* Jaeger: https://www.jaegertracing.io/docs/1.17/getting-started/
* OpenTelemetry Collector: https://opentelemetry.io/docs/collector/about/

Moreover, this demo is built on a Kubernetes cluster running both the OpenTelemetry Collector, and Jaeger.
Therefore, the included OpenTelemetry Collector configuration is a Kubernetes `ConfigMap`.
However, the main concepts should be transferable to any other platform you want to run this demo on.

# Configuring the OTEL Collector

In order to enable our application to send traces to the OpenTelemetry Collector, we need to first configure the OTLP receiver:

```yml
receivers:
      otlp: {}
```

This will create the receiver on the Collector side, and open up port `55680` for receiving traces.
The rest of the configuration is quite standard, with the only mention that we need to create the Jaeger exporter:

```yml
exporters:
      jaeger:
        # Replace with a real endpoint.
        endpoint: "<jaeger-service-url>:14250"
```

After this, apply the configuration to your OpenTelemetry Collector instance (with `kubectl apply -f otel-controller-config.yaml` for k8s users).
You should see that the Collector creates the otlp receiver:

```json
{"level":"info","ts":1589184143.206609,"caller":"builder/receivers_builder.go:79","msg":"Receiver started.","component_kind":"receiver","component_type":"otlp","component_name":"otlp"}
```
and the Jaeger exporter:
```json
{"level":"info","ts":1589184143.1535392,"caller":"builder/exporters_builder.go:94","msg":"Exporter started.","component_kind":"exporter","component_type":"jaeger","component_name":"jaeger"}
```

# Writing the demo

Having the OpenTelemetry Collector started with the OTLP port open for traces, and connected to Jaeger, let's look at the go app that will send traces to the Collector.

First, we need to create an exporter using the otlp package:
```go
exp, _ := otlp.NewExporter(otlp.WithInsecure(),
        // For debug purposes
        otlp.WithGRPCDialOption(grpc.WithBlock()))
defer func() {
		_ = exp.Stop()
	}()
```
This will initialize the exporter with the default configuration.
In this configuration, it will try to connect to an OTLP receiver at the address `localhost:55680`.
If your OpenTelemetry Collector is running at a different address, use the [`otlp.WithAddress`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#WithAddress) function to change the default address.

Feel free to remove the blocking operation, but it might come in handy when testing the connection. Also, make sure to close the exporter before the app exits.

The next steps are the same as for all other otel-go sdk uses:
1) Create a trace provider from the `otlp` exporter: 
```go
tp, _ := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exp, // add following two options to ensure flush
			sdktrace.WithScheduleDelayMillis(5),
			sdktrace.WithMaxExportBatchSize(2),
        ))
```

2) Start sending traces:
```go
tracer := tp.Tracer("test-tracer")
ctx, span := tracer.Start(context.Background(), "CollectorExporter-Example")
	defer span.End()
```

The traces should now be visible from the Jaeger UI (if you have it installed).

# Notes

* There is an issue with the exporter/Collector which causes Jaeger to throw errors when receiving spans from the OpenTelemetry Collector: https://github.com/open-telemetry/opentelemetry-collector/issues/815
