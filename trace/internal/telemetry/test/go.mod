module go.opentelemetry.io/otel/trace/internal/telemetry/test

go 1.25.0

require (
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/collector/pdata v1.65.0
	go.opentelemetry.io/otel/trace v1.47.0-rc.1
)

require (
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	go.opentelemetry.io/collector/featuregate v1.65.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../..

replace go.opentelemetry.io/otel => ../../../..

replace go.opentelemetry.io/otel/metric => ../../../../metric

replace go.opentelemetry.io/otel/log => ../../../../log
