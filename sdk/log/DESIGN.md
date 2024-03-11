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
is defined as follows:

```go
type LoggerProvider struct {
	embedded.LoggerProvider
}

// NewLoggerProvider returns a new and configured LoggerProvider.
//
// By default, the returned LoggerProvider is configured with the default
// Resource and no Processors. Processors cannot be added after a LoggerProvider is
// created. This means the returned LoggerProvider, one created with no
// Processors, will perform no operations.
func NewLoggerProvider(...LoggerProviderOption) *LoggerProvider

// Logger returns a new log.Logger with the provided name and configuration.
//
// This method can be called concurrently.
//
// Logger implements the log.LoggerProvider interface.
func (*LoggerProvider) Logger(name string, options ...log.LoggerOption) log.Logger

type LoggerProviderOption interface { /* ... */ }

// WithResource associates a Resource with a LoggerProvider. This Resource
// represents the entity producing telemetry and is associated with all Loggers
// the LoggerProvider will create.
//
// By default, if this Option is not used, the default Resource from the
// go.opentelemetry.io/otel/sdk/resource package will be used.
func WithResource(*resource.Resource) LoggerProviderOption
```

## LogRecord limits

The [LogRecord limits](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecord-limits)
can be configured using following options:

```go
// WithAttributeCountLimit sets the maximum allowed log record attribute count.
// Any attribute added to a log record once this limit is reached will be dropped.
//
// Setting this to zero means no attributes will be recorded.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no limit 128 will be used.
func WithAttributeCountLimit(limit int) LoggerProviderOption

// AttributeValueLengthLimit sets the maximum allowed attribute value length.
//
// This limit only applies to string and string slice attribute values.
// Any string longer than this value will be truncated to this length.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no limit (-1) will be used.
func WithAttributeValueLengthLimit(limit int) LoggerProviderOption
```

