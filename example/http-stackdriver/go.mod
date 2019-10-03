module github.com/open-telemetry/opentelemetry-go/example/http-stackdriver

go 1.13

replace go.opentelemetry.io => ../..

require (
	go.opentelemetry.io v0.0.0
	google.golang.org/grpc v1.24.0
)
