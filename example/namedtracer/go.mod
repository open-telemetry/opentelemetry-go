module go.opentelemetry.io/otel/example/namedtracer

go 1.16

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	github.com/go-logr/stdr v1.2.2
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace
