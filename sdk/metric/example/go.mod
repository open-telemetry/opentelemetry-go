module go.opentelemetry.io/otel/sdk/metric/example

go 1.20

require (
	go.opentelemetry.io/otel v1.18.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.41.0
	go.opentelemetry.io/otel/sdk v1.18.0
	go.opentelemetry.io/otel/sdk/metric v0.41.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.18.0 // indirect
	go.opentelemetry.io/otel/trace v1.18.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/sdk => ../../

replace go.opentelemetry.io/otel/sdk/metric => ../
