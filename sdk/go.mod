module go.opentelemetry.io/otel/sdk

go 1.14

replace (
	go.opentelemetry.io/otel => ../
	go.opentelemetry.io/otel/codes => ../codes
	go.opentelemetry.io/otel/label => ../label
)

require (
	github.com/benbjohnson/clock v1.0.3
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/codes v0.1.0
	go.opentelemetry.io/otel/label v0.1.0
)
