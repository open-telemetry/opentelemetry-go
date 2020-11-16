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

package minmaxsumcount

import (
	"errors"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
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
)

func TestMinMaxSumCountAbsolute(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		minMaxSumCount(t, profile, positiveOnly)
	})
}

func TestMinMaxSumCountNegativeOnly(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		minMaxSumCount(t, profile, negativeOnly)
	})
}

func TestMinMaxSumCountPositiveAndNegative(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		minMaxSumCount(t, profile, positiveAndNegative)
	})
}

func new2(desc *metric.Descriptor) (_, _ *Aggregator) {
	alloc := New(2, desc)
	return &alloc[0], &alloc[1]
}

func new4(desc *metric.Descriptor) (_, _, _, _ *Aggregator) {
	alloc := New(4, desc)
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func checkZero(t *testing.T, agg *Aggregator, desc *metric.Descriptor) {
	kind := desc.NumberKind()

	sum, err := agg.Sum()
	require.NoError(t, err)
	require.Equal(t, kind.Zero(), sum)

	count, err := agg.Count()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	max, err := agg.Max()
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, kind.Zero(), max)

	min, err := agg.Min()
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, kind.Zero(), min)
}

// Validates min, max, sum and count for a given profile and policy
func minMaxSumCount(t *testing.T, profile aggregatortest.Profile, policy policy) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

	agg, ckpt := new2(descriptor)

	all := aggregatortest.NewNumbers(profile.NumberKind)

	for i := 0; i < count; i++ {
		x := profile.Random(policy.sign())
		all.Append(x)
		aggregatortest.CheckedUpdate(t, agg, x, descriptor)
	}

	require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

	checkZero(t, agg, descriptor)

	all.Sort()

	aggSum, err := ckpt.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InEpsilon(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		aggSum.CoerceToFloat64(profile.NumberKind),
		0.000000001,
		"Same sum - "+policy.name)

	count, err := ckpt.Count()
	require.Equal(t, all.Count(), count, "Same count -"+policy.name)
	require.Nil(t, err)

	min, err := ckpt.Min()
	require.Nil(t, err)
	require.Equal(t,
		all.Min(),
		min,
		"Same min -"+policy.name)

	max, err := ckpt.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max -"+policy.name)
}

func TestMinMaxSumCountMerge(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		agg1, agg2, ckpt1, ckpt2 := new4(descriptor)

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

		checkZero(t, agg1, descriptor)
		checkZero(t, agg2, descriptor)

		aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

		all.Sort()

		aggSum, err := ckpt1.Sum()
		require.Nil(t, err)
		allSum := all.Sum()
		require.InEpsilon(t,
			(&allSum).CoerceToFloat64(profile.NumberKind),
			aggSum.CoerceToFloat64(profile.NumberKind),
			0.000000001,
			"Same sum - absolute")

		count, err := ckpt1.Count()
		require.Equal(t, all.Count(), count, "Same count - absolute")
		require.Nil(t, err)

		min, err := ckpt1.Min()
		require.Nil(t, err)
		require.Equal(t,
			all.Min(),
			min,
			"Same min - absolute")

		max, err := ckpt1.Max()
		require.Nil(t, err)
		require.Equal(t,
			all.Max(),
			max,
			"Same max - absolute")
	})
}

func TestMaxSumCountNotSet(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		alloc := New(2, descriptor)
		agg, ckpt := &alloc[0], &alloc[1]

		require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

		asum, err := ckpt.Sum()
		require.Equal(t, number.Number(0), asum, "Empty checkpoint sum = 0")
		require.Nil(t, err)

		count, err := ckpt.Count()
		require.Equal(t, int64(0), count, "Empty checkpoint count = 0")
		require.Nil(t, err)

		max, err := ckpt.Max()
		require.Equal(t, aggregation.ErrNoData, err)
		require.Equal(t, number.Number(0), max)
	})
}
