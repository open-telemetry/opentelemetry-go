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

package internal // import "go.opentelemetry.io/otel/bridge/opencensus/opencensusmetric/internal"

import (
	"errors"
	"testing"
	"time"

	ocmetricdata "go.opencensus.io/metric/metricdata"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestConvertMetrics(t *testing.T) {
	endTime1 := time.Now()
	endTime2 := endTime1.Add(-time.Millisecond)
	startTime := endTime2.Add(-time.Minute)
	for _, tc := range []struct {
		desc        string
		input       []*ocmetricdata.Metric
		expected    []metricdata.Metrics
		expectedErr error
	}{
		{
			desc:     "empty",
			expected: []metricdata.Metrics{},
		},
		{
			desc: "normal Histogram, gauges, and sums",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/histogram-a",
						Description: "a testing histogram",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeDistribution,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "a"},
							{Key: "b"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{

							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "hello",
									Present: true,
								}, {
									Value:   "world",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewDistributionPoint(endTime1, &ocmetricdata.Distribution{
									Count: 8,
									Sum:   100.0,
									BucketOptions: &ocmetricdata.BucketOptions{
										Bounds: []float64{1.0, 2.0, 3.0},
									},
									Buckets: []ocmetricdata.Bucket{
										{Count: 1},
										{Count: 2},
										{Count: 5},
									},
								}),
								ocmetricdata.NewDistributionPoint(endTime2, &ocmetricdata.Distribution{
									Count: 10,
									Sum:   110.0,
									BucketOptions: &ocmetricdata.BucketOptions{
										Bounds: []float64{1.0, 2.0, 3.0},
									},
									Buckets: []ocmetricdata.Bucket{
										{Count: 1},
										{Count: 4},
										{Count: 5},
									},
								}),
							},
							StartTime: startTime,
						},
					},
				}, {
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/gauge-a",
						Description: "an int testing gauge",
						Unit:        ocmetricdata.UnitBytes,
						Type:        ocmetricdata.TypeGaugeInt64,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "c"},
							{Key: "d"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "foo",
									Present: true,
								}, {
									Value:   "bar",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewInt64Point(endTime1, 123),
								ocmetricdata.NewInt64Point(endTime2, 1236),
							},
						},
					},
				}, {
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/gauge-b",
						Description: "a float testing gauge",
						Unit:        ocmetricdata.UnitBytes,
						Type:        ocmetricdata.TypeGaugeFloat64,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "cf"},
							{Key: "df"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "foof",
									Present: true,
								}, {
									Value:   "barf",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewFloat64Point(endTime1, 123.4),
								ocmetricdata.NewFloat64Point(endTime2, 1236.7),
							},
						},
					},
				}, {
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/sum-a",
						Description: "an int testing sum",
						Unit:        ocmetricdata.UnitMilliseconds,
						Type:        ocmetricdata.TypeCumulativeInt64,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "e"},
							{Key: "f"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "zig",
									Present: true,
								}, {
									Value:   "zag",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewInt64Point(endTime1, 13),
								ocmetricdata.NewInt64Point(endTime2, 14),
							},
						},
					},
				}, {
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/sum-b",
						Description: "a float testing sum",
						Unit:        ocmetricdata.UnitMilliseconds,
						Type:        ocmetricdata.TypeCumulativeFloat64,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "e"},
							{Key: "f"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "zig",
									Present: true,
								}, {
									Value:   "zag",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewFloat64Point(endTime1, 12.3),
								ocmetricdata.NewFloat64Point(endTime2, 123.4),
							},
						},
					},
				},
			},
			expected: []metricdata.Metrics{
				{
					Name:        "foo.com/histogram-a",
					Description: "a testing histogram",
					Unit:        "1",
					Data: metricdata.Histogram[float64]{
						DataPoints: []metricdata.HistogramDataPoint[float64]{
							{
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("a"),
									Value: attribute.StringValue("hello"),
								}, attribute.KeyValue{
									Key:   attribute.Key("b"),
									Value: attribute.StringValue("world"),
								}),
								StartTime:    startTime,
								Time:         endTime1,
								Count:        8,
								Sum:          100.0,
								Bounds:       []float64{1.0, 2.0, 3.0},
								BucketCounts: []uint64{1, 2, 5},
							}, {
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("a"),
									Value: attribute.StringValue("hello"),
								}, attribute.KeyValue{
									Key:   attribute.Key("b"),
									Value: attribute.StringValue("world"),
								}),
								StartTime:    startTime,
								Time:         endTime2,
								Count:        10,
								Sum:          110.0,
								Bounds:       []float64{1.0, 2.0, 3.0},
								BucketCounts: []uint64{1, 4, 5},
							},
						},
						Temporality: metricdata.CumulativeTemporality,
					},
				}, {
					Name:        "foo.com/gauge-a",
					Description: "an int testing gauge",
					Unit:        "By",
					Data: metricdata.Gauge[int64]{
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("c"),
									Value: attribute.StringValue("foo"),
								}, attribute.KeyValue{
									Key:   attribute.Key("d"),
									Value: attribute.StringValue("bar"),
								}),
								Time:  endTime1,
								Value: 123,
							}, {
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("c"),
									Value: attribute.StringValue("foo"),
								}, attribute.KeyValue{
									Key:   attribute.Key("d"),
									Value: attribute.StringValue("bar"),
								}),
								Time:  endTime2,
								Value: 1236,
							},
						},
					},
				}, {
					Name:        "foo.com/gauge-b",
					Description: "a float testing gauge",
					Unit:        "By",
					Data: metricdata.Gauge[float64]{
						DataPoints: []metricdata.DataPoint[float64]{
							{
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("cf"),
									Value: attribute.StringValue("foof"),
								}, attribute.KeyValue{
									Key:   attribute.Key("df"),
									Value: attribute.StringValue("barf"),
								}),
								Time:  endTime1,
								Value: 123.4,
							}, {
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("cf"),
									Value: attribute.StringValue("foof"),
								}, attribute.KeyValue{
									Key:   attribute.Key("df"),
									Value: attribute.StringValue("barf"),
								}),
								Time:  endTime2,
								Value: 1236.7,
							},
						},
					},
				}, {
					Name:        "foo.com/sum-a",
					Description: "an int testing sum",
					Unit:        "ms",
					Data: metricdata.Sum[int64]{
						IsMonotonic: true,
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("e"),
									Value: attribute.StringValue("zig"),
								}, attribute.KeyValue{
									Key:   attribute.Key("f"),
									Value: attribute.StringValue("zag"),
								}),
								Time:  endTime1,
								Value: 13,
							}, {
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("e"),
									Value: attribute.StringValue("zig"),
								}, attribute.KeyValue{
									Key:   attribute.Key("f"),
									Value: attribute.StringValue("zag"),
								}),
								Time:  endTime2,
								Value: 14,
							},
						},
					},
				}, {
					Name:        "foo.com/sum-b",
					Description: "a float testing sum",
					Unit:        "ms",
					Data: metricdata.Sum[float64]{
						IsMonotonic: true,
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.DataPoint[float64]{
							{
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("e"),
									Value: attribute.StringValue("zig"),
								}, attribute.KeyValue{
									Key:   attribute.Key("f"),
									Value: attribute.StringValue("zag"),
								}),
								Time:  endTime1,
								Value: 12.3,
							}, {
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("e"),
									Value: attribute.StringValue("zig"),
								}, attribute.KeyValue{
									Key:   attribute.Key("f"),
									Value: attribute.StringValue("zag"),
								}),
								Time:  endTime2,
								Value: 123.4,
							},
						},
					},
				},
			},
		}, {
			desc: "histogram without data points",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/histogram-a",
						Description: "a testing histogram",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeDistribution,
					},
				},
			},
			expected: []metricdata.Metrics{
				{
					Name:        "foo.com/histogram-a",
					Description: "a testing histogram",
					Unit:        "1",
					Data: metricdata.Histogram[float64]{
						Temporality: metricdata.CumulativeTemporality,
						DataPoints:  []metricdata.HistogramDataPoint[float64]{},
					},
				},
			},
		}, {
			desc: "sum without data points",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/sum-a",
						Description: "a testing sum",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeFloat64,
					},
				},
			},
			expected: []metricdata.Metrics{
				{
					Name:        "foo.com/sum-a",
					Description: "a testing sum",
					Unit:        "1",
					Data: metricdata.Sum[float64]{
						IsMonotonic: true,
						Temporality: metricdata.CumulativeTemporality,
						DataPoints:  []metricdata.DataPoint[float64]{},
					},
				},
			},
		}, {
			desc: "gauge without data points",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/gauge-a",
						Description: "a testing gauge",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeGaugeInt64,
					},
				},
			},
			expected: []metricdata.Metrics{
				{
					Name:        "foo.com/gauge-a",
					Description: "a testing gauge",
					Unit:        "1",
					Data: metricdata.Gauge[int64]{
						DataPoints: []metricdata.DataPoint[int64]{},
					},
				},
			},
		}, {
			desc: "histogram with negative count",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/histogram-a",
						Description: "a testing histogram",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeDistribution,
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							Points: []ocmetricdata.Point{
								ocmetricdata.NewDistributionPoint(endTime1, &ocmetricdata.Distribution{
									Count: -8,
								}),
							},
							StartTime: startTime,
						},
					},
				},
			},
			expectedErr: errConversion,
		}, {
			desc: "histogram with negative bucket count",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/histogram-a",
						Description: "a testing histogram",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeDistribution,
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							Points: []ocmetricdata.Point{
								ocmetricdata.NewDistributionPoint(endTime1, &ocmetricdata.Distribution{
									Buckets: []ocmetricdata.Bucket{
										{Count: -1},
										{Count: 2},
										{Count: 5},
									},
								}),
							},
							StartTime: startTime,
						},
					},
				},
			},
			expectedErr: errConversion,
		}, {
			desc: "histogram with non-histogram datapoint type",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/bad-point",
						Description: "a bad type",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeDistribution,
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							Points: []ocmetricdata.Point{
								ocmetricdata.NewFloat64Point(endTime1, 1.0),
							},
							StartTime: startTime,
						},
					},
				},
			},
			expectedErr: errConversion,
		}, {
			desc: "sum with non-sum datapoint type",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/bad-point",
						Description: "a bad type",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeCumulativeFloat64,
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							Points: []ocmetricdata.Point{
								ocmetricdata.NewDistributionPoint(endTime1, &ocmetricdata.Distribution{}),
							},
							StartTime: startTime,
						},
					},
				},
			},
			expectedErr: errConversion,
		}, {
			desc: "gauge with non-gauge datapoint type",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/bad-point",
						Description: "a bad type",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeGaugeFloat64,
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							Points: []ocmetricdata.Point{
								ocmetricdata.NewDistributionPoint(endTime1, &ocmetricdata.Distribution{}),
							},
							StartTime: startTime,
						},
					},
				},
			},
			expectedErr: errConversion,
		}, {
			desc: "unsupported Gauge Distribution type",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/bad-point",
						Description: "a bad type",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeGaugeDistribution,
					},
				},
			},
			expectedErr: errConversion,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			output, err := ConvertMetrics(tc.input)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("convertAggregation(%+v) = err(%v), want err(%v)", tc.input, err, tc.expectedErr)
			}
			metricdatatest.AssertEqual[metricdata.ScopeMetrics](t,
				metricdata.ScopeMetrics{Metrics: tc.expected},
				metricdata.ScopeMetrics{Metrics: output})
		})
	}
}

