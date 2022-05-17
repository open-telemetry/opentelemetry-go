module go.opentelemetry.io/otel

go 1.16

require (
	github.com/go-logr/logr v1.2.3
	github.com/go-logr/stdr v1.2.2
	github.com/google/go-cmp v0.5.8
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel/trace => ./trace
