[![Circle CI](https://circleci.com/gh/open-telemetry/opentelemetry-go.svg?style=svg)](https://circleci.com/gh/open-telemetry/opentelemetry-go)
[![Docs](https://godoc.org/github.com/open-telemetry/opentelemetry-go?status.svg)](http://godoc.org/github.com/open-telemetry/opentelemetry-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-telemetry/opentelemetry-go)](https://goreportcard.com/report/github.com/open-telemetry/opentelemetry-go)

This is a prototype *intended to be modified* into the opentelemetry-go implementation. The `api` directory here should be used as a starting point to introduce a new OpenTelemetry exporter, whereas the existing `exporter/observer` streaming model should be help verify the api 

To run the examples, first build the stderr tracer plugin (requires Linux or OS X):

```
(cd ./exporter/stdout/plugin && make)
(cd ./exporter/spanlog/plugin && make)
```

Then set the `OPENTELEMETRY_LIB` environment variable to the .so file in that directory, e.g., 

```
OPENTELEMETRY_LIB=./exporter/stderr/plugin/stderr.so go run ./example/server/server.go
```

and

```
OPENTELEMETRY_LIB=./exporter/spanlog/plugin/spanlog.so go run ./example/client/client.go
```
