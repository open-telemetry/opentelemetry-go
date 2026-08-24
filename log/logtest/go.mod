module go.opentelemetry.io/otel/log/logtest

go 1.25.0

require (
	github.com/google/go-cmp v0.7.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/log v0.21.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel => ../../

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/log => ../
