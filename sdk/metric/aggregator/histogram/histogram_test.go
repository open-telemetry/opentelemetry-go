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

package histogram_test

import (
	"context"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
)

const count = 100

type policy struct {
	name     string
	absolute bool
	sign     func() int
}

var (
	positiveOnly = policy{
		name:     "absolute",
		absolute: true,
		sign:     func() int { return +1 },
	}
	negativeOnly = policy{
		name:     "negative",
		absolute: false,
		sign:     func() int { return -1 },
	}
	positiveAndNegative = policy{
		name:     "positiveAndNegative",
		absolute: false,
		sign: func() int {
			if rand.Uint32() > math.MaxUint32/2 {
				return -1
			}
			return 1
		},
	}

	testBoundaries = []float64{500, 250, 750}
)

func new2(desc *metric.Descriptor, options ...histogram.Option) (_, _ *histogram.Aggregator) {
	alloc := histogram.New(2, desc, options...)
	return &alloc[0], &alloc[1]
}

func new4(desc *metric.Descriptor, options ...histogram.Option) (_, _, _, _ *histogram.Aggregator) {
	alloc := histogram.New(4, desc, options...)
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func checkZero(t *testing.T, agg *histogram.Aggregator, desc *metric.Descriptor) {
	asum, err := agg.Sum()
	require.Equal(t, number.Number(0), asum, "Empty checkpoint sum = 0")
	require.NoError(t, err)

	count, err := agg.Count()
	require.Equal(t, uint64(0), count, "Empty checkpoint count = 0")
	require.NoError(t, err)

	buckets, err := agg.Histogram()
	require.NoError(t, err)

	require.Equal(t, len(buckets.Counts), len(testBoundaries)+1, "There should be b + 1 counts, where b is the number of boundaries")
	for i, bCount := range buckets.Counts {
		require.Equal(t, uint64(0), uint64(bCount), "Bucket #%d must have 0 observed values", i)
	}
}

func TestHistogramAbsolute(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		testHistogram(t, profile, positiveOnly)
	})
}

func TestHistogramNegativeOnly(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		testHistogram(t, profile, negativeOnly)
	})
}

func TestHistogramPositiveAndNegative(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		testHistogram(t, profile, positiveAndNegative)
	})
}

// Validates count, sum and buckets for a given profile and policy
func testHistogram(t *testing.T, profile aggregatortest.Profile, policy policy) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

	agg, ckpt := new2(descriptor, histogram.WithExplicitBoundaries(testBoundaries))

	// This needs to repeat at least 3 times to uncover a failure to reset
	// for the overall sum and count fields, since the third time through
	// is the first time a `histogram.state` object is reused.
	for repeat := 0; repeat < 3; repeat++ {
		all := aggregatortest.NewNumbers(profile.NumberKind)

		for i := 0; i < count; i++ {
			x := profile.Random(policy.sign())
			all.Append(x)
			aggregatortest.CheckedUpdate(t, agg, x, descriptor)
		}

		require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

		checkZero(t, agg, descriptor)

		checkHistogram(t, all, profile, ckpt)
	}
}

func TestHistogramInitial(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		agg := &histogram.New(1, descriptor, histogram.WithExplicitBoundaries(testBoundaries))[0]
		buckets, err := agg.Histogram()

		require.NoError(t, err)
		require.Equal(t, len(buckets.Counts), len(testBoundaries)+1)
		require.Equal(t, len(buckets.Boundaries), len(testBoundaries))
	})
}

func TestHistogramMerge(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		agg1, agg2, ckpt1, ckpt2 := new4(descriptor, histogram.WithExplicitBoundaries(testBoundaries))

		all := aggregatortest.NewNumbers(profile.NumberKind)

		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all.Append(x)
			aggregatortest.CheckedUpdate(t, agg1, x, descriptor)
		}
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all.Append(x)
			aggregatortest.CheckedUpdate(t, agg2, x, descriptor)
		}

		require.NoError(t, agg1.SynchronizedMove(ckpt1, descriptor))
		require.NoError(t, agg2.SynchronizedMove(ckpt2, descriptor))

		aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

		checkHistogram(t, all, profile, ckpt1)
	})
}

