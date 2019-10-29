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

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
	"go.opentelemetry.io/sdk/metric/aggregator/test"
)

const count = 100

func TestCounterMonotonic(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		batcher, record := test.NewAggregatorTest(export.CounterMetricKind, profile.NumberKind, false)

		sum := core.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			agg.Update(ctx, x, record)
		}

		agg.Collect(ctx, record, batcher)

		require.Equal(t, sum, agg.AsNumber(), "Same sum - monotonic")
	})
}

func TestCounterMonotonicNegative(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		batcher, record := test.NewAggregatorTest(export.CounterMetricKind, profile.NumberKind, false)

		for i := 0; i < count; i++ {
			agg.Update(ctx, profile.Random(-1), record)
		}

		sum := profile.Random(+1)
		agg.Update(ctx, sum, record)
		agg.Collect(ctx, record, batcher)

		require.Equal(t, sum, agg.AsNumber(), "Same sum - monotonic")
	})
}

func TestCounterNonMonotonic(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		batcher, record := test.NewAggregatorTest(export.CounterMetricKind, profile.NumberKind, true)

		sum := core.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			y := profile.Random(-1)
			sum.AddNumber(profile.NumberKind, x)
			sum.AddNumber(profile.NumberKind, y)
			agg.Update(ctx, x, record)
			agg.Update(ctx, y, record)
		}

		agg.Collect(ctx, record, batcher)

		require.Equal(t, sum, agg.AsNumber(), "Same sum - monotonic")
	})
}
