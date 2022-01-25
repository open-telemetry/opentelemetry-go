module go.opentelemetry.io/otel/exporters/otlp/otlpmetric

go 1.17

require (
	github.com/google/go-cmp v0.5.7
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.3.0
	go.opentelemetry.io/otel/metric v0.26.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/sdk/metric v0.26.0
	go.opentelemetry.io/proto/otlp v0.12.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.3.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	golang.org/x/text v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk => ../../../sdk

replace go.opentelemetry.io/otel/metric => ../../../metric

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/sdk/export/metric => ../../../sdk/export/metric

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/bridge/opencensus => ../../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opentracing => ../../../bridge/opentracing

replace go.opentelemetry.io/otel/example/jaeger => ../../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../../example/namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../../../example/opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../../example/passthrough

replace go.opentelemetry.io/otel/example/prometheus => ../../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../../example/zipkin

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../internal/retry

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric => ./

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ./otlpmetricgrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/internal/tools => ../../../internal/tools

replace go.opentelemetry.io/otel/internal/metric => ../../../internal/metric

replace go.opentelemetry.io/otel/exporters/jaeger => ../../jaeger

replace go.opentelemetry.io/otel/exporters/prometheus => ../../prometheus

replace go.opentelemetry.io/otel/exporters/zipkin => ../../zipkin

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp => ./otlpmetrichttp

replace go.opentelemetry.io/otel/bridge/opencensus/test => ../../../bridge/opencensus/test

replace go.opentelemetry.io/otel/example/fib => ../../../example/fib

replace go.opentelemetry.io/otel/schema => ../../../schema
