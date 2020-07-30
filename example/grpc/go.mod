module go.opentelemetry.io/otel/example/grpc

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/stdout => ../../exporters/stdout
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	github.com/golang/protobuf v1.4.2
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/stdout v0.9.0
	go.opentelemetry.io/otel/sdk v0.9.0
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	google.golang.org/grpc v1.30.0
)
