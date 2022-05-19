module go.opentelemetry.io/otel/exporters/jaeger

go 1.16

require (
	github.com/google/go-cmp v0.5.8
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/bridge/opentracing => ../../bridge/opentracing

replace go.opentelemetry.io/otel/example/fib => ../../example/fib

replace go.opentelemetry.io/otel/example/jaeger => ../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../example/namedtracer

replace go.opentelemetry.io/otel/example/otel-collector => ../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../example/passthrough

replace go.opentelemetry.io/otel/example/zipkin => ../../example/zipkin

replace go.opentelemetry.io/otel/exporters/jaeger => ./

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../otlp/internal/retry

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/zipkin => ../zipkin

replace go.opentelemetry.io/otel/internal/tools => ../../internal/tools

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/schema => ../../schema

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric
