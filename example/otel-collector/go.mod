module go.opentelemetry.io/otel/example/otel-collector

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/baggage => ../../baggage
	go.opentelemetry.io/otel/exporters/otlp => ../../exporters/otlp
	go.opentelemetry.io/otel/label => ../../label
	go.opentelemetry.io/otel/sdk => ../../sdk
	go.opentelemetry.io/otel/semconv => ../../semconv
	go.opentelemetry.io/otel/trace => ../../trace
)

require (
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/otlp v0.16.0
	go.opentelemetry.io/otel/label v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.opentelemetry.io/otel/semconv v0.16.0
	go.opentelemetry.io/otel/trace v0.16.0
	google.golang.org/grpc v1.35.0
)
