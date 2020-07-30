module go.opentelemetry.io/otel/example/http

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/stdout => ../../exporters/stdout
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/stdout v0.9.0
	go.opentelemetry.io/otel/sdk v0.9.0
)
