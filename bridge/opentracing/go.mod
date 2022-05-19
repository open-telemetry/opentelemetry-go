module go.opentelemetry.io/otel/bridge/opentracing

go 1.16

replace go.opentelemetry.io/otel => ../..

require (
	github.com/opentracing/opentracing-go v1.2.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel/trace => ../../trace
