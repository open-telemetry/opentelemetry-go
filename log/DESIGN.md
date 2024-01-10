# Logs Bridge API

## Abstract

`go.opentelemetry.io/otel/log` provides
[Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/).

The initial version of the design and the prototype
was created in [#4725](https://github.com/open-telemetry/opentelemetry-go/pull/4725).

## Background

The key challenge is to create a performant API compliant with the [specification](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/)
with an intuitive and user friendly design.
Performance is seen as one of the most important characteristics of logging libraries in Go.

## Design

This proposed design aims to:

- be specification compliant,
- be similar to Trace and Metrics API,
- take advantage of both OpenTelemetry and `slog` experience to achieve acceptable performance.

### Module structure

The API is published as a single `go.opentelemetry.io/otel/log` Go module.

The module name is compliant with
[Artifact Naming](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/bridge-api.md#artifact-naming)
and the package structure is the same as for Trace API and Metrics API.

The Go module consits of the following packages:

- `go.opentelemetry.io/otel/log`
- `go.opentelemetry.io/otel/log/embedded`
- `go.opentelemetry.io/otel/log/noop`

Rejected alternative:

- [Reuse slog](#reuse-slog)

### LoggerProvider

The [`LoggerProvider` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#loggerprovider)
is defined as an interface:

```go
type LoggerProvider interface {
	embedded.LoggerProvider
	Logger(name string, options ...LoggerOption) Logger
}
```

The specification may add new operations to `LoggerProvider`.
The interface may have methods added without a package major version bump.
This embeds `embedded.LoggerProvider` to help inform an API implementation
author about this non-standard API evolution.
This approach is already used in Trace API and Metrics API.

#### LoggerProvider.Logger

The `Logger` method implements the [`Get a Logger` operation](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#get-a-logger).

The required `name` parameter is accepted as a `string` method argument.

The following options are defined to support optional parameters:

```go
func WithInstrumentationVersion(version string) LoggerOption
func WithInstrumentationAttributes(attr ...attribute.KeyValue) LoggerOption
func WithSchemaURL(schemaURL string) LoggerOption
```

Implementation requirements:

- The [specification requires](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#concurrency-requirements)
  the method to be safe to be called concurrently.

- The method should use some default name if the passed name is empty
  in order to meet the [specification's SDK requirement](https://opentelemetry.io/docs/specs/otel/logs/sdk/#logger-creation)
  to return a working logger when an invalid name is passed
  as well as to resemble the behavior of getting tracers and meters.

`Logger` can be extended by adding new `LoggerOption` options
and adding new exported fields to the `LoggerConfig` struct.
This design is already used in Trace API for getting tracers
and in Metrics API for getting meters.

Rejected alternative:

- [Passing struct as parameter to LoggerProvider.Logger](#passing-struct-as-parameter-to-loggerproviderlogger).

### Logger

The [`Logger` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#logger)
is defined as an interface:

```go
type Logger interface {
	embedded.Logger
	Emit(ctx context.Context, record Record)
}
```

The specification may add new operations to `Logger`.
The interface may have methods added without a package major version bump.
This embeds `embedded.Logger` to help inform an API implementation
author about this non-standard API evolution.
This approach is already used in Trace API and Metrics API.

### Logger.Emit

The `Emit` method implements the [`Emit a LogRecord` operation](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#emit-a-logrecord).

[`Context` associated with the `LogRecord`](https://opentelemetry.io/docs/specs/otel/context/)
is accepted as a `context.Context` method argument.

Calls to `Emit` are supposed to be on the hot path.
Therefore, in order to reduce the number of heap allocations,
the [`LogRecord` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#emit-a-logrecord),
is defined as a `struct`:

```go
type Record struct {
	Timestamp         time.Time
	ObservedTimestamp time.Time
	Severity          Severity
	SeverityText      string
	Body              string
	Attributes        []attribute.KeyValue
}
```

[`Timestamp`](https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-timestamp)
is defined as `time.Time` type.

[`ObservedTimestamp`](https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-observedtimestamp)
is defined as `time.Time` type.

[`SeverityNumber`](https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-severitynumber)
is defined as a type and constants:

```go
type Severity int

const (
	SeverityTrace Severity = iota + 1
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

[`SeverityText`](https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-severitytext)
is defined as `string` type.

[`Body`](https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-body)
is defined as `string` type as the specification says:

> First-party Applications SHOULD use a string message.

[Log record attributes](https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-attributes)
are defined a regular slice of `attribute.KeyValue`.
The users can use [`sync.Pool`](https://pkg.go.dev/sync#Pool)
for reducing the number of allocations when passing attributes.

Implementation requirements:

- The [specification requires](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#concurrency-requirements)
  the method to be safe to be called concurrently.

- The method should not interrupt the record processing if the context is canceled
  per ["ignoring context cancellation" guideline](../CONTRIBUTING.md#ignoring-context-cancellation).

- The method handle the trace context passed via `ctx` argument in order to meet the
  [specification's SDK requirement](https://opentelemetry.io/docs/specs/otel/logs/sdk/#readablelogrecord)
  to populate the trace context fields from the resolved context.

- The method should not modify the record's attribute as it may be reused by caller.

- The method should copy the record's attributes in case of asynchronous processing
  as otherwise it may lead to data races as the user may modify the passed attributes
  e.g. because of using `sync.Pool` to reduce the number of allocation.

- The [specification requires](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#emit-a-logrecord)
  use the current time as observed timestamp if the passed is empty.

`Emit` can be extended by adding new exported fields to the `Record` struct.

Rejected alternatives:

- [Record as interface](#record-as-interface)
- [Options as parameter to Logger.Emit](#options-as-parameter-to-loggeremit)
- [Passing record as pointer to Logger.Emit](#passing-record-as-pointer-to-loggeremit)
- [Logger.WithAttributes](#loggerwithattributes)
- [Record attributes like in slog.Record](#record-attributes-like-in-slogrecord)
- [Record.Body as any](#recordbody-as-any)

### noop package

The `go.opentelemetry.io/otel/log/noop` package provides
[Logs Bridge API No-Op Implementation](https://opentelemetry.io/docs/specs/otel/logs/noop/).

### Trace context correlation

The bridge implementation should do its best to pass
the `ctx` containing the trace context from the caller
so it can later be passed via `Logger.Emit`.
Re-constructing a `context.Context` with [`trace.ContextWithSpanContext`](https://pkg.go.dev/go.opentelemetry.io/otel/trace#ContextWithSpanContext)
and [`trace.NewSpanContext`](https://pkg.go.dev/go.opentelemetry.io/otel/trace#NewSpanContext)
would usually involve more memory allocations.

The logging libraries which have recording methods that accepts `context.Context`,
such us [`slog`](https://pkg.go.dev/log/slog),
[`logrus`](https://pkg.go.dev/github.com/sirupsen/logrus)
[`zerolog`](https://pkg.go.dev/github.com/rs/zerolog),
makes passing the trace context trivial.

However, some libraries do not accept a `context.Context` in their recording methods.
Structured logging libraries,
such as [`logr`](https://pkg.go.dev/github.com/go-logr/logr)
and [`zap`](https://pkg.go.dev/go.uber.org/zap),
offer passing `any` type as a log attribute/field.
Therefore, their bridge implementations can define a "special" log attributes/field
that will be used to capture the trace context.

## Benchmarking

The benchmarks take inspiration from [`slog`](https://pkg.go.dev/log/slog),
because for the Go team it was also critical to create API that would be fast
and interoperable with existing logging packages.[^1][^2]

The benchmark results can be found in [the prototype](https://github.com/open-telemetry/opentelemetry-go/pull/4725).

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

Using `struct` instead of `interface` improves the performance as e.g.
indirect calls are less optimized,
usage of interfaces tend to increase heap allocations.[^2]

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
type LoggerProvider interface{
	embedded.LoggerProvider
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
This would allow having different optimization strategies.

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
when the records are processed asynchronously,
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
Unlike `zerolog`, we are not using a `sync.Pool` in the logger implementation,
which would expose to the issue described above.
Instead, `sync.Pool` is used by bridges, which are the users,
rather than implementers, of the API.

### Record.Body as any

[Logs Data Model](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/logs/data-model.md#field-body)
defines Body to be `any`.
We could define `Body` as `any` instead of `string`.
This could result in a more flexible API
that could support logging libraries which model log records as any objects.

However, using `any` decreases the performance.[^6]

Additionally, we are not aware of any popular Go logging library
that uses any objects as log records.

At last, the specification recommends using `string` in Bridge API:

> First-party Applications SHOULD use a string message.

As according to [the specification](https://opentelemetry.io/docs/specs/otel/logs/#new-first-party-application-logs)
First-party Applications are supposed to use the log bridges
that integrates with the SDK through Bridge API.

We can always add an additional `StructuredBody any` field in a future release
if we would need to support structured bodies. This approach would be
backwards-compatible and should not have such a negative impact on performance
for the most common scenario where `Body` is a `string`.

[^1]: Jonathan Amsterdam, [The Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
[^2]: Jonathan Amsterdam, [GopherCon Europe 2023: A Fast Structured Logging Package](https://www.youtube.com/watch?v=tC4Jt3i62ns)
[^3]: [Emit definition discussion with benchmarks](https://github.com/open-telemetry/opentelemetry-go/pull/4725#discussion_r1400869566)
[^4]: [Logger.WithAttributes analysis](https://github.com/pellared/opentelemetry-go/pull/3)
[^5]: [Record attributes as field and use sync.Pool for reducing allocations analysis](https://github.com/pellared/opentelemetry-go/pull/4)
[^6]: [Record.Body as any](https://github.com/pellared/opentelemetry-go/pull/5)
