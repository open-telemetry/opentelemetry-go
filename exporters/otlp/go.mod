module go.opentelemetry.io/otel/exporters/otlp

replace go.opentelemetry.io/otel => ../..

require (
	github.com/golang/protobuf v1.3.4
	github.com/google/go-cmp v0.4.0
	github.com/open-telemetry/opentelemetry-proto v0.0.0-20200313210948-2e3afbfffa38
	github.com/stretchr/testify v1.4.0
	go.opentelemetry.io/otel v0.3.0
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.27.1
)

go 1.13
