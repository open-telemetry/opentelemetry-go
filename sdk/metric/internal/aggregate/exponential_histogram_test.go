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
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
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
	t.Run("float64 MinMaxSum", testExpoHistogramDataPointRecordMinMaxSum[float64])
	t.Run("float64-2", testExpoHistogramDataPointRecordFloat64)
	t.Run("int64", testExpoHistogramDataPointRecord[int64])
	t.Run("int64 MinMaxSum", testExpoHistogramDataPointRecordMinMaxSum[int64])
}

// TODO: This can be defined in the test after we drop support for go1.19.
type expoHistogramDataPointRecordTestCase[N int64 | float64] struct {
	maxSize         int
	values          []N
	expectedBuckets expoBucket
	expectedScale   int
}

func testExpoHistogramDataPointRecord[N int64 | float64](t *testing.T) {
	testCases := []expoHistogramDataPointRecordTestCase[N]{
		{
			maxSize: 4,
			values:  []N{2, 4, 1},
			expectedBuckets: expoBucket{
				startBin: -1,

				counts: []uint64{1, 1, 1},
			},
			expectedScale: 0,
		},
		{
			maxSize: 4,
			values:  []N{4, 4, 4, 2, 16, 1},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{1, 4, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{1, 2, 4},
			expectedBuckets: expoBucket{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{1, 4, 2},
			expectedBuckets: expoBucket{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{2, 4, 1},
			expectedBuckets: expoBucket{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{2, 1, 4},
			expectedBuckets: expoBucket{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{4, 1, 2},
			expectedBuckets: expoBucket{
				startBin: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{4, 2, 1},
			expectedBuckets: expoBucket{
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

			dp := newExpoHistogramDataPoint[N](tt.maxSize, 20, 0.0)
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

func testExpoHistogramDataPointRecordMinMaxSum[N int64 | float64](t *testing.T) {
	testCases := []expoHistogramDataPointRecordMinMaxSumTestCase[N]{
		{
			values:   []N{2, 4, 1},
			expected: expectedMinMaxSum[N]{1, 4, 7, 3},
		},
		{
			values:   []N{4, 4, 4, 2, 16, 1},
			expected: expectedMinMaxSum[N]{1, 16, 31, 6},
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			dp := newExpoHistogramDataPoint[N](4, 20, 0.0)
			for _, v := range tt.values {
				dp.record(v)
			}

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
		expectedBuckets expoBucket
		expectedScale   int
	}

	testCases := []TestCase{
		{
			maxSize: 4,
			values:  []float64{2, 2, 2, 1, 8, 0.5},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{2, 3, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{1, 0.5, 2},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{1, 2, 0.5},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{2, 0.5, 1},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{2, 1, 0.5},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{0.5, 1, 2},
			expectedBuckets: expoBucket{
				startBin: -1,
				counts:   []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{0.5, 2, 1},
			expectedBuckets: expoBucket{
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

			dp := newExpoHistogramDataPoint[float64](tt.maxSize, 20, 0.0)
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
	// floor( log2( value) * 2^20 )

	fdp := newExpoHistogramDataPoint[float64](4, 20, 0.0)
	fdp.record(math.MaxFloat64)

	if fdp.posBuckets.startBin != 1073741824 {
		t.Errorf("Expected startBin to be 1073741824, got %d", fdp.posBuckets.startBin)
	}

	// Subnormal Numbers do not work with the above formula.
	SmallestNonZeroNormalFloat64 := 0x1p-1022
	fdp = newExpoHistogramDataPoint[float64](4, 20, 0.0)
	fdp.record(SmallestNonZeroNormalFloat64)

	if fdp.posBuckets.startBin != -1071644673 {
		t.Errorf("Expected startBin to be -1071644673, got %d", fdp.posBuckets.startBin)
	}

	idp := newExpoHistogramDataPoint[int64](4, 20, 0.0)
	idp.record(math.MaxInt64)

	if idp.posBuckets.startBin != 66060287 {
		t.Errorf("Expected startBin to be 66060287, got %d", idp.posBuckets.startBin)
	}
}

func TestExpoBucketDownscale(t *testing.T) {
	tests := []struct {
		name   string
		bucket *expoBucket
		scale  int
		want   *expoBucket
	}{
		{
			name:   "Empty bucket",
			bucket: &expoBucket{},
			scale:  3,
			want:   &expoBucket{},
		},
		{
			name: "1 size bucket",
			bucket: &expoBucket{
				startBin: 50,
				counts:   []uint64{7},
			},
			scale: 4,
			want: &expoBucket{
				startBin: 3,
				counts:   []uint64{7},
			},
		},
		{
			name: "zero scale",
			bucket: &expoBucket{
				startBin: 50,
				counts:   []uint64{7, 5},
			},
			scale: 0,
			want: &expoBucket{
				startBin: 50,
				counts:   []uint64{7, 5},
			},
		},
		{
			name: "aligned bucket scale 1",
			bucket: &expoBucket{
				startBin: 0,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 1,
			want: &expoBucket{
				startBin: 0,
				counts:   []uint64{3, 7, 11},
			},
		},
		{
			name: "aligned bucket scale 2",
			bucket: &expoBucket{
				startBin: 0,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 2,
			want: &expoBucket{
				startBin: 0,
				counts:   []uint64{10, 11},
			},
		},
		{
			name: "aligned bucket scale 3",
			bucket: &expoBucket{
				startBin: 0,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 3,
			want: &expoBucket{
				startBin: 0,
				counts:   []uint64{21},
			},
		},
		{
			name: "unaligned bucket scale 1",
			bucket: &expoBucket{
				startBin: 5,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,1,2,3,4,5,6]
			scale: 1,
			want: &expoBucket{
				startBin: 2,
				counts:   []uint64{1, 5, 9, 6},
			}, // This is equivalent to [0,0,1,5,9,6]
		},
		{
			name: "unaligned bucket scale 2",
			bucket: &expoBucket{
				startBin: 7,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,0,0,1,2,3,4,5,6]
			scale: 2,
			want: &expoBucket{
				startBin: 1,
				counts:   []uint64{1, 14, 6},
			}, // This is equivalent to [0,1,14,6]
		},
		{
			name: "unaligned bucket scale 3",
			bucket: &expoBucket{
				startBin: 3,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,1,2,3,4,5,6]
			scale: 3,
			want: &expoBucket{
				startBin: 0,
				counts:   []uint64{15, 6},
			}, // This is equivalent to [0,15,6]
		},
		{
			name: "unaligned bucket scale 1",
			bucket: &expoBucket{
				startBin: 1,
				counts:   []uint64{1, 0, 1},
			},
			scale: 1,
			want: &expoBucket{
				startBin: 0,
				counts:   []uint64{1, 1},
			},
		},
		{
			name: "negative startBin",
			bucket: &expoBucket{
				startBin: -1,
				counts:   []uint64{1, 0, 3},
			},
			scale: 1,
			want: &expoBucket{
				startBin: -1,
				counts:   []uint64{1, 3},
			},
		},
		{
			name: "negative startBin 2",
			bucket: &expoBucket{
				startBin: -4,
				counts:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			scale: 1,
			want: &expoBucket{
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
		bucket *expoBucket
		bin    int
		want   *expoBucket
	}{
		{
			name:   "Empty Bucket creates first count",
			bucket: &expoBucket{},
			bin:    -5,
			want: &expoBucket{
				startBin: -5,
				counts:   []uint64{1},
			},
		},
		{
			name: "Bin is in the bucket",
			bucket: &expoBucket{
				startBin: 3,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			bin: 5,
			want: &expoBucket{
				startBin: 3,
				counts:   []uint64{1, 2, 4, 4, 5, 6},
			},
		},
		{
			name: "Bin is before the start of the bucket",
			bucket: &expoBucket{
				startBin: 1,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			bin: -2,
			want: &expoBucket{
				startBin: -2,
				counts:   []uint64{1, 0, 0, 1, 2, 3, 4, 5, 6},
			},
		},
		{
			name: "Bin is after the end of the bucket",
			bucket: &expoBucket{
				startBin: -2,
				counts:   []uint64{1, 2, 3, 4, 5, 6},
			},
			bin: 4,
			want: &expoBucket{
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

func Test_needRescale(t *testing.T) {
	type args struct {
		bin      int
		startBin int
		length   int
		maxSize  int
	}
	tests := []struct {
		name string
		args args
		want bool
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
			want: false,
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
			want: false,
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
			want: true,
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
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needRescale(tt.args.bin, tt.args.startBin, tt.args.length, tt.args.maxSize); got != tt.want {
				t.Errorf("needRescale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPrepend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		agg := newExpoHistogramDataPoint[float64](1024, 20, 0.0)
		n := math.MaxFloat64
		for j := 0; j < 1024; j++ {
			agg.record(n)
			n = n / 2
		}
	}
}

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		agg := newExpoHistogramDataPoint[float64](1024, 20, 0.0)
		n := smallestNonZeroNormalFloat64
		for j := 0; j < 1024; j++ {
			agg.record(n)
			n = n * 2
		}
	}
}

var expoHistConf = aggregation.DefaultExponentialHistogram()

func BenchmarkExponentialHistogram(b *testing.B) {
	b.Run("Int64", benchmarkExponentialHistogram[int64])
	b.Run("Float64", benchmarkExponentialHistogram[float64])
}

func benchmarkExponentialHistogram[N int64 | float64](b *testing.B) {
	factory := func() aggregator[N] { return newDeltaExponentialHistogram[N](expoHistConf) }
	b.Run("Delta", benchmarkAggregator(factory))
	factory = func() aggregator[N] { return newCumulativeExponentialHistogram[N](expoHistConf) }
	b.Run("Cumulative", benchmarkAggregator(factory))
}

func TestSubNormal(t *testing.T) {
	want := &expoHistogramDataPoint[float64]{
		maxSize: 4,
		count:   3,
		min:     math.SmallestNonzeroFloat64,
		max:     math.SmallestNonzeroFloat64,
		sum:     3 * math.SmallestNonzeroFloat64,

		scale: 20,
		posBuckets: expoBucket{
			startBin: -1071644673, // This is the offset for the smallest normal float64 floor(Log_2(2.2250738585072014e-308)*2^20)
			counts:   []uint64{3},
		},
	}

	ehdp := newExpoHistogramDataPoint[float64](4, 20, 0.0)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)

	assert.Equal(t, want, ehdp)
}

func TestZeroThresholdInt64(t *testing.T) {
	restore := withHandler(t)
	defer restore()

	ehdp := newExpoHistogramDataPoint[int64](4, 20, 3.0)
	ehdp.record(1)
	ehdp.record(2)
	ehdp.record(3)

	assert.Len(t, ehdp.posBuckets.counts, 0)
	assert.Len(t, ehdp.negBuckets.counts, 0)
	assert.Equal(t, uint64(3), ehdp.zeroCount)
}

func TestZeroThresholdFloat64(t *testing.T) {
	restore := withHandler(t)
	defer restore()

	ehdp := newExpoHistogramDataPoint[float64](4, 20, 0.3)
	ehdp.record(0.1)
	ehdp.record(0.2)
	ehdp.record(0.3)

	assert.Len(t, ehdp.posBuckets.counts, 0)
	assert.Len(t, ehdp.negBuckets.counts, 0)
	assert.Equal(t, uint64(3), ehdp.zeroCount)
}

func TestExponentialHistogramAggregation(t *testing.T) {
	t.Run("Int64", testExponentialHistogramAggregation[int64])
	t.Run("Float64", testExponentialHistogramAggregation[float64])
	t.Run("Int64 Empty", testEmptyExponentialHistogramAggregation[int64])
	t.Run("Float64 Empty", testEmptyExponentialHistogramAggregation[float64])
}

// TODO: This can be defined in the test after we drop support for go1.19.
type exponentialHistogramAggregationTestCase[N int64 | float64] struct {
	name       string
	aggregator aggregator[N]
	input      [][]N
	want       metricdata.ExponentialHistogram[N]
}

func testExponentialHistogramAggregation[N int64 | float64](t *testing.T) {
	cfg := aggregation.ExponentialHistogram{
		MaxSize:  4,
		MaxScale: 20,
	}

	tests := []exponentialHistogramAggregationTestCase[N]{
		{
			name:       "Delta Single",
			aggregator: newDeltaExponentialHistogram[N](cfg),
			input: [][]N{
				{4, 4, 4, 2, 16, 1},
			},
			want: metricdata.ExponentialHistogram[N]{
				Temporality: metricdata.DeltaTemporality,
				DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
					{
						Count: 6,
						Min:   metricdata.NewExtrema[N](1),
						Max:   metricdata.NewExtrema[N](16),
						Sum:   31,
						Scale: -1,
						PositiveBucket: metricdata.ExponentialBucket{
							Offset: -1,
							Counts: []uint64{1, 4, 1},
						},
					},
				},
			},
		},
		{
			name:       "Cumulative Single",
			aggregator: newCumulativeExponentialHistogram[N](cfg),
			input: [][]N{
				{4, 4, 4, 2, 16, 1},
			},
			want: metricdata.ExponentialHistogram[N]{
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
					{
						Count: 6,
						Min:   metricdata.NewExtrema[N](1),
						Max:   metricdata.NewExtrema[N](16),
						Sum:   31,
						Scale: -1,
						PositiveBucket: metricdata.ExponentialBucket{
							Offset: -1,
							Counts: []uint64{1, 4, 1},
						},
					},
				},
			},
		},
		{
			name:       "Delta Multiple",
			aggregator: newDeltaExponentialHistogram[N](cfg),
			input: [][]N{
				{2, 3, 8},
				{4, 4, 4, 2, 16, 1},
			},
			want: metricdata.ExponentialHistogram[N]{
				Temporality: metricdata.DeltaTemporality,
				DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
					{
						Count: 6,
						Min:   metricdata.NewExtrema[N](1),
						Max:   metricdata.NewExtrema[N](16),
						Sum:   31,
						Scale: -1,
						PositiveBucket: metricdata.ExponentialBucket{
							Offset: -1,
							Counts: []uint64{1, 4, 1},
						},
					},
				},
			},
		},
		{
			name:       "Cumulative Multiple ",
			aggregator: newCumulativeExponentialHistogram[N](cfg),
			input: [][]N{
				{2, 3, 8},
				{4, 4, 4, 2, 16, 1},
			},
			want: metricdata.ExponentialHistogram[N]{
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.ExponentialHistogramDataPoint[N]{
					{
						Count: 9,
						Min:   metricdata.NewExtrema[N](1),
						Max:   metricdata.NewExtrema[N](16),
						Sum:   44,
						Scale: -1,
						PositiveBucket: metricdata.ExponentialBucket{
							Offset: -1,
							Counts: []uint64{1, 6, 2},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			var got metricdata.Aggregation
			for _, n := range tt.input {
				for _, v := range n {
					tt.aggregator.Aggregate(v, *attribute.EmptySet())
				}
				got = tt.aggregator.Aggregation()
			}

			metricdatatest.AssertAggregationsEqual(t, tt.want, got, metricdatatest.IgnoreTimestamp())
		})
	}
}

func testEmptyExponentialHistogramAggregation[N int64 | float64](t *testing.T) {
	cfg := aggregation.ExponentialHistogram{
		MaxSize:  4,
		MaxScale: 20,
	}
	var want metricdata.Aggregation

	c := newCumulativeExponentialHistogram[N](cfg)
	assert.Equal(t, want, c.Aggregation())

	d := newDeltaExponentialHistogram[N](cfg)
	assert.Equal(t, want, d.Aggregation())
}

func TestNormalizeConfig(t *testing.T) {
	type TestCase struct {
		name string
		cfg  aggregation.ExponentialHistogram
		want aggregation.ExponentialHistogram
	}
	testCases := []TestCase{
		{
			name: "Normal",
			cfg: aggregation.ExponentialHistogram{
				MaxSize:       13,
				MaxScale:      7,
				ZeroThreshold: 0.5,
			},
			want: aggregation.ExponentialHistogram{
				MaxSize:       13,
				MaxScale:      7,
				ZeroThreshold: 0.5,
			},
		},
		{
			name: "MaxScale too small",
			cfg: aggregation.ExponentialHistogram{
				MaxSize:       24,
				MaxScale:      -15,
				ZeroThreshold: 0.46,
			},
			want: aggregation.ExponentialHistogram{
				MaxSize:       24,
				MaxScale:      -10,
				ZeroThreshold: 0.46,
			},
		},
		{
			name: "MaxScale too large",
			cfg: aggregation.ExponentialHistogram{
				MaxSize:       31,
				MaxScale:      25,
				ZeroThreshold: 0.8,
			},
			want: aggregation.ExponentialHistogram{
				MaxSize:       31,
				MaxScale:      20,
				ZeroThreshold: 0.8,
			},
		},
		{
			name: "MaxSize too small",
			cfg: aggregation.ExponentialHistogram{
				MaxSize:       -1,
				MaxScale:      10,
				ZeroThreshold: 0.5,
			},
			want: aggregation.ExponentialHistogram{
				MaxSize:       160,
				MaxScale:      10,
				ZeroThreshold: 0.5,
			},
		},
		{
			name: "ZeroThreshold negative",
			cfg: aggregation.ExponentialHistogram{
				MaxSize:       13,
				MaxScale:      7,
				ZeroThreshold: -0.5,
			},
			want: aggregation.ExponentialHistogram{
				MaxSize:       13,
				MaxScale:      7,
				ZeroThreshold: 0.5,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			restore := withHandler(t)
			defer restore()

			got := normalizeConfig(tt.cfg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getNormalBase2(t *testing.T) {
	tests := []struct {
		value float64
		want  int
	}{
		// TODO: Add test cases.
		{0x1p0, 0}, // 1
		{0x1.0000000000001p0, 0},
		{0x1.fffffffffffffp0, 0},
		{0x1p1, 1}, // 2
		{0x1.0000000000001p1, 1},
		{0x1.fffffffffffffp1, 1},
		{0x1p2, 2}, // 4
		{math.MaxFloat64, 1023},
		{smallestNonZeroNormalFloat64, -1022},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%x", tt.value), func(t *testing.T) {
			if got := getNormalBase2(tt.value); got != tt.want {
				t.Errorf("getNormalBase2() = %v, want %v", got, tt.want)
			}
		})
	}
}
