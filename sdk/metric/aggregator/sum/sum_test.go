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
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	ottest "go.opentelemetry.io/otel/internal/testing"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
)

const count = 100

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []ottest.FieldOffset{
		{
			Name:   "Aggregator.value",
			Offset: unsafe.Offsetof(Aggregator{}.value),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func new2() (_, _ *Aggregator) {
	alloc := New(2)
	return &alloc[0], &alloc[1]
}

func new4() (_, _, _, _ *Aggregator) {
	alloc := New(4)
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func checkZero(t *testing.T, agg *Aggregator, desc *otel.Descriptor) {
	kind := desc.NumberKind()

	sum, err := agg.Sum()
	require.NoError(t, err)
	require.Equal(t, kind.Zero(), sum)
}

func TestCounterSum(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg, ckpt := new2()

		descriptor := aggregatortest.NewAggregatorTest(otel.CounterInstrumentKind, profile.NumberKind)

		sum := otel.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			aggregatortest.CheckedUpdate(t, agg, x, descriptor)
		}

		err := agg.SynchronizedMove(ckpt, descriptor)
		require.NoError(t, err)

		checkZero(t, agg, descriptor)

		asum, err := ckpt.Sum()
		require.Equal(t, sum, asum, "Same sum - monotonic")
		require.Nil(t, err)
	})
}

func TestValueRecorderSum(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg, ckpt := new2()

		descriptor := aggregatortest.NewAggregatorTest(otel.ValueRecorderInstrumentKind, profile.NumberKind)

		sum := otel.Number(0)

		for i := 0; i < count; i++ {
			r1 := profile.Random(+1)
			r2 := profile.Random(-1)
			aggregatortest.CheckedUpdate(t, agg, r1, descriptor)
			aggregatortest.CheckedUpdate(t, agg, r2, descriptor)
			sum.AddNumber(profile.NumberKind, r1)
			sum.AddNumber(profile.NumberKind, r2)
		}

		require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))
		checkZero(t, agg, descriptor)

		asum, err := ckpt.Sum()
		require.Equal(t, sum, asum, "Same sum - monotonic")
		require.Nil(t, err)
	})
}

func TestCounterMerge(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg1, agg2, ckpt1, ckpt2 := new4()

		descriptor := aggregatortest.NewAggregatorTest(otel.CounterInstrumentKind, profile.NumberKind)

		sum := otel.Number(0)
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			sum.AddNumber(profile.NumberKind, x)
			aggregatortest.CheckedUpdate(t, agg1, x, descriptor)
			aggregatortest.CheckedUpdate(t, agg2, x, descriptor)
		}

		require.NoError(t, agg1.SynchronizedMove(ckpt1, descriptor))
		require.NoError(t, agg2.SynchronizedMove(ckpt2, descriptor))

		checkZero(t, agg1, descriptor)
		checkZero(t, agg2, descriptor)

		aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

		sum.AddNumber(descriptor.NumberKind(), sum)

		asum, err := ckpt1.Sum()
		require.Equal(t, sum, asum, "Same sum - monotonic")
		require.Nil(t, err)
	})
}
