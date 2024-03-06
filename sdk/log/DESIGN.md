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

func NewLoggerProvider(...Option) *LoggerProvider

func (*LoggerProvider) Logger(name string, options ...log.LoggerOption) log.Logger

type Option interface { /* ... */ }

func WithResource(*resource.Resource) Option
```

## LogRecord limits

The [LogRecord limits](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecord-limits)
are defined as follows:

```go
func WithLimits(Limits) Option

type Limits struct {
	// AttributeValueLengthLimit is the maximum allowed attribute value length.
	//
	// This limit only applies to string and string slice attribute values.
	// Any string longer than this value will be truncated to this length.
	//
	// Setting this to zero means to use the default setting.
	// See also AttributeValueLengthLimitZero.
	//
	// Setting this to a negative value means no limit is applied.
	AttributeValueLengthLimit int

	// AttributeValueLengthLimitZero is to set a zero value of AttributeValueLengthLimit.
	// Setting this means string and string slice attribute values will be empty.
	// It overrides any value set via AttributeValueLengthLimit.
	AttributeValueLengthLimitZero bool

	// AttributeCountLimit is the maximum allowed span attribute count. Any
	// attribute added to a span once this limit is reached will be dropped.
	//
	// Setting this to zero means to use the default setting.
	// See also AttributeCountLimitZero.
	//
	// Setting this to a negative value means no limit is applied.
	AttributeCountLimit int

	// AttributeCountLimitZero is to set a zero value of AttributeCountLimit.
	// Setting this means no attributes will be recorded.
	// It overrides any value set via AttributeCountLimit.
	AttributeCountLimitZero bool
}
```

### LogRecordProcessor and LogRecordExporter  

Both [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor)
and [LogRecordExporter](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordexporter)
are defined via an `Exporter` interface:

```go
func WithExporter(Exporter) Option {
	return nil
}

type Exporter interface {
	Export(ctx context.Context, records []*Record) error
	Shutdown(ctx context.Context) error
	ForceFlush(ctx context.Context) error
}
```

The `Record` struct represents the [ReadWriteLogRecord](https://opentelemetry.io/docs/specs/otel/logs/sdk/#readwritelogrecord).

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

// SetAttributes sets and overrides the attributes of the log record.
func (r *Record) SetAttributes(attrs ...log.KeyValue)

func (r *Record) TraceID() trace.TraceID

func (r *Record) SpanID() trace.SpanID

func (r *Record) TraceFlags() trace.TraceFlags

func (r *Record) Resource() resource.Resource

func (r *Record) InstrumentationScope() instrumentation.Scope

func (r *Record) AttributeValueLengthLimit() int

func (r *Record) AttributeCountLimit() int
```

The user can implement a custom [LogRecordProcessor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logrecordprocessor)
by implementing a `Exporter` decorator.

This is similar to the design of HTTP server middleware
which is a wrapper of `http.Handler`.[^2]

[Simple processor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#simple-processor)
is a achieved by simply passing a bare-exporter.

[Batching processor](https://opentelemetry.io/docs/specs/otel/logs/sdk/#batching-processor)
is a achieved by wrapping an expoter with `Batcher`:

```go
type Batcher struct { /* ... */ }

func NewBatchingExporter(exporter Exporter, opts ...BatchingOption) *Batcher

func (b *Batcher) Export(ctx context.Context, records []*Record)

func (b *Batcher) Shutdown(ctx context.Context) error

func (b *Batcher) ForceFlush(ctx context.Context) error

type BatchingOption interface { /* ... */ }

func WithInterval(d time.Duration) BatchingOption

func WithTimeout(d time.Duration) BatchingOption

func WithQueueSize(max int) BatchingOption

func WithBatchSize(max int) BatchingOption
```

## Benchmarking

The benchmarks are supposed to test end-to-end scenarios.

However, in order avoid I/O that could affect the stability of the results,
the benchmarks are using an stdout exporter using `io.Discard`.

The benchmark results can be found in [the prototype](https://github.com/open-telemetry/opentelemetry-go/pull/4955).

## Rejected alternatives

## Open issues

The Logs SDK NOT be released as stable before all issues below are closed:

- [Fix what can be modified via ReadWriteLogRecord](https://github.com/open-telemetry/opentelemetry-specification/pull/3907)

[^1]: [OpenTelemetry Logging](https://opentelemetry.io/docs/specs/otel/logs)
[^2]: [Middleware - Go Web Examples](https://gowebexamples.com/basic-middleware/)
