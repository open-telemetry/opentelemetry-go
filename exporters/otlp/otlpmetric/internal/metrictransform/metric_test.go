// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrictransform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/reader"

	// "go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"

	"go.opentelemetry.io/otel/sdk/metric/number"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

func TestStringKeyValues(t *testing.T) {
	tests := []struct {
		kvs      []attribute.KeyValue
		expected []*commonpb.KeyValue
	}{
		{
			nil,
			nil,
		},
		{
			[]attribute.KeyValue{},
			nil,
		},
		{
			[]attribute.KeyValue{
				attribute.Bool("true", true),
				attribute.Int64("one", 1),
				attribute.Int64("two", 2),
				attribute.Float64("three", 3),
				attribute.Int("four", 4),
				attribute.Int("five", 5),
				attribute.Float64("six", 6),
				attribute.Int("seven", 7),
				attribute.Int("eight", 8),
				attribute.String("the", "final word"),
			},
			[]*commonpb.KeyValue{
				{Key: "eight", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 8}}},
				{Key: "five", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 5}}},
				{Key: "four", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 4}}},
				{Key: "one", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 1}}},
				{Key: "seven", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 7}}},
				{Key: "six", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 6.0}}},
				{Key: "the", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "final word"}}},
				{Key: "three", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 3.0}}},
				{Key: "true", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_BoolValue{BoolValue: true}}},
				{Key: "two", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 2}}},
			},
		},
	}

	for _, test := range tests {
		labels := attribute.NewSet(test.kvs...)
		assert.Equal(t, test.expected, Iterator(labels.Iter()))
	}
}

func TestSumPoints(t *testing.T) {
	// desc := sdkinstrument.NewDescriptor("", sdkinstrument.HistogramKind, number.Int64Kind)
	// labels := attribute.NewSet(attribute.String("one", "1"))

	testcases := []struct {
		name        string
		points      []reader.Point
		kind        number.Kind
		temporality aggregation.Temporality
		want        *metricpb.Metric_Sum
	}{
		{
			name:   "no points",
			points: []reader.Point{},
			want:   nil,
		},
		{
			name: "incorrect aggregation",
			points: []reader.Point{
				{
					Aggregation: gauge.NewInt64(1),
				},
			},
			want: nil, // Error is in the error handler
		},
		{
			name: "int data",
			points: []reader.Point{
				{
					Aggregation: sum.NewInt64Monotonic(2),
				},
			},
			kind: number.Int64Kind,
			want: &metricpb.Metric_Sum{
				Sum: &metricpb.Sum{
					DataPoints: []*metricpb.NumberDataPoint{
						{
							Value: &metricpb.NumberDataPoint_AsInt{
								AsInt: 2,
							},
						},
					},
					IsMonotonic: true,
				},
			},
		},
		{
			name: "float data",
			points: []reader.Point{
				{
					Aggregation: sum.NewFloat64NonMonotonic(5),
				},
			},
			kind: number.Float64Kind,
			want: &metricpb.Metric_Sum{
				Sum: &metricpb.Sum{
					DataPoints: []*metricpb.NumberDataPoint{
						{
							Value: &metricpb.NumberDataPoint_AsDouble{
								AsDouble: 5,
							},
						},
					},
					IsMonotonic: false,
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got := sumPoints(tt.points, tt.kind, tt.temporality)
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestGaugePoints(t *testing.T) {
	// desc := sdkinstrument.NewDescriptor("", sdkinstrument.HistogramKind, number.Int64Kind)
	// labels := attribute.NewSet(attribute.String("one", "1"))

	testcases := []struct {
		name        string
		points      []reader.Point
		kind        number.Kind
		temporality aggregation.Temporality
		want        *metricpb.Metric_Gauge
	}{
		{
			name:   "no points",
			points: []reader.Point{},
			want:   nil,
		},
		{
			name: "incorrect aggregation",
			points: []reader.Point{
				{
					Aggregation: sum.NewInt64Monotonic(1),
				},
			},
			want: nil, // Error is in the error handler
		},
		{
			name: "int data",
			points: []reader.Point{
				{
					Aggregation: gauge.NewInt64(2),
				},
			},
			kind: number.Int64Kind,
			want: &metricpb.Metric_Gauge{
				Gauge: &metricpb.Gauge{
					DataPoints: []*metricpb.NumberDataPoint{
						{
							Value: &metricpb.NumberDataPoint_AsInt{
								AsInt: 2,
							},
						},
					},
				},
			},
		},
		{
			name: "float data",
			points: []reader.Point{
				{
					Aggregation: gauge.NewFloat64(5),
				},
			},
			kind: number.Float64Kind,
			want: &metricpb.Metric_Gauge{
				Gauge: &metricpb.Gauge{
					DataPoints: []*metricpb.NumberDataPoint{
						{
							Value: &metricpb.NumberDataPoint_AsDouble{
								AsDouble: 5,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got := gaugePoints(tt.points, tt.kind)
			assert.Equal(t, tt.want, got)
		})
	}

}
func TestHistogramPoints(t *testing.T) {
	boundaries := []float64{2.0, 5.0, 8.0}
	testSum := float64(11)
	testcases := []struct {
		name        string
		points      []reader.Point
		kind        number.Kind
		temporality aggregation.Temporality
		want        *metricpb.Metric_Histogram
	}{
		{
			name:   "no points",
			points: []reader.Point{},
			want:   nil,
		},
		{
			name: "incorrect aggregation",
			points: []reader.Point{
				{
					Aggregation: sum.NewInt64Monotonic(1),
				},
			},
			want: nil, // Error is in the error handler
		},
		{
			name: "int data",
			points: []reader.Point{
				{
					Aggregation: histogram.NewInt64(boundaries, 1, 10),
				},
			},
			kind: number.Int64Kind,
			want: &metricpb.Metric_Histogram{
				Histogram: &metricpb.Histogram{
					DataPoints: []*metricpb.HistogramDataPoint{
						{
							Count:          2,
							Sum:            &testSum,
							BucketCounts:   []uint64{1, 0, 0, 1},
							ExplicitBounds: boundaries,
						},
					},
				},
			},
		},
		{
			name: "float data",
			points: []reader.Point{
				{
					Aggregation: histogram.NewFloat64(boundaries, 1, 10),
				},
			},
			kind: number.Float64Kind,
			want: &metricpb.Metric_Histogram{
				Histogram: &metricpb.Histogram{
					DataPoints: []*metricpb.HistogramDataPoint{
						{
							Count:          2,
							Sum:            &testSum,
							BucketCounts:   []uint64{1, 0, 0, 1},
							ExplicitBounds: boundaries,
						},
					},
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got := histogramPoints(tt.points, tt.kind, tt.temporality)
			assert.Equal(t, tt.want, got)
		})
	}

}
