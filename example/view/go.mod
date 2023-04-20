module go.opentelemetry.io/otel/example/view

go 1.19

require (
	github.com/prometheus/client_golang v1.15.0
	go.opentelemetry.io/otel v1.15.0-rc.2
	go.opentelemetry.io/otel/exporters/prometheus v0.38.0-rc.2
	go.opentelemetry.io/otel/metric v1.15.0-rc.2
	go.opentelemetry.io/otel/sdk v1.15.0-rc.2
	go.opentelemetry.io/otel/sdk/metric v0.38.0-rc.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	go.opentelemetry.io/otel/trace v1.15.0-rc.2 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/exporters/prometheus => ../../exporters/prometheus

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/trace => ../../trace
