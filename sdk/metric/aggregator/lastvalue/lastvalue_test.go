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

package lastvalue

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	ottest "go.opentelemetry.io/otel/internal/testing"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 100

var _ export.Aggregator = &Aggregator{}

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []ottest.FieldOffset{
		{
			Name:   "lastValueData.value",
			Offset: unsafe.Offsetof(lastValueData{}.value),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestLastValueUpdate(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		record := test.NewAggregatorTest(metric.ValueObserverKind, profile.NumberKind)

		var last metric.Number
		for i := 0; i < count; i++ {
			x := profile.Random(rand.Intn(1)*2 - 1)
			last = x
			test.CheckedUpdate(t, agg, x, record)
		}

		agg.Checkpoint(ctx, record)

		lv, _, err := agg.LastValue()
		require.Equal(t, last, lv, "Same last value - non-monotonic")
		require.Nil(t, err)
	})
}

func TestLastValueMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg1 := New()
		agg2 := New()

		descriptor := test.NewAggregatorTest(metric.ValueObserverKind, profile.NumberKind)

		first1 := profile.Random(+1)
		first2 := profile.Random(+1)
		first1.AddNumber(profile.NumberKind, first2)

		test.CheckedUpdate(t, agg1, first1, descriptor)
		test.CheckedUpdate(t, agg2, first2, descriptor)

		agg1.Checkpoint(ctx, descriptor)
		agg2.Checkpoint(ctx, descriptor)

		_, t1, err := agg1.LastValue()
		require.Nil(t, err)
		_, t2, err := agg2.LastValue()
		require.Nil(t, err)
		require.True(t, t1.Before(t2))

		test.CheckedMerge(t, agg1, agg2, descriptor)

		lv, ts, err := agg1.LastValue()
		require.Nil(t, err)
		require.Equal(t, t2, ts, "Merged timestamp - non-monotonic")
		require.Equal(t, first2, lv, "Merged value - non-monotonic")
	})
}

func TestLastValueNotSet(t *testing.T) {
	descriptor := test.NewAggregatorTest(metric.ValueObserverKind, metric.Int64NumberKind)

	g := New()
	g.Checkpoint(context.Background(), descriptor)

	value, timestamp, err := g.LastValue()
	require.Equal(t, aggregator.ErrNoData, err)
	require.True(t, timestamp.IsZero())
	require.Equal(t, metric.Number(0), value)
}
