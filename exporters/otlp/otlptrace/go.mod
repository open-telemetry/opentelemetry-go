module go.opentelemetry.io/otel/exporters/otlp/otlptrace

go 1.20

require (
	github.com/google/go-cmp v0.6.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.23.0-rc.1
	go.opentelemetry.io/otel/sdk v1.23.0-rc.1
	go.opentelemetry.io/otel/trace v1.23.0-rc.1
	go.opentelemetry.io/proto/otlp v1.1.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	go.opentelemetry.io/otel/metric v1.23.0-rc.1 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/proto/otlp => ../../../../opentelemetry-proto-go/otlp/
