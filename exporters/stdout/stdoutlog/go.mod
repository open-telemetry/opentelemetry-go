module go.opentelemetry.io/otel/exporters/stdout/stdoutlog

go 1.18

replace (
	go.opentelemetry.io/otel => ../../..
	go.opentelemetry.io/otel/sdk => ../../../sdk
)

require go.opentelemetry.io/otel/sdk v1.11.2

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/log v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
)

replace go.opentelemetry.io/otel/log => ../../../log

replace go.opentelemetry.io/otel/sdk/log => ../../../sdk/log
