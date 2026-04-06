// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/internal/x"
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
		maxSize        int
		values         []N
		expectedStart  int32
		expectedCounts []uint64
		expectedScale  int32
	}{
		{
			maxSize:        4,
			values:         []N{2, 4, 1},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 1, 1},
			expectedScale:  0,
		},
		{
			maxSize:        4,
			values:         []N{4, 4, 4, 2, 16, 1},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 4, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []N{1, 2, 4},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 2},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []N{1, 4, 2},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 2},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []N{2, 4, 1},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 2},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []N{2, 1, 4},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 2},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []N{4, 1, 2},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 2},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []N{4, 2, 1},
			expectedStart:  -1,
			expectedCounts: []uint64{1, 2},
			expectedScale:  -1,
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			dp := newExpoHistogramDataPoint[N](alice, tt.maxSize, 20, false, false)
			for _, v := range tt.values {
				dp.record(v)
				dp.record(-v)
			}

			assertBuckets(t, tt.expectedStart, tt.expectedCounts, dp.posBuckets, "positive buckets")
			assertBuckets(t, tt.expectedStart, tt.expectedCounts, dp.negBuckets, "negative buckets")
			assert.Equal(t, tt.expectedScale, dp.scale.Load(), "scale")
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

			h := newCumulativeExpoHistogram[int64](4, 20, false, false, 0, dropExemplars[int64])
			for _, v := range tt.values {
				h.measure(t.Context(), v, alice, nil)
			}
			v, _ := h.values.Load(alice.Equivalent())
			dp := v.(*cumulativePoint[int64]).points[0]

			assert.Equal(t, tt.expected.max, dp.minMax.maximum.Load())
			assert.Equal(t, tt.expected.min, dp.minMax.minimum.Load())
			assert.InDelta(t, tt.expected.sum, dp.sum.load(), 0.01)
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

			h := newCumulativeExpoHistogram[float64](4, 20, false, false, 0, dropExemplars[float64])
			for _, v := range tt.values {
				h.measure(t.Context(), v, alice, nil)
			}
			v, _ := h.values.Load(alice.Equivalent())
			dp := v.(*cumulativePoint[float64]).points[0]

			assert.Equal(t, tt.expected.max, dp.minMax.maximum.Load())
			assert.Equal(t, tt.expected.min, dp.minMax.minimum.Load())
			assert.InDelta(t, tt.expected.sum, dp.sum.load(), 0.01)
		})
	}
}

func testExpoHistogramDataPointRecordFloat64(t *testing.T) {
	type TestCase struct {
		maxSize        int
		values         []float64
		expectedStart  int32
		expectedCounts []uint64
		expectedScale  int32
	}

	testCases := []TestCase{
		{
			maxSize:        4,
			values:         []float64{2, 2, 2, 1, 8, 0.5},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 3, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []float64{1, 0.5, 2},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []float64{1, 2, 0.5},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []float64{2, 0.5, 1},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []float64{2, 1, 0.5},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []float64{0.5, 1, 2},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 1},
			expectedScale:  -1,
		},
		{
			maxSize:        2,
			values:         []float64{0.5, 2, 1},
			expectedStart:  -1,
			expectedCounts: []uint64{2, 1},
			expectedScale:  -1,
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			dp := newExpoHistogramDataPoint[float64](alice, tt.maxSize, 20, false, false)
			for _, v := range tt.values {
				dp.record(v)
				dp.record(-v)
			}

			assertBuckets(t, tt.expectedStart, tt.expectedCounts, dp.posBuckets, "positive buckets")
			assertBuckets(t, tt.expectedStart, tt.expectedCounts, dp.negBuckets, "negative buckets")
			assert.Equal(t, tt.expectedScale, dp.scale.Load(), "scale")
			assert.Equal(t, tt.expectedScale, dp.scale.Load())
		})
	}
}

