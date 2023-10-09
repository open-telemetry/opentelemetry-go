module go.opentelemetry.io/otel

go 1.20

require (
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr v1.2.2
	github.com/google/go-cmp v0.5.9
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/otel/metric v1.19.0
	go.opentelemetry.io/otel/trace v1.19.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/tools v0.12.1-0.20230815132531-74c255bcf846 // indirect
	golang.org/x/vuln v1.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.opentelemetry.io/otel/trace => ./trace

replace go.opentelemetry.io/otel/metric => ./metric
