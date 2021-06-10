module go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc

go 1.15

require (
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/proto/otlp v0.9.0
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace go.opentelemetry.io/otel => ../../../..

replace go.opentelemetry.io/otel/sdk => ../../../../sdk

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../

replace go.opentelemetry.io/otel/metric => ../../../../metric

replace go.opentelemetry.io/otel/oteltest => ../../../../oteltest

replace go.opentelemetry.io/otel/trace => ../../../../trace

replace go.opentelemetry.io/otel/bridge/opencensus => ../../../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opentracing => ../../../../bridge/opentracing

replace go.opentelemetry.io/otel/example/jaeger => ../../../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../../../example/namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../../../../example/opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../../../../example/otel-collector

replace go.opentelemetry.io/otel/example/prom-collector => ../../../../example/prom-collector

replace go.opentelemetry.io/otel/example/prometheus => ../../../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../../../example/zipkin

replace go.opentelemetry.io/otel/exporters/metric/prometheus => ../../../metric/prometheus

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ./

replace go.opentelemetry.io/otel/exporters/stdout => ../../../stdout

replace go.opentelemetry.io/otel/exporters/trace/jaeger => ../../../trace/jaeger

replace go.opentelemetry.io/otel/exporters/trace/zipkin => ../../../trace/zipkin

replace go.opentelemetry.io/otel/internal/tools => ../../../../internal/tools

replace go.opentelemetry.io/otel/sdk/export/metric => ../../../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../../../sdk/metric

replace go.opentelemetry.io/otel/example/passthrough => ../../../../example/passthrough

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../otlptracehttp

replace go.opentelemetry.io/otel/internal/metric => ../../../../internal/metric
