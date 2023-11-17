# Logs Bridge API Design

Author: Robert Pająk

Tracking issue at [#4696](https://github.com/open-telemetry/opentelemetry-go/issues/4696).

## Abstract

<!-- A short summary of the proposal. -->

We propose adding a `go.opentelemetry.io/otel/log` Go module which will provide
[Logs Data Model](https://opentelemetry.io/docs/specs/otel/logs/data-model/)
and [Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/data-model/).

## Background

<!-- An introduction of the necessary background and the problem being solved by the proposed change. -->

They key challenge is to create API which will be complaint with the specification
and be as performant as possible.

Performance is seen as one of the most imporatant charactristics of logging libraries in Go.

## Proposal

The design and benchmarks takes inspiration from [`slog`](https://pkg.go.dev/log/slog),
because for the Go team it was also critical to create API that would be fast
and interoperable with existing logging packages. [^1] [^2]

<!-- A precise statement of the proposed change. -->

## Rationale

<!-- A discussion of alternate approaches and the trade offs, advantages, and disadvantages of the specified approach. -->

## Compatibility

The backwards compatibility is achieved using the `embedded` design pattern
that is already used in Trace API and Metrics API.

## Implementation

<!-- A description of the steps in the implementation, who will do them, and when. -->

## Open issues (if applicable)

<!-- A discussion of issues relating to this proposal for which the author does not
know the solution. This section may be omitted if there are none. -->

[^1]: Jonathan Amsterdam, [The Go Blog: Structured Logging with slog](https://go.dev/blog/slog)
[^2]: Jonathan Amsterdam, [GopherCon Europe 2023: A Fast Structured Logging Package](https://www.youtube.com/watch?v=tC4Jt3i62ns)
