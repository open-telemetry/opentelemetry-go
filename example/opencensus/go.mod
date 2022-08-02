module go.opentelemetry.io/otel/example/opencensus

go 1.17

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/bridge/opencensus => ../../bridge/opencensus
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opencensus.io v0.22.6-0.20201102222123-380f4078db9f
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/bridge/opencensus v0.31.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.31.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.9.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/sdk/metric v0.31.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace
