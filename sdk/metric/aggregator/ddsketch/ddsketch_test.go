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

package ddsketch

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
)

const count = 1000

type updateTest struct {
}

func new2(desc *metric.Descriptor) (_, _ *Aggregator) {
	alloc := New(2, desc, NewDefaultConfig())
	return &alloc[0], &alloc[1]
}

func new4(desc *metric.Descriptor) (_, _, _, _ *Aggregator) {
	alloc := New(4, desc, NewDefaultConfig())
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func checkZero(t *testing.T, agg *Aggregator, desc *metric.Descriptor) {
	kind := desc.NumberKind()

	sum, err := agg.Sum()
	require.NoError(t, err)
	require.Equal(t, kind.Zero(), sum)

	count, err := agg.Count()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)

	max, err := agg.Max()
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, kind.Zero(), max)

	median, err := agg.Quantile(0.5)
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, kind.Zero(), median)

	min, err := agg.Min()
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, kind.Zero(), min)
}

func (ut *updateTest) run(t *testing.T, profile aggregatortest.Profile) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)
	agg, ckpt := new2(descriptor)

	all := aggregatortest.NewNumbers(profile.NumberKind)
	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		aggregatortest.CheckedUpdate(t, agg, x, descriptor)

		y := profile.Random(-1)
		all.Append(y)
		aggregatortest.CheckedUpdate(t, agg, y, descriptor)
	}

	err := agg.SynchronizedMove(ckpt, descriptor)
	require.NoError(t, err)

	checkZero(t, agg, descriptor)

	all.Sort()

	sum, err := ckpt.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InDelta(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		1,
		"Same sum")

	count, err := ckpt.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	max, err := ckpt.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max")

	median, err := ckpt.Quantile(0.5)
	require.Nil(t, err)
	allMedian := all.Median()
	require.InDelta(t,
		(&allMedian).CoerceToFloat64(profile.NumberKind),
		median.CoerceToFloat64(profile.NumberKind),
		10,
		"Same median")
}

func TestDDSketchUpdate(t *testing.T) {
	ut := updateTest{}
	aggregatortest.RunProfiles(t, ut.run)
}

type mergeTest struct {
	absolute bool
}

func (mt *mergeTest) run(t *testing.T, profile aggregatortest.Profile) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

	agg1, agg2, ckpt1, ckpt2 := new4(descriptor)

	all := aggregatortest.NewNumbers(profile.NumberKind)
	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		aggregatortest.CheckedUpdate(t, agg1, x, descriptor)

		if !mt.absolute {
			y := profile.Random(-1)
			all.Append(y)
			aggregatortest.CheckedUpdate(t, agg1, y, descriptor)
		}
	}

	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		aggregatortest.CheckedUpdate(t, agg2, x, descriptor)

		if !mt.absolute {
			y := profile.Random(-1)
			all.Append(y)
			aggregatortest.CheckedUpdate(t, agg2, y, descriptor)
		}
	}

	require.NoError(t, agg1.SynchronizedMove(ckpt1, descriptor))
	require.NoError(t, agg2.SynchronizedMove(ckpt2, descriptor))

	checkZero(t, agg1, descriptor)
	checkZero(t, agg1, descriptor)

	aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

	all.Sort()

	aggSum, err := ckpt1.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InDelta(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		aggSum.CoerceToFloat64(profile.NumberKind),
		1,
		"Same sum")

	count, err := ckpt1.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	max, err := ckpt1.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max")

	median, err := ckpt1.Quantile(0.5)
	require.Nil(t, err)
	allMedian := all.Median()
	require.InDelta(t,
		(&allMedian).CoerceToFloat64(profile.NumberKind),
		median.CoerceToFloat64(profile.NumberKind),
		10,
		"Same median")
}

func TestDDSketchMerge(t *testing.T) {
	// Test absolute and non-absolute
	for _, absolute := range []bool{false, true} {
		t.Run(fmt.Sprint("Absolute=", absolute), func(t *testing.T) {
			mt := mergeTest{
				absolute: absolute,
			}
			// Test integer and floating point
			aggregatortest.RunProfiles(t, mt.run)
		})
	}
}
