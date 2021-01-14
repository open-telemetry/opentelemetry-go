module go.opentelemetry.io/otel/example/prom-collector

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/metric/prometheus => ../../exporters/metric/prometheus
	go.opentelemetry.io/otel/exporters/otlp => ../../exporters/otlp
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.16.0
	go.opentelemetry.io/otel/exporters/otlp v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	google.golang.org/grpc v1.35.0
)
