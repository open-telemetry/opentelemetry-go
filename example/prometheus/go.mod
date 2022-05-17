module go.opentelemetry.io/otel/example/prometheus

go 1.16

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/prometheus => ../../exporters/prometheus
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/prometheus v0.30.0
	go.opentelemetry.io/otel/metric v0.30.0
	go.opentelemetry.io/otel/sdk/metric v0.30.0
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace
