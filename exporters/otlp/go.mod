module go.opentelemetry.io/otel/exporters/otlp

replace go.opentelemetry.io/otel => ../..

replace github.com/open-telemetry/opentelemetry-proto => ../../../github.com/open-telemetry/opentelemetry-proto

require (
	github.com/google/go-cmp v0.5.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/open-telemetry/opentelemetry-proto v0.4.0
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.7.0
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.30.0
)

go 1.13
