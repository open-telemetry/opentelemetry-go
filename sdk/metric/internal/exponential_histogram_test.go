package internal

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

type noErrorHandler struct{ t *testing.T }

func (h *noErrorHandler) Handle(e error) {
	require.NoError(h.t, e)
}

func withHandler(t *testing.T) func() {
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

			dp := NewExpoHistogramDataPoint[N](tt.maxSize, 20, 0.0)
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

			dp := NewExpoHistogramDataPoint[N](4, 20, 0.0)
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

			dp := NewExpoHistogramDataPoint[float64](tt.maxSize, 20, 0.0)
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

	fdp := NewExpoHistogramDataPoint[float64](4, 20, 0.0)
	fdp.record(math.MaxFloat64)

	if fdp.posBuckets.startIndex != 1073741824 {
		t.Errorf("Expected startIndex to be 1073741824, got %d", fdp.posBuckets.startIndex)
	}

	// Subnormal Numbers do not work with the above formula.
	SmallestNonZeroNormalFloat64 := 0x1p-1022
	fdp = NewExpoHistogramDataPoint[float64](4, 20, 0.0)
	fdp.record(SmallestNonZeroNormalFloat64)

	if fdp.posBuckets.startIndex != -1071644673 {
		t.Errorf("Expected startIndex to be -1071644673, got %d", fdp.posBuckets.startIndex)
	}

	idp := NewExpoHistogramDataPoint[int64](4, 20, 0.0)
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

func BenchmarkLinear(b *testing.B) {
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := NewExpoHistogramDataPoint[float64](1024, 20, 0.0)
	for i := 0; i < b.N; i++ {
		x := 2 - rnd.Float64()
		agg.record(x)
	}
}

// Benchmarks the Update() function for values in the range (0, MaxValue].
func BenchmarkExponential(b *testing.B) {
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := NewExpoHistogramDataPoint[float64](1024, 20, 0.0)
	for i := 0; i < b.N; i++ {
		x := rnd.ExpFloat64()
		agg.record(x)
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
