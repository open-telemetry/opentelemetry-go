# Logs Bridge API

OpenTelemetry Logs tracking issue at [#4696](https://github.com/open-telemetry/opentelemetry-go/issues/4696).

## Abstract

We propose adding a `go.opentelemetry.io/otel/log` Go module which will provide
[Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/).

## Background

They key challenge is to create a well-performant API compliant with the specification.
Performance is seen as one of the most imporatant charactristics of logging libraries in Go.

## Design

This proposed design aims to:

- be specification compliant,
- have similar API to Trace and Metrics API,
- take advantage of both OpenTelemetry and `slog` experience to achieve acceptable performance.

### Module structure

The Go module consits of the following packages:

- `go.opentelemetry.io/otel/log`
- `go.opentelemetry.io/otel/log/embedded`
- `go.opentelemetry.io/otel/log/noop`

### LoggerProvider

The [`LoggerProvider` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#loggerprovider)
is defined as an interface.

```go
type LoggerProvider interface{
	embedded.LoggerProvider
    Logger(name string, options ...LoggerOption) Logger
}
```

### Logger

The [`Logger` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#logger)
is defined as an interface.

```go
type Logger interface{
	embedded.Logger
    Emit(ctx context.Context, options ...RecordOption)
}
```

The `Logger` has `Emit(context.Context, options ...RecordOption` method.

### Record

The [`LogRecord` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#logger)
is defined as a struct.

```go
type Record struct {
	Timestamp         time.Time
	ObservedTimestamp time.Time
	Severity          Severity
	SeverityText      string
	Body              string

	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of Attrs.
	front [nAttrsInline]attribute.KeyValue

	// The number of Attrs in front.
	nFront int

	// The list of Attrs except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []attribute.KeyValue
}

const nAttrsInline = 5

type Severity int

const (
	SeverityUndefined Severity = iota
	SeverityTrace
	SeverityTrace2
	SeverityTrace3
	SeverityTrace4
	SeverityDebug
	SeverityDebug2
	SeverityDebug3
	SeverityDebug4
	SeverityInfo
	SeverityInfo2
	SeverityInfo3
	SeverityInfo4
	SeverityWarn
	SeverityWarn2
	SeverityWarn3
	SeverityWarn4
	SeverityError
	SeverityError2
	SeverityError3
	SeverityError4
	SeverityFatal
	SeverityFatal2
	SeverityFatal3
	SeverityFatal4
)
```

`Record` has `Attributes` and `AddAttributes` methods,
like [`slog.Record.Attrs`](https://pkg.go.dev/log/slog#Record.Attrs)
and [`slog.Record.AddAttrs`](https://pkg.go.dev/log/slog#Record.AddAttrs),
in order to achieve high-performance when accessing and setting attributes efficiently.

The `NewRecord(...RecordOption) (Record, error)` is a factory function
used to create records using provided options.

`Record` has a `Clone` method to allow copying records
so that the SDK can offer concurrency safety.

## Compatibility

The backwards compatibility is achieved using the `embedded` design pattern
that is already used in Trace API and Metrics API.

## Benchmarking

The benchmarks takes inspiration from [`slog`](https://pkg.go.dev/log/slog),
because for the Go team it was also critical to create API that would be fast
and interoperable with existing logging packages.[^1]

## Rationale

### Logger.Emit definition

One of the ideas was to have:

```go
type Logger interface{
    Emit(ctx context.Context, record Record)
}
```

This gives the advantage that the SDK would not need to call `NewRecord(options...)`.

The user can still easily create a helper that could be easier to use:

```go
func log(ctx context.Context, l Logger, options ...RecordOption) {
    r := log.NewRecord(options...)
    l.Emit(ctx, r)
}
```

The main reasons against this defintion are that following:

1. The existing design is similar to the [Meter API](https://pkg.go.dev/go.opentelemetry.io/otel/metric#Meter)
for creating instruments.
2. It is unsure if anyone would like to reuse a record.
3. Just passing options should be more-user friendly API.

### Record as struct

`Record` is defined as a `struct` because of the following reasons.

Log record is a value object without any behavior.
It is used as data input for Logger methods.

The log record resembles the instrument config structs like [metric.Float64CounterConfig](https://pkg.go.dev/go.opentelemetry.io/otel/metric#Float64CounterConfig).

Using `struct` instead of `interface` should have better the performance as e.g.
indirect calls are less optimized,
usage of intefaces tend to increase heap allocations.[^2]

The `Record` design is inspired by [`slog.Record`](https://pkg.go.dev/log/slog#Record).

## Open issues (if applicable)

<!-- A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none. -->

[^1]: Jonathan Amsterdam, [The Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
[^2]: Jonathan Amsterdam, [GopherCon Europe 2023: A Fast Structured Logging Package](https://www.youtube.com/watch?v=tC4Jt3i62ns)
