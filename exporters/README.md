# Exporters

Included in this directory are exporters that export both metric and trace telemetry.

- [stdout](./stdout): Writes telemetry to a specified local output as structured JSON.
- [otlp](./otlp): Sends telemetry to an OpenTelemetry collector as OTLP.

Additionally, there are [metric](./metric) and [trace](./trace) only exporters.

## Metric Telemetry Only

- [prometheus](./metric/prometheus): Exposes metric telemetry as Prometheus metrics.

## Trace Telemetry Only

- [jaeger](./trace/jaeger): Sends properly transformed trace telemetry to a Jaeger endpoint.
- [zipkin](./trace/zipkin): Sends properly transformed trace telemetry to a Zipkin endpoint.
