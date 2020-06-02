module go.opentelemetry.io/otel/exporters/otlp

replace go.opentelemetry.io/otel => ../..

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/google/go-cmp v0.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3 // indirect
	github.com/open-telemetry/opentelemetry-proto v0.3.0
	github.com/stretchr/testify v1.4.0
	go.opentelemetry.io/otel v0.6.0
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.27.1
)

go 1.13
