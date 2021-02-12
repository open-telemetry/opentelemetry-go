module go.opentelemetry.io/otel/exporters/trace/jaeger

go 1.14

require (
	github.com/apache/thrift v0.13.0
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.opentelemetry.io/otel/trace v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.39.0
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/exporters/trace/jaeger => ./

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/trace => ../../../trace
