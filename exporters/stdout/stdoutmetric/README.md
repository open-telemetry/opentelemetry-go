# STDOUT Metric Exporter

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/stdout/stdoutmetric)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/stdout/stdoutmetric)

The default JSON output uses the [OpenTelemetry attribute representation for non-OTLP protocols](https://opentelemetry.io/docs/specs/otel/common/#attribute-representation-for-non-otlp).
Resource, instrumentation scope, and data point attributes are JSON objects, such as `{"service.name":"example","rate":1.5}`.
Exemplar attributes are arrays of single-entry objects, such as `[{"request.id":"abc"}]`.
Empty attribute sets are `{}`; a missing resource is `null`.
Non-finite metric and attribute values are encoded as `"NaN"`, `"Infinity"`, or `"-Infinity"`, including values nested in arrays and maps.

Custom encoders supplied with `WithEncoder` receive the original metric data and choose their own representation.
The output is intended for testing and debugging and has no stability or compatibility guarantees.
Use an OTLP exporter when an interchange format is required.
