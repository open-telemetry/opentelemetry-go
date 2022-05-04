module go.opentelemetry.io/otel/sdk

go 1.16

replace go.opentelemetry.io/otel => ../

require (
	github.com/go-logr/logr v1.2.3
	github.com/google/go-cmp v0.5.8
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7
)

replace go.opentelemetry.io/otel/trace => ../trace
