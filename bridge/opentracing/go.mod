module go.opentelemetry.io/otel/bridge/opentracing

go 1.14

require (
	github.com/opentracing/opentracing-go v1.2.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/trace v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/bridge/opentracing => ./

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/trace => ../../trace
