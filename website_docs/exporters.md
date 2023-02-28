---
title: Exporters
aliases: [/docs/instrumentation/go/exporting_data]
weight: 4
---

In order to visualize and analyze your [traces](/docs/concepts/signals/traces/)
and metrics, you will need to export them to a backend.

## OTLP endpoint

To send trace data to an OTLP endpoint (like the [collector](/docs/collector) or
Jaeger >= v1.35.0) you'll want to configure an OTLP exporter that sends to your endpoint.

### Using HTTP

```go
import (
  	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
  	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func installExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}
  	/* … */
}
```

To learn more on how to use the OTLP HTTP exporter, try out the [otel-collector](https://github.com/open-telemetry/opentelemetry-go/tree/main/example/otel-collector)

### Jaeger

To try out the OTLP exporter, since v1.35.0 you can run
[Jaeger](https://www.jaegertracing.io/)  as an OTLP endpoint and for trace
visualization in a docker container:

```shell
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```

## Prometheus

Prometheus export is available in the `go.opentelemetry.io/otel/exporters/prometheus` package.

Please find more documentation on [GitHub](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/prometheus)
