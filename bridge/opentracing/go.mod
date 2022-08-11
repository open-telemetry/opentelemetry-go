module go.opentelemetry.io/otel/bridge/opentracing

go 1.17

replace go.opentelemetry.io/otel => ../..

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.7.2
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace
