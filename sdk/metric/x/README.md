# Experimental SDK Metric Features

<!--
Adding/uncommenting the following after initial PR merge:
[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric/x)
-->

This package contains experimental SDK-level features for the OpenTelemetry metric SDK.
These types are currently under development and not part of the stable API.
They may be changed in backwards-incompatible ways, or removed entirely.

## Features

- [MeterConfigurator](#meterconfigurator)

### MeterConfigurator

`WithMeterConfigurator` enables or disables Meters by instrumentation scope, following the specification's [MeterConfigurator](https://opentelemetry.io/docs/specs/otel/metrics/sdk/#meterconfigurator).
The configuration can be updated at runtime with `MeterConfiguratorHandle.Set`.

#### Disabled meter behavior

- A disabled Meter drops new measurements and skips callback invocations.
- Instrument creation still validates arguments.
  The SDK keeps callbacks that users register while the Meter is disabled.
  These callbacks can run when the Meter becomes enabled again.
- The SDK continues to export earlier cumulative data while the Meter is disabled.

This behavior differs from the specification's current requirement that a disabled Meter behaves as a no-op, and it can change.
That difference needs to be resolved before stabilization; see [specification issue #5301](https://github.com/open-telemetry/opentelemetry-specification/issues/5301).

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.
