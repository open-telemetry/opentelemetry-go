module go.opentelemetry.io/otel

go 1.25.0

require (
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/go-logr/logr v1.4.4
	github.com/go-logr/stdr v1.2.2
	github.com/google/go-cmp v0.7.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/auto/sdk v1.2.1
	go.opentelemetry.io/otel/metric v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

replace go.opentelemetry.io/otel/trace => ./trace

replace go.opentelemetry.io/otel/metric => ./metric
