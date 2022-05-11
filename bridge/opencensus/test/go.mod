module go.opentelemetry.io/otel/bridge/opencensus/test

go 1.16

require (
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/bridge/opencensus v0.30.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/bridge/opencensus => ../

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace
