module go.opentelemetry.io/otel/exporters/trace/zipkin

go 1.14

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require (
	github.com/google/go-cmp v0.5.3
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
)
