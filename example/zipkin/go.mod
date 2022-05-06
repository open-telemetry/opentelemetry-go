module go.opentelemetry.io/otel/example/zipkin

go 1.16

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/zipkin => ../../exporters/zipkin
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/zipkin v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel/trace => ../../trace