The limits can be also configured using the `OTEL_LOGRECORD_*` environment variables as
[defined by the specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#logrecord-limits).

### Processor

The [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor)
is defined as follows:

```go
// WithProcessor associates Processor with a LoggerProvider.
//
// By default, if this option is not used, the LoggerProvider will perform no
// operations; no data will be exported without a processor.
//
// Each WithProcessor creates a separate pipeline. Use custom decorators
// for advanced scenarios such as enriching with attributes.
//
// Use NewBatchingProcessor to batch log records before they are exported.
// Use NewSimpleProcessor to synchronously export log records.
func WithProcessor(processor Processor) LoggerProviderOption

// Processor handles the processing of log records.
//
// Any of the Exporter's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Exporter to manage
// this concurrency.
type Processor interface {
	// OnEmit is called when a Record is emitted.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// All retry logic must be contained in this function. The SDK does not
	// implement any retry logic. All errors returned by this function are
	// considered unrecoverable and will be reported to a configured error
	// Handler.
	//
	// Before modifying a Record, the implementation must use Record.Clone
	// to create a copy that shares no state with the original.
	OnEmit(ctx context.Context, record Record) error

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the exporter should be done in this call.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	Shutdown(ctx context.Context) error

	// ForceFlush exports log records to the configured Exporter that have not yet
	// been exported.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	ForceFlush(ctx context.Context) error
}
```

The user can configure custom processors and decorate built-in processors.

### SimpleProcessor

The [Simple processor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#simple-processor)
is defined as follows:

```go
// SimpleProcessor implements Processor.
type SimpleProcessor struct { /* ... */ }

// NewBatchingProcessor decorates the provided exporter
// so that the log records are batched before exporting.
func NewSimpleProcessor(exporter Exporter) *SimpleProcessor
```

### BatchingProcessor

The [Batching processor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#batching-processor)
is defined as follows:

```go
// BatchingProcessor implements Processor.
type BatchingProcessor struct { /* ... */ }

// NewBatchingProcessor decorates the provided exporter
// so that the log records are batched before exporting.
func NewBatchingProcessor(exporter Exporter, opts ...BatchingOption) *BatchingProcessor 

// BatchingOption applies a configuration to a Batcher.
type BatchingOption interface { /* ... */ }

// WithMaxQueueSize sets the maximum queue size used by the Batcher.
// After the size is reached log records are dropped.
//
// If the OTEL_BLRP_MAX_QUEUE_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BLRP_MAX_QUEUE_SIZE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 2048 will be used.
// The default value is also used when the provided value is not a positive value.
func WithMaxQueueSize(max int) BatchingOption

// WithExportInterval sets the maximum duration between batched exports.
//
// If the OTEL_BSP_SCHEDULE_DELAY environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BSP_SCHEDULE_DELAY will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 1s will be used.
// The default value is also used when the provided value is not a positive value.
func WithExportInterval(d time.Duration) BatchingOption

// WithExportTimeout sets the duration after which a batched export is canceled.
//
// If the OTEL_BSP_EXPORT_TIMEOUT environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BSP_EXPORT_TIMEOUT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 30s will be used.
// The default value is also used when the provided value is not a positive value.
func WithExportTimeout(d time.Duration) BatchingOption

// WithExportMaxBatchSize sets the maximum batch size of every export.
//
// If the OTEL_BSP_MAX_EXPORT_BATCH_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BSP_MAX_EXPORT_BATCH_SIZE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 512 will be used.
// The default value is also used when the provided value is not a positive value.
func WithExportMaxBatchSize(max int) BatchingOption
```

The `Batcher` can be also configured using the `OTEL_BLRP_*` environment variables as
[defined by the specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#batch-logrecord-processor).

### Exporter

The [LogRecordExporter](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordexporter)
is defined as follows:

```go
// Exporter handles the delivery of log records to external receivers.
//
// Any of the Exporter's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Exporter to manage
// this concurrency.
type Exporter interface {
	// Export transmits log records to a receiver.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// All retry logic must be contained in this function. The SDK does not
	// implement any retry logic. All errors returned by this function are
	// considered unrecoverable and will be reported to a configured error
	// Handler.
	//
	// Implementations must not retain the records slice.
	//
	// Before modifying a Record, the implementation must use Record.Clone
	// to create a copy that shares no state with the original.
	Export(ctx context.Context, records []Record) error

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the exporter should be done in this call.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	//
	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	Shutdown(ctx context.Context) error

	// ForceFlush exports log records to the configured Exporter that have not yet
	// been exported.
	//
	// The deadline or cancellation of the passed context must be honored. An
	// appropriate error should be returned in these situations.
	ForceFlush(ctx context.Context) error
}
```

The slice passed to `Export` must not be retained by the implementation
(like e.g. [`io.Writer`](https://pkg.go.dev/io#Writer))
so that the caller can reuse the passed slice
(e.g. using [`sync.Pool`](https://pkg.go.dev/sync#Pool))
to avoid heap allocations on each call.

### Record

The [ReadWriteLogRecord](https://opentelemetry.io/docs/specs/otel/logs/sdk/#readwritelogrecord)
is defined as follows:

```go
type Record struct { /* ... */ }

func (r *Record) Timestamp()

func (r *Record) SetTimestamp(t time.Time)

func (r *Record) ObservedTimestamp() time.Time

func (r *Record) SetObservedTimestamp(t time.Time)

func (r *Record) Severity() log.Severity

func (r *Record) SetSeverity(level log.Severity)

func (r *Record) SeverityText() string

func (r *Record) SetSeverityText(text string)

func (r *Record) Body() log.Value

func (r *Record) SetBody(v log.Value)

func (r *Record) WalkAttributes(f func(log.KeyValue) bool)

func (r *Record) AddAttributes(attrs ...log.KeyValue)

// SetAttributes sets and overrides the attributes of the log record.
func (r *Record) SetAttributes(attrs ...log.KeyValue)

func (r *Record) TraceID() trace.TraceID

func (r *Record) SpanID() trace.SpanID

func (r *Record) TraceFlags() trace.TraceFlags

// Resource returns the entity that collected the log.
func (r *Record) Resource() resource.Resource

// InstrumentationScope returns the scope that the Logger was created with.
func (r *Record) InstrumentationScope() instrumentation.Scope

// AttributeValueLengthLimit is the maximum allowed attribute value length.
//
// This limit only applies to string and string slice attribute values.
// Any string longer than this value should be truncated to this length.
//
// Negative value means no limit should be applied.
func (r *Record) AttributeValueLengthLimit() int

// AttributeCountLimit is the maximum allowed log record attribute count. Any
// attribute added to a log record once this limit is reached should be dropped.
//
// Zero means no attributes should be recorded.
//
// Negative value means no limit should be applied.
func (r *Record) AttributeCountLimit() int

// Clone returns a copy of the record with no shared state. The original record
// and the clone can both be modified without interfering with each other.
func (r *Record) Clone() Record
```

The `Record` is designed similarly to [`log.Record`](https://pkg.go.dev/go.opentelemetry.io/otel/log#Record)
in order to reduce the number of heap allocations when processing attributes.

The SDK does have not have an additional definition of
[ReadableLogRecord](https://opentelemetry.io/docs/specs/otel/logs/sdk/#readablelogrecord)
as the specification does not say that the exporters must not be able to modify
the log records. It simply requires them to be able to read the log records.
Having less abstractions reduces the API surface and makes the design simpler.

## Benchmarking

The benchmarks are supposed to test end-to-end scenarios
and avoid I/O that could affect the stability of the results,

The benchmark results can be found in [the prototype](https://github.com/open-telemetry/opentelemetry-go/pull/4955).

## Rejected alternatives

### Represent both LogRecordProcessor and LogRecordExporter as Expoter

Because the [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor)
and the [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordexporter)
abstractions are so similar, there was a proposal to unify them under
single `Expoter` interface.[^2]

However, introducing a `Processor` interface makes it easier
to create custom processor decorators[^3]
and makes the design more aligned with the specifiation.

### Embedd log.Record

Because [`Record`](#record) and [`log.Record`](https://pkg.go.dev/go.opentelemetry.io/otel/log#Record)
are very similar, there was a proposal to embedd `log.Record` in `Record` definition.

[`log.Record`](https://pkg.go.dev/go.opentelemetry.io/otel/log#Record)
supports only adding attributes.
In the SDK, we also need to be able to modify the attributes (e.g. removal)
provided via API.

Moreover it is safer to have these abstraction decoupled.
E.g. there can be a need for some fields that can be set via API and cannot be modified by the processors.

## Open issues

The Logs SDK NOT be released as stable before all issues below are closed:

- [Redefine ReadableLogRecord and ReadWriteLogRecord](https://github.com/open-telemetry/opentelemetry-specification/pull/3898)
- [Fix what can be modified via ReadWriteLogRecord](https://github.com/open-telemetry/opentelemetry-specification/pull/3907)
- [logs: Allow duplicate keys](https://github.com/open-telemetry/opentelemetry-specification/issues/3931)
- [Add an Enabled method to Logger](https://github.com/open-telemetry/opentelemetry-specification/issues/3917)

[^1]: [OpenTelemetry Logging](https://opentelemetry.io/docs/specs/otel/logs)
[^2]: [Conversation on representing LogRecordProcessor and LogRecordExporter via a single Expoter interface](https://github.com/open-telemetry/opentelemetry-go/pull/4954#discussion_r1515050480)
[^3]: [Introduce Processor](https://github.com/pellared/opentelemetry-go/pull/9)
