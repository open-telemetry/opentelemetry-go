# OpenTelemetry/OpenCensus Bridge

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/bridge/opencensus)](https://pkg.go.dev/go.opentelemetry.io/otel/bridge/opencensus)
[Example](https://github.com/open-telemetry/opentelemetry-go/blob/main/example/opencensus/main.go)

## Getting started

Assuming you have configured an OpenTelemetry `TracerProvider`, these will be
the steps to follow to wire up the bridge:

```go
import (
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
)

func main() {
	/* Create tracerProvider and configure OpenTelemetry ... */

	tracer := otel.Tracer("opencensus")
	octrace.DefaultTracer = opencensus.NewTracer(tracer)

	/* ... */
}
```
