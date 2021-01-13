module go.opentelemetry.io/otel/exporters/otlp

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
	google.golang.org/grpc v1.34.0
)
