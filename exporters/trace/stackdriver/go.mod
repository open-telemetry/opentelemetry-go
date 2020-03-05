module go.opentelemetry.io/otel/exporters/trace/stackdriver

go 1.13

replace go.opentelemetry.io/otel => ../../..

require (
	cloud.google.com/go v0.53.0
	github.com/golang/protobuf v1.3.4
	github.com/stretchr/testify v1.4.0
	go.opentelemetry.io/otel v0.2.3
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200303153909-beee998c1893
	google.golang.org/grpc v1.27.1
)
