module go.opentelemetry.io/otel/example/grpc

go 1.13

replace go.opentelemetry.io/otel => ../..

require (
	github.com/golang/protobuf v1.4.2
	go.opentelemetry.io/otel v0.7.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.30.0
)
