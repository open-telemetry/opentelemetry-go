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

package aggregate

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type noErrorHandler struct{ t *testing.T }

func (h *noErrorHandler) Handle(e error) {
	require.NoError(h.t, e)
}

func withHandler(t *testing.T) func() {
	t.Helper()
	h := &noErrorHandler{t: t}
	original := global.GetErrorHandler()
	global.SetErrorHandler(h)
	return func() { global.SetErrorHandler(original) }
}

func TestExpoHistogramDataPointRecord(t *testing.T) {
	t.Run("float64", testExpoHistogramDataPointRecord[float64])
	t.Run("float64 MinMaxSum", testExpoHistogramMinMaxSumFloat64)
	t.Run("float64-2", testExpoHistogramDataPointRecordFloat64)
	t.Run("int64", testExpoHistogramDataPointRecord[int64])
	t.Run("int64 MinMaxSum", testExpoHistogramMinMaxSumInt64)
}

func testExpoHistogramDataPointRecord[N int64 | float64](t *testing.T) {
	testCases := []struct {
		maxSize         int
		values          []N
		expectedBuckets expoBuckets
		expectedScale   int
	}{
		{
			maxSize: 4,
			values:  []N{2, 4, 1},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 1, 1},
			},
			expectedScale: 0,
		},
		{
			maxSize: 4,
			values:  []N{4, 4, 4, 2, 16, 1},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{1, 4, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{1, 2, 4},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{1, 4, 2},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{2, 4, 1},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{2, 1, 4},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{4, 1, 2},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{4, 2, 1},
			expectedBuckets: expoBuckets{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			dp := newExpoHistogramDataPoint[N](tt.maxSize, 20, false, false)
			for _, v := range tt.values {
				dp.record(v)
				dp.record(-v)
			}

			assert.Equal(t, tt.expectedBuckets, dp.posBuckets, "positive buckets")
			assert.Equal(t, tt.expectedBuckets, dp.negBuckets, "negative buckets")
			assert.Equal(t, tt.expectedScale, dp.scale, "scale")
		})
	}
}

// TODO: This can be defined in the test after we drop support for go1.19.
type expectedMinMaxSum[N int64 | float64] struct {
	min   N
	max   N
	sum   N
	count uint
}
type expoHistogramDataPointRecordMinMaxSumTestCase[N int64 | float64] struct {
	values   []N
	expected expectedMinMaxSum[N]
}

func testExpoHistogramMinMaxSumInt64(t *testing.T) {
	testCases := []expoHistogramDataPointRecordMinMaxSumTestCase[int64]{
		{
			values:   []int64{2, 4, 1},
			expected: expectedMinMaxSum[int64]{1, 4, 7, 3},
		},
		{
			values:   []int64{4, 4, 4, 2, 16, 1},
			expected: expectedMinMaxSum[int64]{1, 16, 31, 6},
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			h := newExponentialHistogram[int64](4, 20, false, false, 0)
			for _, v := range tt.values {
				h.measure(context.Background(), v, alice)
			}
			dp := h.values[alice]

			assert.Equal(t, tt.expected.max, dp.max)
			assert.Equal(t, tt.expected.min, dp.min)
			assert.Equal(t, tt.expected.sum, dp.sum)
		})
	}
}

func testExpoHistogramMinMaxSumFloat64(t *testing.T) {
	testCases := []expoHistogramDataPointRecordMinMaxSumTestCase[float64]{
		{
			values:   []float64{2, 4, 1},
			expected: expectedMinMaxSum[float64]{1, 4, 7, 3},
		},
		{
			values:   []float64{2, 4, 1, math.Inf(1)},
			expected: expectedMinMaxSum[float64]{1, 4, 7, 4},
		},
		{
			values:   []float64{2, 4, 1, math.Inf(-1)},
			expected: expectedMinMaxSum[float64]{1, 4, 7, 4},
		},
		{
			values:   []float64{2, 4, 1, math.NaN()},
			expected: expectedMinMaxSum[float64]{1, 4, 7, 4},
		},
		{
			values:   []float64{4, 4, 4, 2, 16, 1},
			expected: expectedMinMaxSum[float64]{1, 16, 31, 6},
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			h := newExponentialHistogram[float64](4, 20, false, false, 0)
			for _, v := range tt.values {
				h.measure(context.Background(), v, alice)
			}
			dp := h.values[alice]

			assert.Equal(t, tt.expected.max, dp.max)
			assert.Equal(t, tt.expected.min, dp.min)
			assert.Equal(t, tt.expected.sum, dp.sum)
		})
	}
}

