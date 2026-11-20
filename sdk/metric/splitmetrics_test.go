// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestSplitResourceMetrics(t *testing.T) {
	points := func(n int) []metricdata.DataPoint[int64] {
		var dps []metricdata.DataPoint[int64]
		for i := range n {
			dps = append(dps, metricdata.DataPoint[int64]{Value: int64(i)})
		}
		return dps
	}

	tests := []struct {
		name     string
		size     int
		input    *metricdata.ResourceMetrics
		expected [][][]int // Expected representation of the batching structure
	}{
		{
			name: "no splitting needed",
			size: 10,
			input: &metricdata.ResourceMetrics{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(3)}},
							{Data: metricdata.Gauge[int64]{DataPoints: points(2)}},
						},
					},
				},
			},
			expected: [][][]int{{{3, 2}}},
		},
		{
			name: "split on metric boundary",
			size: 3,
			input: &metricdata.ResourceMetrics{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(3)}},
							{Data: metricdata.Gauge[int64]{DataPoints: points(2)}},
						},
					},
				},
			},
			expected: [][][]int{
				{{3}},
				{{2}},
			},
		},
		{
			name: "split inside a single metric",
			size: 2,
			input: &metricdata.ResourceMetrics{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(5)}},
						},
					},
				},
			},
			expected: [][][]int{
				{{2}},
				{{2}},
				{{1}},
			},
		},
		{
			name: "split across scopes",
			size: 4,
			input: &metricdata.ResourceMetrics{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(2)}},
							{Data: metricdata.Gauge[int64]{DataPoints: points(3)}},
						},
					},
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(2)}},
						},
					},
				},
			},
			expected: [][][]int{
				{{2, 2}},
				{
					{1},
					{2},
				}, // The 3-point metric's 3rd point overflowed into the second batch, filling 1, leaving 3 left for the new scope
			},
		},
		{
			name: "zero points input",
			size: 5,
			input: &metricdata.ResourceMetrics{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(0)}},
						},
					},
				},
			},
			expected: [][][]int{{{0}}},
		},
		{
			name: "size zero",
			size: 0,
			input: &metricdata.ResourceMetrics{
				ScopeMetrics: []metricdata.ScopeMetrics{
					{
						Metrics: []metricdata.Metrics{
							{Data: metricdata.Gauge[int64]{DataPoints: points(1)}},
						},
					},
				},
			},
			expected: [][][]int{{{1}}},
		},
		{
			name:     "empty scope metrics",
			size:     10,
			input:    &metricdata.ResourceMetrics{},
			expected: [][][]int{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := batcher{size: tt.size}
			batches := b.splitResourceMetrics(tt.input)

			var actual [][][]int
			for _, batch := range batches {
				var scopes [][]int
				for _, sm := range batch.ScopeMetrics {
					var metrics []int
					for _, m := range sm.Metrics {
						metrics = append(metrics, metricDPC(m))
					}
					scopes = append(scopes, metrics)
				}
				actual = append(actual, scopes)
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCopyMetricData(t *testing.T) {
	tests := []struct {
		name           string
		data           metricdata.Aggregation
		expectedPoints int
	}{
		{
			name:           "Gauge[int64]",
			data:           metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{}, {}}},
			expectedPoints: 2,
		},
		{
			name:           "Gauge[float64]",
			data:           metricdata.Gauge[float64]{DataPoints: []metricdata.DataPoint[float64]{{}, {}}},
			expectedPoints: 2,
		},
		{
			name: "Sum[int64]",
			data: metricdata.Sum[int64]{
				DataPoints:  []metricdata.DataPoint[int64]{{}, {}},
				Temporality: metricdata.DeltaTemporality,
				IsMonotonic: true,
			},
			expectedPoints: 2,
		},
		{
			name: "Sum[float64]",
			data: metricdata.Sum[float64]{
				DataPoints:  []metricdata.DataPoint[float64]{{}, {}},
				Temporality: metricdata.DeltaTemporality,
				IsMonotonic: true,
			},
			expectedPoints: 2,
		},
		{
			name: "Histogram[int64]",
			data: metricdata.Histogram[int64]{
				DataPoints:  []metricdata.HistogramDataPoint[int64]{{}, {}},
				Temporality: metricdata.DeltaTemporality,
			},
			expectedPoints: 2,
		},
		{
			name: "Histogram[float64]",
			data: metricdata.Histogram[float64]{
				DataPoints:  []metricdata.HistogramDataPoint[float64]{{}, {}},
				Temporality: metricdata.CumulativeTemporality,
			},
			expectedPoints: 2,
		},
		{
			name: "ExponentialHistogram[int64]",
			data: metricdata.ExponentialHistogram[int64]{
				DataPoints:  []metricdata.ExponentialHistogramDataPoint[int64]{{}, {}},
				Temporality: metricdata.DeltaTemporality,
			},
			expectedPoints: 2,
		},
		{
			name: "ExponentialHistogram[float64]",
			data: metricdata.ExponentialHistogram[float64]{
				DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{{}, {}},
				Temporality: metricdata.CumulativeTemporality,
			},
			expectedPoints: 2,
		},
		{
			name:           "Summary",
			data:           metricdata.Summary{DataPoints: []metricdata.SummaryDataPoint{{}, {}}},
			expectedPoints: 2,
		},
		{
			name:           "Unknown Type",
			data:           nil,
			expectedPoints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := metricdata.Metrics{
				Name:        "test",
				Description: "desc",
				Unit:        "1",
				Data:        tt.data,
			}
			assert.Equal(t, tt.expectedPoints, metricDPC(m))

			if tt.expectedPoints == 0 {
				return
			}
			// Test copying 1 element out of 2
			copied := copyMetricData(m, 0, 1)
			assert.Equal(t, 1, metricDPC(copied))
			assert.Equal(t, "test", copied.Name)
			assert.Equal(t, "desc", copied.Description)
			assert.Equal(t, "1", copied.Unit)

			switch expectedData := tt.data.(type) {
			case metricdata.Sum[int64]:
				copiedData := copied.Data.(metricdata.Sum[int64])
				assert.Equal(t, expectedData.Temporality, copiedData.Temporality)
				assert.Equal(t, expectedData.IsMonotonic, copiedData.IsMonotonic)
			case metricdata.Sum[float64]:
				copiedData := copied.Data.(metricdata.Sum[float64])
				assert.Equal(t, expectedData.Temporality, copiedData.Temporality)
				assert.Equal(t, expectedData.IsMonotonic, copiedData.IsMonotonic)
			case metricdata.Histogram[int64]:
				copiedData := copied.Data.(metricdata.Histogram[int64])
				assert.Equal(t, expectedData.Temporality, copiedData.Temporality)
			case metricdata.Histogram[float64]:
				copiedData := copied.Data.(metricdata.Histogram[float64])
				assert.Equal(t, expectedData.Temporality, copiedData.Temporality)
			case metricdata.ExponentialHistogram[int64]:
				copiedData := copied.Data.(metricdata.ExponentialHistogram[int64])
				assert.Equal(t, expectedData.Temporality, copiedData.Temporality)
			case metricdata.ExponentialHistogram[float64]:
				copiedData := copied.Data.(metricdata.ExponentialHistogram[float64])
				assert.Equal(t, expectedData.Temporality, copiedData.Temporality)
			}
		})
	}
}
