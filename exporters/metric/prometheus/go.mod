module go.opentelemetry.io/otel/exporters/metric/prometheus

go 1.13

replace go.opentelemetry.io/otel => ../../..

require (
	github.com/prometheus/client_golang v1.6.0
	github.com/stretchr/testify v1.4.0
	go.opentelemetry.io/otel v0.6.0
)
