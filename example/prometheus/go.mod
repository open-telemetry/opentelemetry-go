module go.opentelemetry.io/otel/example/prometheus

go 1.15

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/metric/prometheus => ../../exporters/metric/prometheus
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.20.0
	go.opentelemetry.io/otel/metric v0.20.0
	go.opentelemetry.io/otel/sdk/export/metric v0.20.0
	go.opentelemetry.io/otel/sdk/metric v0.20.0
)

replace go.opentelemetry.io/otel/bridge/opencensus => ../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opentracing => ../../bridge/opentracing

replace go.opentelemetry.io/otel/example/jaeger => ../jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../otel-collector

replace go.opentelemetry.io/otel/example/prom-collector => ../prom-collector

replace go.opentelemetry.io/otel/example/prometheus => ./

replace go.opentelemetry.io/otel/example/zipkin => ../zipkin

replace go.opentelemetry.io/otel/exporters/stdout => ../../exporters/stdout

replace go.opentelemetry.io/otel/exporters/trace/jaeger => ../../exporters/trace/jaeger

replace go.opentelemetry.io/otel/exporters/trace/zipkin => ../../exporters/trace/zipkin

replace go.opentelemetry.io/otel/internal/tools => ../../internal/tools

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/oteltest => ../../oteltest

replace go.opentelemetry.io/otel/sdk/export/metric => ../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/example/passthrough => ../passthrough

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../../exporters/otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../../exporters/otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../../exporters/otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/internal/metric => ../../internal/metric
