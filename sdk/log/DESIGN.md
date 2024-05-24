# Logs SDK

## Abstract

`go.opentelemetry.io/otel/sdk/log` provides Logs SDK compliant with the
[specification](https://opentelemetry.io/docs/specs/otel/logs/sdk/).

The main and recommended use case is to configure the SDK to use an OTLP
exporter with a batch processor.[^1] Therefore, the design aims to be
high-performant in this scenario.

The prototype was created in
[#4955](https://github.com/open-telemetry/opentelemetry-go/pull/4955).

## Modules structure

The SDK is published as a single `go.opentelemetry.io/otel/sdk/log` Go module.

The exporters are going to be published as following Go modules:

- `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc`
- `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp`
- `go.opentelemetry.io/otel/exporters/stdout/stdoutlog`

## LoggerProvider

The [LoggerProvider](https://opentelemetry.io/docs/specs/otel/logs/sdk/#loggerprovider)
is implemented as `LoggerProvider` struct in [provider.go](provider.go).

## LogRecord limits

The [LogRecord limits](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecord-limits)
can be configured using following options:

```go
func WithAttributeCountLimit(limit int) LoggerProviderOption
func WithAttributeValueLengthLimit(limit int) LoggerProviderOption
```

The limits can be also configured using the `OTEL_LOGRECORD_*` environment variables as
[defined by the specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#logrecord-limits).

### Processor

The [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor)
is defined as `Processor` interface in [processor.go](processor.go).

The user set processors for the `LoggerProvider` using
`func WithProcessor(processor Processor) LoggerProviderOption`.

The user can configure custom processors and decorate built-in processors.

The specification may add new operations to the
[LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor).
If it happens, [CONTRIBUTING.md](../../CONTRIBUTING.md#how-to-change-other-interfaces)
describes how the SDK can be extended in a backwards-compatible way.

### SimpleProcessor

The [Simple processor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#simple-processor)
is implemented as `SimpleProcessor` struct in [simple.go](simple.go).

### BatchProcessor

The [Batching processor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#batching-processor)
is implemented as `BatchProcessor` struct in [batch.go](batch.go).

The `Batcher` can be also configured using the `OTEL_BLRP_*` environment variables as
[defined by the specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#batch-logrecord-processor).

### Exporter

The [LogRecordExporter](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordexporter)
is defined as `Exporter` interface in [exporter.go](exporter.go).

The slice passed to `Export` must not be retained by the implementation
(like e.g. [`io.Writer`](https://pkg.go.dev/io#Writer))
so that the caller can reuse the passed slice
(e.g. using [`sync.Pool`](https://pkg.go.dev/sync#Pool))
to avoid heap allocations on each call.

The specification may add new operations to the
[LogRecordExporter](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordexporter).
If it happens, [CONTRIBUTING.md](../../CONTRIBUTING.md#how-to-change-other-interfaces)
describes how the SDK can be extended in a backwards-compatible way.

### Record

The [ReadWriteLogRecord](https://opentelemetry.io/docs/specs/otel/logs/sdk/#readwritelogrecord)
is defined as `Record` struct in [record.go](record.go).

The `Record` is designed similarly to [`log.Record`](https://pkg.go.dev/go.opentelemetry.io/otel/log#Record)
in order to reduce the number of heap allocations when processing attributes.

The SDK does not have have an additional definition of
[ReadableLogRecord](https://opentelemetry.io/docs/specs/otel/logs/sdk/#readablelogrecord)
as the specification does not say that the exporters must not be able to modify
the log records. It simply requires them to be able to read the log records.
Having less abstractions reduces the API surface and makes the design simpler.

## Benchmarking

The benchmarks are supposed to test end-to-end scenarios
and avoid I/O that could affect the stability of the results.

The benchmark results can be found in [the prototype](https://github.com/open-telemetry/opentelemetry-go/pull/4955).

## Rejected alternatives

### Represent both LogRecordProcessor and LogRecordExporter as Expoter

Because the [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor)
and the [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordexporter)
abstractions are so similar, there was a proposal to unify them under
single `Expoter` interface.[^2]

However, introducing a `Processor` interface makes it easier
to create custom processor decorators[^3]
and makes the design more aligned with the specification.

### Embed log.Record

Because [`Record`](#record) and [`log.Record`](https://pkg.go.dev/go.opentelemetry.io/otel/log#Record)
are very similar, there was a proposal to embed `log.Record` in `Record` definition.

[`log.Record`](https://pkg.go.dev/go.opentelemetry.io/otel/log#Record)
supports only adding attributes.
In the SDK, we also need to be able to modify the attributes (e.g. removal)
provided via API.

Moreover it is safer to have these abstraction decoupled.
E.g. there can be a need for some fields that can be set via API and cannot be modified by the processors.

[^1]: [OpenTelemetry Logging](https://opentelemetry.io/docs/specs/otel/logs)
[^2]: [Conversation on representing LogRecordProcessor and LogRecordExporter via a single Expoter interface](https://github.com/open-telemetry/opentelemetry-go/pull/4954#discussion_r1515050480)
[^3]: [Introduce Processor](https://github.com/pellared/opentelemetry-go/pull/9)
