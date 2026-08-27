# Experimental Metric SDK

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric/x)

This module is a proof of concept for bound metric instruments and explicit
series lifetime management. It is an alternative provider implementation, not
an add-on to `go.opentelemetry.io/otel/sdk/metric`.

The proof of concept implements `Int64Counter` with cumulative and delta Sum
aggregation. Other instrument kinds return no-op instruments. Stable SDK
Readers cannot be passed to this provider, and data from a stable provider is
not merged with this provider's collection.

Calling `Bind` performs attribute processing and cardinality selection once.
Recording without additional attributes then uses a direct aggregation handle.
Record-time attributes remain valid but use the normal lookup path. `Finish`
ends an exact concrete series and exports its value once. Shared overflow series
cannot be finished. A bound handle used after `Finish` lazily starts a new
series lifetime.

The bound Sum store is generated from `internal/shared/metricx`. Expanding the
template boundary to the full stable SDK runtime is intentionally deferred
until this vertical slice validates the API and lifecycle behavior.
See [DESIGN.md](./DESIGN.md) for the evaluated package boundaries, limitations,
and rejected alternatives.
