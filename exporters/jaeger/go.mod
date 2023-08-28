// Deprecated: This module is no longer supported.
// OpenTelemetry dropped support for Jaeger exporter in July 2023.
// Jaeger officially accepts and recommends using OTLP.
// Use [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp]
// or [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc] instead.
module go.opentelemetry.io/otel/exporters/jaeger

go 1.19

require (
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr v1.2.2
	github.com/google/go-cmp v0.5.9
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.17.0
	go.opentelemetry.io/otel/sdk v1.17.0
	go.opentelemetry.io/otel/trace v1.17.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	go.opentelemetry.io/otel/metric v1.17.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/metric => ../../metric
