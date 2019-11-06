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

package counter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 100

func TestCounterMonotonic(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		descriptor := test.NewAggregatorTest(export.CounterKind, profile.NumberKind, false)

		sum := core.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			test.CheckedUpdate(ctx, agg, x, descriptor)
		}

		agg.Checkpoint(ctx, descriptor)

		require.Equal(t, sum, agg.Sum(), "Same sum - monotonic")
	})
}

func TestCounterMonotonicNegative(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		descriptor := test.NewAggregatorTest(export.CounterKind, profile.NumberKind, false)

		for i := 0; i < count; i++ {
			test.CheckedUpdate(ctx, agg, profile.Random(-1), descriptor)
		}

		sum := profile.Random(+1)
		test.CheckedUpdate(ctx, agg, sum, descriptor)
		agg.Checkpoint(ctx, descriptor)

		require.Equal(t, sum, agg.Sum(), "Same sum - monotonic")
	})
}

func TestCounterNonMonotonic(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		descriptor := test.NewAggregatorTest(export.CounterKind, profile.NumberKind, true)

		sum := core.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			y := profile.Random(-1)
			sum.AddNumber(profile.NumberKind, x)
			sum.AddNumber(profile.NumberKind, y)
			test.CheckedUpdate(ctx, agg, x, descriptor)
			test.CheckedUpdate(ctx, agg, y, descriptor)
		}

		agg.Checkpoint(ctx, descriptor)

		require.Equal(t, sum, agg.Sum(), "Same sum - monotonic")
	})
}

func TestCounterMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg1 := New()
		agg2 := New()

		descriptor := test.NewAggregatorTest(export.CounterKind, profile.NumberKind, false)

		sum := core.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			test.CheckedUpdate(ctx, agg1, x, descriptor)
			test.CheckedUpdate(ctx, agg2, x, descriptor)
		}

		agg1.Checkpoint(ctx, descriptor)
		agg2.Checkpoint(ctx, descriptor)

		agg1.Merge(agg2, descriptor)

		sum.AddNumber(descriptor.NumberKind(), sum)

		require.Equal(t, sum, agg1.Sum(), "Same sum - monotonic")
	})
}
