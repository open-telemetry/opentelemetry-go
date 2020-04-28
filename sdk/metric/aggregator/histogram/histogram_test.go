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

package histogram

import (
	"context"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
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

	boundaries = map[core.NumberKind][]core.Number{
		core.Float64NumberKind: {core.NewFloat64Number(500), core.NewFloat64Number(250), core.NewFloat64Number(750)},
		core.Int64NumberKind:   {core.NewInt64Number(500), core.NewInt64Number(250), core.NewInt64Number(750)},
	}
)

func TestHistogramAbsolute(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		histogram(t, profile, positiveOnly)
	})
}

func TestHistogramNegativeOnly(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		histogram(t, profile, negativeOnly)
	})
}

func TestHistogramPositiveAndNegative(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		histogram(t, profile, positiveAndNegative)
	})
}

// Validates count, sum and buckets for a given profile and policy
func histogram(t *testing.T, profile test.Profile, policy policy) {
	ctx := context.Background()
	descriptor := test.NewAggregatorTest(metric.MeasureKind, profile.NumberKind)

	agg := New(descriptor, boundaries[profile.NumberKind])

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
	require.Nil(t, err)

	count, err := agg.Count()
	require.Equal(t, all.Count(), count, "Same count -"+policy.name)
	require.Nil(t, err)

	require.Equal(t, len(agg.checkpoint.buckets.Counts), len(boundaries[profile.NumberKind])+1, "There should be b + 1 counts, where b is the number of boundaries")

	counts := calcBuckets(all.Points(), profile)
	for i, v := range counts {
		bCount := agg.checkpoint.buckets.Counts[i].AsUint64()
		require.Equal(t, v, bCount, "Wrong bucket #%d count: %v != %v", i, counts, agg.checkpoint.buckets.Counts)
	}
}

func TestHistogramMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(metric.MeasureKind, profile.NumberKind)

		agg1 := New(descriptor, boundaries[profile.NumberKind])
		agg2 := New(descriptor, boundaries[profile.NumberKind])

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
		require.Nil(t, err)

		count, err := agg1.Count()
		require.Equal(t, all.Count(), count, "Same count - absolute")
		require.Nil(t, err)

		require.Equal(t, len(agg1.checkpoint.buckets.Counts), len(boundaries[profile.NumberKind])+1, "There should be b + 1 counts, where b is the number of boundaries")

		counts := calcBuckets(all.Points(), profile)
		for i, v := range counts {
			bCount := agg1.checkpoint.buckets.Counts[i].AsUint64()
			require.Equal(t, v, bCount, "Wrong bucket #%d count: %v != %v", i, counts, agg1.checkpoint.buckets.Counts)
		}
	})
}

func TestHistogramNotSet(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(metric.MeasureKind, profile.NumberKind)

		agg := New(descriptor, boundaries[profile.NumberKind])
		agg.Checkpoint(ctx, descriptor)

		asum, err := agg.Sum()
		require.Equal(t, core.Number(0), asum, "Empty checkpoint sum = 0")
		require.Nil(t, err)

		count, err := agg.Count()
		require.Equal(t, int64(0), count, "Empty checkpoint count = 0")
		require.Nil(t, err)

		require.Equal(t, len(agg.checkpoint.buckets.Counts), len(boundaries[profile.NumberKind])+1, "There should be b + 1 counts, where b is the number of boundaries")
		for i, bCount := range agg.checkpoint.buckets.Counts {
			require.Equal(t, uint64(0), bCount.AsUint64(), "Bucket #%d must have 0 observed values", i)
		}
	})
}

func calcBuckets(points []core.Number, profile test.Profile) []uint64 {
	sortedBoundaries := numbers{
		numbers: make([]core.Number, len(boundaries[profile.NumberKind])),
		kind:    profile.NumberKind,
	}

	copy(sortedBoundaries.numbers, boundaries[profile.NumberKind])
	sort.Sort(&sortedBoundaries)
	boundaries := sortedBoundaries.numbers

	counts := make([]uint64, len(boundaries)+1)
	idx := 0
	for _, p := range points {
		for idx < len(boundaries) && p.CompareNumber(profile.NumberKind, boundaries[idx]) != -1 {
			idx++
		}
		counts[idx]++
	}

	return counts
}
