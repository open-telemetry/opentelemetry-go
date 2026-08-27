# Bound Instrument PoC Design Notes

This module tests whether experimental metric SDK behavior can be developed
without adding provisional methods to `go.opentelemetry.io/otel/sdk/metric`.
It intentionally has its own MeterProvider, Reader contract, pipelines, and
aggregation state. A stable SDK Reader cannot implement the local Reader
interface because its registration methods are package-private.

## Implemented seam

The public façade reuses stable `Instrument`, `Stream`, `View`, aggregation,
exemplar, resource, and `metricdata` types. The bound Sum state is generated
from `internal/shared/metricx`, so it can be rendered into another module
without importing a parent module's `internal` packages. The generated state
owns cardinality, cumulative and delta collection, exemplars, direct handles,
finished-series tombstones, and lazy reactivation.

The stable SDK package is not changed and its instrument method sets remain
unchanged. This is the most important compatibility property of the
experiment.

## Deliberate vertical-slice limits

- Only `Int64Counter` with Sum aggregation is implemented. Other instrument
  methods are inherited from the API no-op implementation.
- Stable PeriodicReader, exporter, and external Producer integration are not
  included. ManualReader is sufficient to exercise cumulative and delta
  pipelines in the same provider.
- A shared overflow point cannot be finished because its value may contain
  contributions from several original attribute sets.
- If a finished set is reactivated before collection, its final lifetime is
  exported in that collection and its new lifetime starts exporting in the
  next collection. This avoids duplicate data points with identical attributes.
- Expanding templates to the stable provider, reader, meter, and all
  aggregations would duplicate most of the SDK source in this PoC. That work is
  deferred until the API and lifecycle experiment demonstrates sufficient
  value to justify the larger refactor.

## Rejected alternatives

- Importing `sdk/metric/internal` from this module would compile because of the
  directory layout, but would couple independently versioned modules to an
  unstable internal implementation.
- Wrapping a stable MeterProvider cannot provide bound handles because the
  stable Reader and pipeline contracts are sealed by package-private methods.
- Running a second provider for bound measurements would require merging
  batches and would apply views, cardinality, and collection lifecycles
  independently.
- Reflection, `unsafe`, and linkname-based bridges would make the prototype
  fragile and unsuitable as a migration path.