func testExpoHistogramDataPointRecordFloat64(t *testing.T) {
	type TestCase struct {
		maxSize         int
		values          []float64
		expectedBuckets expoBuckets
		expectedScale   int
	}

	testCases := []TestCase{
		{
			maxSize: 4,
			values:  []float64{2, 2, 2, 1, 8, 0.5},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 3, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{1, 0.5, 2},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{1, 2, 0.5},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{2, 0.5, 1},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{2, 1, 0.5},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{0.5, 1, 2},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{0.5, 2, 1},
			expectedBuckets: expoBuckets{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			dp := newExpoHistogramDataPoint[float64](tt.maxSize, 20, false, false)
			for _, v := range tt.values {
				dp.record(v)
				dp.record(-v)
			}

			assert.Equal(t, tt.expectedBuckets, dp.posBuckets)
			assert.Equal(t, tt.expectedBuckets, dp.negBuckets)
			assert.Equal(t, tt.expectedScale, dp.scale)
		})
	}
}

func TestExponentialHistogramDataPointRecordLimits(t *testing.T) {
	// These bins are calculated from the following formula:
	// floor( log2( value) * 2^20 ) using an arbitrary precision calculator.

	fdp := newExpoHistogramDataPoint[float64](4, 20, false, false)
	fdp.record(math.MaxFloat64)

	if fdp.posBuckets.startBin != 1073741823 {
		t.Errorf("Expected startBin to be 1073741823, got %d", fdp.posBuckets.startBin)
	}

	fdp = newExpoHistogramDataPoint[float64](4, 20, false, false)
	fdp.record(math.SmallestNonzeroFloat64)

	if fdp.posBuckets.startBin != -1126170625 {
		t.Errorf("Expected startBin to be -1126170625, got %d", fdp.posBuckets.startBin)
	}

	idp := newExpoHistogramDataPoint[int64](4, 20, false, false)
	idp.record(math.MaxInt64)

	if idp.posBuckets.startBin != 66060287 {
		t.Errorf("Expected startBin to be 66060287, got %d", idp.posBuckets.startBin)
	}
}

