module go.opentelemetry.io/example/jaeger

go 1.13

replace (
	go.opentelemetry.io => ../..
	go.opentelemetry.io/exporter/trace/jaeger => ../../exporter/trace/jaeger
)

require (
	go.opentelemetry.io v0.0.0-20191025183852-68310ab97435
	go.opentelemetry.io/exporter/trace/jaeger v0.0.0-20191025183852-68310ab97435
)
