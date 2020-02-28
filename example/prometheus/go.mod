module go.opentelemetry.io/otel/example/prometheus

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/metric/prometheus => ../../exporters/metric/prometheus
)

require (
	go.opentelemetry.io/otel v0.2.1
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.2.1
)
