module go.opentelemetry.io/otel/example/fib

go 1.20

require (
	go.opentelemetry.io/otel v1.18.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.18.0
	go.opentelemetry.io/otel/sdk v1.18.0
	go.opentelemetry.io/otel/trace v1.18.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.18.0 // indirect
	go.opentelemetry.io/otel/schema v0.0.6 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/schema => ../../schema
