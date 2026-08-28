# Experimental Metric SDK

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric/x)

This module is a proof of concept for bound metric instruments and explicit
series lifetime management. It is a thin experimental facade over
`go.opentelemetry.io/otel/sdk/metric`.

The stable provider uses its regular readers, views, pipelines, exporters, and
aggregation state. Enable binding, finishing, or both independently by passing
`sdkmetricx.WithBinding()` and `sdkmetricx.WithFinish()` to
`sdkmetric.NewMeterProvider`. The options currently affect only `Int64Counter`;
their instrument-agnostic names leave room for additional synchronous
instrument support. This module does not re-export the stable SDK surface.

Calling `Bind` performs attribute processing and cardinality selection once.
Recording without additional attributes then uses a direct aggregation handle.
Record-time attributes remain valid but use the normal lookup path. `Finish`
ends an exact concrete series and exports its value once. Shared overflow series
cannot be finished. A bound handle used after `Finish` lazily starts a new
series lifetime. For compatible non-Sum views, bound recording uses the normal
fixed-attribute measurement path; the direct-path performance claim applies
only to Sum.

No SDK runtime is copied or generated. Each option embeds a stable option and
adds an experimental construction hook that the stable SDK recognizes
structurally, without the stable module importing this module. See
[DESIGN.md](./DESIGN.md) for the package boundary and limitations.
