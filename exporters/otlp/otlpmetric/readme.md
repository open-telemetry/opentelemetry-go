# OpenTelemetry-Go OTLP Metric Exporter

[![Go Reference](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/otlp/otlpmetric.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric)

[OpenTelemetry Protocol Exporter](https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/protocol/exporter.md) implementation.

## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc
```

or 

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp
```

## Examples

- [HTTP Exporter setup and examples](./otlpmetrichttp/example_test.go)
- [Full example of gRPC Exporter sending telemetry to a local collector](../../../example/otel-collector)

## [`otlpmetric`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric)

The `otlptrace` package provides common functions used within the [otlpmetricgrpc](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc) and [otlpmetrichttp](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp) packages.

Exporters should be created using the New functions from the respective packages.

## Configuration

### Environment Variables

The following environment variables can be used (instead of options objects) to override the default configuration.
For more information about how each of these environment variables is interpreted, see [the OpenTelemetry
specification](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/exporter.md).

Note: all names are prefixed with `OTEL_EXPORTER_OTLP_` and are uppercased.

| Name | Description | Default | Override with |
|------|-------------|---------|---------------|
| `ENDPOINT` | Endpoint to which the exporter is going to send traces or metrics. If the scheme is `http` or `unix` this will also set `INSECURE` to true |`http://localhost:4317`| `METRICS_ENDPOINT`, [WithEndpoint()] |
| `METRICS_ENDPOINT` | Endpoint to which the exporter is going to send metrics. If the scheme is `http` or `unix` this will also set `INSECURE` to true |`http://localhost:4317`| [WithEndpoint()] |
| `INSECURE` | If set to true the connection will not attempt to use tls when connecting | http: insecure, grpc: secure | `METRICS_INSECURE`, [WithInsecure()] |
| `METRICS_INSECURE` | If set to true the connection will not attempt to use tls when connecting | http: insecure, grpc: secure | [WithInsecure()] |
| `HEADERS` | A list of headers to send with each request | none | `METRICS_HEADERS`, [WithHeaders()] |
| `METRICS_HEADERS` | A list of headers to send with each request | none | [WithHeaders()] |
| `COMPRESSION` | Sets the compressions used in the connection. Supports `none` and `gzip`.  Must import the compressor for gzip to work. | none | `METRICS_COMPRESSION`, http: [WithCompression()] grpc: [WithCompressor()] |
| `METRICS_COMPRESSION` | Sets the compressions used in the connection. Supports `none` and `gzip`.  Must import the compressor for gzip to work. | none | http: [WithCompression()] grpc: [WithCompressor()] |
| `TIMEOUT` | Sets the max amount of time an Exporter will attempt an export | 10s | `METRICS_TIMEOUT`, [WithTimeout()] |
| `METRICS_TIMEOUT` | Sets the max amount of time an Exporter will attempt an export | 10s | [WithTimeout()] |
| `METRICS_TEMPORALITY_PREFERENCE` | The aggregation temporality to use on the basis of instrument kind. Available values are `Cumulative`, `Delta`, and `LowMemory`.  See [The OTLP Exporter Specification] for more details| Cumulative | [WithTemporalitySelector()] |

[WithEndpoint()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithEndpoint
[WithInsecure()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithInsecure
[WithHeader()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithHeaders
[WithCompression()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp#WithCompression
[WithCompressor()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithCompressor
[WithTimeout()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithTimeout
[The OTLP Exporter Specification]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/sdk_exporters/otlp.md#additional-configuration
[WithTemporalitySelector()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithTemporalitySelector