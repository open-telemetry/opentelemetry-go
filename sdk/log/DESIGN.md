# Logs SDK

## Abstract

`go.opentelemetry.io/otel/sdk/log` provides Logs SDK compliant with the
[specification](https://opentelemetry.io/docs/specs/otel/logs/sdk/).

The main and recommended use case is to configure the SDK to use an OTLP
exporter with a batch processor.[^1] Therefore, the design aims to be
high-performant in this scenario.

The prototype was created in TODO.

## Module structure

The SDK is published as a single `go.opentelemetry.io/otel/sdk/log` Go module.

The Go module consists of the following packages:

- `go.opentelemetry.io/otel/sdk/log`
- `go.opentelemetry.io/otel/sdk/log/logtest`

The exporters are published as following Go modules:

- `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc`
- `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp`
- `go.opentelemetry.io/otel/exporters/stdout/stdoutlog`

## Rejected alternatives

## Open issues

The Logs SDK NOT be released as stable before all issues below are closed:

- TBD

[^1]: [OpenTelemetry Logging](https://opentelemetry.io/docs/specs/otel/logs)
