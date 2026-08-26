// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x contains experimental metric SDK options and features.
//
// These features are experimental and under development. They may be changed
// in backwards-incompatible ways or removed in future releases.
//
// # Composable View Matching
//
// Composable view matching (configured via [WithViewMatchingMode] with
// [ViewMatchingModeComposable]) merges matching views for an instrument rather
// than creating independent metric streams:
//
//   - Views matching an instrument are grouped by target stream name. Non-renaming
//     views apply across each stream group for the instrument.
//   - A view that returns the instrument's own Name, Description, or Unit is
//     treated as leaving that property unspecified. A direct View function can
//     explicitly alter or zero these fields by returning different values.
//   - Scalar properties (Description, Unit, ExemplarReservoirProviderSelector)
//     follow last-wins precedence among matching views.
//   - Aggregations follow last-wins precedence. If an invalid or incompatible
//     aggregation is specified, the SDK logs a warning and falls back to
//     preceding matching views or the reader default.
//   - Per-view cardinality limits (aggregation_cardinality_limit) are not yet
//     supported on [go.opentelemetry.io/otel/sdk/metric.Stream]; cardinality
//     limits continue to be configured at the Reader or MeterProvider level.
//   - Attribute filters across matching views in a group are combined using logical
//     AND. If no matching view in the group specifies an AttributeFilter, any
//     instrument-level default attributes (e.g. from WithDefaultAttributes) act as
//     the baseline filter. An explicit AttributeFilter in a matching view overrides
//     the baseline.
package x