func TestExponentialHistogramDataPointRecordLimits(t *testing.T) {
	// These bins are calculated from the following formula:
	// floor( log2( value) * 2^20 ) using an arbitrary precision calculator.

	fdp := newExpoHistogramDataPoint[float64](alice, 4, 20, false, false)
	fdp.record(math.MaxFloat64)

	if fdp.posBuckets.startBin != 1073741823 {
		t.Errorf("Expected startBin to be 1073741823, got %d", fdp.posBuckets.startBin)
	}

	fdp = newExpoHistogramDataPoint[float64](alice, 4, 20, false, false)
	fdp.record(math.SmallestNonzeroFloat64)

	if fdp.posBuckets.startBin != -1126170625 {
		t.Errorf("Expected startBin to be -1126170625, got %d", fdp.posBuckets.startBin)
	}

	idp := newExpoHistogramDataPoint[int64](alice, 4, 20, false, false)
	idp.record(math.MaxInt64)

	if idp.posBuckets.startBin != 66060287 {
		t.Errorf("Expected startBin to be 66060287, got %d", idp.posBuckets.startBin)
	}
}

func newBucket(startBin int32, counts []uint64) *expoBuckets {
	b := &expoBuckets{startBin: startBin, counts: make([]atomic.Uint64, len(counts))}
	for i, v := range counts {
		b.counts[i].Store(v)
	}
	return b
}

func assertBuckets(t *testing.T, expectedStart int32, expectedCounts []uint64, actual expoBuckets, msg string) {
	t.Helper()
	assert.Equal(t, expectedStart, actual.startBin, "%s: startBin", msg)
	var actualCounts []uint64
	if len(actual.counts) > 0 {
		actualCounts = make([]uint64, len(actual.counts))
		for i := range actual.counts {
			actualCounts[i] = actual.counts[i].Load()
		}
	}
	assert.Equal(t, expectedCounts, actualCounts, "%s: counts", msg)
}

func TestExpoBucketDownscale(t *testing.T) {
	tests := []struct {
		name       string
		bucket     *expoBuckets
		scale      int32
		wantStart  int32
		wantCounts []uint64
	}{
		{
			name:       "Empty bucket",
			bucket:     newBucket(0, nil),
			scale:      3,
			wantStart:  0,
			wantCounts: nil,
		},
		{
			name:       "1 size bucket",
			bucket:     newBucket(50, []uint64{7}),
			scale:      4,
			wantStart:  3,
			wantCounts: []uint64{7},
		},
		{
			name:       "zero scale",
			bucket:     newBucket(50, []uint64{7, 5}),
			scale:      0,
			wantStart:  50,
			wantCounts: []uint64{7, 5},
		},
		{
			name:       "aligned bucket scale 1",
			bucket:     newBucket(0, []uint64{1, 2, 3, 4, 5, 6}),
			scale:      1,
			wantStart:  0,
			wantCounts: []uint64{3, 7, 11},
		},
		{
			name:       "aligned bucket scale 2",
			bucket:     newBucket(0, []uint64{1, 2, 3, 4, 5, 6}),
			scale:      2,
			wantStart:  0,
			wantCounts: []uint64{10, 11},
		},
		{
			name:       "aligned bucket scale 3",
			bucket:     newBucket(0, []uint64{1, 2, 3, 4, 5, 6}),
			scale:      3,
			wantStart:  0,
			wantCounts: []uint64{21},
		},
		{
			name:       "unaligned bucket scale 1",
			bucket:     newBucket(5, []uint64{1, 2, 3, 4, 5, 6}), // This is equivalent to [0,0,0,0,0,1,2,3,4,5,6]
			scale:      1,
			wantStart:  2,
			wantCounts: []uint64{1, 5, 9, 6}, // This is equivalent to [0,0,1,5,9,6]
		},
		{
			name:       "unaligned bucket scale 2",
			bucket:     newBucket(7, []uint64{1, 2, 3, 4, 5, 6}), // This is equivalent to [0,0,0,0,0,0,0,1,2,3,4,5,6]
			scale:      2,
			wantStart:  1,
			wantCounts: []uint64{1, 14, 6}, // This is equivalent to [0,1,14,6]
		},
		{
			name:       "unaligned bucket scale 3",
			bucket:     newBucket(3, []uint64{1, 2, 3, 4, 5, 6}), // This is equivalent to [0,0,0,1,2,3,4,5,6]
			scale:      3,
			wantStart:  0,
			wantCounts: []uint64{15, 6}, // This is equivalent to [0,15,6]
		},
		{
			name:       "unaligned bucket scale 1",
			bucket:     newBucket(1, []uint64{1, 0, 1}),
			scale:      1,
			wantStart:  0,
			wantCounts: []uint64{1, 1},
		},
		{
			name:       "negative startBin",
			bucket:     newBucket(-1, []uint64{1, 0, 3}),
			scale:      1,
			wantStart:  -1,
			wantCounts: []uint64{1, 3},
		},
		{
			name:       "negative startBin 2",
			bucket:     newBucket(-4, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			scale:      1,
			wantStart:  -2,
			wantCounts: []uint64{3, 7, 11, 15, 19},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bucket.downscale(tt.scale)

			assertBuckets(t, tt.wantStart, tt.wantCounts, *tt.bucket, tt.name)
		})
	}
}

