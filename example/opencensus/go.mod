module go.opentelemetry.io/otel/example/opencensus

go 1.19

replace (
	go.opentelemetry.io/otel => ../..
	go.opentelemetry.io/otel/bridge/opencensus => ../../bridge/opencensus
	go.opentelemetry.io/otel/sdk => ../../sdk
)

require (
	go.opencensus.io v0.24.0
	go.opentelemetry.io/otel v1.15.0-rc.2
	go.opentelemetry.io/otel/bridge/opencensus v0.38.0-rc.2
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.38.0-rc.2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.15.0-rc.2
	go.opentelemetry.io/otel/sdk v1.15.0-rc.2
	go.opentelemetry.io/otel/sdk/metric v0.38.0-rc.2
)

require (
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/kr/text v0.2.0 // indirect
	go.opentelemetry.io/otel/metric v1.15.0-rc.2 // indirect
	go.opentelemetry.io/otel/schema v0.0.4 // indirect
	go.opentelemetry.io/otel/trace v1.15.0-rc.2 // indirect
	golang.org/x/sys v0.7.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../metric

replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../trace

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/schema => ../../schema
