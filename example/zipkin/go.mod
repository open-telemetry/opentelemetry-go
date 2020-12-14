module go.opentelemetry.io/otel/example/zipkin

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/zipkin => ../../exporters/trace/zipkin
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
)