func TestExpoBucketRecord(t *testing.T) {
	tests := []struct {
		name       string
		bucket     *expoBuckets
		bin        int32
		wantStart  int32
		wantCounts []uint64
	}{
		{
			name:       "Empty Bucket creates first count",
			bucket:     newBucket(0, nil),
			bin:        -5,
			wantStart:  -5,
			wantCounts: []uint64{1},
		},
		{
			name:       "Bin is in the bucket",
			bucket:     newBucket(3, []uint64{1, 2, 3, 4, 5, 6}),
			bin:        5,
			wantStart:  3,
			wantCounts: []uint64{1, 2, 4, 4, 5, 6},
		},
		{
			name:       "Bin is before the start of the bucket",
			bucket:     newBucket(1, []uint64{1, 2, 3, 4, 5, 6}),
			bin:        -2,
			wantStart:  -2,
			wantCounts: []uint64{1, 0, 0, 1, 2, 3, 4, 5, 6},
		},
		{
			name:       "Bin is after the end of the bucket",
			bucket:     newBucket(-2, []uint64{1, 2, 3, 4, 5, 6}),
			bin:        4,
			wantStart:  -2,
			wantCounts: []uint64{1, 2, 3, 4, 5, 6, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bucket.record(tt.bin)

			assertBuckets(t, tt.wantStart, tt.wantCounts, *tt.bucket, tt.name)
		})
	}
}

func TestScaleChange(t *testing.T) {
	type args struct {
		bin      int32
		startBin int32
		length   int
		maxSize  int
	}
	tests := []struct {
		name string
		args args
		want int32
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
			p := newExpoHistogramDataPoint[float64](alice, tt.args.maxSize, 20, false, false)
			got := p.scaleChange(tt.args.bin, tt.args.startBin, tt.args.length)
			if got != tt.want {
				t.Errorf("scaleChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPrepend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		agg := newExpoHistogramDataPoint[float64](alice, 1024, 20, false, false)
		n := math.MaxFloat64
		for range 1024 {
			agg.record(n)
			n /= 2
		}
	}
}

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		agg := newExpoHistogramDataPoint[float64](alice, 1024, 20, false, false)
		n := smallestNonZeroNormalFloat64
		for range 1024 {
			agg.record(n)
			n *= 2
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
		attrs:   alice,
		maxSize: 4,
	}
	want.minMax.Update(math.SmallestNonzeroFloat64)
	want.sum.add(3 * math.SmallestNonzeroFloat64)
	want.scale.Store(20)
	want.posBuckets = *newBucket(-1126170625, []uint64{3})

	ehdp := newExpoHistogramDataPoint[float64](alice, 4, 20, false, false)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)

	want.startTime = ehdp.startTime

	assert.Equal(t, want, ehdp)
}

