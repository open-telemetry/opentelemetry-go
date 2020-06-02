module go.opentelemetry.io/otel/example/otel-test

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/otlp => ../../exporters/otlp
)

require (
	go.opentelemetry.io/otel v0.6.0
	go.opentelemetry.io/otel/exporters/otlp v0.6.0
)
