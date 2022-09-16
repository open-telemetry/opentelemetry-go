module go.opentelemetry.io/otel/bridge/opencensus

go 1.18

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
)

replace go.opentelemetry.io/otel => ../..

replace go.opentelemetry.io/otel/trace => ../../trace
