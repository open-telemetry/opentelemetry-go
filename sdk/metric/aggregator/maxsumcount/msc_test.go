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

package maxsumcount

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 100

func TestMaxSumCountAbsolute(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		record := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, false)

		agg := New()

		all := test.NewNumbers(profile.NumberKind)

		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all.Append(x)
			test.CheckedUpdate(t, agg, x, record)
		}

		agg.Checkpoint(ctx, record)

		all.Sort()

		asum, err := agg.Sum()
		require.InEpsilon(t,
			all.Sum().CoerceToFloat64(profile.NumberKind),
			asum.CoerceToFloat64(profile.NumberKind),
			0.000000001,
			"Same sum - absolute")
		require.Nil(t, err)

		count, err := agg.Count()
		require.Equal(t, all.Count(), count, "Same count - absolute")
		require.Nil(t, err)

		max, err := agg.Max()
		require.Nil(t, err)
		require.Equal(t,
			all.Max(),
			max,
			"Same max - absolute")
	})
}

func TestMaxSumCountMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		descriptor := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, false)

		agg1 := New()
		agg2 := New()

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

		max, err := agg1.Max()
		require.Nil(t, err)
		require.Equal(t,
			all.Max(),
			max,
			"Same max - absolute")
	})
}
