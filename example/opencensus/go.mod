module go.opentelemetry.io/otel/example/opencensus

go 1.16

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/bridge/opencensus => ../../bridge/opencensus
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/bridge/opencensus v0.30.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.30.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/sdk/metric v0.30.0
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace
