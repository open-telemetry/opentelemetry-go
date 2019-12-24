// Copyright 2019, OpenTelemetry Authors
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
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
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
	descriptor := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, !policy.absolute)

	agg := New(descriptor, []float64{250, 500, 700})

	all := test.NewNumbers(profile.NumberKind)

	for i := 0; i < count; i++ {
		x := profile.Random(policy.sign())
		all.Append(x)
		test.CheckedUpdate(t, agg, x, descriptor)
	}

	agg.Checkpoint(ctx, descriptor)

	all.Sort()

	asum, err := agg.Sum()
	require.InEpsilon(t,
		all.Sum().CoerceToFloat64(profile.NumberKind),
		asum.CoerceToFloat64(profile.NumberKind),
		0.000000001,
		"Same sum - "+policy.name)
	require.Nil(t, err)

	count, err := agg.Count()
	require.Equal(t, all.Count(), count, "Same count -"+policy.name)
	require.Nil(t, err)

	for _, p := range all.Points() {
		fmt.Print(p.Emit(profile.NumberKind), " ")
	}
	fmt.Println()
	fmt.Println(agg.checkpoint)
}

func TestHistogramMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, false)

		agg1 := New(descriptor, []float64{250, 500, 700})
		agg2 := New(descriptor, []float64{250, 500, 700})

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
		require.InEpsilon(t,
			all.Sum().CoerceToFloat64(profile.NumberKind),
			asum.CoerceToFloat64(profile.NumberKind),
			0.000000001,
			"Same sum - absolute")
		require.Nil(t, err)

		count, err := agg1.Count()
		require.Equal(t, all.Count(), count, "Same count - absolute")
		require.Nil(t, err)

	})
}

func TestHistogramNotSet(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, false)

		agg := New(descriptor, []float64{250, 500, 700})
		agg.Checkpoint(ctx, descriptor)

		asum, err := agg.Sum()
		require.Equal(t, core.Number(0), asum, "Empty checkpoint sum = 0")
		require.Nil(t, err)

		count, err := agg.Count()
		require.Equal(t, int64(0), count, "Empty checkpoint count = 0")
		require.Nil(t, err)

	})
}
