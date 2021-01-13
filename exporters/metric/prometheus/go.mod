module go.opentelemetry.io/otel/exporters/metric/prometheus

go 1.14

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require (
	github.com/prometheus/client_golang v1.9.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
)
