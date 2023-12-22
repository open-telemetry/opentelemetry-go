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

## Usage examples

### Log Bridge implementation

A naive implementation of
the [slog.Handler](https://pkg.go.dev/log/slog#Handler) interface
is in [benchmark/slog_test.go](benchmark/slog_test.go).

A naive implementation of
the [logr.LogSink](https://pkg.go.dev/github.com/go-logr/logr#LogSink) interface
is in [benchmark/logr_test.go](benchmark/slog_test.go).

The log bridges can use [`sync.Pool`](https://pkg.go.dev/sync#Pool)
for reducing the number of allocations when mapping attributes.

### Direct API usage

The users may also chose to use the API directly.

```go
package app

var logger = otel.Logger("my-service")

// In some function:
logger.Emit(ctx, Record{Severity: log.SeverityInfo, Body: "Application started."})
```

### API implementation

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

If the record is processed asynchronously,
then the processor has to copy record attributes,
in order to avoid use after free bugs and race condition.

Canceling the context should not affect record processing.
Among other things, log messages may be necessary to debug a
cancellation-related problem.
The context is used to pass request-scoped values such as Trace ID and Span ID.

## Compatibility

The backwards compatibility is achieved using the `embedded` design pattern
that is already used in Trace API and Metrics API.

Additionally, the `Logger.Emit` functionality can be extended by
adding new exported fields to the `Record` struct.

## Benchmarking

The benchmarks takes inspiration from [`slog`](https://pkg.go.dev/log/slog),
because for the Go team it was also critical to create API that would be fast
and interoperable with existing logging packages.[^1][^2]

## Rejected Alternatives

### Reuse slog

The API must not be coupled to [`slog`](https://pkg.go.dev/log/slog),
nor any other logging library.

The API needs to evolve orthogonally to `slog`.

`slog` is not compliant with the [Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/).
and we cannot expect the Go team to make `slog` compliant with it.

The interoperabilty can be achieved using [a log bridge](https://opentelemetry.io/docs/specs/otel/glossary/#log-appender--bridge).

You can read more about OpenTelemetry Logs design on [opentelemetry.io](https://opentelemetry.io/docs/concepts/signals/logs/).

### Record as interface

`Record` is defined as a `struct` because of the following reasons.

Log record is a value object without any behavior.
It is used as data input for Logger methods.

The log record resembles the instrument config structs like [metric.Float64CounterConfig](https://pkg.go.dev/go.opentelemetry.io/otel/metric#Float64CounterConfig).

Using `struct` instead of `interface` should have better the performance as e.g.
indirect calls are less optimized,
usage of interfaces tend to increase heap allocations.[^2]

The `Record` design is inspired by [`slog.Record`](https://pkg.go.dev/log/slog#Record).

### Options as parameter to Logger.Emit

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

### Passing record as pointer to Logger.Emit

So far the benchmarks do not show differences that would
favor passing the record via pointer (and vice versa).

Passing via value feels safer because of the following reasons.

It follows the design of [`slog.Handler`](https://pkg.go.dev/log/slog#Handler).

It should reduce the possibility of a heap allocation.

The user would not be able to pass `nil`.
Therefore, it reduces the possiblity to have a nil pointer dereference.

### Passing struct as parameter to LoggerProvider.Logger

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

### Logger.WithAttributes

We could add `WithAttributes` to the `Logger` interface.
Then `Record` could be a simple struct with only exported fields.
The idea was that the SDK would implement the performance improvements
instead of doing it in the API.
This would allow having different optimisation strategies.

During the analysis[^4], it occurred that the main problem of this proposal
is that the variadic slice passed to an interface method is always heap allocated.

Moreover, the logger returned by `WithAttribute` was allocated on the heap.

At last, the proposal was not specification compliant.

### Record attributes like in slog.Record

To reduce the number of allocations of the attributes,
the `Record` could be modeled similarly to [`slog.Record`](https://pkg.go.dev/log/slog#Record).
`Record` could have `WalkAttributes` and `AddAttributes` methods,
like [`slog.Record.Attrs`](https://pkg.go.dev/log/slog#Record.Attrs)
and [`slog.Record.AddAttrs`](https://pkg.go.dev/log/slog#Record.AddAttrs),
in order to achieve high-performance when accessing and setting attributes efficiently.
`Record` would have a `AttributesLen` method that returns
the number of attributes to allow slice preallocation
when converting records to a different representation.

However, during the analysis[^5] we decided that having
a simple slice in `Record` is more flexible.

It is possible to achieve better performance, by using [`sync.Pool`](https://pkg.go.dev/sync#Pool).

Having a simple `Record` without any logic makes it possible
that the optimisations can be done in API implementation
and bridge implementations.
For instance, in order to reduce the heap allocations of attributes,
the bridge implementation can use a `sync.Pool`.
In such case, the API implementation (SDK) would need to copy the attributes
when the records are processed asynchrounsly,
in order to avoid use after free bugs and race conditions.

For reference, here is the reason why `slog` does not use `sync.Pool`[^2]:

> We can use a sync pool for records though we decided not to.
You can but it's a bad idea for us. Why?
Because users have control of Records.
Handler writers can get their hands on a record
and we'd have to ask them to free it
or try to free it magically at some some point.
But either way, they could get themselves in trouble by freeing it twice
or holding on to one after they free it.
That's a use after free bug and that's why `zerolog` was problematic for us.
`zerolog` as as part of its speed exposes a pool allocated value to users
if you use `zerolog` the normal way, that you'll see in all the examples,
you will never encounter a problem.
But if you do something a little out of the ordinary you can get
use after free bugs and we just didn't want to put that in the standard library.

We took a different decision, because the key difference is that `slog`
is a logging library and Logs Bridge API is only a logging abstraction.
We want to provide more flexibility and offer better speed.

## Open issues (if applicable)

<!-- A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none. -->

[^1]: Jonathan Amsterdam, [The Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
[^2]: Jonathan Amsterdam, [GopherCon Europe 2023: A Fast Structured Logging Package](https://www.youtube.com/watch?v=tC4Jt3i62ns)
[^3]: [Emit definition discussion with benchmarks](https://github.com/open-telemetry/opentelemetry-go/pull/4725#discussion_r1400869566)
[^4]: [Logger.WithAttributes analysis](https://github.com/pellared/opentelemetry-go/pull/3)
[^5]: [Record attributes as field and use sync.Pool for reducing allocations analysis](https://github.com/pellared/opentelemetry-go/pull/4)
