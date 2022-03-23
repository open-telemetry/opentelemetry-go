module go.opentelemetry.io/otel/exporters/otlp/internal

go 1.16

require (
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.5.0
	google.golang.org/grpc v1.45.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ./retry

replace go.opentelemetry.io/otel/bridge/opencensus => ../../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opencensus/test => ../../../bridge/opencensus/test

replace go.opentelemetry.io/otel/bridge/opentracing => ../../../bridge/opentracing

replace go.opentelemetry.io/otel/example/fib => ../../../example/fib

replace go.opentelemetry.io/otel/example/jaeger => ../../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../../example/namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../../../example/opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../../example/passthrough

replace go.opentelemetry.io/otel/example/prometheus => ../../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../../example/zipkin

replace go.opentelemetry.io/otel/exporters/jaeger => ../../jaeger

replace go.opentelemetry.io/otel/exporters/otlp/internal => ./

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric => ../otlpmetric

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ../otlpmetric/otlpmetricgrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp => ../otlpmetric/otlpmetrichttp

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/exporters/prometheus => ../../prometheus

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/zipkin => ../../zipkin

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/internal/tools => ../../../internal/tools

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/schema => ../../../schema

replace go.opentelemetry.io/otel/sdk/export/metric => ../../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace
