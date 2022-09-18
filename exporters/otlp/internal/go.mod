module go.opentelemetry.io/otel/exporters/otlp/internal

go 1.18

require (
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0
	google.golang.org/grpc v1.49.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/bridge/opencensus => ../../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opencensus/opencensusmetric => ../../../bridge/opencensus/opencensusmetric

replace go.opentelemetry.io/otel/bridge/opencensus/test => ../../../bridge/opencensus/test

replace go.opentelemetry.io/otel/bridge/opentracing => ../../../bridge/opentracing

replace go.opentelemetry.io/otel/example/fib => ../../../example/fib

replace go.opentelemetry.io/otel/example/jaeger => ../../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../../example/namedtracer

replace go.opentelemetry.io/otel/example/otel-collector => ../../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../../example/passthrough

replace go.opentelemetry.io/otel/example/prometheus => ../../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../../example/zipkin

replace go.opentelemetry.io/otel/exporters/jaeger => ../../jaeger

replace go.opentelemetry.io/otel/exporters/otlp/internal => ./

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ./retry

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

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace
