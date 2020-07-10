module go.opentelemetry.io/otel/exporters/trace/zipkin

go 1.13

replace go.opentelemetry.io/otel => ../../..

require (
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.8.0
	google.golang.org/grpc v1.30.0
)
