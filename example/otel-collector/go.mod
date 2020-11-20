module go.opentelemetry.io/otel/example/otel-collector

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/otlp => ../../exporters/otlp
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/otlp v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
	google.golang.org/grpc v1.32.0
)
