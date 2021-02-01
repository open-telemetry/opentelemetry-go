module go.opentelemetry.io/otel/example/zipkin

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/zipkin => ../../exporters/trace/zipkin
	go.opentelemetry.io/otel/label => ../../label
	go.opentelemetry.io/otel/sdk => ../../sdk
	go.opentelemetry.io/otel/semconv => ../../semconv
	go.opentelemetry.io/otel/trace => ../../trace
)

require (
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
)
