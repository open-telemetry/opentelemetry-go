module go.opentelemetry.io/otel/log/logtest

go 1.26.0

require (
	github.com/google/go-cmp v0.7.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.47.0-rc.1
	go.opentelemetry.io/otel/log v1.47.0-rc.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel => ../../

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/log => ../
