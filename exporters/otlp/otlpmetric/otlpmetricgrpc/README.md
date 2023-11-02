# OpenTelemetry-Go OTLP Metric gRPC Exporter

[![Go Reference](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc)

[OpenTelemetry Protocol Exporter](https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/protocol/exporter.md) implementation using gRPC.

## Installation

```
go get -u go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc
```

## Usage

Exporters should be created using the [New](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#New) and used with a [PeriodicReader].

## Environment Variables

The following environment variables can be used (instead of options objects) to override the default configuration.

| Name | Description | Default | Override with |
|------|-------------|---------|---------------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Endpoint to which the exporter is going to send traces or metrics. If the scheme is `http` or `unix` this will also set `OTEL_EXPORTER_OTLP_INSECURE` to true | `https://localhost:4317` | `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`, [WithEndpoint()] |
| `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` | Endpoint to which the exporter is going to send metrics. If the scheme is `http` or `unix` this will also set `OTEL_EXPORTER_OTLP_INSECURE` to true | `https://localhost:4317` | [WithEndpoint()] |
| `OTEL_EXPORTER_OTLP_INSECURE` | If set to true the connection will not attempt to use TLS when connecting | `false` | `OTEL_EXPORTER_OTLP_METRICS_INSECURE`, [WithInsecure()] |
| `OTEL_EXPORTER_OTLP_METRICS_INSECURE` | If set to true the connection will not attempt to use TLS when connecting | `false` | [WithInsecure()] |
| `OTEL_EXPORTER_OTLP_HEADERS` | A list of headers to send with each request | none | `OTEL_EXPORTER_OTLP_METRICS_HEADERS`, [WithHeaders()] |
| `OTEL_EXPORTER_OTLP_METRICS_HEADERS` | A list of headers to send with each request | none | [WithHeaders()] |
| `OTEL_EXPORTER_OTLP_COMPRESSION` | Sets the compressions used in the connection. Supports `none` and `gzip`.  Must import the compressor for gzip to work. | none | `OTEL_EXPORTER_OTLP_METRICS_COMPRESSION`, [WithCompressor()] |
| `OTEL_EXPORTER_OTLP_METRICS_COMPRESSION` | Sets the compressions used in the connection. Supports `none` and `gzip`.  Must import the compressor for gzip to work. | none | [WithCompressor()] |
| `OTEL_EXPORTER_OTLP_TIMEOUT` | Sets the max amount of time (as milliseconds) an Exporter will attempt an export | `10000` | `OTEL_EXPORTER_OTLP_METRICS_TIMEOUT`, [WithTimeout()] |
| `OTEL_EXPORTER_OTLP_METRICS_TIMEOUT` | Sets the max amount of time (as milliseconds) an Exporter will attempt an export | `10000` | [WithTimeout()] |
| `OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE` | The aggregation temporality to use on the basis of instrument kind. Available values are `Cumulative`, `Delta`, and `LowMemory`.  See [The OTLP Exporter Specification] for more details | `Cumulative` | [WithTemporalitySelector()] |

[PeriodicReader]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#NewPeriodicReader
[WithEndpoint()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithEndpoint
[WithInsecure()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithInsecure
[WithHeaders()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithHeaders
[WithCompressor()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithCompressor
[WithTimeout()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithTimeout
[The OTLP Exporter Specification]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/sdk_exporters/otlp.md#additional-configuration
[WithTemporalitySelector()]: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc#WithTemporalitySelector
