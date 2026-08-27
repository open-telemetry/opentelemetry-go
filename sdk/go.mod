module go.opentelemetry.io/otel/sdk

go 1.26.0

replace go.opentelemetry.io/otel => ../

require (
	github.com/go-logr/logr v1.4.4
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.47.0-rc.1
	go.opentelemetry.io/otel/metric v1.47.0-rc.1
	go.opentelemetry.io/otel/sdk/metric v1.47.0-rc.1
	go.opentelemetry.io/otel/trace v1.47.0-rc.1
	go.uber.org/goleak v1.3.0
	golang.org/x/sys v0.47.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/log v1.47.0-rc.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace go.opentelemetry.io/otel/trace => ../trace

replace go.opentelemetry.io/otel/metric => ../metric

replace go.opentelemetry.io/otel/sdk/metric => ./metric

replace go.opentelemetry.io/otel/metric/x => ../metric/x

replace go.opentelemetry.io/otel/log => ../log
