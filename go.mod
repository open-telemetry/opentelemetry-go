module go.opentelemetry.io/otel

go 1.14

replace (
	go.opentelemetry.io/otel/label => ./label
	go.opentelemetry.io/otel/codes => ./codes
)

require (
	github.com/google/go-cmp v0.5.4
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel/codes v0.1.0
	go.opentelemetry.io/otel/label v0.1.0
)
