# OpenTelemetry Exporters

Once the OpenTelemetry SDK has created and processed telemetry, it needs to be exported.
This package contains exporters for this purpose.

## Exporter Packages

The following exporter packages are provided with the following OpenTelemetry signal support.

| Exporter Package                                                                | Metrics | Traces |
| :-----------------------------------------------------------------------------: | :-----: | :----: |
| [go.opentelemetry.io/otel/exporters/jaeger](./jaeger)                           |         | ✓      |
| [go.opentelemetry.io/otel/exporters/otlp/otlpmetric](./otlp/otlpmetric)         | ✓       |        |
| [go.opentelemetry.io/otel/exporters/otlp/otlptrace](./otlp/otlptrace)           |         | ✓      |
| [go.opentelemetry.io/otel/exporters/prometheus](./prometheus)                   | ✓       |        |
| [go.opentelemetry.io/otel/exporters/stdout/stdoutmetric](./stdout/stdoutmetric) | ✓       |        |
| [go.opentelemetry.io/otel/exporters/stdout/stdouttrace](./stdout/stdouttrace)   |         | ✓      |
| [go.opentelemetry.io/otel/exporters/zipkin](./zipkin)                           |         | ✓      |

See the [OpenTelemetry registry] for 3rd-part exporters compatible with this project.

[OpenTelemetry registry]: https://opentelemetry.io/registry/?language=go&component=exporter
