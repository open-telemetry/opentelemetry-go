module go.opentelemetry.io/otel/example/namedtracer

go 1.13

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/stdout => ../../exporters/stdout
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	github.com/google/gofuzz v1.1.0 // indirect
	go.opentelemetry.io/otel v0.9.0
	go.opentelemetry.io/otel/exporters/stdout v0.9.0
	go.opentelemetry.io/otel/sdk v0.9.0
)
