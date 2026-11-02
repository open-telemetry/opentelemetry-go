// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	metricFilterAccept = iota
	metricFilterDrop
	metricFilterAcceptPartial
)

const (
	metricFilterAttrAccept = iota
	metricFilterAttrDrop
)

// These values must match AggregationKind in sdk/metric/x.
const (
	aggregationKindUnknown uint8 = iota
	aggregationKindSum
	aggregationKindGauge
	aggregationKindHistogram
	aggregationKindExponentialHistogram
	aggregationKindSummary
)

// metricFilter is an internal, private mirror of the experimental metricFilter
// interface from the 'x' package. It is duplicated here to avoid package
// import cycles and layer violations while leveraging Go's implicit interfaces.
type metricFilter interface {
	Experimental()
	TestMetric(scope instrumentation.Scope, name string, kind uint8, unit string) int
	TestAttributes(
		scope instrumentation.Scope,
		name string,
		kind uint8,
		unit string,
		attrs attribute.Set,
	) int
}

// filterScopeMetrics filters metrics in-place. The caller's array is
// reused; no allocation. A scope left with zero metrics after filtering is
// dropped.
func filterScopeMetrics(sm []metricdata.ScopeMetrics, filter metricFilter) []metricdata.ScopeMetrics {
	i := 0
	for idx := range sm {
		sm[idx].Metrics = filterMetrics(sm[idx].Scope, sm[idx].Metrics, filter)
		if len(sm[idx].Metrics) == 0 {
			continue
		}
		sm[i], sm[idx] = sm[idx], sm[i]
		i++
	}
	return sm[:i]
}

// filterMetrics filters metrics in-place via read/write-index compaction:
// kept metrics move to the front and the returned slice views only them.
// A metric left with zero data points after partial filtering is dropped.
func filterMetrics(
	scope instrumentation.Scope,
	metrics []metricdata.Metrics,
	filter metricFilter,
) []metricdata.Metrics {
	j := 0
	for i := range metrics {
		kind := aggregationKind(metrics[i].Data)
		switch filter.TestMetric(scope, metrics[i].Name, kind, metrics[i].Unit) {
		case metricFilterAccept:
		case metricFilterAcceptPartial:
			keep := func(attrs attribute.Set) bool {
				return filter.TestAttributes(
					scope,
					metrics[i].Name,
					kind,
					metrics[i].Unit,
					attrs,
				) == metricFilterAttrAccept
			}
			metrics[i].Data = filterAggregation(metrics[i].Data, keep)
			if dataPointCount(metrics[i].Data) == 0 {
				continue
			}
		default: // metricFilterDrop
			continue
		}
		metrics[j], metrics[i] = metrics[i], metrics[j]
		j++
	}
	return metrics[:j]
}

func aggregationKind(data metricdata.Aggregation) uint8 {
	switch data.(type) {
	case metricdata.Sum[int64], metricdata.Sum[float64]:
		return aggregationKindSum
	case metricdata.Gauge[int64], metricdata.Gauge[float64]:
		return aggregationKindGauge
	case metricdata.Histogram[int64], metricdata.Histogram[float64]:
		return aggregationKindHistogram
	case metricdata.ExponentialHistogram[int64], metricdata.ExponentialHistogram[float64]:
		return aggregationKindExponentialHistogram
	case metricdata.Summary:
		return aggregationKindSummary
	default:
		return aggregationKindUnknown
	}
}

