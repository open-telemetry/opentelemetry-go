module go.opentelemetry.io/exporter/trace/stackdriver

go 1.12

replace go.opentelemetry.io => ../../..

require (
	cloud.google.com/go v0.47.0
	github.com/golang/protobuf v1.3.2
	go.opentelemetry.io v0.0.0-20191021171549-9b5f5dd13acd
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.11.0
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03
	google.golang.org/grpc v1.24.0
)
