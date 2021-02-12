module go.opentelemetry.io/otel/exporters/otlp

go 1.14

require (
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v0.16.0
	go.opentelemetry.io/otel/sdk/export/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/trace v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.35.0
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/otlp => ./

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/sdk/export/metric => ../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace
