module go.opentelemetry.io/otel/bridge/opentracing

go 1.17

replace go.opentelemetry.io/otel => ../..

require (
	github.com/opentracing/opentracing-go v1.2.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace
