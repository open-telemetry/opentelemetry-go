module go.opentelemetry.io/otel/example/otel-collector

go 1.14

require (
	github.com/open-telemetry/opentelemetry-collector v0.3.0
	go.opentelemetry.io/otel v0.6.0
	go.opentelemetry.io/otel/exporters/otlp v0.6.0
	google.golang.org/grpc v1.29.1
)

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/otlp => ../../exporters/otlp
)