func TestExponentialHistogramAggregation(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64/Delta", testDeltaExpoHist[int64]())
	c.Reset()

	t.Run("Float64/Delta", testDeltaExpoHist[float64]())
	c.Reset()

	t.Run("Int64/Cumulative", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeExpoHist[int64]()(t)
	})
	c.Reset()

	t.Run("Int64/Cumulative/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeExpoHist[int64]()(t)
	})
	c.Reset()

	t.Run("Float64/Cumulative", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeExpoHist[float64]()(t)
	})
	c.Reset()

	t.Run("Float64/Cumulative/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeExpoHist[float64]()(t)
	})
	c.Reset()
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
				{ctx, -1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(3),
							Count:      7,
							Min:        metricdata.NewExtrema[N](-1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        30,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
							NegativeBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1},
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
				{ctx, -1, alice},
			},
			expect: output{
				n: 2,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(7),
							Count:      7,
							Min:        metricdata.NewExtrema[N](-1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        30,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
							NegativeBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1},
							},
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(7),
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

	aliceStartTime := y2kPlus(0)
	overflowStartTime := y2kPlus(0)

	if x.PerSeriesStartTimestamps.Enabled() {
		aliceStartTime = y2kPlus(2)
		overflowStartTime = y2kPlus(6)
	}

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
				{ctx, -1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.ExponentialHistogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       y2kPlus(3),
							Count:      7,
							Min:        metricdata.NewExtrema[N](-1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        30,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 4, 1},
							},
							NegativeBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1},
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
							StartTime:  aliceStartTime,
							Time:       y2kPlus(4),
							Count:      10,
							Min:        metricdata.NewExtrema[N](-1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        43,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 6, 2},
							},
							NegativeBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1},
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
							StartTime:  aliceStartTime,
							Time:       y2kPlus(5),
							Count:      10,
							Min:        metricdata.NewExtrema[N](-1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        43,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 6, 2},
							},
							NegativeBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1},
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
							StartTime:  aliceStartTime,
							Time:       y2kPlus(7),
							Count:      10,
							Min:        metricdata.NewExtrema[N](-1),
							Max:        metricdata.NewExtrema[N](16),
							Sum:        43,
							Scale:      -1,
							PositiveBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1, 6, 2},
							},
							NegativeBucket: metricdata.ExponentialBucket{
								Offset: -1,
								Counts: []uint64{1},
							},
						},
						{
							Attributes: overflowSet,
							StartTime:  overflowStartTime,
							Time:       y2kPlus(7),
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

func TestExponentialHistogramAggregationConcurrentSafe(t *testing.T) {
	t.Run("Int64/Delta", testDeltaExpoHistConcurrentSafe[int64]())
	t.Run("Float64/Delta", testDeltaExpoHistConcurrentSafe[float64]())
	t.Run("Int64/Cumulative", testCumulativeExpoHistConcurrentSafe[int64]())
	t.Run("Float64/Cumulative", testCumulativeExpoHistConcurrentSafe[float64]())
}

func testDeltaExpoHistConcurrentSafe[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.ExponentialBucketHistogram(4, 20, false, false)
	return testAggregationConcurrentSafe[N](in, out, validateExponentialHistogram[N])
}

func testCumulativeExpoHistConcurrentSafe[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.ExponentialBucketHistogram(4, 20, false, false)
	return testAggregationConcurrentSafe[N](in, out, validateExponentialHistogram[N])
}

func validateExponentialHistogram[N int64 | float64](t *testing.T, aggs []metricdata.Aggregation) {
	sums := make(map[attribute.Set]N)
	counts := make(map[attribute.Set]uint64)
	var isDelta bool
	for i, agg := range aggs {
		s, ok := agg.(metricdata.ExponentialHistogram[N])
		require.True(t, ok)
		if s.Temporality == metricdata.DeltaTemporality {
			isDelta = true
		}
		require.LessOrEqual(t, len(s.DataPoints), 3, "AggregationLimit of 3 exceeded in a single cycle")
		for _, dp := range s.DataPoints {
			assert.False(t,
				dp.Time.Before(dp.StartTime),
				"Timestamp %v must not be before start time %v", dp.Time, dp.StartTime,
			)

			if s.Temporality == metricdata.DeltaTemporality {
				sums[dp.Attributes] += dp.Sum
				counts[dp.Attributes] += dp.Count
			} else if i == len(aggs)-1 {
				sums[dp.Attributes] = dp.Sum
				counts[dp.Attributes] = dp.Count
			}

			var totalCount uint64
			for _, bc := range dp.PositiveBucket.Counts {
				totalCount += bc
			}
			for _, bc := range dp.NegativeBucket.Counts {
				totalCount += bc
			}
			assert.Equal(t, totalCount, dp.Count)
		}
	}

	var totalSum N
	var totalCount uint64
	for attr, sum := range sums {
		totalSum += sum
		count := counts[attr]
		totalCount += count

		expectedSingleSum := expectedConcurrentSum[N]() / N(concurrentNumGoroutines)
		expectedSingleCount := expectedConcurrentCount / uint64(concurrentNumGoroutines)

		if !isDelta {
			if attr == overflowSet {
				// The overflow set contains all the goroutines that didn't make the limit of 3
				assert.Equal(t, uint64(0), count%expectedSingleCount)
				assertSumEqual[N](t, N(count/expectedSingleCount)*expectedSingleSum, sum)
			} else {
				// Individual attributes should have exactly one goroutine's worth of data
				assert.Equal(t, expectedSingleSum, sum)
				assert.Equal(t, expectedSingleCount, count)
			}
		}
	}
	assertSumEqual[N](t, expectedConcurrentSum[N](), totalSum)
	assert.Equal(t, expectedConcurrentCount, totalCount)
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
	scales := []int32{0, 15, -5}

	for _, s := range scales {
		for _, v := range values {
			f.Add(v, s)
		}
	}

	f.Fuzz(func(t *testing.T, v float64, scale int32) {
		// GetBin only works on positive values.
		if math.Signbit(v) {
			v *= -1
		}
		// GetBin Doesn't work on zero.
		if v == 0.0 {
			t.Skip("skipping test for zero")
		}

		p := newExpoHistogramDataPoint[float64](alice, 4, 20, false, false)
		// scale range is -10 to 20.
		scaleValue := (scale%31+31)%31 - 10
		p.scale.Store(scaleValue)
		got := p.getBin(v)
		if v <= lowerBound(got, scaleValue) {
			t.Errorf(
				"v=%x scale =%d had bin %d, but was below lower bound %x",
				v,
				scaleValue,
				got,
				lowerBound(got, scaleValue),
			)
		}
		if v > lowerBound(got+1, scaleValue) {
			t.Errorf(
				"v=%x scale =%d had bin %d, but was above upper bound %x",
				v,
				scaleValue,
				got,
				lowerBound(got+1, scaleValue),
			)
		}
	})
}

func lowerBound(index, scale int32) float64 {
	// The lowerBound of the index of Math.SmallestNonzeroFloat64 at any scale
	// is always rounded down to 0.0.
	// For example lowerBound(getBin(Math.SmallestNonzeroFloat64, 7), 7) == 0.0
	// 2 ^ (index * 2 ^ (-scale))
	return math.Exp2(math.Ldexp(float64(index), -int(scale)))
}

func TestExponentialHistogramConcurrentSafeEdgeCases(t *testing.T) {
	t.Run("Int64/Delta", testExpoHistConcurrentSafeEdgeCases[int64](metricdata.DeltaTemporality))
	t.Run("Float64/Delta", testExpoHistConcurrentSafeEdgeCases[float64](metricdata.DeltaTemporality))
	t.Run("Int64/Cumulative", testExpoHistConcurrentSafeEdgeCases[int64](metricdata.CumulativeTemporality))
	t.Run("Float64/Cumulative", testExpoHistConcurrentSafeEdgeCases[float64](metricdata.CumulativeTemporality))
}

func testExpoHistConcurrentSafeEdgeCases[N int64 | float64](temporality metricdata.Temporality) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("ZeroValues", func(t *testing.T) {
			meas, comp := Builder[N]{
				Temporality:      temporality,
				Filter:           attrFltr,
				AggregationLimit: 3,
			}.ExponentialBucketHistogram(160, 20, false, false)

			ctx := t.Context()
			var wg sync.WaitGroup
			const numGoroutines = 10
			const numRecords = 100
			wg.Add(numGoroutines)
			for range numGoroutines {
				go func() {
					defer wg.Done()
					for range numRecords {
						meas(ctx, 0, alice)
					}
				}()
			}
			wg.Wait()

			dest := new(metricdata.Aggregation)
			comp(dest)
			h := (*dest).(metricdata.ExponentialHistogram[N])
			require.Len(t, h.DataPoints, 1)
			assert.Equal(t, uint64(numGoroutines*numRecords), h.DataPoints[0].ZeroCount)
			assert.Equal(t, uint64(numGoroutines*numRecords), h.DataPoints[0].Count)
		})

		t.Run("RescalingStress", func(t *testing.T) {
			meas, comp := Builder[N]{
				Temporality:      temporality,
				Filter:           attrFltr,
				AggregationLimit: 3,
			}.ExponentialBucketHistogram(160, 20, false, false)

			ctx := t.Context()
			var wg sync.WaitGroup
			const numGoroutines = 10
			const numRecords = 100

			// To verify exact outcome, we sequentially record the same to a reference.
			refMeas, refComp := Builder[N]{
				Temporality:      temporality,
				Filter:           attrFltr,
				AggregationLimit: 3,
			}.ExponentialBucketHistogram(160, 20, false, false)
			var m sync.Mutex

			wg.Add(numGoroutines)
			for i := range numGoroutines {
				go func(id int) {
					defer wg.Done()
					for j := range numRecords {
						// generate a mix of very large and very small powers of 2
						valFloat := math.Exp2((float64(j) / float64(numRecords)) * 60.0)
						if id%2 == 0 {
							valFloat = -valFloat
						}
						val := N(valFloat)
						// For integers, values less than 1 will truncate to 0. Mix things up.
						if id%3 == 0 {
							val = N(float64(id+1) * 100.0)
						}

						meas(ctx, val, alice)

						m.Lock()
						refMeas(ctx, val, alice)
						m.Unlock()
					}
				}(i)
			}
			wg.Wait()

			dest := new(metricdata.Aggregation)
			comp(dest)
			h := (*dest).(metricdata.ExponentialHistogram[N])
			require.Len(t, h.DataPoints, 1)

			refDest := new(metricdata.Aggregation)
			refComp(refDest)
			refH := (*refDest).(metricdata.ExponentialHistogram[N])
			require.Len(t, refH.DataPoints, 1)

			// StartTime/Time will differ slightly.
			h.DataPoints[0].StartTime = refH.DataPoints[0].StartTime
			h.DataPoints[0].Time = refH.DataPoints[0].Time

			// Float sums might be slightly slightly off due to summing wildly different magnitudes
			// in different orders concurrently versus sequentially.
			actualSum := float64(h.DataPoints[0].Sum)
			expectedSum := float64(refH.DataPoints[0].Sum)
			if actualSum != expectedSum {
				assert.InEpsilon(t, expectedSum, actualSum, 0.5, "Sum")
				// Force equality for the deep struct comparison below
				h.DataPoints[0].Sum = refH.DataPoints[0].Sum
			}

			// Normalize Exemplars to avoid nil vs empty slice comparison failures
			h.DataPoints[0].Exemplars = nil
			refH.DataPoints[0].Exemplars = nil

			assert.Equal(t, refH, h)
		})
	}
}

