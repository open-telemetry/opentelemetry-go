# Bound Instrument PoC Design Notes

This module evaluates experimental bound instruments without adding their
methods to regular `go.opentelemetry.io/otel/sdk/metric` instruments. It is a
thin facade over one stable SDK provider, not a second SDK implementation.

## Runtime seam

`sdk/metric/x.NewMeterProvider` delegates to `sdk/metric.NewMeterProvider` and
injects an experimental option. The stable SDK discovers that option through a
structural type assertion and uses its wrapper only while constructing
`Int64Counter` instruments. The stable SDK does not import the experimental
module.

Readers, views, pipelines, aggregation selection, resources, cardinality,
exemplars, collection, and shutdown therefore use the regular SDK machinery.
Stable Readers and reader options are accepted directly by this facade. Other
instrument kinds retain their complete stable implementations.

The stable aggregation package models binding and finishing as independent
internal capabilities. An experimental Sum supplies both. Bind resolves and
filters attributes, applies cardinality, and returns a direct measurement
target for each reader and view. The bound Add hot path updates those targets
without attribute processing or map lookup. Finish independently retires the
exact filtered series.

Regular stable providers do not enable the experimental aggregation path.
Their counters do not implement `metric/x.Int64CounterBinder` or
`metric/x.Finisher`, and their existing measurement path is unchanged.

## Deliberate vertical-slice limits

- Only `Int64Counter` receives experimental methods.
- Sum has the direct bound implementation and exact Finish support.
- Compatible non-Sum views record bound values through their normal
  fixed-attribute measurement function, but do not yet remove series on
  Finish.
- A shared overflow point cannot be finished because its value may contain
  contributions from several original attribute sets.
- If a finished set is reactivated before collection, its final lifetime is
  exported in that collection and its new lifetime starts exporting in the
  following collection. This avoids duplicate points with identical
  attributes in one batch.

## Rejected alternatives

- Duplicating a provider, Reader, pipeline, and aggregation stack would test a
  separate SDK whose behavior could drift from the stable implementation.
- Templates with only one generated consumer add indirection without sharing
  implementation. Generating the entire SDK twice would preserve a source of
  truth but produce a much larger experimental surface.
- Importing `sdk/metric/internal` directly would couple independently selected
  module versions to an unstable internal API. The structural option keeps the
  dependency in the normal `sdk/metric/x` to `sdk/metric` direction.
- Running two providers would require merging batches and would apply views,
  cardinality, and collection lifecycles independently.
- Reflection, `unsafe`, and linkname-based bridges would be fragile migration
  paths.
