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

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

var (
	bounds   = []float64{1, 5}
	histConf = aggregation.ExplicitBucketHistogram{
		Boundaries: bounds,
		NoMinMax:   false,
	}
)

func TestHistogram(t *testing.T) {
	t.Cleanup(mockTime(now))
	t.Run("Int64", testHistogram[int64])
	t.Run("Float64", testHistogram[float64])
}

func testHistogram[N int64 | float64](t *testing.T) {
	tester := &aggregatorTester[N]{
		GoroutineN:   defaultGoroutines,
		MeasurementN: defaultMeasurements,
		CycleN:       defaultCycles,
	}

	incr := monoIncr[N]()
	eFunc := deltaHistExpecter[N](incr)
	t.Run("Delta", tester.Run(NewDeltaHistogram[N](histConf), incr, eFunc))
	eFunc = cumuHistExpecter[N](incr)
	t.Run("Cumulative", tester.Run(NewCumulativeHistogram[N](histConf), incr, eFunc))
}

func deltaHistExpecter[N int64 | float64](incr setMap[N]) expectFunc {
	h := metricdata.Histogram[N]{Temporality: metricdata.DeltaTemporality}
	return func(m int) metricdata.Aggregation {
		h.DataPoints = make([]metricdata.HistogramDataPoint[N], 0, len(incr))
		for a, v := range incr {
			h.DataPoints = append(h.DataPoints, hPoint[N](a, v, uint64(m)))
		}
		return h
	}
}

func cumuHistExpecter[N int64 | float64](incr setMap[N]) expectFunc {
	var cycle int
	h := metricdata.Histogram[N]{Temporality: metricdata.CumulativeTemporality}
	return func(m int) metricdata.Aggregation {
		cycle++
		h.DataPoints = make([]metricdata.HistogramDataPoint[N], 0, len(incr))
		for a, v := range incr {
			h.DataPoints = append(h.DataPoints, hPoint[N](a, v, uint64(cycle*m)))
		}
		return h
	}
}

// hPoint returns an HistogramDataPoint that started and ended now with multi
// number of measurements values v. It includes a min and max (set to v).
func hPoint[N int64 | float64](a attribute.Set, v N, multi uint64) metricdata.HistogramDataPoint[N] {
	idx := sort.SearchFloat64s(bounds, float64(v))
	counts := make([]uint64, len(bounds)+1)
	counts[idx] += multi
	return metricdata.HistogramDataPoint[N]{
		Attributes:   a,
		StartTime:    now(),
		Time:         now(),
		Count:        multi,
		Bounds:       bounds,
		BucketCounts: counts,
		Min:          metricdata.NewExtrema(v),
		Max:          metricdata.NewExtrema(v),
		Sum:          v * N(multi),
	}
}

func TestBucketsBin(t *testing.T) {
	t.Run("Int64", testBucketsBin[int64]())
	t.Run("Float64", testBucketsBin[float64]())
}

func testBucketsBin[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		b := newBuckets[N](3)
		assertB := func(counts []uint64, count uint64, sum, min, max N) {
			assert.Equal(t, counts, b.counts)
			assert.Equal(t, count, b.count)
			assert.Equal(t, sum, b.sum)
			assert.Equal(t, min, b.min)
			assert.Equal(t, max, b.max)
		}

		assertB([]uint64{0, 0, 0}, 0, 0, 0, 0)
		b.bin(1, 2)
		assertB([]uint64{0, 1, 0}, 1, 2, 0, 2)
		b.bin(0, -1)
		assertB([]uint64{1, 1, 0}, 2, 1, -1, 2)
	}
}

func testHistImmutableBounds[N int64 | float64](newA func(aggregation.ExplicitBucketHistogram) Aggregator[N], getBounds func(Aggregator[N]) []float64) func(t *testing.T) {
	b := []float64{0, 1, 2}
	cpB := make([]float64, len(b))
	copy(cpB, b)

	a := newA(aggregation.ExplicitBucketHistogram{Boundaries: b})
	return func(t *testing.T) {
		require.Equal(t, cpB, getBounds(a))

		b[0] = 10
		assert.Equal(t, cpB, getBounds(a), "modifying the bounds argument should not change the bounds")

		a.Aggregate(5, alice)
		hdp := a.Aggregation().(metricdata.Histogram[N]).DataPoints[0]
		hdp.Bounds[1] = 10
		assert.Equal(t, cpB, getBounds(a), "modifying the Aggregation bounds should not change the bounds")
	}
}

func TestHistogramImmutableBounds(t *testing.T) {
	t.Run("Delta", testHistImmutableBounds(
		NewDeltaHistogram[int64],
		func(a Aggregator[int64]) []float64 {
			deltaH := a.(*deltaHistogram[int64])
			return deltaH.bounds
		},
	))

	t.Run("Cumulative", testHistImmutableBounds(
		NewCumulativeHistogram[int64],
		func(a Aggregator[int64]) []float64 {
			cumuH := a.(*cumulativeHistogram[int64])
			return cumuH.bounds
		},
	))
}

func TestCumulativeHistogramImutableCounts(t *testing.T) {
	a := NewCumulativeHistogram[int64](histConf)
	a.Aggregate(5, alice)
	hdp := a.Aggregation().(metricdata.Histogram[int64]).DataPoints[0]

	cumuH := a.(*cumulativeHistogram[int64])
	require.Equal(t, hdp.BucketCounts, cumuH.values[alice].counts)

	cpCounts := make([]uint64, len(hdp.BucketCounts))
	copy(cpCounts, hdp.BucketCounts)
	hdp.BucketCounts[0] = 10
	assert.Equal(t, cpCounts, cumuH.values[alice].counts, "modifying the Aggregator bucket counts should not change the Aggregator")
}

func TestDeltaHistogramReset(t *testing.T) {
	t.Cleanup(mockTime(now))

	a := NewDeltaHistogram[int64](histConf)
	assert.Nil(t, a.Aggregation())

	a.Aggregate(1, alice)
	expect := metricdata.Histogram[int64]{Temporality: metricdata.DeltaTemporality}
	expect.DataPoints = []metricdata.HistogramDataPoint[int64]{hPoint[int64](alice, 1, 1)}
	metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

	// The attr set should be forgotten once Aggregations is called.
	expect.DataPoints = nil
	assert.Nil(t, a.Aggregation())

	// Aggregating another set should not affect the original (alice).
	a.Aggregate(1, bob)
	expect.DataPoints = []metricdata.HistogramDataPoint[int64]{hPoint[int64](bob, 1, 1)}
	metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())
}

func TestEmptyHistogramNilAggregation(t *testing.T) {
	assert.Nil(t, NewCumulativeHistogram[int64](histConf).Aggregation())
	assert.Nil(t, NewCumulativeHistogram[float64](histConf).Aggregation())
	assert.Nil(t, NewDeltaHistogram[int64](histConf).Aggregation())
	assert.Nil(t, NewDeltaHistogram[float64](histConf).Aggregation())
}

func BenchmarkHistogram(b *testing.B) {
	b.Run("Int64", benchmarkHistogram[int64])
	b.Run("Float64", benchmarkHistogram[float64])
}

func benchmarkHistogram[N int64 | float64](b *testing.B) {
	factory := func() Aggregator[N] { return NewDeltaHistogram[N](histConf) }
	b.Run("Delta", benchmarkAggregator(factory))
	factory = func() Aggregator[N] { return NewCumulativeHistogram[N](histConf) }
	b.Run("Cumulative", benchmarkAggregator(factory))
}
