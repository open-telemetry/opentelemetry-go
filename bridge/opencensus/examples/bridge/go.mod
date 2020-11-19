module go.opentelemetry.io/otel/bridge/opencensus/examples/bridge

go 1.14

replace (
	go.opentelemetry.io/otel => ../../../..
	go.opentelemetry.io/otel/bridge/opencensus => ../..
	go.opentelemetry.io/otel/exporters/stdout => ../../../../exporters/stdout
	go.opentelemetry.io/otel/sdk => ../../../../sdk
)

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/bridge/opencensus v0.0.0-20201117180221-c857a3da18cb
	go.opentelemetry.io/otel/exporters/stdout v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
)
