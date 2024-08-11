module go.opentelemetry.io/otel/example/zipkin

go 1.21

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/zipkin => ../../exporters/zipkin
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/zipkin v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/metric => ../../metric
