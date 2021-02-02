module go.opentelemetry.io/otel/bridge/opentracing

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/trace => ../../trace
)

require (
	github.com/opentracing/opentracing-go v1.2.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/trace v0.16.0
)
