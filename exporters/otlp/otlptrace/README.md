# Otlptrace Exporter

To constructs a new Otlptrace Exporter and starts it:
```
exp, err := otlptrace.New(ctx, opts...)
```
# Otlptrace Client
Otlptrace package define a trace exporter that uses a `otlptrace.Client` .

`otlptrace.Client` manages connections to the collector, handles the transformation of data into wire format, 
and the transmission of that data to the collector.

# Otlptracegrpc and Otlptracehttp
The otlptracegrpc and otlptracehttp implements a gRPC `otlptrace.Client` and 
an HTTP `otlptrace.Client`respectively,
both offering convenience functions .


## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlptrace
```

