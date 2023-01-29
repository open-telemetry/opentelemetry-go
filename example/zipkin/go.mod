module go.opentelemetry.io/otel/example/zipkin

go 1.18

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/exporters/zipkin => ../../exporters/zipkin
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.12.0
	go.opentelemetry.io/otel/exporters/zipkin v1.12.0
	go.opentelemetry.io/otel/sdk v1.12.0
	go.opentelemetry.io/otel/trace v1.12.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	golang.org/x/sys v0.0.0-20221010170243-090e33056c14 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace
