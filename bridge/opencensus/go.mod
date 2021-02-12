module go.opentelemetry.io/otel/bridge/opencensus

go 1.14

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/oteltest v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/trace v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/bridge/opencensus => ./

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/oteltest => ../../oteltest

replace go.opentelemetry.io/otel/trace => ../../trace
