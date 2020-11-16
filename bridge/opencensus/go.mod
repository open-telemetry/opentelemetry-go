module go.opentelemetry.io/opentelemetry-go/bridge/opencensus

go 1.15

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v0.13.0
)

replace go.opentelemetry.io/otel => ../..
