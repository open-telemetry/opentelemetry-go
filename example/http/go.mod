module github.com/open-telemetry/opentelemetry-go/example/http

go 1.12

replace go.opentelemetry.io => ../..

require (
	go.opentelemetry.io v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.24.0
)
