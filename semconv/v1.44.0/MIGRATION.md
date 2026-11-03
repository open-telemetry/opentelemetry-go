<!-- Generated. DO NOT MODIFY. -->
# Migration from v1.43.0 to v1.44.0

The `go.opentelemetry.io/otel/semconv/v1.44.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.43.0`, with the following exceptions.

## Metric API changes

Some metric instruments in this release have changed names or types to match the updated semantic conventions. These changes may require updates to applications using the affected APIs.

- `containerconv.MemoryPagingFaults` and `containerconv.MemoryPagingFaultsObservable` have been renamed to `containerconv.PagingFaults` and `containerconv.PagingFaultsObservable`, respectively.
- `k8sconv.NodeMemoryPagingFaults` and `k8sconv.NodeMemoryPagingFaultsObservable` have been renamed to `k8sconv.NodePagingFaults` and `k8sconv.NodePagingFaultsObservable`, respectively.
- `k8sconv.PodMemoryPagingFaults` and `k8sconv.PodMemoryPagingFaultsObservable` have been renamed to `k8sconv.PodPagingFaults` and `k8sconv.PodPagingFaultsObservable`, respectively.
- `containerconv.MemoryUsage` and `containerconv.MemoryUsageObservable` changed from `Int64Counter` and `Int64ObservableCounter` to `Int64UpDownCounter` and `Int64ObservableUpDownCounter`, respectively. Update code that supplies instrument-specific options or depends on the underlying instrument types.
- `k8sconv.NodeMemoryUsage`, `k8sconv.NodeMemoryUsageObservable`, `k8sconv.PodMemoryUsage`, and `k8sconv.PodMemoryUsageObservable` changed from gauge instruments to up-down-counter instruments. Update the instrument-specific option types and replace `Record` with `Add` and `RecordSet` with `AddSet` where applicable.
- `hwconv.Errors.AttrNetworkIODirection` and `hwconv.ErrorsObservable.AttrNetworkIODirection` were removed. Review callers of these helpers and use the appropriate remaining semantic-convention attributes where applicable.
