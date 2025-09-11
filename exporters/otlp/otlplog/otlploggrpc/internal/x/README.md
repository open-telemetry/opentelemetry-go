Experimental Features

The `otlploggrpc` exporter contains features that have not yet stabilized in the OpenTelemetry specification.
These features are added to the `otlploggrpc` exporter prior to stabilization in the specification so that users can start experimenting with them and provide feedback.

These features may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Observability](#observability)

### Observability

The `otlploggrpc` exporter provides a self-observability feature that allows you to monitor the exporter itself.



To opt-in, set the environment variable `OTEL_GO_X_OBSERVABILITY` to `true`.

When enabled, the exporter will create the following metrics using the global `MeterProvider`:

- `otel.sdk.exporter.log.inflight`
- `otel.sdk.exporter.log.exported`
- `otel.sdk.exporter.operation.duration`

Please see the [Semantic conventions for OpenTelemetry SDK metrics] documentation for more details on these metrics.