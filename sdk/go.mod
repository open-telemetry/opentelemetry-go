module go.opentelemetry.io/otel/sdk

go 1.14

replace (
	go.opentelemetry.io/otel => ../
	go.opentelemetry.io/otel/label => ../label
	go.opentelemetry.io/otel/semconv => ../semconv
)

require (
	github.com/benbjohnson/clock v1.0.3
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/label v0.16.0
	go.opentelemetry.io/otel/semconv v0.16.0
)
