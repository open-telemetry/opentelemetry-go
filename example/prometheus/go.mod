module go.opentelemetry.io/otel/example/prometheus

go 1.14

require (
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.16.0
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/metric/prometheus => ../../exporters/metric/prometheus

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/sdk/export/metric => ../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace
