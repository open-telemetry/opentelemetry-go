module go.opentelemetry.io/otel/example/otel-collector

go 1.14

require (
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/otlp v0.9.0
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	google.golang.org/grpc v1.30.0
)

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/otlp => ../../exporters/otlp
)
