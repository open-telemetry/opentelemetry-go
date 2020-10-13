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

package array

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
)

type updateTest struct {
	count int
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

	min, err := agg.Min()
	require.True(t, errors.Is(err, aggregation.ErrNoData))
	require.Equal(t, kind.Zero(), min)
}

func new2() (_, _ *Aggregator) {
	alloc := New(2)
	return &alloc[0], &alloc[1]
}

func new4() (_, _, _, _ *Aggregator) {
	alloc := New(4)
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func (ut *updateTest) run(t *testing.T, profile aggregatortest.Profile) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)
	agg, ckpt := new2()

	all := aggregatortest.NewNumbers(profile.NumberKind)

	for i := 0; i < ut.count; i++ {
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
	require.InEpsilon(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		0.0000001,
		"Same sum")
	count, err := ckpt.Count()
	require.Nil(t, err)
	require.Equal(t, all.Count(), count, "Same count")

	min, err := ckpt.Min()
	require.Nil(t, err)
	require.Equal(t, all.Min(), min, "Same min")

	max, err := ckpt.Max()
	require.Nil(t, err)
	require.Equal(t, all.Max(), max, "Same max")

	qx, err := ckpt.Quantile(0.5)
	require.Nil(t, err)
	require.Equal(t, all.Median(), qx, "Same median")
}

func TestArrayUpdate(t *testing.T) {
	// Test with an odd an even number of measurements
	for count := 999; count <= 1000; count++ {
		t.Run(fmt.Sprint("Odd=", count%2 == 1), func(t *testing.T) {
			ut := updateTest{
				count: count,
			}

			// Test integer and floating point
			aggregatortest.RunProfiles(t, ut.run)
		})
	}
}

type mergeTest struct {
	count    int
	absolute bool
}

func (mt *mergeTest) run(t *testing.T, profile aggregatortest.Profile) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)
	agg1, agg2, ckpt1, ckpt2 := new4()

	all := aggregatortest.NewNumbers(profile.NumberKind)

	for i := 0; i < mt.count; i++ {
		x1 := profile.Random(+1)
		all.Append(x1)
		aggregatortest.CheckedUpdate(t, agg1, x1, descriptor)

		x2 := profile.Random(+1)
		all.Append(x2)
		aggregatortest.CheckedUpdate(t, agg2, x2, descriptor)

		if !mt.absolute {
			y1 := profile.Random(-1)
			all.Append(y1)
			aggregatortest.CheckedUpdate(t, agg1, y1, descriptor)

			y2 := profile.Random(-1)
			all.Append(y2)
			aggregatortest.CheckedUpdate(t, agg2, y2, descriptor)
		}
	}

	require.NoError(t, agg1.SynchronizedMove(ckpt1, descriptor))
	require.NoError(t, agg2.SynchronizedMove(ckpt2, descriptor))

	checkZero(t, agg1, descriptor)
	checkZero(t, agg2, descriptor)

	aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

	all.Sort()

	sum, err := ckpt1.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InEpsilon(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		0.0000001,
		"Same sum - absolute")
	count, err := ckpt1.Count()
	require.Nil(t, err)
	require.Equal(t, all.Count(), count, "Same count - absolute")

	min, err := ckpt1.Min()
	require.Nil(t, err)
	require.Equal(t, all.Min(), min, "Same min - absolute")

	max, err := ckpt1.Max()
	require.Nil(t, err)
	require.Equal(t, all.Max(), max, "Same max - absolute")

	qx, err := ckpt1.Quantile(0.5)
	require.Nil(t, err)
	require.Equal(t, all.Median(), qx, "Same median - absolute")
}

func TestArrayMerge(t *testing.T) {
	// Test with an odd an even number of measurements
	for count := 999; count <= 1000; count++ {
		t.Run(fmt.Sprint("Odd=", count%2 == 1), func(t *testing.T) {
			// Test absolute and non-absolute
			for _, absolute := range []bool{false, true} {
				t.Run(fmt.Sprint("Absolute=", absolute), func(t *testing.T) {
					mt := mergeTest{
						count:    count,
						absolute: absolute,
					}

					// Test integer and floating point
					aggregatortest.RunProfiles(t, mt.run)
				})
			}
		})
	}
}

func TestArrayErrors(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg, ckpt := new2()

		_, err := ckpt.Max()
		require.Error(t, err)
		require.Equal(t, err, aggregation.ErrNoData)

		_, err = ckpt.Min()
		require.Error(t, err)
		require.Equal(t, err, aggregation.ErrNoData)

		_, err = ckpt.Quantile(0.1)
		require.Error(t, err)
		require.Equal(t, err, aggregation.ErrNoData)

		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		aggregatortest.CheckedUpdate(t, agg, metric.Number(0), descriptor)

		if profile.NumberKind == metric.Float64NumberKind {
			aggregatortest.CheckedUpdate(t, agg, metric.NewFloat64Number(math.NaN()), descriptor)
		}
		require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

		count, err := ckpt.Count()
		require.Equal(t, int64(1), count, "NaN value was not counted")
		require.Nil(t, err)

		num, err := ckpt.Quantile(0)
		require.Nil(t, err)
		require.Equal(t, num, metric.Number(0))

		_, err = ckpt.Quantile(-0.0001)
		require.Error(t, err)
		require.True(t, errors.Is(err, aggregation.ErrInvalidQuantile))

		_, err = agg.Quantile(1.0001)
		require.Error(t, err)
		require.True(t, errors.Is(err, aggregation.ErrNoData))
	})
}

func TestArrayFloat64(t *testing.T) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, metric.Float64NumberKind)

	fpsf := func(sign int) []float64 {
		// Check behavior of a bunch of odd floating
		// points except for NaN, which is invalid.
		return []float64{
			0,
			1 / math.Inf(sign),
			1,
			2,
			1e100,
			math.MaxFloat64,
			math.SmallestNonzeroFloat64,
			math.MaxFloat32,
			math.SmallestNonzeroFloat32,
			math.E,
			math.Pi,
			math.Phi,
			math.Sqrt2,
			math.SqrtE,
			math.SqrtPi,
			math.SqrtPhi,
			math.Ln2,
			math.Log2E,
			math.Ln10,
			math.Log10E,
		}
	}

	all := aggregatortest.NewNumbers(metric.Float64NumberKind)

	agg, ckpt := new2()

	for _, f := range fpsf(1) {
		all.Append(metric.NewFloat64Number(f))
		aggregatortest.CheckedUpdate(t, agg, metric.NewFloat64Number(f), descriptor)
	}

	for _, f := range fpsf(-1) {
		all.Append(metric.NewFloat64Number(f))
		aggregatortest.CheckedUpdate(t, agg, metric.NewFloat64Number(f), descriptor)
	}

	require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

	all.Sort()

	sum, err := ckpt.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InEpsilon(t, (&allSum).AsFloat64(), sum.AsFloat64(), 0.0000001, "Same sum")

	count, err := ckpt.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	min, err := ckpt.Min()
	require.Nil(t, err)
	require.Equal(t, all.Min(), min, "Same min")

	max, err := ckpt.Max()
	require.Nil(t, err)
	require.Equal(t, all.Max(), max, "Same max")

	qx, err := ckpt.Quantile(0.5)
	require.Nil(t, err)
	require.Equal(t, all.Median(), qx, "Same median")

	po, err := ckpt.Points()
	require.Nil(t, err)
	require.Equal(t, all.Len(), len(po), "Points() must have same length of updates")
	for i := 0; i < len(po); i++ {
		require.Equal(t, all.Points()[i], po[i], "Wrong point at position %d", i)
	}
}
