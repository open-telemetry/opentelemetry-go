// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/bridge/opencensus/opencensusmetric/internal"

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	ocmetricdata "go.opencensus.io/metric/metricdata"
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestConvertMetrics(t *testing.T) {
	endTime1 := time.Now()
	exemplarTime := endTime1.Add(-10 * time.Second)
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
			desc: "normal Histogram, summary, gauges, and sums",
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
										{
											Count: 1,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     0.8,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: octrace.SpanContext{
														TraceID: octrace.TraceID([16]byte{1}),
														SpanID:  octrace.SpanID([8]byte{2}),
													},
													"bool": true,
												},
											},
										},
										{
											Count: 2,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     1.5,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: octrace.SpanContext{
														TraceID: octrace.TraceID([16]byte{3}),
														SpanID:  octrace.SpanID([8]byte{4}),
													},
												},
											},
										},
										{
											Count: 5,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     2.6,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: octrace.SpanContext{
														TraceID: octrace.TraceID([16]byte{5}),
														SpanID:  octrace.SpanID([8]byte{6}),
													},
												},
											},
										},
									},
								}),
								ocmetricdata.NewDistributionPoint(endTime2, &ocmetricdata.Distribution{
									Count: 10,
									Sum:   110.0,
									BucketOptions: &ocmetricdata.BucketOptions{
										Bounds: []float64{1.0, 2.0, 3.0},
									},
									Buckets: []ocmetricdata.Bucket{
										{
											Count: 1,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     0.9,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: octrace.SpanContext{
														TraceID: octrace.TraceID([16]byte{7}),
														SpanID:  octrace.SpanID([8]byte{8}),
													},
												},
											},
										},
										{
											Count: 4,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     1.1,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: octrace.SpanContext{
														TraceID: octrace.TraceID([16]byte{9}),
														SpanID:  octrace.SpanID([8]byte{10}),
													},
												},
											},
										},
										{
											Count: 5,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     2.7,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: octrace.SpanContext{
														TraceID: octrace.TraceID([16]byte{11}),
														SpanID:  octrace.SpanID([8]byte{12}),
													},
												},
											},
										},
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
				}, {
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/summary-a",
						Description: "a testing summary",
						Unit:        ocmetricdata.UnitMilliseconds,
						Type:        ocmetricdata.TypeSummary,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "g"},
							{Key: "h"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "ding",
									Present: true,
								}, {
									Value:   "dong",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewSummaryPoint(endTime1, &ocmetricdata.Summary{
									Count:          10,
									Sum:            13.2,
									HasCountAndSum: true,
									Snapshot: ocmetricdata.Snapshot{
										Percentiles: map[float64]float64{
											50.0:  1.0,
											0.0:   0.1,
											100.0: 10.4,
										},
									},
								}),
								ocmetricdata.NewSummaryPoint(endTime2, &ocmetricdata.Summary{
									Count: 12,
									Snapshot: ocmetricdata.Snapshot{
										Percentiles: map[float64]float64{
											0.0:   0.2,
											50.0:  1.1,
											100.0: 10.5,
										},
									},
								}),
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
								Exemplars: []metricdata.Exemplar[float64]{
									{
										Time:    exemplarTime,
										Value:   0.8,
										TraceID: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
										SpanID:  []byte{2, 0, 0, 0, 0, 0, 0, 0},
										FilteredAttributes: []attribute.KeyValue{
											attribute.Bool("bool", true),
										},
									},
									{
										Time:    exemplarTime,
										Value:   1.5,
										TraceID: []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
										SpanID:  []byte{4, 0, 0, 0, 0, 0, 0, 0},
									},
									{
										Time:    exemplarTime,
										Value:   2.6,
										TraceID: []byte{5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
										SpanID:  []byte{6, 0, 0, 0, 0, 0, 0, 0},
									},
								},
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
								Exemplars: []metricdata.Exemplar[float64]{
									{
										Time:    exemplarTime,
										Value:   0.9,
										TraceID: []byte{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
										SpanID:  []byte{8, 0, 0, 0, 0, 0, 0, 0},
									},
									{
										Time:    exemplarTime,
										Value:   1.1,
										TraceID: []byte{9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
										SpanID:  []byte{10, 0, 0, 0, 0, 0, 0, 0},
									},
									{
										Time:    exemplarTime,
										Value:   2.7,
										TraceID: []byte{11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
										SpanID:  []byte{12, 0, 0, 0, 0, 0, 0, 0},
									},
								},
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
				}, {
					Name:        "foo.com/summary-a",
					Description: "a testing summary",
					Unit:        "ms",
					Data: metricdata.Summary{
						DataPoints: []metricdata.SummaryDataPoint{
							{
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("g"),
									Value: attribute.StringValue("ding"),
								}, attribute.KeyValue{
									Key:   attribute.Key("h"),
									Value: attribute.StringValue("dong"),
								}),
								Time:  endTime1,
								Count: 10,
								Sum:   13.2,
								QuantileValues: []metricdata.QuantileValue{
									{
										Quantile: 0.0,
										Value:    0.1,
									},
									{
										Quantile: 0.5,
										Value:    1.0,
									},
									{
										Quantile: 1.0,
										Value:    10.4,
									},
								},
							}, {
								Attributes: attribute.NewSet(attribute.KeyValue{
									Key:   attribute.Key("g"),
									Value: attribute.StringValue("ding"),
								}, attribute.KeyValue{
									Key:   attribute.Key("h"),
									Value: attribute.StringValue("dong"),
								}),
								Time:  endTime2,
								Count: 12,
								QuantileValues: []metricdata.QuantileValue{
									{
										Quantile: 0.0,
										Value:    0.2,
									},
									{
										Quantile: 0.5,
										Value:    1.1,
									},
									{
										Quantile: 1.0,
										Value:    10.5,
									},
								},
							},
						},
					},
				},
			},
		},
		{
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
		},
		{
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
		},
		{
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
		},
		{
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
			expectedErr: errNegativeCount,
		},
		{
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
			expectedErr: errNegativeBucketCount,
		},
		{
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
			expectedErr: errMismatchedValueTypes,
		},
		{
			desc: "summary with mismatched attributes",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/summary-mismatched",
						Description: "a mismatched summary",
						Unit:        ocmetricdata.UnitMilliseconds,
						Type:        ocmetricdata.TypeSummary,
						LabelKeys: []ocmetricdata.LabelKey{
							{Key: "g"},
						},
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							LabelValues: []ocmetricdata.LabelValue{
								{
									Value:   "ding",
									Present: true,
								}, {
									Value:   "dong",
									Present: true,
								},
							},
							Points: []ocmetricdata.Point{
								ocmetricdata.NewSummaryPoint(endTime1, &ocmetricdata.Summary{
									Count:          10,
									Sum:            13.2,
									HasCountAndSum: true,
									Snapshot: ocmetricdata.Snapshot{
										Percentiles: map[float64]float64{
											0.0: 0.1,
											0.5: 1.0,
											1.0: 10.4,
										},
									},
								}),
							},
						},
					},
				},
			},
			expectedErr: errMismatchedAttributeKeyValues,
		},
		{
			desc: "summary with negative count",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/summary-negative",
						Description: "a negative count summary",
						Unit:        ocmetricdata.UnitMilliseconds,
						Type:        ocmetricdata.TypeSummary,
					},
					TimeSeries: []*ocmetricdata.TimeSeries{
						{
							Points: []ocmetricdata.Point{
								ocmetricdata.NewSummaryPoint(endTime1, &ocmetricdata.Summary{
									Count:          -10,
									Sum:            13.2,
									HasCountAndSum: true,
									Snapshot: ocmetricdata.Snapshot{
										Percentiles: map[float64]float64{
											0.0: 0.1,
											0.5: 1.0,
											1.0: 10.4,
										},
									},
								}),
							},
						},
					},
				},
			},
			expectedErr: errNegativeCount,
		},
		{
			desc: "histogram with invalid span context exemplar",
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
									Count: 8,
									Sum:   100.0,
									BucketOptions: &ocmetricdata.BucketOptions{
										Bounds: []float64{1.0, 2.0, 3.0},
									},
									Buckets: []ocmetricdata.Bucket{
										{
											Count: 1,
											Exemplar: &ocmetricdata.Exemplar{
												Value:     0.8,
												Timestamp: exemplarTime,
												Attachments: map[string]interface{}{
													ocmetricdata.AttachmentKeySpanContext: "notaspancontext",
												},
											},
										},
									},
								}),
							},
							StartTime: startTime,
						},
					},
				},
			},
			expectedErr: errInvalidExemplarSpanContext,
		},
		{
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
			expectedErr: errMismatchedValueTypes,
		},
		{
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
			expectedErr: errMismatchedValueTypes,
		},
		{
			desc: "summary with non-summary datapoint type",
			input: []*ocmetricdata.Metric{
				{
					Descriptor: ocmetricdata.Descriptor{
						Name:        "foo.com/bad-point",
						Description: "a bad type",
						Unit:        ocmetricdata.UnitDimensionless,
						Type:        ocmetricdata.TypeSummary,
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
			expectedErr: errMismatchedValueTypes,
		},
		{
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
			expectedErr: errAggregationType,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			output, err := ConvertMetrics(tc.input)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("ConvertMetrics(%+v) = err(%v), want err(%v)", tc.input, err, tc.expectedErr)
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

type fakeStringer string

func (f fakeStringer) String() string {
	return string(f)
}

func TestConvertKV(t *testing.T) {
	key := "foo"
	for _, tt := range []struct {
		value    any
		expected attribute.Value
	}{
		{
			value:    bool(true),
			expected: attribute.BoolValue(true),
		},
		{
			value:    []bool{true, false},
			expected: attribute.BoolSliceValue([]bool{true, false}),
		},
		{
			value:    int(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []int{10, 20},
			expected: attribute.IntSliceValue([]int{10, 20}),
		},
		{
			value:    int8(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []int8{10, 20},
			expected: attribute.IntSliceValue([]int{10, 20}),
		},
		{
			value:    int16(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []int16{10, 20},
			expected: attribute.IntSliceValue([]int{10, 20}),
		},
		{
			value:    int32(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []int32{10, 20},
			expected: attribute.IntSliceValue([]int{10, 20}),
		},
		{
			value:    int64(10),
			expected: attribute.Int64Value(10),
		},
		{
			value:    []int64{10, 20},
			expected: attribute.Int64SliceValue([]int64{10, 20}),
		},
		{
			value:    uint(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    uint(math.MaxUint),
			expected: attribute.StringValue(fmt.Sprintf("%v", uint(math.MaxUint))),
		},
		{
			value:    []uint{10, 20},
			expected: attribute.StringSliceValue([]string{"10", "20"}),
		},
		{
			value:    uint8(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []uint8{10, 20},
			expected: attribute.StringSliceValue([]string{"10", "20"}),
		},
		{
			value:    uint16(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []uint16{10, 20},
			expected: attribute.StringSliceValue([]string{"10", "20"}),
		},
		{
			value:    uint32(10),
			expected: attribute.IntValue(10),
		},
		{
			value:    []uint32{10, 20},
			expected: attribute.StringSliceValue([]string{"10", "20"}),
		},
		{
			value:    uint64(10),
			expected: attribute.Int64Value(10),
		},
		{
			value:    uint64(math.MaxUint64),
			expected: attribute.StringValue("18446744073709551615"),
		},
		{
			value:    []uint64{10, 20},
			expected: attribute.StringSliceValue([]string{"10", "20"}),
		},
		{
			value:    uintptr(10),
			expected: attribute.Int64Value(10),
		},
		{
			value:    []uintptr{10, 20},
			expected: attribute.StringSliceValue([]string{"10", "20"}),
		},
		{
			value:    float32(10),
			expected: attribute.Float64Value(10),
		},
		{
			value:    []float32{10, 20},
			expected: attribute.Float64SliceValue([]float64{10, 20}),
		},
		{
			value:    float64(10),
			expected: attribute.Float64Value(10),
		},
		{
			value:    []float64{10, 20},
			expected: attribute.Float64SliceValue([]float64{10, 20}),
		},
		{
			value:    complex64(10),
			expected: attribute.StringValue("(10+0i)"),
		},
		{
			value:    []complex64{10, 20},
			expected: attribute.StringSliceValue([]string{"(10+0i)", "(20+0i)"}),
		},
		{
			value:    complex128(10),
			expected: attribute.StringValue("(10+0i)"),
		},
		{
			value:    []complex128{10, 20},
			expected: attribute.StringSliceValue([]string{"(10+0i)", "(20+0i)"}),
		},
		{
			value:    "string",
			expected: attribute.StringValue("string"),
		},
		{
			value:    []string{"string", "slice"},
			expected: attribute.StringSliceValue([]string{"string", "slice"}),
		},
		{
			value:    fakeStringer("stringer"),
			expected: attribute.StringValue("stringer"),
		},
		{
			value:    metricdata.Histogram[float64]{},
			expected: attribute.StringValue("unhandled attribute value: {DataPoints:[] Temporality:undefinedTemporality}"),
		},
	} {
		t.Run(fmt.Sprintf("%v(%+v)", reflect.TypeOf(tt.value), tt.value), func(t *testing.T) {
			got := convertKV(key, tt.value)
			assert.Equal(t, key, string(got.Key))
			assert.Equal(t, tt.expected, got.Value)
		})
	}
}

func BenchmarkConvertExemplar(b *testing.B) {
	const attchmentsN = 10
	data := make([]*ocmetricdata.Exemplar, b.N)
	for i := range data {
		a := make(ocmetricdata.Attachments, attchmentsN)
		for j := 0; j < attchmentsN; j++ {
			a[strconv.Itoa(j)] = rand.Int63()
		}
		data[i] = &ocmetricdata.Exemplar{
			Value:       rand.NormFloat64(),
			Timestamp:   time.Now(),
			Attachments: a,
		}
	}

	var out metricdata.Exemplar[float64]

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, _ = convertExemplar(data[n])
	}

	_ = out
}

func BenchmarkConvertQuantiles(b *testing.B) {
	const percentileN = 20
	data := make([]ocmetricdata.Snapshot, b.N)
	for i := range data {
		p := make(map[float64]float64, percentileN)
		for j := 0; j < percentileN; j++ {
			v := rand.Float64()
			for v == 0 {
				// Convert from [0, 1) interval to (0, 1).
				v = rand.Float64()
			}
			v *= 100 // Convert from (0, 1) interval to (0, 100).
			p[v] = rand.ExpFloat64()
		}
		data[i] = ocmetricdata.Snapshot{Percentiles: p}
	}

	var out []metricdata.QuantileValue

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out = convertQuantiles(data[n])
	}

	_ = out
}
