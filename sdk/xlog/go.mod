module go.opentelemetry.io/otel/sdk/xlog

go 1.22.0

require (
	go.opentelemetry.io/otel/log v0.10.0
	go.opentelemetry.io/otel/sdk v1.34.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace go.opentelemetry.io/otel/log => ../../log

replace go.opentelemetry.io/otel/sdk => ../../sdk
