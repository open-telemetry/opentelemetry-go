module go.opentelemetry.io/otel/sdk

go 1.20

replace go.opentelemetry.io/otel => ../

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/go-logr/logr v1.2.4
	github.com/google/go-cmp v0.5.9
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.17.0
	go.opentelemetry.io/otel/schema v0.0.5
	go.opentelemetry.io/otel/trace v1.17.0
	golang.org/x/sys v0.12.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ../trace

replace go.opentelemetry.io/otel/metric => ../metric
