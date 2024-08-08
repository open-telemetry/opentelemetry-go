# Experimental Features

The log SDK contains features that have not yet stabilized.
These features are added to the OpenTelemetry Go log SDK prior to stabilization so that users can start experimenting with them and provide feedback.

These feature may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for more information.

## Features

- [Filter Processors](#filter-processor)

### Filter Processor

Users of logging APIs often want to know if a log `Record` will be processed or dropped before they perform complex operations to construct the `Record`.
The [`Logger`] in the Logs Bridge API provides the `Enabled` method for just this use-case.
In order for the Logs Bridge SDK to effectively implement this API, it needs to be known if the registered [`Processor`]s are enabled for the `Record` within a context.
A [`Processor`] that knows, and can identify, what `Record` it will process or drop when it is passed to `OnEmit` can communicate this to the SDK `Logger` by implementing the `FilterProcessor`.

The SDK `Logger` will check all of the registered [`Processor`]s that implement the `FilterProcessor` interface by calling `Enabled` when the `Logger.Enabled` method is called.

[`Logger`]: https://pkg.go.dev/go.opentelemetry.io/otel/log#Logger
[`Processor`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/log#Processor

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go versioning and stability [policy](../../../../VERSIONING.md).
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
