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

	"go.opentelemetry.io/otel/sdk/export"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 100

func TestMaxSumCountAbsolute(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		batcher, record := test.NewAggregatorTest(export.MeasureMetricKind, profile.NumberKind, false)

		agg := New()

		var all test.Numbers
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all = append(all, x)
			agg.Update(ctx, x, record)
		}

		agg.Collect(ctx, record, batcher)

		all.Sort()

		require.InEpsilon(t,
			all.Sum(profile.NumberKind).CoerceToFloat64(profile.NumberKind),
			agg.Sum().CoerceToFloat64(profile.NumberKind),
			0.000000001,
			"Same sum - absolute")
		require.Equal(t, all.Count(), agg.Count(), "Same count - absolute")
		require.Equal(t,
			all[len(all)-1],
			agg.Max(),
			"Same max - absolute")
	})
}

func TestMaxSumCountMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		batcher, record := test.NewAggregatorTest(export.MeasureMetricKind, profile.NumberKind, false)

		agg1 := New()
		agg2 := New()

		var all test.Numbers

		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all = append(all, x)
			agg1.Update(ctx, x, record)
		}
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all = append(all, x)
			agg2.Update(ctx, x, record)
		}

		agg1.Collect(ctx, record, batcher)
		agg2.Collect(ctx, record, batcher)

		agg1.Merge(agg2, record.Descriptor())

		all.Sort()

		require.InEpsilon(t,
			all.Sum(profile.NumberKind).CoerceToFloat64(profile.NumberKind),
			agg1.Sum().CoerceToFloat64(profile.NumberKind),
			0.000000001,
			"Same sum - absolute")
		require.Equal(t, all.Count(), agg1.Count(), "Same count - absolute")
		require.Equal(t,
			all[len(all)-1],
			agg1.Max(),
			"Same max - absolute")
	})
}
