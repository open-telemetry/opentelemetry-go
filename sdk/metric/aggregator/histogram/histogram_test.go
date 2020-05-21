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

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
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

	boundaries = []float64{500, 250, 750}
)

func TestHistogramAbsolute(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		testHistogram(t, profile, positiveOnly)
	})
}

func TestHistogramNegativeOnly(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		testHistogram(t, profile, negativeOnly)
	})
}

func TestHistogramPositiveAndNegative(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		testHistogram(t, profile, positiveAndNegative)
	})
}

// Validates count, sum and buckets for a given profile and policy
func testHistogram(t *testing.T, profile test.Profile, policy policy) {
	ctx := context.Background()
	descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

	agg := histogram.New(descriptor, boundaries)

	all := test.NewNumbers(profile.NumberKind)

	for i := 0; i < count; i++ {
		x := profile.Random(policy.sign())
		all.Append(x)
		test.CheckedUpdate(t, agg, x, descriptor)
	}

	agg.Checkpoint(ctx, descriptor)

	all.Sort()

	asum, err := agg.Sum()
	sum := all.Sum()
	require.InEpsilon(t,
		sum.CoerceToFloat64(profile.NumberKind),
		asum.CoerceToFloat64(profile.NumberKind),
		0.000000001,
		"Same sum - "+policy.name)
	require.NoError(t, err)

	count, err := agg.Count()
	require.Equal(t, all.Count(), count, "Same count -"+policy.name)
	require.NoError(t, err)

	buckets, err := agg.Histogram()
	require.NoError(t, err)

	require.Equal(t, len(buckets.Counts), len(boundaries)+1, "There should be b + 1 counts, where b is the number of boundaries")

	counts := calcBuckets(all.Points(), profile)
	for i, v := range counts {
		bCount := uint64(buckets.Counts[i])
		require.Equal(t, v, bCount, "Wrong bucket #%d count: %v != %v", i, counts, buckets.Counts)
	}
}

func TestHistogramInitial(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

		agg := histogram.New(descriptor, boundaries)
		buckets, err := agg.Histogram()

		require.NoError(t, err)
		require.Equal(t, len(buckets.Counts), len(boundaries)+1)
		require.Equal(t, len(buckets.Boundaries), len(boundaries))
	})
}

func TestHistogramMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

		agg1 := histogram.New(descriptor, boundaries)
		agg2 := histogram.New(descriptor, boundaries)

		all := test.NewNumbers(profile.NumberKind)

		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all.Append(x)
			test.CheckedUpdate(t, agg1, x, descriptor)
		}
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all.Append(x)
			test.CheckedUpdate(t, agg2, x, descriptor)
		}

		agg1.Checkpoint(ctx, descriptor)
		agg2.Checkpoint(ctx, descriptor)

		test.CheckedMerge(t, agg1, agg2, descriptor)

		all.Sort()

		asum, err := agg1.Sum()
		sum := all.Sum()
		require.InEpsilon(t,
			sum.CoerceToFloat64(profile.NumberKind),
			asum.CoerceToFloat64(profile.NumberKind),
			0.000000001,
			"Same sum - absolute")
		require.NoError(t, err)

		count, err := agg1.Count()
		require.Equal(t, all.Count(), count, "Same count - absolute")
		require.NoError(t, err)

		buckets, err := agg1.Histogram()
		require.NoError(t, err)

		require.Equal(t, len(buckets.Counts), len(boundaries)+1, "There should be b + 1 counts, where b is the number of boundaries")

		counts := calcBuckets(all.Points(), profile)
		for i, v := range counts {
			bCount := uint64(buckets.Counts[i])
			require.Equal(t, v, bCount, "Wrong bucket #%d count: %v != %v", i, counts, buckets.Counts)
		}
	})
}

func TestHistogramNotSet(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

		agg := histogram.New(descriptor, boundaries)
		agg.Checkpoint(ctx, descriptor)

		asum, err := agg.Sum()
		require.Equal(t, metric.Number(0), asum, "Empty checkpoint sum = 0")
		require.NoError(t, err)

		count, err := agg.Count()
		require.Equal(t, int64(0), count, "Empty checkpoint count = 0")
		require.NoError(t, err)

		buckets, err := agg.Histogram()
		require.NoError(t, err)

		require.Equal(t, len(buckets.Counts), len(boundaries)+1, "There should be b + 1 counts, where b is the number of boundaries")
		for i, bCount := range buckets.Counts {
			require.Equal(t, uint64(0), uint64(bCount), "Bucket #%d must have 0 observed values", i)
		}
	})
}

func calcBuckets(points []metric.Number, profile test.Profile) []uint64 {
	sortedBoundaries := make([]float64, len(boundaries))

	copy(sortedBoundaries, boundaries)
	sort.Float64s(sortedBoundaries)

	counts := make([]uint64, len(sortedBoundaries)+1)
	idx := 0
	for _, p := range points {
		for idx < len(sortedBoundaries) && p.CoerceToFloat64(profile.NumberKind) >= sortedBoundaries[idx] {
			idx++
		}
		counts[idx]++
	}

	return counts
}
