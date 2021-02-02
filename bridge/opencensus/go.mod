module go.opentelemetry.io/otel/bridge/opencensus

go 1.14

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/label v0.16.0
	go.opentelemetry.io/otel/trace v0.16.0
)

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/label => ../../label
	go.opentelemetry.io/otel/trace => ../../trace
	go.opentelemetry.io/otel/baggage => ../../baggage
)
