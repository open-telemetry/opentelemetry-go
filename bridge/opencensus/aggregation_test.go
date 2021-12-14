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

package opencensus

import (
	"errors"
	"testing"
	"time"

	"go.opencensus.io/metric/metricdata"

	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

func TestNewAggregationFromPoints(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		desc         string
		input        []metricdata.Point
		expectedKind aggregation.Kind
		expectedErr  error
	}{
		{
			desc:        "no points",
			expectedErr: errEmpty,
		},
		{
			desc: "int point",
			input: []metricdata.Point{
				{
					Time:  now,
					Value: int64(23),
				},
			},
			expectedKind: aggregation.LastValueKind,
		},
		{
			desc: "float point",
			input: []metricdata.Point{
				{
					Time:  now,
					Value: float64(23),
				},
			},
			expectedKind: aggregation.LastValueKind,
		},
		{
			desc: "distribution point",
			input: []metricdata.Point{
				{
					Time: now,
					Value: &metricdata.Distribution{
						Count: 2,
						Sum:   55,
						BucketOptions: &metricdata.BucketOptions{
							Bounds: []float64{20, 30},
						},
						Buckets: []metricdata.Bucket{
							{Count: 1},
							{Count: 1},
						},
					},
				},
			},
			expectedKind: aggregation.HistogramKind,
		},
		{
			desc: "bad distribution bucket count",
			input: []metricdata.Point{
				{
					Time: now,
					Value: &metricdata.Distribution{
						Count: 2,
						Sum:   55,
						BucketOptions: &metricdata.BucketOptions{
							Bounds: []float64{20, 30},
						},
						Buckets: []metricdata.Bucket{
							// negative bucket
							{Count: -1},
							{Count: 1},
						},
					},
				},
			},
			expectedErr: errBadPoint,
		},
		{
			desc: "bad distribution count",
			input: []metricdata.Point{
				{
					Time: now,
					Value: &metricdata.Distribution{
						// negative count
						Count: -2,
						Sum:   55,
						BucketOptions: &metricdata.BucketOptions{
							Bounds: []float64{20, 30},
						},
						Buckets: []metricdata.Bucket{
							{Count: 1},
							{Count: 1},
						},
					},
				},
			},
			expectedErr: errBadPoint,
		},
		{
			desc: "incompatible point type bool",
			input: []metricdata.Point{
				{
					Time:  now,
					Value: true,
				},
			},
			expectedErr: errIncompatibleType,
		},
		{
			desc: "dist is incompatible with raw points",
			input: []metricdata.Point{
				{
					Time:  now,
					Value: int64(23),
				},
				{
					Time: now,
					Value: &metricdata.Distribution{
						Count: 2,
						Sum:   55,
						BucketOptions: &metricdata.BucketOptions{
							Bounds: []float64{20, 30},
						},
						Buckets: []metricdata.Bucket{
							{Count: 1},
							{Count: 1},
						},
					},
				},
			},
			expectedErr: errIncompatibleType,
		},
		{
			desc: "int point is incompatible with dist",
			input: []metricdata.Point{
				{
					Time: now,
					Value: &metricdata.Distribution{
						Count: 2,
						Sum:   55,
						BucketOptions: &metricdata.BucketOptions{
							Bounds: []float64{20, 30},
						},
						Buckets: []metricdata.Bucket{
							{Count: 1},
							{Count: 1},
						},
					},
				},
				{
					Time:  now,
					Value: int64(23),
				},
			},
			expectedErr: errBadPoint,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			var output []aggregation.Aggregation
			err := recordAggregationsFromPoints(tc.input, func(agg aggregation.Aggregation, ts time.Time) error {
				last := tc.input[len(tc.input)-1]
				if ts != last.Time {
					t.Errorf("incorrect timestamp %v != %v", ts, last.Time)
				}
				output = append(output, agg)
				return nil
			})
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("newAggregationFromPoints(%v) = err(%v), want err(%v)", tc.input, err, tc.expectedErr)
			}
			for _, out := range output {
				if tc.expectedErr == nil && out.Kind() != tc.expectedKind {
					t.Errorf("newAggregationFromPoints(%v) = %v, want %v", tc.input, out.Kind(), tc.expectedKind)
				}
			}
		})
	}
}

