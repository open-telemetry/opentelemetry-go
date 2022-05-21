module go.opentelemetry.io/otel/sdk/metric

go 1.16

require (
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v1.7.0
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/sdk => ../
