module go.opentelemetry.io/otel/example/basic

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/stdout => ../../exporters/stdout
)

require (
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/stdout v0.9.0
)
