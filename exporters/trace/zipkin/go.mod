module go.opentelemetry.io/otel/exporters/trace/zipkin

go 1.15

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require (
	github.com/google/go-cmp v0.5.6
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
)

replace go.opentelemetry.io/otel/bridge/opencensus => ../../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opentracing => ../../../bridge/opentracing

replace go.opentelemetry.io/otel/example/jaeger => ../../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../../example/namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../../../example/opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../../../example/otel-collector

replace go.opentelemetry.io/otel/example/prom-collector => ../../../example/prom-collector

replace go.opentelemetry.io/otel/example/prometheus => ../../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../../example/zipkin

replace go.opentelemetry.io/otel/exporters/metric/prometheus => ../../metric/prometheus

replace go.opentelemetry.io/otel/exporters/stdout => ../../stdout

replace go.opentelemetry.io/otel/exporters/trace/jaeger => ../jaeger

replace go.opentelemetry.io/otel/exporters/trace/zipkin => ./

replace go.opentelemetry.io/otel/internal/tools => ../../../internal/tools

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/oteltest => ../../../oteltest

replace go.opentelemetry.io/otel/sdk/export/metric => ../../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/example/passthrough => ../../../example/passthrough

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../../otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../../otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../../otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/internal/metric => ../../../internal/metric

replace go.opentelemetry.io/otel/exporters/jaeger => ../../jaeger

replace go.opentelemetry.io/otel/exporters/prometheus => ../../prometheus

replace go.opentelemetry.io/otel/exporters/zipkin => ../../zipkin

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric => ../../otlp/otlpmetric

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ../../otlp/otlpmetric/otlpmetricgrpc
