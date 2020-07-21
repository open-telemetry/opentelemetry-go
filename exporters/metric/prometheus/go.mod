module go.opentelemetry.io/otel/exporters/metric/prometheus

go 1.13

replace go.opentelemetry.io/otel => ../../..

require (
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.9.0
)
