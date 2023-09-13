module go.opentelemetry.io/otel/exporters/stdout/stdoutmetric

go 1.20

require (
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.18.0
	go.opentelemetry.io/otel/sdk v1.18.0
	go.opentelemetry.io/otel/sdk/metric v0.41.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.18.0 // indirect
	go.opentelemetry.io/otel/schema v0.0.6 // indirect
	go.opentelemetry.io/otel/trace v1.18.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/schema => ../../../schema
