package internal

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

func testExpoHistogramDataPointRecord[N int64 | float64](t *testing.T) {
	type TestCase struct {
		maxSize         int
		values          []N
		expectedBuckets expoBucket
		expectedScale   int
	}

	testCases := []TestCase{
		{
			maxSize: 4,
			values:  []N{2, 4, 1},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 1, 1},
			},
			expectedScale: 0,
		},
		{
			maxSize: 4,
			values:  []N{4, 4, 4, 2, 16, 1},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{1, 4, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{1, 2, 4},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{1, 4, 2},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{2, 4, 1},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{2, 1, 4},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{4, 1, 2},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []N{4, 2, 1},
			expectedBuckets: expoBucket{
				startIndex: -1,

				counts: []uint64{1, 2},
			},
			expectedScale: -1,
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			defer withHandler(t)()

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

func testExpoHistogramDataPointRecordMinMaxSum[N int64 | float64](t *testing.T) {
	type Expected struct {
		min, max, sum N
		count         uint
	}
	type TestCase struct {
		values   []N
		expected Expected
	}
	testCases := []TestCase{
		{
			values:   []N{2, 4, 1},
			expected: Expected{1, 4, 7, 3},
		},
		{
			values:   []N{4, 4, 4, 2, 16, 1},
			expected: Expected{1, 16, 31, 6},
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			defer withHandler(t)()

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
				startIndex: -1,
				counts:     []uint64{2, 3, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{1, 0.5, 2},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{1, 2, 0.5},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{2, 0.5, 1},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{2, 1, 0.5},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{0.5, 1, 2},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{2, 1},
			},
			expectedScale: -1,
		},
		{
			maxSize: 2,
			values:  []float64{0.5, 2, 1},
			expectedBuckets: expoBucket{
				startIndex: -1,
				counts:     []uint64{2, 1},
			},
			expectedScale: -1,
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprint(tt.values), func(t *testing.T) {
			defer withHandler(t)()

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
	// These indexes are calculated from the following formula:
	// floor( log2( value) * 2^20 )

	fdp := newExpoHistogramDataPoint[float64](4, 20, 0.0)
	fdp.record(math.MaxFloat64)

	if fdp.posBuckets.startIndex != 1073741824 {
		t.Errorf("Expected startIndex to be 1073741824, got %d", fdp.posBuckets.startIndex)
	}

	// Subnormal Numbers do not work with the above formula.
	SmallestNonZeroNormalFloat64 := 0x1p-1022
	fdp = newExpoHistogramDataPoint[float64](4, 20, 0.0)
	fdp.record(SmallestNonZeroNormalFloat64)

	if fdp.posBuckets.startIndex != -1071644673 {
		t.Errorf("Expected startIndex to be -1071644673, got %d", fdp.posBuckets.startIndex)
	}

	idp := newExpoHistogramDataPoint[int64](4, 20, 0.0)
	idp.record(math.MaxInt64)

	if idp.posBuckets.startIndex != 66060287 {
		t.Errorf("Expected startIndex to be 66060287, got %d", idp.posBuckets.startIndex)
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
				startIndex: 50,
				counts:     []uint64{7},
			},
			scale: 4,
			want: &expoBucket{
				startIndex: 3,
				counts:     []uint64{7},
			},
		},
		{
			name: "zero scale",
			bucket: &expoBucket{
				startIndex: 50,
				counts:     []uint64{7, 5},
			},
			scale: 0,
			want: &expoBucket{
				startIndex: 50,
				counts:     []uint64{7, 5},
			},
		},
		{
			name: "aligned bucket scale 1",
			bucket: &expoBucket{
				startIndex: 0,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 1,
			want: &expoBucket{
				startIndex: 0,
				counts:     []uint64{3, 7, 11},
			},
		},
		{
			name: "aligned bucket scale 2",
			bucket: &expoBucket{
				startIndex: 0,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 2,
			want: &expoBucket{
				startIndex: 0,
				counts:     []uint64{10, 11},
			},
		},
		{
			name: "aligned bucket scale 3",
			bucket: &expoBucket{
				startIndex: 0,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			},
			scale: 3,
			want: &expoBucket{
				startIndex: 0,
				counts:     []uint64{21},
			},
		},
		{
			name: "unaligned bucket scale 1",
			bucket: &expoBucket{
				startIndex: 5,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,1,2,3,4,5,6]
			scale: 1,
			want: &expoBucket{
				startIndex: 2,
				counts:     []uint64{1, 5, 9, 6},
			}, // This is equivalent to [0,0,1,5,9,6]
		},
		{
			name: "unaligned bucket scale 2",
			bucket: &expoBucket{
				startIndex: 7,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,0,0,1,2,3,4,5,6]
			scale: 2,
			want: &expoBucket{
				startIndex: 1,
				counts:     []uint64{1, 14, 6},
			}, // This is equivalent to [0,1,14,6]
		},
		{
			name: "unaligned bucket scale 3",
			bucket: &expoBucket{
				startIndex: 3,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,1,2,3,4,5,6]
			scale: 3,
			want: &expoBucket{
				startIndex: 0,
				counts:     []uint64{15, 6},
			}, // This is equivalent to [0,15,6]
		},
		{
			name: "unaligned bucket scale 1",
			bucket: &expoBucket{
				startIndex: 1,
				counts:     []uint64{1, 0, 1},
			},
			scale: 1,
			want: &expoBucket{
				startIndex: 0,
				counts:     []uint64{1, 1},
			},
		},
		{
			name: "negative startIndex",
			bucket: &expoBucket{
				startIndex: -1,
				counts:     []uint64{1, 0, 3},
			},
			scale: 1,
			want: &expoBucket{
				startIndex: -1,
				counts:     []uint64{1, 3},
			},
		},
		{
			name: "negative startIndex 2",
			bucket: &expoBucket{
				startIndex: -4,
				counts:     []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			scale: 1,
			want: &expoBucket{
				startIndex: -2,
				counts:     []uint64{3, 7, 11, 15, 19},
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
		index  int
		want   *expoBucket
	}{
		{
			name:   "Empty Bucket creates first count",
			bucket: &expoBucket{},
			index:  -5,
			want: &expoBucket{
				startIndex: -5,
				counts:     []uint64{1},
			},
		},
		{
			name: "Index is in the bucket",
			bucket: &expoBucket{
				startIndex: 3,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			},
			index: 5,
			want: &expoBucket{
				startIndex: 3,
				counts:     []uint64{1, 2, 4, 4, 5, 6},
			},
		},
		{
			name: "Index is before the start of the bucket",
			bucket: &expoBucket{
				startIndex: 1,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			},
			index: -2,
			want: &expoBucket{
				startIndex: -2,
				counts:     []uint64{1, 0, 0, 1, 2, 3, 4, 5, 6},
			},
		},
		{
			name: "Index is after the end of the bucket",
			bucket: &expoBucket{
				startIndex: -2,
				counts:     []uint64{1, 2, 3, 4, 5, 6},
			},
			index: 4,
			want: &expoBucket{
				startIndex: -2,
				counts:     []uint64{1, 2, 3, 4, 5, 6, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bucket.record(tt.index)

			assert.Equal(t, tt.want, tt.bucket)
		})
	}
}

func Test_needRescale(t *testing.T) {
	type args struct {
		index      int
		startIndex int
		length     int
		maxSize    int
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
				index:      5,
				startIndex: 0,
				length:     0,
				maxSize:    4,
			},
			want: false,
		},
		{
			name: "if index is between start, and the end, no rescale needed",
			// [-1, ..., 8] Length 10 -> [-1, ..., 5, ..., 8] Length 10
			args: args{
				index:      5,
				startIndex: -1,
				length:     10,
				maxSize:    20,
			},
			want: false,
		},
		{
			name: "if len([index,... end]) > maxSize, rescale needed",
			// [8,9,10] Length 3 -> [5, ..., 10] Length 6
			args: args{
				index:      5,
				startIndex: 8,
				length:     3,
				maxSize:    5,
			},
			want: true,
		},
		{
			name: "if len([start, ..., index]) > maxSize, rescale needed",
			// [2,3,4] Length 3 -> [2, ..., 7] Length 6
			args: args{
				index:      7,
				startIndex: 2,
				length:     3,
				maxSize:    5,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needRescale(tt.args.index, tt.args.startIndex, tt.args.length, tt.args.maxSize); got != tt.want {
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
	factory := func() Aggregator[N] { return NewDeltaExponentialHistogram[N](expoHistConf) }
	b.Run("Delta", benchmarkAggregator(factory))
	factory = func() Aggregator[N] { return NewCumulativeExponentialHistogram[N](expoHistConf) }
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
			startIndex: -1071644673, // This is the offset for the smallest normal float64 floor(Log_2(2.2250738585072014e-308)*2^20)
			counts:     []uint64{3},
		},
	}

	ehdp := newExpoHistogramDataPoint[float64](4, 20, 0.0)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)
	ehdp.record(math.SmallestNonzeroFloat64)

	assert.Equal(t, want, ehdp)
}

func TestZeroThresholdInt64(t *testing.T) {
	defer withHandler(t)()

	ehdp := newExpoHistogramDataPoint[int64](4, 20, 3.0)
	ehdp.record(1)
	ehdp.record(2)
	ehdp.record(3)

	assert.Len(t, ehdp.posBuckets.counts, 0)
	assert.Len(t, ehdp.negBuckets.counts, 0)
	assert.Equal(t, uint64(3), ehdp.zeroCount)
}

func TestZeroThresholdFloat64(t *testing.T) {
	defer withHandler(t)()

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

func testExponentialHistogramAggregation[N int64 | float64](t *testing.T) {
	cfg := aggregation.ExponentialHistogram{
		MaxSize:  4,
		MaxScale: 20,
	}

	type TestCase struct {
		name       string
		aggregator Aggregator[N]
		input      [][]N
		want       metricdata.ExponentialHistogram[N]
	}

	tests := []TestCase{
		{
			name:       "Delta Single",
			aggregator: NewDeltaExponentialHistogram[N](cfg),
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
			aggregator: NewCumulativeExponentialHistogram[N](cfg),
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
			aggregator: NewDeltaExponentialHistogram[N](cfg),
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
			aggregator: NewCumulativeExponentialHistogram[N](cfg),
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
			defer withHandler(t)()

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

	c := NewCumulativeExponentialHistogram[N](cfg)
	assert.Equal(t, want, c.Aggregation())

	d := NewDeltaExponentialHistogram[N](cfg)
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
			defer withHandler(t)()

			got := normalizeConfig(tt.cfg)
			assert.Equal(t, tt.want, got)
		})
	}
}
