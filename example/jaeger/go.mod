module go.opentelemetry.io/otel/example/jaeger

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporter/trace/jaeger => ../../exporter/trace/jaeger
)

require (
	go.opentelemetry.io/otel v0.2.0
	go.opentelemetry.io/otel/exporter/trace/jaeger v0.2.0
)
