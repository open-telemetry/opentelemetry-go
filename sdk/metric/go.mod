module go.opentelemetry.io/otel/sdk/metric

go 1.14

require (
	github.com/benbjohnson/clock v1.0.3 // do not upgrade to v1.1.x because it would require Go >= 1.15
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk/export/metric v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk => ../

replace go.opentelemetry.io/otel/sdk/export/metric => ../export/metric

replace go.opentelemetry.io/otel/sdk/metric => ./

replace go.opentelemetry.io/otel/trace => ../../trace
