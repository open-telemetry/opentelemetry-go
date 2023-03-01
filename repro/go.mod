module repro

go 1.20

require (
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/contrib/instrumentation/runtime v0.40.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.37.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.37.0
	go.opentelemetry.io/otel/metric v0.37.0
	go.opentelemetry.io/otel/sdk/metric v0.37.0
	go.opentelemetry.io/proto/otlp v0.19.0
	google.golang.org/grpc v1.53.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.14.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/sdk => ../sdk

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../exporters/otlp/internal/retry

replace go.opentelemetry.io/otel => ../

replace go.opentelemetry.io/otel/trace => ../trace

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric => ../exporters/otlp/otlpmetric

replace go.opentelemetry.io/otel/metric => ../metric

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => ../exporters/otlp/otlpmetric/otlpmetricgrpc

replace go.opentelemetry.io/otel/sdk/metric => ../sdk/metric

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../exporters/stdout/stdoutmetric
