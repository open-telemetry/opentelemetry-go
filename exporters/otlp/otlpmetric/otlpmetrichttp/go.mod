module go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp

go 1.18

retract v0.32.2 // Contains unresolvable dependencies.

require (
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/otel v1.11.0
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.32.3
	go.opentelemetry.io/otel/metric v0.32.3
	go.opentelemetry.io/otel/sdk/metric v0.32.3
	go.opentelemetry.io/proto/otlp v0.19.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.12.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/sdk v1.11.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.0 // indirect
	golang.org/x/net v0.0.0-20221017152216-f25eb7ecb193 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20221014213838-99cd37c6964a // indirect
	google.golang.org/grpc v1.50.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel => ../../../..

replace go.opentelemetry.io/otel/sdk => ../../../../sdk

replace go.opentelemetry.io/otel/sdk/metric => ../../../../sdk/metric

replace go.opentelemetry.io/otel/exporters/otlp/otlpmetric => ../

replace go.opentelemetry.io/otel/metric => ../../../../metric

replace go.opentelemetry.io/otel/trace => ../../../../trace

replace go.opentelemetry.io/otel/exporters/otlp/internal/retry => ../../internal/retry
