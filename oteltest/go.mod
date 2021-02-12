module go.opentelemetry.io/otel/oteltest

go 1.14

require (
	go.opentelemetry.io/otel v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/trace v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/metric => ../metric

replace go.opentelemetry.io/otel/oteltest => ./

replace go.opentelemetry.io/otel/trace => ../trace
