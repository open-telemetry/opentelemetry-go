module go.opentelemetry.io/otel/example/jaeger

go 1.19

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/jaeger => ../../exporters/jaeger
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/exporters/jaeger v1.15.1
	go.opentelemetry.io/otel/sdk v1.15.1
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v0.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	golang.org/x/sys v0.7.0 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/metric => ../../metric
