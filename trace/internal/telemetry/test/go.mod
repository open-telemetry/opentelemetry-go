module go.opentelemetry.io/otel/trace/internal/telemetry/test

go 1.23.0

toolchain go1.24.0

require (
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/collector/pdata v1.25.0
	go.opentelemetry.io/otel/trace v1.34.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/grpc v1.70.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ../../..

replace go.opentelemetry.io/otel => ../../../..

replace go.opentelemetry.io/otel/metric => ../../../../metric
