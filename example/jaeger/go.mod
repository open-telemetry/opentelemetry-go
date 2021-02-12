module go.opentelemetry.io/otel/example/jaeger

go 1.14

require (
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/trace/jaeger => ../../exporters/trace/jaeger

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/trace => ../../trace
