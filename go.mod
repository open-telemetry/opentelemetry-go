module go.opentelemetry.io/otel

go 1.14

require (
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel/label v0.16.0
	go.opentelemetry.io/otel/trace v0.16.0
)

replace (
	go.opentelemetry.io/otel/label => ./label
	go.opentelemetry.io/otel/trace => ./trace
)
