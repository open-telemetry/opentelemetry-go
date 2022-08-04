module go.opentelemetry.io/otel/exporters/otlp/otlpmetric

go 1.18

require (
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/sdk/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/proto/otlp v0.18.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/sdk v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/sdk => ../../../sdk
