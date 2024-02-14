module go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc

go 1.20

require (
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel v1.23.0-rc.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.0-rc.1
	go.opentelemetry.io/otel/sdk v1.23.0-rc.1
	go.opentelemetry.io/otel/trace v1.23.0-rc.1
	go.opentelemetry.io/proto/otlp v1.1.0
	go.uber.org/goleak v1.3.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe
	google.golang.org/grpc v1.61.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.23.0-rc.1 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240125205218-1f4bbc51befe // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel => ../../../..

replace go.opentelemetry.io/otel/sdk => ../../../../sdk

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../

replace go.opentelemetry.io/otel/trace => ../../../../trace

replace go.opentelemetry.io/otel/metric => ../../../../metric

replace go.opentelemetry.io/proto/otlp => ../../../../../opentelemetry-proto-go/otlp/

replace go.opentelemetry.io/otel/bridge/opencensus => ../../../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opencensus/test => ../../../../bridge/opencensus/test

replace go.opentelemetry.io/otel/bridge/opentracing => ../../../../bridge/opentracing

replace go.opentelemetry.io/otel/bridge/opentracing/test => ../../../../bridge/opentracing/test

replace go.opentelemetry.io/otel/example/dice => ../../../../example/dice

replace go.opentelemetry.io/otel/example/namedtracer => ../../../../example/namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../../../../example/opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ../../../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../../../example/passthrough

replace go.opentelemetry.io/otel/example/prometheus => ../../../../example/prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../../../../example/zipkin

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ../../otlpmetric/otlpmetricgrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp => ../../otlpmetric/otlpmetrichttp

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ./

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../otlptracehttp

replace go.opentelemetry.io/otel/exporters/prometheus => ../../../prometheus

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../../stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../../stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/zipkin => ../../../zipkin

replace go.opentelemetry.io/otel/internal/tools => ../../../../internal/tools

replace go.opentelemetry.io/otel/schema => ../../../../schema

replace go.opentelemetry.io/otel/sdk/metric => ../../../../sdk/metric
