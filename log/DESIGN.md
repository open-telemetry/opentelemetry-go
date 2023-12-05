# Logs Bridge API

OpenTelemetry Logs tracking issue at [#4696](https://github.com/open-telemetry/opentelemetry-go/issues/4696).

## Abstract

We propose adding a `go.opentelemetry.io/otel/log` Go module which will provide
[Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/).

## Background

The key challenge is to create a well-performant API compliant with the specification.
Performance is seen as one of the most important characteristics of logging libraries in Go.

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
is defined as an interface [provider.go](provider.go).

### Logger

The [`Logger` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#logger)
is defined as an interface in [logger.go](logger.go).

### Record

The [`LogRecord` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#logger)
is defined as a struct in [record.go](record.go).

`Record` has `WalkAttributes` and `AddAttributes` methods,
like [`slog.Record.Attrs`](https://pkg.go.dev/log/slog#Record.Attrs)
and [`slog.Record.AddAttrs`](https://pkg.go.dev/log/slog#Record.AddAttrs),
in order to achieve high-performance when accessing and setting attributes efficiently.

`Record` has a `AttributesLen` method that returns
the number of attributes to allow slice preallocation
when converting records to a different representation.

### Usage Example: Log Bridge implementation

A naive implementation of
the [slog.Handler](https://pkg.go.dev/log/slog#Handler) interface
is in [benchmark/slog_test.go](benchmark/slog_test.go).

A naive implementation of
the [logr.LogSink](https://pkg.go.dev/github.com/go-logr/logr#LogSink) interface
is in [benchmark/logr_test.go](benchmark/slog_test.go).

### Usage Example: Direct API usage

The users may also chose to use the API directly.

```go
package app

var logger = otel.Logger("my-service")

// In some function:
logger.Emit(ctx, Record{Severity: log.SeverityInfo, Body: "Application started."})
```

### Usage Example: API implementation

Excerpt of how SDK can implement the `Logger` interface.

```go
type Logger struct {
	scope     instrumentation.Scope
	processor Processor
}

func (l *Logger) Emit(ctx context.Context, r log.Record) {
	// Create log record model.
	record, err := toModel(r)
	if err != nil {
		otel.Handle(err)
		return
	}
	l.processor.Process(ctx, record)
}
```

A test implementation of the the `Logger` interface
is in [benchmark/writer_logger_test.go](benchmark/writer_logger_test.go).

Canceling the context should not affect record processing.
Among other things, log messages may be necessary to debug a
cancellation-related problem.
The context is used to pass request-scoped values such as Trace ID and Span ID.

## Compatibility

The backwards compatibility is achieved using the `embedded` design pattern
that is already used in Trace API and Metrics API.

Additionally, the `Logger.Emit` functionality can be extended by
adding new exported fields and methods to the `Record` struct.

## Benchmarking

The benchmarks takes inspiration from [`slog`](https://pkg.go.dev/log/slog),
because for the Go team it was also critical to create API that would be fast
and interoperable with existing logging packages.[^1][^2]

## Rationale

### Rejected Alternative: Reuse slog

The API must not be coupled to [`slog`](https://pkg.go.dev/log/slog),
nor any other logging library.

The API needs to evolve orthogonally to `slog`.

`slog` is not compliant with the [Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/).
and we cannot expect the Go team to make `slog` compliant with it.

The interoperabilty can be achieved using [a log bridge](https://opentelemetry.io/docs/specs/otel/glossary/#log-appender--bridge).

You can read more about OpenTelemetry Logs design on [opentelemetry.io](https://opentelemetry.io/docs/concepts/signals/logs/).

### Rejected Alternative: Record as interface

`Record` is defined as a `struct` because of the following reasons.

Log record is a value object without any behavior.
It is used as data input for Logger methods.

The log record resembles the instrument config structs like [metric.Float64CounterConfig](https://pkg.go.dev/go.opentelemetry.io/otel/metric#Float64CounterConfig).

Using `struct` instead of `interface` should have better the performance as e.g.
indirect calls are less optimized,
usage of interfaces tend to increase heap allocations.[^2]

The `Record` design is inspired by [`slog.Record`](https://pkg.go.dev/log/slog#Record).

### Rejected Alternative: Options as parameter to Logger.Emit

One of the initial ideas was to have:

```go
type Logger interface{
	embedded.Logger
	Emit(ctx context.Context, options ...RecordOption)
}
```

The main reason was that design would be similar
to the [Meter API](https://pkg.go.dev/go.opentelemetry.io/otel/metric#Meter)
for creating instruments.

However, passing `Record` directly, instead of using options,
is more performant as it reduces heap allocations.[^3]

Another advantage of passing `Record` is that API would not have functions like `NewRecord(options...)`,
which would be used by the SDK and not by the users.

At last, the definition would be similar to [`slog.Handler.Handle`](https://pkg.go.dev/log/slog#Handler)
that was designed to provide optimization opportunities.[^1]

### Rejected Alternative: Passing record as pointer to Logger.Emit

So far the benchmarks do not show differences that would
favor passing the record via pointer (and vice versa).

Passing via value feels safer because of the following reasons.

It follows the design of [`slog.Handler`](https://pkg.go.dev/log/slog#Handler).

It should reduce the possibility of a heap allocation.

The user would not be able to pass `nil`.
Therefore, it reduces the possiblity to have a nil pointer dereference.

### Rejected Alternative: Passing struct as parameter to LoggerProvider.Logger

Similarly to `Logger.Emit`, we could have something like:

```go
type Logger interface{
	embedded.Logger
	Logger(name context.Context, config LoggerConfig)
}
```

The drawback of this idea would be that this would be
a different design from Trace and Metrics API.

The performance of acquiring a logger is not as critical
as the performance of emitting a log record. While a single
HTTP/RPC handler could write hundreds of logs, it should not
create a new logger for each log entry.
The application should reuse loggers whenever possible.

### Rejected Proposal: Export attributesInlineCount

There was a proposal to export `attributesInlineCount`
so the bridge implementation could use it
to reduce the number of heap allocations
when the record has more attribute than 5
(the value of `attributesInlineCount`).

However, according to [^1], only ~5% of code emits log records
with more than 5 attributes.
Moreover, according to
[the benchmarks](https://github.com/open-telemetry/opentelemetry-go/pull/4725#discussion_r1413884476),
it would only save a few allocations
when the number of attributes is greater than 5
and the time execution tend to be slower for 5 attributes or less.

At last, nothing prevents us to export this constant in future
if it will occur that it could be helpful in some scenarios.
However, without a strong reason, we prefer to hide the implementation detail
and have smaller API surface.

## Open issues (if applicable)

<!-- A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none. -->

[^1]: Jonathan Amsterdam, [The Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
[^2]: Jonathan Amsterdam, [GopherCon Europe 2023: A Fast Structured Logging Package](https://www.youtube.com/watch?v=tC4Jt3i62ns)
[^3]: [Emit definition discussion with benchmarks](https://github.com/open-telemetry/opentelemetry-go/pull/4725#discussion_r1400869566)
