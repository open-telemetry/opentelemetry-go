module go.opentelemetry.io/otel/bridge/opencensus/opencensusmetric

go 1.18

require (
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/metric v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/sdk/metric v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/sdk v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/otel/trace v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace
