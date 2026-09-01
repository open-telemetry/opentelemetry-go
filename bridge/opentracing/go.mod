module go.opentelemetry.io/otel/bridge/opentracing

go 1.26.0

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/trace => ../../trace

require (
	github.com/opentracing-contrib/go-grpc v0.1.4
	github.com/opentracing-contrib/go-grpc/test v0.0.0-20260830185921-7bb06b076e19
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.47.0-rc.1
	go.opentelemetry.io/otel/trace v1.47.0-rc.1
	google.golang.org/grpc v1.83.2
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/log v1.47.0-rc.1 // indirect
	go.opentelemetry.io/otel/metric v1.47.0-rc.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260831171406-18b4a7587f8a // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/log => ../../log
