module go.opentelemetry.io/otel/sdk/metric

go 1.16

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk => ../

replace go.opentelemetry.io/otel/trace => ../../trace

require (
	github.com/benbjohnson/clock v1.3.0
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/metric v0.30.0
	go.opentelemetry.io/otel/sdk v1.7.0
)
