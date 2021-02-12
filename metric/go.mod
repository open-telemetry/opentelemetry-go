module go.opentelemetry.io/otel/metric

go 1.14

require (
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/oteltest v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/metric => ./

replace go.opentelemetry.io/otel/oteltest => ../oteltest

replace go.opentelemetry.io/otel/trace => ../trace