func TestConvertAttributes(t *testing.T) {
	setWithMultipleKeys := attribute.NewSet(
		attribute.KeyValue{Key: attribute.Key("first"), Value: attribute.StringValue("1")},
		attribute.KeyValue{Key: attribute.Key("second"), Value: attribute.StringValue("2")},
	)
	for _, tc := range []struct {
		desc        string
		inputKeys   []ocmetricdata.LabelKey
		inputValues []ocmetricdata.LabelValue
		expected    *attribute.Set
		expectedErr error
	}{
		{
			desc:     "no attributes",
			expected: attribute.EmptySet(),
		},
		{
			desc:        "different numbers of keys and values",
			inputKeys:   []ocmetricdata.LabelKey{{Key: "foo"}},
			expected:    attribute.EmptySet(),
			expectedErr: errMismatchedAttributeKeyValues,
		},
		{
			desc:      "multiple keys and values",
			inputKeys: []ocmetricdata.LabelKey{{Key: "first"}, {Key: "second"}},
			inputValues: []ocmetricdata.LabelValue{
				{Value: "1", Present: true},
				{Value: "2", Present: true},
			},
			expected: &setWithMultipleKeys,
		},
		{
			desc:      "multiple keys and values with some not present",
			inputKeys: []ocmetricdata.LabelKey{{Key: "first"}, {Key: "second"}, {Key: "third"}},
			inputValues: []ocmetricdata.LabelValue{
				{Value: "1", Present: true},
				{Value: "2", Present: true},
				{Present: false},
			},
			expected: &setWithMultipleKeys,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			output, err := convertAttrs(tc.inputKeys, tc.inputValues)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("convertAttrs(keys: %v, values: %v) = err(%v), want err(%v)", tc.inputKeys, tc.inputValues, err, tc.expectedErr)
			}
			if !output.Equals(tc.expected) {
				t.Errorf("convertAttrs(keys: %v, values: %v) = %+v, want %+v", tc.inputKeys, tc.inputValues, output.ToSlice(), tc.expected.ToSlice())
			}
		})
	}
}
