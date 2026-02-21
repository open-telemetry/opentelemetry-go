// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// splitMetrics removes metric data points from src and puts up to size metric
// data points in dest. This is adapted from the collector's batch processor:
// https://github.com/open-telemetry/opentelemetry-collector/blob/587b90b9ecc1db959ee9104d5bf993591f80ca43/processor/batchprocessor/splitmetrics.go
func splitMetrics(size int, src, dest *metricdata.ResourceMetrics) {
	totalCopiedDataPoints := 0
	dest.Resource = src.Resource
	i := 0
	for ; i < len(src.ScopeMetrics); i++ {
		// If we are done skip everything else.
		if totalCopiedDataPoints == size {
			break
		}
		srcIlm := src.ScopeMetrics[i]
		// If possible to move all metrics do that.
		srcIlmDataPointCount := scopeMetricsDPC(srcIlm)
		if srcIlmDataPointCount+totalCopiedDataPoints <= size {
			totalCopiedDataPoints += srcIlmDataPointCount
			dest.ScopeMetrics = append(dest.ScopeMetrics, srcIlm)
			continue
		}

		destIlm := metricdata.ScopeMetrics{
			Scope: srcIlm.Scope,
		}
		j := 0
		for ; j < len(srcIlm.Metrics); j++ {
			// If we are done skip everything else.
			if totalCopiedDataPoints == size {
				break
			}
			srcMetric := srcIlm.Metrics[j]
			// If possible to move all points do that.
			srcMetricPointCount := metricDPC(srcMetric)
			if srcMetricPointCount+totalCopiedDataPoints <= size {
				totalCopiedDataPoints += srcMetricPointCount
				destIlm.Metrics = append(destIlm.Metrics, srcMetric)
				continue
			}

			// If the metric has more data points than free slots we should split it.
			newMetrics := metricdata.Metrics{}
			copiedDataPoints := size - totalCopiedDataPoints
			splitMetric(&srcIlm.Metrics[j], &newMetrics, copiedDataPoints)
			destIlm.Metrics = append(destIlm.Metrics, newMetrics)
			totalCopiedDataPoints += copiedDataPoints
			break
		}
		// Delete all of the metrics we fully moved.
		srcIlm.Metrics = srcIlm.Metrics[j:]
		dest.ScopeMetrics = append(dest.ScopeMetrics, destIlm)
		src.ScopeMetrics[i] = srcIlm
		break
	}
	// Delete all of the scope metrics we fully moved.
	src.ScopeMetrics = src.ScopeMetrics[i:]
}

// resourceMetricsDPC calculates the total number of data points in the metricdata.ResourceMetrics.
func resourceMetricsDPC(rs *metricdata.ResourceMetrics) int {
	dataPointCount := 0
	ilms := rs.ScopeMetrics
	for k := 0; k < len(ilms); k++ {
		dataPointCount += scopeMetricsDPC(ilms[k])
	}
	return dataPointCount
}

// scopeMetricsDPC calculates the total number of data points in the metricdata.ScopeMetrics.
func scopeMetricsDPC(ilm metricdata.ScopeMetrics) int {
	dataPointCount := 0
	ms := ilm.Metrics
	for k := 0; k < len(ms); k++ {
		dataPointCount += metricDPC(ms[k])
	}
	return dataPointCount
}

// metricDPC calculates the total number of data points in the metricdata.Metrics.
func metricDPC(ms metricdata.Metrics) int {
	switch a := ms.Data.(type) {
	case metricdata.Gauge[int64]:
		return len(a.DataPoints)
	case metricdata.Gauge[float64]:
		return len(a.DataPoints)
	case metricdata.Sum[int64]:
		return len(a.DataPoints)
	case metricdata.Sum[float64]:
		return len(a.DataPoints)
	case metricdata.Histogram[int64]:
		return len(a.DataPoints)
	case metricdata.Histogram[float64]:
		return len(a.DataPoints)
	case metricdata.ExponentialHistogram[int64]:
		return len(a.DataPoints)
	case metricdata.ExponentialHistogram[float64]:
		return len(a.DataPoints)
	case metricdata.Summary:
		return len(a.DataPoints)
	}
	return 0
}

// splitMetric removes metric points from the input data and moves data of the specified size to destination.
func splitMetric(ms, dest *metricdata.Metrics, size int) {
	dest.Name = ms.Name
	dest.Description = ms.Description
	dest.Unit = ms.Unit

	switch a := ms.Data.(type) {
	case metricdata.Gauge[int64]:
		dest.Data = metricdata.Gauge[int64]{
			DataPoints: a.DataPoints[:size],
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.Gauge[float64]:
		dest.Data = metricdata.Gauge[float64]{
			DataPoints: a.DataPoints[:size],
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.Sum[int64]:
		dest.Data = metricdata.Sum[int64]{
			DataPoints:  a.DataPoints[:size],
			Temporality: a.Temporality,
			IsMonotonic: a.IsMonotonic,
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.Sum[float64]:
		dest.Data = metricdata.Sum[float64]{
			DataPoints:  a.DataPoints[:size],
			Temporality: a.Temporality,
			IsMonotonic: a.IsMonotonic,
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.Histogram[int64]:
		dest.Data = metricdata.Histogram[int64]{
			DataPoints:  a.DataPoints[:size],
			Temporality: a.Temporality,
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.Histogram[float64]:
		dest.Data = metricdata.Histogram[float64]{
			DataPoints:  a.DataPoints[:size],
			Temporality: a.Temporality,
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.ExponentialHistogram[int64]:
		dest.Data = metricdata.ExponentialHistogram[int64]{
			DataPoints:  a.DataPoints[:size],
			Temporality: a.Temporality,
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.ExponentialHistogram[float64]:
		dest.Data = metricdata.ExponentialHistogram[float64]{
			DataPoints:  a.DataPoints[:size],
			Temporality: a.Temporality,
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	case metricdata.Summary:
		dest.Data = metricdata.Summary{
			DataPoints: a.DataPoints[:size],
		}
		a.DataPoints = a.DataPoints[size:]
		ms.Data = a
	}
}
