module go.opentelemetry.io/otel/example/fib

go 1.16

require (
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/trace => ../../trace
