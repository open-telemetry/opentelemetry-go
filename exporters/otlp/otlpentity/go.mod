module go.opentelemetry.io/otel/exporters/otlp/otlpentity

go 1.20

require (
	go.opentelemetry.io/otel v1.23.1
	go.opentelemetry.io/otel/sdk v1.23.1
	go.opentelemetry.io/proto/otlp v1.1.0
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.23.1 // indirect
	go.opentelemetry.io/otel/trace v1.23.0-rc.1 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/proto/otlp => ../../../../opentelemetry-proto-go/otlp/
