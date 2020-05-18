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

package sum

import (
	"context"
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	ottest "go.opentelemetry.io/otel/internal/testing"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 100

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []ottest.FieldOffset{
		{
			Name:   "Aggregator.current",
			Offset: unsafe.Offsetof(Aggregator{}.current),
		},
		{
			Name:   "Aggregator.checkpoint",
			Offset: unsafe.Offsetof(Aggregator{}.checkpoint),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestCounterSum(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		descriptor := test.NewAggregatorTest(metric.CounterKind, profile.NumberKind)

		sum := metric.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			test.CheckedUpdate(t, agg, x, descriptor)
		}

		agg.Checkpoint(ctx, descriptor)

		asum, err := agg.Sum()
		require.Equal(t, sum, asum, "Same sum - monotonic")
		require.Nil(t, err)
	})
}

func TestValueRecorderSum(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

		sum := metric.Number(0)

		for i := 0; i < count; i++ {
			r1 := profile.Random(+1)
			r2 := profile.Random(-1)
			test.CheckedUpdate(t, agg, r1, descriptor)
			test.CheckedUpdate(t, agg, r2, descriptor)
			sum.AddNumber(profile.NumberKind, r1)
			sum.AddNumber(profile.NumberKind, r2)
		}

		agg.Checkpoint(ctx, descriptor)

		asum, err := agg.Sum()
		require.Equal(t, sum, asum, "Same sum - monotonic")
		require.Nil(t, err)
	})
}

func TestCounterMerge(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg1 := New()
		agg2 := New()

		descriptor := test.NewAggregatorTest(metric.CounterKind, profile.NumberKind)

		sum := metric.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			test.CheckedUpdate(t, agg1, x, descriptor)
			test.CheckedUpdate(t, agg2, x, descriptor)
		}

		agg1.Checkpoint(ctx, descriptor)
		agg2.Checkpoint(ctx, descriptor)

		test.CheckedMerge(t, agg1, agg2, descriptor)

		sum.AddNumber(descriptor.NumberKind(), sum)

		asum, err := agg1.Sum()
		require.Equal(t, sum, asum, "Same sum - monotonic")
		require.Nil(t, err)
	})
}