// dataPointCount returns the number of data points in data. Needed after
// AcceptPartial filtering to drop metrics left with zero data points: the
// exported spec disallows emitting a metric with no data points.
func dataPointCount(data metricdata.Aggregation) int {
	switch d := data.(type) {
	case metricdata.Sum[int64]:
		return len(d.DataPoints)
	case metricdata.Sum[float64]:
		return len(d.DataPoints)
	case metricdata.Gauge[int64]:
		return len(d.DataPoints)
	case metricdata.Gauge[float64]:
		return len(d.DataPoints)
	case metricdata.Histogram[int64]:
		return len(d.DataPoints)
	case metricdata.Histogram[float64]:
		return len(d.DataPoints)
	case metricdata.ExponentialHistogram[int64]:
		return len(d.DataPoints)
	case metricdata.ExponentialHistogram[float64]:
		return len(d.DataPoints)
	case metricdata.Summary:
		return len(d.DataPoints)
	default:
		return 0
	}
}

// filterAggregation applies keep to each data point in data. The type
// switch is required because the attribute filter operates on the concrete
// DataPoints slice, which lives on concrete structs behind the Aggregation
// interface. The reconstructed value is returned to the caller since the
// original interface value cannot be mutated in-place.
func filterAggregation(
	data metricdata.Aggregation,
	keep func(attribute.Set) bool,
) metricdata.Aggregation {
	switch d := data.(type) {
	case metricdata.Sum[int64]:
		d.DataPoints = filterDataPoints(d.DataPoints, keep)
		return d
	case metricdata.Sum[float64]:
		d.DataPoints = filterDataPoints(d.DataPoints, keep)
		return d
	case metricdata.Gauge[int64]:
		d.DataPoints = filterDataPoints(d.DataPoints, keep)
		return d
	case metricdata.Gauge[float64]:
		d.DataPoints = filterDataPoints(d.DataPoints, keep)
		return d
	case metricdata.Histogram[int64]:
		d.DataPoints = filterHistogramDataPoints(d.DataPoints, keep)
		return d
	case metricdata.Histogram[float64]:
		d.DataPoints = filterHistogramDataPoints(d.DataPoints, keep)
		return d
	case metricdata.ExponentialHistogram[int64]:
		d.DataPoints = filterExponentialHistogramDataPoints(d.DataPoints, keep)
		return d
	case metricdata.ExponentialHistogram[float64]:
		d.DataPoints = filterExponentialHistogramDataPoints(d.DataPoints, keep)
		return d
	case metricdata.Summary:
		d.DataPoints = filterSummaryDataPoints(d.DataPoints, keep)
		return d
	default:
		return d
	}
}

// filterDataPoints filters in-place. All four filter helpers share the same
// compaction shape. They differ only in the concrete data point type: Go
// generics can't abstract over struct field access.
func filterDataPoints[N int64 | float64](
	pts []metricdata.DataPoint[N],
	keep func(attribute.Set) bool,
) []metricdata.DataPoint[N] {
	j := 0
	for i := range pts {
		if keep(pts[i].Attributes) {
			pts[j], pts[i] = pts[i], pts[j]
			j++
		}
	}
	return pts[:j]
}

func filterHistogramDataPoints[N int64 | float64](
	pts []metricdata.HistogramDataPoint[N],
	keep func(attribute.Set) bool,
) []metricdata.HistogramDataPoint[N] {
	j := 0
	for i := range pts {
		if keep(pts[i].Attributes) {
			pts[j], pts[i] = pts[i], pts[j]
			j++
		}
	}
	return pts[:j]
}

func filterExponentialHistogramDataPoints[N int64 | float64](
	pts []metricdata.ExponentialHistogramDataPoint[N],
	keep func(attribute.Set) bool,
) []metricdata.ExponentialHistogramDataPoint[N] {
	j := 0
	for i := range pts {
		if keep(pts[i].Attributes) {
			pts[j], pts[i] = pts[i], pts[j]
			j++
		}
	}
	return pts[:j]
}

func filterSummaryDataPoints(
	pts []metricdata.SummaryDataPoint,
	keep func(attribute.Set) bool,
) []metricdata.SummaryDataPoint {
	j := 0
	for i := range pts {
		if keep(pts[i].Attributes) {
			pts[j], pts[i] = pts[i], pts[j]
			j++
		}
	}
	return pts[:j]
}
