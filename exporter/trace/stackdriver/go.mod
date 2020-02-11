module go.opentelemetry.io/otel/exporter/trace/stackdriver

go 1.13

replace go.opentelemetry.io/otel => ../../..

require (
	cloud.google.com/go v0.47.0
	github.com/golang/protobuf v1.3.2
	github.com/stretchr/testify v1.4.0
	go.opentelemetry.io/otel v0.2.1
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.11.0
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03
	google.golang.org/grpc v1.24.0
)