func TestExpoBucketDownscale(t *testing.T) {
	tests := []struct {
		name   string
		bucket *expoBuckets
		scale  int
		want   *expoBuckets
	}{
		{
			name:   "Empty bucket",
			bucket: &expoBuckets{},
			scale:  3,
			want:   &expoBuckets{},
		},
		{
			name: "1 size bucket",
			bucket: &expoBuckets{
				startBin: 50,
				counts:   []uint64{7},
			},
			scale: 4,
			want: &expoBuckets{
				startBin: 3,
				counts:   []uint64{7},
			},
		},
		{
			name: "zero scale",
			bucket: &expoBuckets{
				startBin: 50,
				counts:   []uint64{7, 5},
			},
			scale: 0,
			want: &expoBuckets{
				startBin: 50,
				counts:   []uint64{7, 5},
			},
		},
		{
			name: "aligned bucket scale 1",
			bucket: &expoBuckets{
				startBin: 0,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 1,
			want: &expoBuckets{
				startBin: 0,
				counts:   []uint64{3, 7, 11},
			},
		},
		{
			name: "aligned bucket scale 2",
			bucket: &expoBuckets{
				startBin: 0,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 2,
			want: &expoBuckets{
				startBin: 0,
				counts:   []uint64{10, 11},
			},
		},
		{
			name: "aligned bucket scale 3",
			bucket: &expoBuckets{
				startBin: 0,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 3,
			want: &expoBuckets{
				startBin: 0,
				counts:   []uint64{21},
			},
		},
		{
			name: "unaligned bucket scale 1",
			bucket: &expoBuckets{
				startBin: 5,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,1,2,3,4,5,6]
			scale: 1,
			want: &expoBuckets{
				startBin: 2,
				counts:   []uint64{1, 5, 9, 6},
			}, // This is equivalent to [0,0,1,5,9,6]
		},
		{
			name: "unaligned bucket scale 2",
			bucket: &expoBuckets{
				startBin: 7,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,0,0,1,2,3,4,5,6]
			scale: 2,
			want: &expoBuckets{
				startBin: 1,
				counts:   []uint64{1, 14, 6},
			}, // This is equivalent to [0,1,14,6]
		},
		{
			name: "unaligned bucket scale 3",
			bucket: &expoBuckets{
				startBin: 3,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,1,2,3,4,5,6]
			scale: 3,
			want: &expoBuckets{
				startBin: 0,
				counts:   []uint64{15, 6},
			}, // This is equivalent to [0,15,6]
		},
		{
			name: "unaligned bucket scale 1",
			bucket: &expoBuckets{
				startBin: 1,
				counts:   []uint64{1, 0, 1},
			},
			scale: 1,
			want: &expoBuckets{
				startBin: 0,
				counts:   []uint64{1, 1},
			},
		},
		{
			name: "negative startBin",
			bucket: &expoBuckets{
				startBin: -1,
				counts:   []uint64{1, 0, 3},
			},
			scale: 1,
			want: &expoBuckets{
				startBin: -1,
				counts:   []uint64{1, 3},
			},
		},
		{
			name: "negative startBin 2",
			bucket: &expoBuckets{
				startBin: -4,
				counts:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			scale: 1,
			want: &expoBuckets{
				startBin: -2,
				counts:   []uint64{3, 7, 11, 15, 19},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bucket.downscale(tt.scale)

			assert.Equal(t, tt.want, tt.bucket)
		})
	}
}

func TestExpoBucketRecord(t *testing.T) {
	tests := []struct {
		name   string
		bucket *expoBuckets
		bin    int
		want   *expoBuckets
	}{
		{
			name:   "Empty Bucket creates first count",
			bucket: &expoBuckets{},
			bin:    -5,
			want: &expoBuckets{
				startBin: -5,
				counts:   []uint64{1},
			},
		},
		{
			name: "Bin is in the bucket",
			bucket: &expoBuckets{
				startBin: 3,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			bin: 5,
			want: &expoBuckets{
				startBin: 3,
				counts:   []uint64{1, 2, 4, 4, 5, 6},
			},
		},
		{
			name: "Bin is before the start of the bucket",
			bucket: &expoBuckets{
				startBin: 1,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			bin: -2,
			want: &expoBuckets{
				startBin: -2,
				counts:   []uint64{1, 0, 0, 1, 2, 3, 4, 5, 6},
			},
		},
		{
			name: "Bin is after the end of the bucket",
			bucket: &expoBuckets{
				startBin: -2,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			bin: 4,
			want: &expoBuckets{
				startBin: -2,
				counts:   []uint64{1, 2, 3, 4, 5, 6, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bucket.record(tt.bin)

			assert.Equal(t, tt.want, tt.bucket)
		})
	}
}

func TestScaleChange(t *testing.T) {
	type args struct {
		bin      int
		startBin int
		length   int
		maxSize  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "if length is 0, no rescale is needed",
			// [] -> [5] Length 1
			args: args{
				bin:      5,
				startBin: 0,
				length:   0,
				maxSize:  4,
			},
			want: 0,
		},
		{
			name: "if bin is between start, and the end, no rescale needed",
			// [-1, ..., 8] Length 10 -> [-1, ..., 5, ..., 8] Length 10
			args: args{
				bin:      5,
				startBin: -1,
				length:   10,
				maxSize:  20,
			},
			want: 0,
		},
		{
			name: "if len([bin,... end]) > maxSize, rescale needed",
			// [8,9,10] Length 3 -> [5, ..., 10] Length 6
			args: args{
				bin:      5,
				startBin: 8,
				length:   3,
				maxSize:  5,
			},
			want: 1,
		},
		{
			name: "if len([start, ..., bin]) > maxSize, rescale needed",
			// [2,3,4] Length 3 -> [2, ..., 7] Length 6
			args: args{
				bin:      7,
				startBin: 2,
				length:   3,
				maxSize:  5,
			},
			want: 1,
		},
		{
			name: "if len([start, ..., bin]) > maxSize, rescale needed",
			// [2,3,4] Length 3 -> [2, ..., 7] Length 12
			args: args{
				bin:      13,
				startBin: 2,
				length:   3,
				maxSize:  5,
			},
			want: 2,
		},
		{
			name: "It should not hang if it will never be able to rescale",
			args: args{
				bin:      1,
				startBin: -1,
				length:   1,
				maxSize:  1,
			},
			want: 31,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newExpoHistogramDataPoint[float64](tt.args.maxSize, 20, false, false)
			got := p.scaleChange(tt.args.bin, tt.args.startBin, tt.args.length)
			if got != tt.want {
				t.Errorf("scaleChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPrepend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		agg := newExpoHistogramDataPoint[float64](1024, 20, false, false)
		n := math.MaxFloat64
		for j := 0; j < 1024; j++ {
			agg.record(n)
			n = n / 2
		}
	}
}

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		agg := newExpoHistogramDataPoint[float64](1024, 20, false, false)
		n := smallestNonZeroNormalFloat64
		for j := 0; j < 1024; j++ {
			agg.record(n)
			n = n * 2
		}
	}
}

func BenchmarkExponentialHistogram(b *testing.B) {
	const (
		maxSize  = 160
		maxScale = 20
		noMinMax = false
		noSum    = false
	)

	b.Run("Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.ExponentialBucketHistogram(maxSize, maxScale, noMinMax, noSum)
	}))
	b.Run("Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.ExponentialBucketHistogram(maxSize, maxScale, noMinMax, noSum)
	}))
	b.Run("Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.ExponentialBucketHistogram(maxSize, maxScale, noMinMax, noSum)
	}))
	b.Run("Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.ExponentialBucketHistogram(maxSize, maxScale, noMinMax, noSum)
	}))
}

