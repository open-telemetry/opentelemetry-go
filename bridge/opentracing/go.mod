module go.opentelemetry.io/otel/bridge/opentracing

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/label => ../../label
)

require (
	github.com/opentracing/opentracing-go v1.2.0
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/label v0.1.0
)
