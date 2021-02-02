module go.opentelemetry.io/otel/semconv

go 1.15

require (
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel/label v0.16.0
)

replace go.opentelemetry.io/otel/label => ../label
