module github.com/open-telemetry/opentelemetry-go/example/namedtracer

go 1.12

require (
	go.opentelemetry.io v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/example/namedtracer v0.0.0-00010101000000-000000000000
)

replace go.opentelemetry.io => ../..

replace go.opentelemetry.io/example/namedtracer => ./
