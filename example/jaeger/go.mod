// Deprecated: This example is no longer supported as 
[go.opentelemetry.io/otel/exporters/jaeger] is no longer supported.
// OpenTelemetry dropped support for Jaeger exporter in July 2023.
// Jaeger officially accepts and recommends using OTLP.
// Use [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp]
// or [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc] instead.
module go.opentelemetry.io/otel/example/jaeger

go 1.19

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/jaeger => ../../exporters/jaeger
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/exporters/jaeger v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/metric => ../../metric
