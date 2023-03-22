module go.opentelemetry.io/otel/metric

go 1.19

require (
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/otel v1.15.0-rc.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/trace => ../trace
