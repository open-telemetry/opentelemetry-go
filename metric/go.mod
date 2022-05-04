module go.opentelemetry.io/otel/metric

go 1.16

require (
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.7.0
)

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/trace => ../trace
