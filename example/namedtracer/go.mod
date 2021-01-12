module go.opentelemetry.io/otel/example/namedtracer

go 1.14

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/codes => ../../codes
	go.opentelemetry.io/otel/exporters/stdout => ../../exporters/stdout
	go.opentelemetry.io/otel/label => ../../label
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/stdout v0.15.0
	go.opentelemetry.io/otel/label v0.1.0
	go.opentelemetry.io/otel/sdk v0.15.0
)
