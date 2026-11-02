<!-- Generated. DO NOT MODIFY. -->
# Migration from v1.43.0 to v1.44.0

The `go.opentelemetry.io/otel/semconv/v1.44.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.43.0` with the following exceptions.

<!-- This fragment is included in the generated MIGRATION.md. -->
<!-- markdownlint-disable MD041 -->

## Metric API changes

The following changes affect metric subpackages as well as their constructor options and embedded instrument types.

### Paging fault renames

The synchronous and observable helpers have been renamed:

| Subpackage | Previous helpers | New helpers |
| --- | --- | --- |
| `containerconv` | `MemoryPagingFaults`, `MemoryPagingFaultsObservable` | `PagingFaults`, `PagingFaultsObservable` |
| `k8sconv` | `NodeMemoryPagingFaults`, `NodeMemoryPagingFaultsObservable` | `NodePagingFaults`, `NodePagingFaultsObservable` |
| `k8sconv` | `PodMemoryPagingFaults`, `PodMemoryPagingFaultsObservable` | `PodPagingFaults`, `PodPagingFaultsObservable` |

Update calls to their `New...` constructors as well. The emitted metric names also change:

| Previous metric name | New metric name |
| --- | --- |
| `container.memory.paging.faults` | `container.paging.faults` |
| `k8s.node.memory.paging.faults` | `k8s.node.paging.faults` |
| `k8s.pod.memory.paging.faults` | `k8s.pod.paging.faults` |

Update dashboards, queries, recording rules, and alerts that select the previous names.

### Memory instrument changes

- `containerconv.MemoryUsage` and `MemoryUsageObservable` change from `Int64Counter` and `Int64ObservableCounter` to `Int64UpDownCounter` and `Int64ObservableUpDownCounter`. Update instrument-specific options and code that uses the embedded instruments. Synchronous measurements remain deltas; observable measurements remain absolute values, which can now decrease.
- `k8sconv.NodeMemoryUsage` and `PodMemoryUsage` change from `Int64Gauge` to `Int64UpDownCounter`. Replace `Record` with `Add` and `RecordSet` with `AddSet`, and update instrument-specific options. **Convert absolute readings to deltas for each attribute set.** For readings of 100 and then 120, add 100 and then 20; adding 100 and 120 reports 220. Alternatively, use the observable helpers to report absolute readings without keeping delta state.
- `k8sconv.NodeMemoryUsageObservable` and `PodMemoryUsageObservable` change from `Int64ObservableGauge` to `Int64ObservableUpDownCounter`. Update instrument-specific options and embedded instrument references. Continue to observe absolute values; do not convert observable measurements to deltas.

### Floating point measurements

The helpers below now use the corresponding `Float64` instruments instead of `Int64` instruments so they can represent fractional values required by the semantic conventions.
This also applies to their `Observable` variants, where provided.
Update constructor option types, embedded instrument references, and measurement arguments from `int64` to `float64`.
For example, `containerconv.CPUUsage.Record` now accepts a value such as `0.25` CPUs, and `azureconv.CosmosDBClientOperationRequestCharge.Record` accepts a charge such as `46.18` request units.

| Subpackage | Helpers |
| --- | --- |
| `azureconv` | `CosmosDBClientOperationRequestCharge` |
| `containerconv` | `CPUUsage` |
| `hwconv` | `BatteryCharge`, `BatteryChargeLimit`, `Energy`, `FanSpeedRatio`, `GpuMemoryUtilization`, `GpuUtilization`, `HostAmbientTemperature`, `HostHeatingMargin`, `HostPower`, `LogicalDiskUtilization`, `NetworkBandwidthUtilization`, `PhysicalDiskEnduranceUtilization`, `Power`, `PowerSupplyUsage`, `PowerSupplyUtilization`, `Temperature`, `TemperatureLimit`, `Voltage`, `VoltageLimit`, `VoltageNominal` |
| `k8sconv` | `ContainerCPULimitCurrent`, `ContainerCPULimitDesired`, `ContainerCPULimitUtilization`, `ContainerCPURequestCurrent`, `ContainerCPURequestDesired`, `ContainerCPURequestUtilization`, `NodeCPUAllocatable`, `NodeCPUUsage`, `NodeSystemContainerCPUUsage`, `PodCPUUsage`, `ResourceQuotaCPULimitHard`, `ResourceQuotaCPULimitUsed`, `ResourceQuotaCPURequestHard`, `ResourceQuotaCPURequestUsed` |
| `openshiftconv` | `ClusterquotaCPULimitHard`, `ClusterquotaCPULimitUsed`, `ClusterquotaCPURequestHard`, `ClusterquotaCPURequestUsed` |
| `processconv` | `CPUUtilization` |
| `systemconv` | `CPUUtilization`, `FilesystemUtilization`, `MemoryLinuxHugepagesUtilization`, `PagingUtilization` |

### Removed attribute helpers

`hwconv.Errors.AttrNetworkIODirection` and `hwconv.ErrorsObservable.AttrNetworkIODirection` were removed.
Review callers and use the remaining semantic-convention attributes appropriate to each instrument.
