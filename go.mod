module go.opentelemetry.io/otel

go 1.22.0

require (
	github.com/go-logr/logr v1.4.2
	github.com/go-logr/stdr v1.2.2
	github.com/google/go-cmp v0.6.0
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/auto/sdk v1.1.0
	go.opentelemetry.io/otel/metric v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ./trace

replace go.opentelemetry.io/otel/metric => ./metric
