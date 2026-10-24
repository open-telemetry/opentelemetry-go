<!-- Generated. DO NOT MODIFY. -->
# Migration from v1.43.0 to v1.44.0

## Metric API changes

Some metric instruments in this release have changed type to match the updated
semantic conventions. These changes may require updates to applications using
the affected APIs.

- `containerconv.MemoryUsage` and `containerconv.MemoryUsageObservable` changed
  from `Int64Counter` and `Int64ObservableCounter` to `Int64UpDownCounter` and
  `Int64ObservableUpDownCounter`, respectively. Update code that supplies
  instrument-specific options or depends on the underlying instrument types.
- `k8sconv.NodeMemoryUsage` changed from `Int64Gauge` to
  `Int64UpDownCounter`. Its recording method changed from `Record` to `Add`,
  and the constructor option types changed accordingly.
- `hwconv.Errors.AttrNetworkIODirection` and
  `hwconv.ErrorsObservable.AttrNetworkIODirection` were removed. Review callers
  of these helpers and use the appropriate remaining semantic-convention
  attributes where applicable.

The `go.opentelemetry.io/otel/semconv/v1.44.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.43.0`.
