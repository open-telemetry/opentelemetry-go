module go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp

go 1.22

require (
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/google/go-cmp v0.6.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/log v0.6.0
	go.opentelemetry.io/otel/sdk v1.30.0
	go.opentelemetry.io/otel/sdk/log v0.6.0
	go.opentelemetry.io/otel/trace v1.30.0
	go.opentelemetry.io/proto/otlp v1.3.1
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240924160255-9d4c2d233b61 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240924160255-9d4c2d233b61 // indirect
	google.golang.org/grpc v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel => ../../../..

replace go.opentelemetry.io/otel/sdk/log => ../../../../sdk/log

replace go.opentelemetry.io/otel/trace => ../../../../trace

replace go.opentelemetry.io/otel/sdk => ../../../../sdk

replace go.opentelemetry.io/otel/metric => ../../../../metric

replace go.opentelemetry.io/otel/log => ../../../../log
