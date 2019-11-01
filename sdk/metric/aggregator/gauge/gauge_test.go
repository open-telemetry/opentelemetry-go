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

package gauge

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/export"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 100

var _ export.MetricAggregator = &Aggregator{}

func TestGaugeNonMonotonic(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		batcher, record := test.NewAggregatorTest(export.GaugeMetricKind, profile.NumberKind, false)

		var last core.Number
		for i := 0; i < count; i++ {
			x := profile.Random(rand.Intn(1)*2 - 1)
			last = x
			agg.Update(ctx, x, record)
		}

		agg.Collect(ctx, record, batcher)

		require.Equal(t, last, agg.AsNumber(), "Same last value - non-monotonic")
	})
}

func TestGaugeMonotonic(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		batcher, record := test.NewAggregatorTest(export.GaugeMetricKind, profile.NumberKind, true)

		small := profile.Random(+1)
		last := small
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			last.AddNumber(profile.NumberKind, x)
			agg.Update(ctx, last, record)
		}

		agg.Collect(ctx, record, batcher)

		require.Equal(t, last, agg.AsNumber(), "Same last value - monotonic")
	})
}

func TestGaugeMonotonicDescending(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		batcher, record := test.NewAggregatorTest(export.GaugeMetricKind, profile.NumberKind, true)

		first := profile.Random(+1)
		agg.Update(ctx, first, record)

		for i := 0; i < count; i++ {
			x := profile.Random(-1)
			agg.Update(ctx, x, record)
		}

		agg.Collect(ctx, record, batcher)

		require.Equal(t, first, agg.AsNumber(), "Same last value - monotonic")
	})
}

func TestGaugeNormalMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg1 := New()
		agg2 := New()

		batcher, record := test.NewAggregatorTest(export.GaugeMetricKind, profile.NumberKind, false)

		first1 := profile.Random(+1)
		first2 := profile.Random(+1)
		first1.AddNumber(profile.NumberKind, first2)

		agg1.Update(ctx, first1, record)
		agg2.Update(ctx, first2, record)

		agg1.Collect(ctx, record, batcher)
		agg2.Collect(ctx, record, batcher)

		t1 := agg1.Timestamp()
		t2 := agg2.Timestamp()
		require.True(t, t1.Before(t2))

		agg1.Merge(agg2, record.Descriptor())

		require.Equal(t, t2, agg1.Timestamp(), "Merged timestamp - non-monotonic")
		require.Equal(t, first2, agg1.AsNumber(), "Merged value - non-monotonic")
	})
}

func TestGaugeMonotonicMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg1 := New()
		agg2 := New()

		batcher, record := test.NewAggregatorTest(export.GaugeMetricKind, profile.NumberKind, true)

		first1 := profile.Random(+1)
		agg1.Update(ctx, first1, record)

		first2 := profile.Random(+1)
		first2.AddNumber(profile.NumberKind, first1)
		agg2.Update(ctx, first2, record)

		agg1.Collect(ctx, record, batcher)
		agg2.Collect(ctx, record, batcher)

		agg1.Merge(agg2, record.Descriptor())

		require.Equal(t, first2, agg1.AsNumber(), "Merged value - monotonic")
		require.Equal(t, agg2.Timestamp(), agg1.Timestamp(), "Merged timestamp - monotonic")
	})
}
