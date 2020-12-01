module go.opentelemetry.io/otel/exporters/otlp

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884 // indirect
	google.golang.org/grpc v1.32.0
)