func TestHistogramNotSet(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		agg, ckpt := new2(descriptor, histogram.WithExplicitBoundaries(testBoundaries))

		err := agg.SynchronizedMove(ckpt, descriptor)
		require.NoError(t, err)

		checkZero(t, agg, descriptor)
		checkZero(t, ckpt, descriptor)
	})
}

// checkHistogram ensures the correct aggregated state between `all`
// (test aggregator) and `agg` (code under test).
func checkHistogram(t *testing.T, all aggregatortest.Numbers, profile aggregatortest.Profile, agg *histogram.Aggregator) {

	all.Sort()

	asum, err := agg.Sum()
	require.NoError(t, err)

	sum := all.Sum()
	require.InEpsilon(t,
		sum.CoerceToFloat64(profile.NumberKind),
		asum.CoerceToFloat64(profile.NumberKind),
		0.000000001)

	count, err := agg.Count()
	require.NoError(t, err)
	require.Equal(t, all.Count(), count)

	buckets, err := agg.Histogram()
	require.NoError(t, err)

	require.Equal(t, len(buckets.Counts), len(testBoundaries)+1,
		"There should be b + 1 counts, where b is the number of boundaries")

	sortedBoundaries := make([]float64, len(testBoundaries))
	copy(sortedBoundaries, testBoundaries)

	sort.Float64s(sortedBoundaries)

	require.EqualValues(t, sortedBoundaries, buckets.Boundaries)

	counts := make([]uint64, len(sortedBoundaries)+1)
	idx := 0
	for _, p := range all.Points() {
		for idx < len(sortedBoundaries) && p.CoerceToFloat64(profile.NumberKind) >= sortedBoundaries[idx] {
			idx++
		}
		counts[idx]++
	}
	for i, v := range counts {
		bCount := uint64(buckets.Counts[i])
		require.Equal(t, v, bCount, "Wrong bucket #%d count: %v != %v", i, counts, buckets.Counts)
	}
}

func TestSynchronizedMoveReset(t *testing.T) {
	aggregatortest.SynchronizedMoveResetTest(
		t,
		metric.ValueRecorderInstrumentKind,
		func(desc *metric.Descriptor) export.Aggregator {
			return &histogram.New(1, desc, histogram.WithExplicitBoundaries(testBoundaries))[0]
		},
	)
}

func TestHistogramDefaultBoundaries(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		ctx := context.Background()
		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		agg, ckpt := new2(descriptor)

		bounds := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10} // len 11
		values := append(bounds, 100)                                         // len 12
		expect := []uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}                // len 12

		for _, value := range values {
			var num number.Number

			value -= .001 // Avoid exact boundaries

			if descriptor.NumberKind() == number.Int64Kind {
				value *= 1e6
				num = number.NewInt64Number(int64(value))
			} else {
				num = number.NewFloat64Number(value)
			}

			require.NoError(t, agg.Update(ctx, num, descriptor))
		}

		bucks, err := agg.Histogram()
		require.NoError(t, err)

		// Check for proper lengths, 1 count in each bucket.
		require.Equal(t, len(values), len(bucks.Counts))
		require.Equal(t, len(bounds), len(bucks.Boundaries))
		require.EqualValues(t, expect, bucks.Counts)

		require.Equal(t, expect, bucks.Counts)

		// Move and repeat the test on `ckpt`.
		err = agg.SynchronizedMove(ckpt, descriptor)
		require.NoError(t, err)

		bucks, err = ckpt.Histogram()
		require.NoError(t, err)

		require.Equal(t, len(values), len(bucks.Counts))
		require.Equal(t, len(bounds), len(bucks.Boundaries))
		require.EqualValues(t, expect, bucks.Counts)
	})
}
