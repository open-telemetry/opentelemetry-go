module go.opentelemetry.io/otel/exporters/otlp

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.1
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/sdk v0.9.0
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.30.0
)
