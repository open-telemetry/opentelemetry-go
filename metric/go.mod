module go.opentelemetry.io/otel/metric

go 1.18

require (
	github.com/stretchr/testify v1.8.1
	go.opentelemetry.io/otel v1.11.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/trace => ../trace
