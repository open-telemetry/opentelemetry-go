module go.opentelemetry.io/otel/exporters/metric/prometheus

go 1.13

replace go.opentelemetry.io/otel => ../../..

require (
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/prometheus/client_golang v1.5.0
	github.com/prometheus/procfs v0.0.10 // indirect
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.6.0
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
)