func TestLastValueAggregation(t *testing.T) {
	now := time.Now()
	input := []metricdata.Point{
		{Value: int64(15), Time: now.Add(-time.Minute)},
		{Value: int64(-23), Time: now},
	}
	idx := 0
	err := recordAggregationsFromPoints(input, func(agg aggregation.Aggregation, end time.Time) error {
		if agg.Kind() != aggregation.LastValueKind {
			t.Errorf("recordAggregationsFromPoints(%v) = %v, want %v", input, agg.Kind(), aggregation.LastValueKind)
		}
		if end != input[idx].Time {
			t.Errorf("recordAggregationsFromPoints(%v).end() = %v, want %v", input, end, input[idx].Time)
		}
		pointsLV, ok := agg.(aggregation.LastValue)
		if !ok {
			t.Errorf("recordAggregationsFromPoints(%v) = %v does not implement the aggregation.LastValue interface", input, agg)
		}
		lv, ts, _ := pointsLV.LastValue()
		if lv.AsInt64() != input[idx].Value {
			t.Errorf("recordAggregationsFromPoints(%v) = %v, want %v", input, lv.AsInt64(), input[idx].Value)
		}
		if ts != input[idx].Time {
			t.Errorf("recordAggregationsFromPoints(%v) = %v, want %v", input, ts, input[idx].Time)
		}
		idx++
		return nil
	})
	if err != nil {
		t.Errorf("recordAggregationsFromPoints(%v) = unexpected error %v", input, err)
	}
}

func TestHistogramAggregation(t *testing.T) {
	now := time.Now()
	input := []metricdata.Point{
		{
			Value: &metricdata.Distribution{
				Count: 0,
				Sum:   0,
				BucketOptions: &metricdata.BucketOptions{
					Bounds: []float64{20, 30},
				},
				Buckets: []metricdata.Bucket{
					{Count: 0},
					{Count: 0},
				},
			},
		},
		{
			Time: now,
			Value: &metricdata.Distribution{
				Count: 2,
				Sum:   55,
				BucketOptions: &metricdata.BucketOptions{
					Bounds: []float64{20, 30},
				},
				Buckets: []metricdata.Bucket{
					{Count: 1},
					{Count: 1},
				},
			},
		},
	}
	var output aggregation.Aggregation
	var end time.Time
	err := recordAggregationsFromPoints(input, func(argAgg aggregation.Aggregation, argEnd time.Time) error {
		output = argAgg
		end = argEnd
		return nil
	})
	if err != nil {
		t.Fatalf("recordAggregationsFromPoints(%v) = err(%v), want <nil>", input, err)
	}
	if output.Kind() != aggregation.HistogramKind {
		t.Errorf("recordAggregationsFromPoints(%v) = %v, want %v", input, output.Kind(), aggregation.HistogramKind)
	}
	if end != now {
		t.Errorf("recordAggregationsFromPoints(%v).end() = %v, want %v", input, end, now)
	}
	distAgg, ok := output.(aggregation.Histogram)
	if !ok {
		t.Errorf("recordAggregationsFromPoints(%v) = %v does not implement the aggregation.Points interface", input, output)
	}
	sum, err := distAgg.Sum()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if sum.AsFloat64() != float64(55) {
		t.Errorf("recordAggregationsFromPoints(%v).Sum() = %v, want %v", input, sum.AsFloat64(), float64(55))
	}
	count, err := distAgg.Count()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if count != 2 {
		t.Errorf("recordAggregationsFromPoints(%v).Count() = %v, want %v", input, count, 2)
	}
	hist, err := distAgg.Histogram()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	inputBucketBoundaries := []float64{20, 30}
	if len(hist.Boundaries) != len(inputBucketBoundaries) {
		t.Fatalf("recordAggregationsFromPoints(%v).Histogram() produced %d boundaries, want %d boundaries", input, len(hist.Boundaries), len(inputBucketBoundaries))
	}
	for i, b := range hist.Boundaries {
		if b != inputBucketBoundaries[i] {
			t.Errorf("recordAggregationsFromPoints(%v).Histogram().Boundaries[%d] = %v, want %v", input, i, b, inputBucketBoundaries[i])
		}
	}
	inputBucketCounts := []uint64{1, 1}
	if len(hist.Counts) != len(inputBucketCounts) {
		t.Fatalf("recordAggregationsFromPoints(%v).Histogram() produced %d buckets, want %d buckets", input, len(hist.Counts), len(inputBucketCounts))
	}
	for i, c := range hist.Counts {
		if c != inputBucketCounts[i] {
			t.Errorf("recordAggregationsFromPoints(%v).Histogram().Counts[%d] = %d, want %d", input, i, c, inputBucketCounts[i])
		}
	}
}
