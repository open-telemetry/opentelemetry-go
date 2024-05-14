module go.opentelemetry.io/otel/example/prometheus

go 1.21

require (
	github.com/prometheus/client_golang v1.19.1
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/exporters/prometheus v0.48.0
	go.opentelemetry.io/otel/metric v1.26.0
	go.opentelemetry.io/otel/sdk/metric v1.26.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.15.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/prometheus => ../../exporters/prometheus

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/trace => ../../trace
