module go.opentelemetry.io/otel/bridge/opentracing/test

go 1.18

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/bridge/opentracing => ../

replace go.opentelemetry.io/otel/trace => ../../../trace

require (
	github.com/opentracing-contrib/go-grpc v0.0.0-20210225150812-73cb765af46e
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/otel v1.13.0
	go.opentelemetry.io/otel/bridge/opentracing v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.53.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
