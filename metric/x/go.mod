module go.opentelemetry.io/otel/metric/x

go 1.25.0

require (
	go.opentelemetry.io/otel v1.47.0-rc.1
	go.opentelemetry.io/otel/metric v1.47.0-rc.1
)

require github.com/cespare/xxhash/v2 v2.3.0 // indirect

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/log => ../../log
