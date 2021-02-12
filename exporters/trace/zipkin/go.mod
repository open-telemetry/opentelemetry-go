module go.opentelemetry.io/otel/exporters/trace/zipkin

go 1.14

require (
	github.com/google/go-cmp v0.5.4
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.opentelemetry.io/otel/trace v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/trace => ../../../trace
