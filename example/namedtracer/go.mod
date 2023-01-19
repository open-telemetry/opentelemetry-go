module go.opentelemetry.io/otel/example/namedtracer

go 1.18

replace go.opentelemetry.io/otel => ../..

require (
	github.com/go-logr/stdr v1.2.2
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.11.2
	go.opentelemetry.io/otel/sdk v1.11.2
	go.opentelemetry.io/otel/trace v1.11.2
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/log v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/log => ../../log

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/stdout/stdoutlog => ../../exporters/stdout/stdoutlog
