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

	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
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
			expectedKind: aggregation.ExactKind,
		},
		{
			desc: "float point",
			input: []metricdata.Point{
				{
					Time:  now,
					Value: float64(23),
				},
			},
			expectedKind: aggregation.ExactKind,
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
			desc: "dist is incompatible with exact",
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
			output, err := newAggregationFromPoints(tc.input)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("newAggregationFromPoints(%v) = err(%v), want err(%v)", tc.input, err, tc.expectedErr)
			}
			if tc.expectedErr == nil && output.Kind() != tc.expectedKind {
				t.Errorf("newAggregationFromPoints(%v) = %v, want %v", tc.input, output.Kind(), tc.expectedKind)
			}
		})
	}
}

func TestPointsAggregation(t *testing.T) {
	now := time.Now()
	input := []metricdata.Point{
		{Value: int64(15)},
		{Value: int64(-23), Time: now},
	}
	output, err := newAggregationFromPoints(input)
	if err != nil {
		t.Fatalf("newAggregationFromPoints(%v) = err(%v), want <nil>", input, err)
	}
	if output.Kind() != aggregation.ExactKind {
		t.Errorf("newAggregationFromPoints(%v) = %v, want %v", input, output.Kind(), aggregation.ExactKind)
	}
	if output.end() != now {
		t.Errorf("newAggregationFromPoints(%v).end() = %v, want %v", input, output.end(), now)
	}
	pointsAgg, ok := output.(aggregation.Points)
	if !ok {
		t.Errorf("newAggregationFromPoints(%v) = %v does not implement the aggregation.Points interface", input, output)
	}
	points, err := pointsAgg.Points()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if len(points) != len(input) {
		t.Fatalf("newAggregationFromPoints(%v) resulted in %d points, want %d points", input, len(points), len(input))
	}
	for i := range points {
		inputPoint := input[i]
		outputPoint := points[i]
		if inputPoint.Value != outputPoint.AsInt64() {
			t.Errorf("newAggregationFromPoints(%v)[%d] = %v, want %v", input, i, outputPoint.AsInt64(), inputPoint.Value)
		}
	}
}

func TestLastValueAggregation(t *testing.T) {
	now := time.Now()
	input := []metricdata.Point{
		{Value: int64(15)},
		{Value: int64(-23), Time: now},
	}
	output, err := newAggregationFromPoints(input)
	if err != nil {
		t.Fatalf("newAggregationFromPoints(%v) = err(%v), want <nil>", input, err)
	}
	if output.Kind() != aggregation.ExactKind {
		t.Errorf("newAggregationFromPoints(%v) = %v, want %v", input, output.Kind(), aggregation.ExactKind)
	}
	if output.end() != now {
		t.Errorf("newAggregationFromPoints(%v).end() = %v, want %v", input, output.end(), now)
	}
	lvAgg, ok := output.(aggregation.LastValue)
	if !ok {
		t.Errorf("newAggregationFromPoints(%v) = %v does not implement the aggregation.Points interface", input, output)
	}
	num, endTime, err := lvAgg.LastValue()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if endTime != now {
		t.Errorf("newAggregationFromPoints(%v).LastValue() = endTime: %v, want %v", input, endTime, now)
	}
	if num.AsInt64() != int64(-23) {
		t.Errorf("newAggregationFromPoints(%v).LastValue() = number: %v, want %v", input, num.AsInt64(), int64(-23))
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
	output, err := newAggregationFromPoints(input)
	if err != nil {
		t.Fatalf("newAggregationFromPoints(%v) = err(%v), want <nil>", input, err)
	}
	if output.Kind() != aggregation.HistogramKind {
		t.Errorf("newAggregationFromPoints(%v) = %v, want %v", input, output.Kind(), aggregation.HistogramKind)
	}
	if output.end() != now {
		t.Errorf("newAggregationFromPoints(%v).end() = %v, want %v", input, output.end(), now)
	}
	distAgg, ok := output.(aggregation.Histogram)
	if !ok {
		t.Errorf("newAggregationFromPoints(%v) = %v does not implement the aggregation.Points interface", input, output)
	}
	sum, err := distAgg.Sum()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if sum.AsFloat64() != float64(55) {
		t.Errorf("newAggregationFromPoints(%v).Sum() = %v, want %v", input, sum.AsFloat64(), float64(55))
	}
	count, err := distAgg.Count()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if count != 2 {
		t.Errorf("newAggregationFromPoints(%v).Count() = %v, want %v", input, count, 2)
	}
	hist, err := distAgg.Histogram()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	inputBucketBoundaries := []float64{20, 30}
	if len(hist.Boundaries) != len(inputBucketBoundaries) {
		t.Fatalf("newAggregationFromPoints(%v).Histogram() produced %d boundaries, want %d boundaries", input, len(hist.Boundaries), len(inputBucketBoundaries))
	}
	for i, b := range hist.Boundaries {
		if b != inputBucketBoundaries[i] {
			t.Errorf("newAggregationFromPoints(%v).Histogram().Boundaries[%d] = %v, want %v", input, i, b, inputBucketBoundaries[i])
		}
	}
	inputBucketCounts := []uint64{1, 1}
	if len(hist.Counts) != len(inputBucketCounts) {
		t.Fatalf("newAggregationFromPoints(%v).Histogram() produced %d buckets, want %d buckets", input, len(hist.Counts), len(inputBucketCounts))
	}
	for i, c := range hist.Counts {
		if c != inputBucketCounts[i] {
			t.Errorf("newAggregationFromPoints(%v).Histogram().Counts[%d] = %d, want %d", input, i, c, inputBucketCounts[i])
		}
	}
}
