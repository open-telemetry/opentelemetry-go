module go.opentelemetry.io/otel/example/grpc

go 1.13

replace go.opentelemetry.io/otel => ../..

require (
	github.com/golang/protobuf v1.3.2
	go.opentelemetry.io/otel v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.24.0
)
