module go.opentelemetry.io/otel/example/http-stackdriver

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/trace/stackdriver => ../../exporters/trace/stackdriver
)

require (
	go.opentelemetry.io/otel v0.2.1
	go.opentelemetry.io/otel/exporters/trace/stackdriver v0.2.1
	google.golang.org/grpc v1.24.0
)
