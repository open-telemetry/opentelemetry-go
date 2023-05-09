module go.opentelemetry.io/otel/bridge/opentracing/test

go 1.19

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/bridge/opentracing => ../

replace go.opentelemetry.io/otel/trace => ../../../trace

require (
	github.com/opentracing-contrib/go-grpc v0.0.0-20210225150812-73cb765af46e
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/otel v1.16.0-rc.1
	go.opentelemetry.io/otel/bridge/opentracing v1.16.0-rc.1
	google.golang.org/grpc v1.55.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0-rc.1 // indirect
	go.opentelemetry.io/otel/trace v1.16.0-rc.1 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../../metric
