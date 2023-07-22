# Prometheus Bridge

Status: Experimental

The Prometheus Bridge allows using the Prometheus Golang client library
(github.com/prometheus/client_golang) with the OpenTelemetry SDK.

## Usage

```golang
// Make a Periodic Reader to periodically gather metrics from the Prometheus
// client library, and push to an OpenTelemetry exporter.
reader := metric.NewPeriodicReader(otelExporter)
// Register the Prometheus bridge to add metrics from the Prometheus
// DefaultGatherer to the output. Add the WithGatherer(registry) option to add
// your own registries.
reader.RegisterProducer(prombridge.NewMetricProducer())
// Create an OTel MeterProvider with our reader. Metrics from OpenTelemetry
// instruments are combined with metrics from Prometheus instruments in
// exported batches of metrics.
mp := metric.NewMeterProvider(metric.WithReader(reader))
```

## Limitations

* Summary metrics are dropped by the bridge.
* Start times for histograms and counters are set to the process start time.