func TestExpoHistogramRecordUnderflow(t *testing.T) {
	var errs []error
	original := global.GetErrorHandler()
	global.SetErrorHandler(otel.ErrorHandlerFunc(func(e error) {
		errs = append(errs, e)
	}))
	t.Cleanup(func() {
		global.SetErrorHandler(original)
	})

	dp := newExpoHistogramDataPoint[float64](attribute.NewSet(), 1, 20, false, false)
	// Force scale to a low value
	dp.scale.Store(-10)
	dp.record(1)
	dp.record(math.MaxFloat64)
	require.Len(t, errs, 1)
	assert.EqualError(t, errs[0], "exponential histogram scale underflow")
}

func TestDeltaExpoHistogramMeasureNaNAndInf(t *testing.T) {
	h := newExponentialHistogram[float64](4, 20, false, false, 0, dropExemplars[float64])
	ctx := t.Context()

	h.measure(ctx, math.NaN(), attribute.NewSet(), nil)
	h.measure(ctx, math.Inf(1), attribute.NewSet(), nil)
	h.measure(ctx, math.Inf(-1), attribute.NewSet(), nil)

	var dest metricdata.Aggregation
	h.delta(&dest)
	eh := dest.(metricdata.ExponentialHistogram[float64])
	assert.Empty(t, eh.DataPoints)
}
