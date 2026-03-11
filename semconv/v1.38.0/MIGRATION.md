<!-- Generated. DO NOT MODIFY. -->
# Migration from v1.37.0 to v1.38.0

The `go.opentelemetry.io/otel/semconv/v1.38.0` package should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.37.0` with the following exceptions.

## Removed

The following declarations have been removed.
Refer to the [OpenTelemetry Semantic Conventions documentation] for deprecation instructions.

If the type is not listed in the documentation as deprecated, it has been removed in this version due to lack of applicability or use.
If you use any of these non-deprecated declarations in your Go application, please [open an issue] describing your use-case.

- `ProcessPagingFaultTypeKey`
- `ProcessPagingFaultTypeMajor`
- `ProcessPagingFaultTypeMinor`
- `SystemCPULogicalNumber`
- `SystemCPULogicalNumberKey`
- `SystemPagingTypeKey`
- `SystemPagingTypeMajor`
- `SystemPagingTypeMinor`
- `SystemProcessStatusDefunct`
- `SystemProcessStatusKey`
- `SystemProcessStatusRunning`
- `SystemProcessStatusSleeping`
- `SystemProcessStatusStopped`

[OpenTelemetry Semantic Conventions documentation]: https://github.com/open-telemetry/semantic-conventions
[open an issue]: https://github.com/open-telemetry/opentelemetry-go/issues/new?template=Blank+issue
