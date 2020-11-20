module go.opentelemetry.io/otel/example/zipkin

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/zipkin => ../../exporters/trace/zipkin
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
)
