# Logs Bridge API

Author: Robert Pająk

Tracking issue at [#4696](https://github.com/open-telemetry/opentelemetry-go/issues/4696).

## Abstract

<!-- A short summary of the proposal. -->

We propose adding a `go.opentelemetry.io/otel/log` Go module which will provide
[Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/).

## Background

<!-- An introduction of the necessary background and the problem being solved by the proposed change. -->

They key challenge is to create a well-performant API compliant with the specification.
Performance is seen as one of the most imporatant charactristics of logging libraries in Go.

## Design

This proposed design aims to:

- be specification compliant,
- have similar API to Trace and Metrics API,
- take advantage of both OpenTelemetry and `slog` experience to achieve acceptable performance.

### LoggerProvider

The [`LoggerProvider` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#loggerprovider)
is defined as an interface.

```go
type LoggerProvider interface{
    Logger(name string, options ...LoggerOption) Logger
}
```

### Logger

The [`Logger` abstraction](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/#logger)
is defined as an interface.

```go
type Logger interface{
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

The `NewRecord(...RecordOption) Record` is a factory function used to create records using provided options.

`Record` has a `Clone` method to allow copying records so that the SDK can offer concurrency safety.

## Compatibility

The backwards compatibility is achieved using the `embedded` design pattern
that is already used in Trace API and Metrics API.

## Benchmarking

The benchmarks takes inspiration from [`slog`](https://pkg.go.dev/log/slog),
because for the Go team it was also critical to create API that would be fast
and interoperable with existing logging packages. [^1] [^2]

## Rationale

<!-- A discussion of alternate approaches and the trade offs, advantages, and disadvantages of the specified approach. -->

## Implementation

<!-- A description of the steps in the implementation, who will do them, and when. -->

## Open issues (if applicable)

<!-- A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none. -->

[^1]: Jonathan Amsterdam, [The Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
[^2]: Jonathan Amsterdam, [GopherCon Europe 2023: A Fast Structured Logging Package](https://www.youtube.com/watch?v=tC4Jt3i62ns)
