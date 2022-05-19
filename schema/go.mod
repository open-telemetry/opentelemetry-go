module go.opentelemetry.io/otel/schema

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/stretchr/testify v1.7.1
	gopkg.in/yaml.v2 v2.4.0
)

replace go.opentelemetry.io/otel/bridge/opentracing => ../bridge/opentracing

replace go.opentelemetry.io/otel/example/fib => ../example/fib

replace go.opentelemetry.io/otel/example/jaeger => ../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../example/namedtracer

replace go.opentelemetry.io/otel/example/otel-collector => ../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../example/passthrough

replace go.opentelemetry.io/otel/example/zipkin => ../example/zipkin

replace go.opentelemetry.io/otel/exporters/jaeger => ../exporters/jaeger

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../exporters/otlp/internal/retry

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../exporters/otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../exporters/otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../exporters/otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/zipkin => ../exporters/zipkin

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/internal/tools => ../internal/tools

replace go.opentelemetry.io/otel/metric => ../metric

replace go.opentelemetry.io/otel/schema => ./

replace go.opentelemetry.io/otel/sdk => ../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../sdk/metric

replace go.opentelemetry.io/otel/trace => ../trace
