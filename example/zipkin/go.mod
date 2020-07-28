module go.opentelemetry.go/otel/example/zipkin

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/zipkin => ../../exporters/trace/zipkin
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.9.0
	go.opentelemetry.io/otel/sdk v0.9.0
)
