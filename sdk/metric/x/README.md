# Experimental Features

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric/x)

This module contains experimental Metric SDK features.

## MetricFilter

`WithMetricFilter` configures a reader to filter aggregated metric streams during collection.
`MetricFilter.TestMetric` receives `TestMetricParameters` and returns `Accept`, `Drop`, or `AcceptPartial`.
For `AcceptPartial`, `MetricFilter.TestAttributes` receives `TestAttributesParameters` for each data point and returns `AttrAccept` or `AttrDrop`.
Parameters contain the instrumentation scope, metric name, aggregation kind, and unit.
Attribute parameters also contain the data point's attributes.
The aggregation is passed as `metricdata.Aggregation`, matching the collected data, including View-selected aggregations.

## Compatibility and Stability

Experimental features are outside the OpenTelemetry Go [versioning and stability policy](../../../VERSIONING.md).
These features may be removed or modified in successive releases, including patch releases.
# Experimental Metric SDK

<!--
Adding/uncommenting the following after initial PR merge:
[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric/x)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric/x)
-->

This module contains experimental extensions to the OpenTelemetry Go Metric SDK.
Its APIs may change incompatibly between releases.
