module github.com/open-telemetry/opentelemetry-go/example/http-stackdriver

go 1.13

replace (
	go.opentelemetry.io => ../..
	go.opentelemetry.io/exporter/trace/stackdriver => ../../exporter/trace/stackdriver
)

require (
	go.opentelemetry.io v0.0.0-20191031063502-886243699327
	go.opentelemetry.io/exporter/trace/stackdriver v0.0.0-20191025183852-68310ab97435
	google.golang.org/grpc v1.24.0
)
