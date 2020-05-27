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
	"context"
	"fmt"
	"math"
	"os"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	ottest "go.opentelemetry.io/otel/internal/testing"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []ottest.FieldOffset{
		{
			Name:   "Aggregator.ckptSum",
			Offset: unsafe.Offsetof(Aggregator{}.ckptSum),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

type updateTest struct {
	count int
}

func (ut *updateTest) run(t *testing.T, profile test.Profile) {
	descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

	agg := New()

	all := test.NewNumbers(profile.NumberKind)

	for i := 0; i < ut.count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		test.CheckedUpdate(t, agg, x, descriptor)

		y := profile.Random(-1)
		all.Append(y)
		test.CheckedUpdate(t, agg, y, descriptor)
	}

	ctx := context.Background()
	agg.Checkpoint(ctx, descriptor)

	all.Sort()

	sum, err := agg.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InEpsilon(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		0.0000001,
		"Same sum")
	count, err := agg.Count()
	require.Nil(t, err)
	require.Equal(t, all.Count(), count, "Same count")

	min, err := agg.Min()
	require.Nil(t, err)
	require.Equal(t, all.Min(), min, "Same min")

	max, err := agg.Max()
	require.Nil(t, err)
	require.Equal(t, all.Max(), max, "Same max")

	qx, err := agg.Quantile(0.5)
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
			test.RunProfiles(t, ut.run)
		})
	}
}

type mergeTest struct {
	count    int
	absolute bool
}

func (mt *mergeTest) run(t *testing.T, profile test.Profile) {
	ctx := context.Background()

	descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

	agg1 := New()
	agg2 := New()

	all := test.NewNumbers(profile.NumberKind)

	for i := 0; i < mt.count; i++ {
		x1 := profile.Random(+1)
		all.Append(x1)
		test.CheckedUpdate(t, agg1, x1, descriptor)

		x2 := profile.Random(+1)
		all.Append(x2)
		test.CheckedUpdate(t, agg2, x2, descriptor)

		if !mt.absolute {
			y1 := profile.Random(-1)
			all.Append(y1)
			test.CheckedUpdate(t, agg1, y1, descriptor)

			y2 := profile.Random(-1)
			all.Append(y2)
			test.CheckedUpdate(t, agg2, y2, descriptor)
		}
	}

	agg1.Checkpoint(ctx, descriptor)
	agg2.Checkpoint(ctx, descriptor)

	test.CheckedMerge(t, agg1, agg2, descriptor)

	all.Sort()

	sum, err := agg1.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InEpsilon(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		0.0000001,
		"Same sum - absolute")
	count, err := agg1.Count()
	require.Nil(t, err)
	require.Equal(t, all.Count(), count, "Same count - absolute")

	min, err := agg1.Min()
	require.Nil(t, err)
	require.Equal(t, all.Min(), min, "Same min - absolute")

	max, err := agg1.Max()
	require.Nil(t, err)
	require.Equal(t, all.Max(), max, "Same max - absolute")

	qx, err := agg1.Quantile(0.5)
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
					test.RunProfiles(t, mt.run)
				})
			}
		})
	}
}

func TestArrayErrors(t *testing.T) {
	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		agg := New()

		_, err := agg.Max()
		require.Error(t, err)
		require.Equal(t, err, aggregator.ErrNoData)

		_, err = agg.Min()
		require.Error(t, err)
		require.Equal(t, err, aggregator.ErrNoData)

		_, err = agg.Quantile(0.1)
		require.Error(t, err)
		require.Equal(t, err, aggregator.ErrNoData)

		ctx := context.Background()

		descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

		test.CheckedUpdate(t, agg, metric.Number(0), descriptor)

		if profile.NumberKind == metric.Float64NumberKind {
			test.CheckedUpdate(t, agg, metric.NewFloat64Number(math.NaN()), descriptor)
		}
		agg.Checkpoint(ctx, descriptor)

		count, err := agg.Count()
		require.Equal(t, int64(1), count, "NaN value was not counted")
		require.Nil(t, err)

		num, err := agg.Quantile(0)
		require.Nil(t, err)
		require.Equal(t, num, metric.Number(0))

		_, err = agg.Quantile(-0.0001)
		require.Error(t, err)
		require.Equal(t, err, aggregator.ErrInvalidQuantile)

		_, err = agg.Quantile(1.0001)
		require.Error(t, err)
		require.Equal(t, err, aggregator.ErrInvalidQuantile)
	})
}

func TestArrayFloat64(t *testing.T) {
	descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, metric.Float64NumberKind)

	fpsf := func(sign int) []float64 {
		// Check behavior of a bunch of odd floating
		// points except for NaN, which is invalid.
		return []float64{
			0,
			math.Inf(sign),
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

	all := test.NewNumbers(metric.Float64NumberKind)

	ctx := context.Background()
	agg := New()

	for _, f := range fpsf(1) {
		all.Append(metric.NewFloat64Number(f))
		test.CheckedUpdate(t, agg, metric.NewFloat64Number(f), descriptor)
	}

	for _, f := range fpsf(-1) {
		all.Append(metric.NewFloat64Number(f))
		test.CheckedUpdate(t, agg, metric.NewFloat64Number(f), descriptor)
	}

	agg.Checkpoint(ctx, descriptor)

	all.Sort()

	sum, err := agg.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InEpsilon(t, (&allSum).AsFloat64(), sum.AsFloat64(), 0.0000001, "Same sum")

	count, err := agg.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	min, err := agg.Min()
	require.Nil(t, err)
	require.Equal(t, all.Min(), min, "Same min")

	max, err := agg.Max()
	require.Nil(t, err)
	require.Equal(t, all.Max(), max, "Same max")

	qx, err := agg.Quantile(0.5)
	require.Nil(t, err)
	require.Equal(t, all.Median(), qx, "Same median")

	po, err := agg.Points()
	require.Nil(t, err)
	require.Equal(t, all.Len(), len(po), "Points() must have same length of updates")
	for i := 0; i < len(po); i++ {
		require.Equal(t, all.Points()[i], po[i], "Wrong point at position %d", i)
	}
}
