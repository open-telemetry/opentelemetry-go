module go.opentelemetry.io/otel/exporters/prometheus

go 1.17

require (
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/metric v0.26.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/sdk/metric v0.26.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.3.0 // indirect
	golang.org/x/sys v0.0.0-20210603081109-ebe580a85c40 // indirect
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/bridge/opencensus => ../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opentracing => ../../bridge/opentracing

replace go.opentelemetry.io/otel/example/jaeger => ../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../example/namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../../example/opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../example/passthrough

replace go.opentelemetry.io/otel/example/prom-collector => ../../example/prom-collector

replace go.opentelemetry.io/otel/example/prometheus => ../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../example/zipkin

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/exporters/prometheus => ./

replace go.opentelemetry.io/otel/exporters/jaeger => ../jaeger

replace go.opentelemetry.io/otel/exporters/zipkin => ../zipkin

replace go.opentelemetry.io/otel/internal/metric => ../../internal/metric

replace go.opentelemetry.io/otel/internal/tools => ../../internal/tools

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/sdk/export/metric => ../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric => ../otlp/otlpmetric

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ../otlp/otlpmetric/otlpmetricgrpc

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp => ../otlp/otlpmetric/otlpmetrichttp

replace go.opentelemetry.io/otel/bridge/opencensus/test => ../../bridge/opencensus/test

replace go.opentelemetry.io/otel/example/fib => ../../example/fib

replace go.opentelemetry.io/otel/schema => ../../schema

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../otlp/internal/retry
