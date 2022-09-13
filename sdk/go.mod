module go.opentelemetry.io/otel/sdk

go 1.17

replace go.opentelemetry.io/otel => ../

require (
	github.com/go-logr/logr v1.2.3
	github.com/google/go-cmp v0.5.8
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace go.opentelemetry.io/otel/trace => ../trace
