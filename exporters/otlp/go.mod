module go.opentelemetry.io/otel/exporters/otlp

replace go.opentelemetry.io/otel => ../..

require (
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.8.0
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.30.0
)

go 1.13
