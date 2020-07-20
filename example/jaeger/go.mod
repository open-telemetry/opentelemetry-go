module go.opentelemetry.io/otel/example/jaeger

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/jaeger => ../../exporters/trace/jaeger
)

require (
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.9.0
)
