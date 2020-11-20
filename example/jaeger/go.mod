module go.opentelemetry.io/otel/example/jaeger

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/jaeger => ../../exporters/trace/jaeger
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
)
