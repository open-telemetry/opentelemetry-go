module go.opentelemetry.io/otel/example/dice

go 1.20

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0
	go.opentelemetry.io/otel v1.23.0-rc.1
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.23.0-rc.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.23.0-rc.1
	go.opentelemetry.io/otel/metric v1.23.0-rc.1
	go.opentelemetry.io/otel/sdk v1.23.0-rc.1
	go.opentelemetry.io/otel/sdk/metric v1.23.0-rc.1
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/trace v1.23.0-rc.1 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/sdk => ../../sdk
