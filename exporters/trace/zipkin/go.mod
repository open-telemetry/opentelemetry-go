module go.opentelemetry.io/otel/exporters/trace/zipkin

go 1.14

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require (
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.10.0
	go.opentelemetry.io/otel/sdk v0.10.0
	google.golang.org/grpc v1.30.0
)