func TestSubNormal(t *testing.T) {
	want := &expoHistogramDataPoint[float64]{
		maxSize: 4,
		count:   3,
		min:     math.SmallestNonzeroFloat64,
		max:     math.SmallestNonzeroFloat64,
		sum:     3 * math.SmallestNonzeroFloat64,

		scale: 20,
		posBuckets: expoBuckets{
			startBin: -1126170625,
			counts:   []uint64{3},
		},
	}

	ehdp := newExpoHistogramDataPoint[float64](4, 20, false, false)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)

	assert.Equal(t, want, ehdp)
}

func TestExponentialHistogramAggregation(t *testing.T) {
	t.Cleanup(mockTime(now))

	t.Run("Int64/Delta", testDeltaExpoHist[int64]())
	t.Run("Float64/Delta", testDeltaExpoHist[float64]())
	t.Run("Int64/Cumulative", testCumulativeExpoHist[int64]())
	t.Run("Float64/Cumulative", testCumulativeExpoHist[float64]())
}

func testDeltaExpoHist[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 2,
	}.ExponentialBucketHistogram(4, 20, false, false)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.ExponentialHistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 4, alice},
				{ctx, 4, alice},
				{ctx, 4, alice},
				{ctx, 2, alice},
				{ctx, 16, alice},
				{ctx, 1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      6,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        31,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
						},
					},
				},
			},
		},
		{
			// Delta sums are expected to reset.
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.ExponentialHistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 4, alice},
				{ctx, 4, alice},
				{ctx, 4, alice},
				{ctx, 2, alice},
				{ctx, 16, alice},
				{ctx, 1, alice},
				// These will exceed the cardinality limit.
				{ctx, 4, bob},
				{ctx, 4, bob},
				{ctx, 4, bob},
				{ctx, 2, carol},
				{ctx, 16, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 2,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      6,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        31,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
						},
						{
							Attributes: overflowSet,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      6,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        31,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
						},
					},
				},
			},
		},
	})
}

func testCumulativeExpoHist[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 2,
	}.ExponentialBucketHistogram(4, 20, false, false)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.ExponentialHistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 4, alice},
				{ctx, 4, alice},
				{ctx, 4, alice},
				{ctx, 2, alice},
				{ctx, 16, alice},
				{ctx, 1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      6,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        31,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 2, alice},
				{ctx, 3, alice},
				{ctx, 8, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      9,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        44,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 6, 2},
							},
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			expect: output{
				n: 1,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      9,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        44,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 6, 2},
							},
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				// These will exceed the cardinality limit.
				{ctx, 4, bob},
				{ctx, 4, bob},
				{ctx, 4, bob},
				{ctx, 2, carol},
				{ctx, 16, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 2,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      9,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        44,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 6, 2},
							},
						},
						{
							Attributes: overflowSet,
							StartTime:  staticTime,
							Time:       staticTime,
							Count:      6,
							Min:        metricdata.NewExtrema[N](1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        31,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
						},
					},
				},
			},
		},
	})
}

func FuzzGetBin(f *testing.F) {
	values := []float64{
		2.0,
		0x1p35,
		0x1.0000000000001p35,
		0x1.fffffffffffffp34,
		0x1p300,
		0x1.0000000000001p300,
		0x1.fffffffffffffp299,
	}
	scales := []int{0, 15, -5}

	for _, s := range scales {
		for _, v := range values {
			f.Add(v, s)
		}
	}

	f.Fuzz(func(t *testing.T, v float64, scale int) {
		// GetBin only works on positive values.
		if math.Signbit(v) {
			v = v * -1
		}
		// GetBin Doesn't work on zero.
		if v == 0.0 {
			t.Skip("skipping test for zero")
		}

		p := newExpoHistogramDataPoint[float64](4, 20, false, false)
		// scale range is -10 to 20.
		p.scale = (scale%31+31)%31 - 10
		got := p.getBin(v)
		if v <= lowerBound(got, p.scale) {
			t.Errorf("v=%x scale =%d had bin %d, but was below lower bound %x", v, p.scale, got, lowerBound(got, p.scale))
		}
		if v > lowerBound(got+1, p.scale) {
			t.Errorf("v=%x scale =%d had bin %d, but was above upper bound %x", v, p.scale, got, lowerBound(got+1, p.scale))
		}
	})
}

func lowerBound(index int, scale int) float64 {
	// The lowerBound of the index of Math.SmallestNonzeroFloat64 at any scale
	// is always rounded down to 0.0.
	// For example lowerBound(getBin(Math.SmallestNonzeroFloat64, 7), 7) == 0.0
	// 2 ^ (index * 2 ^ (-scale))
	return math.Exp2(math.Ldexp(float64(index), -scale))
}
