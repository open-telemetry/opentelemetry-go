module go.opentelemetry.io/otel/example/otel-collector

go 1.20

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opentelemetry.io/otel v1.23.0-rc.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.23.0-rc.1
	go.opentelemetry.io/otel/sdk v1.23.0-rc.1
	go.opentelemetry.io/otel/trace v1.23.0-rc.1
	google.golang.org/grpc v1.61.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.0-rc.1 // indirect
	go.opentelemetry.io/otel/metric v1.23.0-rc.1 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../../exporters/otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../../exporters/otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/proto/otlp => ../../../opentelemetry-proto-go/otlp/

replace go.opentelemetry.io/otel/bridge/opencensus => ../../bridge/opencensus

replace go.opentelemetry.io/otel/bridge/opencensus/test => ../../bridge/opencensus/test

replace go.opentelemetry.io/otel/bridge/opentracing => ../../bridge/opentracing

replace go.opentelemetry.io/otel/bridge/opentracing/test => ../../bridge/opentracing/test

replace go.opentelemetry.io/otel/example/dice => ../dice

replace go.opentelemetry.io/otel/example/namedtracer => ../namedtracer

replace go.opentelemetry.io/otel/example/opencensus => ../opencensus

replace go.opentelemetry.io/otel/example/otel-collector => ./

replace go.opentelemetry.io/otel/example/passthrough => ../passthrough

replace go.opentelemetry.io/otel/example/prometheus => ../prometheus

replace go.opentelemetry.io/otel/example/zipkin => ../zipkin

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ../../exporters/otlp/otlpmetric/otlpmetricgrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp => ../../exporters/otlp/otlpmetric/otlpmetrichttp

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../../exporters/otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/exporters/prometheus => ../../exporters/prometheus

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/zipkin => ../../exporters/zipkin

replace go.opentelemetry.io/otel/internal/tools => ../../internal/tools

replace go.opentelemetry.io/otel/schema => ../../schema

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric
