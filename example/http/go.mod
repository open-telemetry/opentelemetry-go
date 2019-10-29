module github.com/open-telemetry/opentelemetry-go/example/http

go 1.13

replace go.opentelemetry.io => ../..

require (
	go.opentelemetry.io v0.0.0-20191025183852-68310ab97435
	google.golang.org/grpc v1.24.0
)
