module go.opentelemetry.io/otel/internal/tools

go 1.16

require (
	github.com/client9/misspell v0.3.4
	github.com/gogo/protobuf v1.3.2
	github.com/golangci/golangci-lint v1.45.3-0.20220330010013-d2ccc6d2bbbb
	github.com/itchyny/gojq v0.12.7
	github.com/jcchavezs/porto v0.4.0
	github.com/wadey/gocovmerge v0.0.0-20160331181800-b5bfa59ec0ad
	go.opentelemetry.io/build-tools/crosslink v0.0.0-20220502161954-e2bf744925c0
	go.opentelemetry.io/build-tools/dbotconf v0.0.0-20220321164008-b8e03aff061a
	go.opentelemetry.io/build-tools/multimod v0.0.0-20210920164323-2ceabab23375
	go.opentelemetry.io/build-tools/semconvgen v0.0.0-20210920164323-2ceabab23375
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3
	golang.org/x/tools v0.1.10
)

replace go.opentelemetry.io/otel/bridge/opentracing => ../../bridge/opentracing

replace go.opentelemetry.io/otel/example/fib => ../../example/fib

replace go.opentelemetry.io/otel/example/jaeger => ../../example/jaeger

replace go.opentelemetry.io/otel/example/namedtracer => ../../example/namedtracer

replace go.opentelemetry.io/otel/example/otel-collector => ../../example/otel-collector

replace go.opentelemetry.io/otel/example/passthrough => ../../example/passthrough

replace go.opentelemetry.io/otel/example/zipkin => ../../example/zipkin

replace go.opentelemetry.io/otel/exporters/jaeger => ../../exporters/jaeger

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../../exporters/otlp/internal/retry

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace => ../../exporters/otlp/otlptrace

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => ../../exporters/otlp/otlptrace/otlptracegrpc

replace go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp => ../../exporters/otlp/otlptrace/otlptracehttp

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/zipkin => ../../exporters/zipkin

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/internal/tools => ./

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/schema => ../../schema

replace go.opentelemetry.io/otel/sdk => ../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace
