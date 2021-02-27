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
	"errors"
	"math/rand"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
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

func new2() (_, _ *Aggregator) {
	alloc := New(2)
	return &alloc[0], &alloc[1]
}

func new4() (_, _, _, _ *Aggregator) {
	alloc := New(4)
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func checkZero(t *testing.T, agg *Aggregator) {
	lv, ts, err := agg.LastValue()
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, time.Time{}, ts)
	require.Equal(t, number.Number(0), lv)
}

func TestLastValueUpdate(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg, ckpt := new2()

		record := aggregatortest.NewAggregatorTest(metric.ValueObserverInstrumentKind, profile.NumberKind)

		var last number.Number
		for i := 0; i < count; i++ {
			x := profile.Random(rand.Intn(1)*2 - 1)
			last = x
			aggregatortest.CheckedUpdate(t, agg, x, record)
		}

		err := agg.SynchronizedMove(ckpt, record)
		require.NoError(t, err)

		lv, _, err := ckpt.LastValue()
		require.Equal(t, last, lv, "Same last value - non-monotonic")
		require.Nil(t, err)
	})
}

func TestLastValueMerge(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg1, agg2, ckpt1, ckpt2 := new4()

		descriptor := aggregatortest.NewAggregatorTest(metric.ValueObserverInstrumentKind, profile.NumberKind)

		first1 := profile.Random(+1)
		first2 := profile.Random(+1)
		first1.AddNumber(profile.NumberKind, first2)

		aggregatortest.CheckedUpdate(t, agg1, first1, descriptor)
		// Ensure these should not have the same timestamp.
		time.Sleep(time.Nanosecond)
		aggregatortest.CheckedUpdate(t, agg2, first2, descriptor)

		require.NoError(t, agg1.SynchronizedMove(ckpt1, descriptor))
		require.NoError(t, agg2.SynchronizedMove(ckpt2, descriptor))

		checkZero(t, agg1)
		checkZero(t, agg2)

		_, t1, err := ckpt1.LastValue()
		require.Nil(t, err)
		_, t2, err := ckpt2.LastValue()
		require.Nil(t, err)
		require.True(t, t1.Before(t2))

		aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

		lv, ts, err := ckpt1.LastValue()
		require.Nil(t, err)
		require.Equal(t, t2, ts, "Merged timestamp - non-monotonic")
		require.Equal(t, first2, lv, "Merged value - non-monotonic")
	})
}

func TestLastValueNotSet(t *testing.T) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueObserverInstrumentKind, number.Int64Kind)

	g, ckpt := new2()
	require.NoError(t, g.SynchronizedMove(ckpt, descriptor))

	checkZero(t, g)
}

func TestSynchronizedMoveReset(t *testing.T) {
	aggregatortest.SynchronizedMoveResetTest(
		t,
		metric.ValueObserverInstrumentKind,
		func(desc *metric.Descriptor) export.Aggregator {
			return &New(1)[0]
		},
	)
}
