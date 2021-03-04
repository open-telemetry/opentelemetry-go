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

package exact

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
)

type updateTest struct {
	count int
}

func requireNotAfter(t *testing.T, t1, t2 time.Time) {
	require.False(t, t1.After(t2), "expected %v ≤ %v", t1, t2)
}

func checkZero(t *testing.T, agg *Aggregator, desc *metric.Descriptor) {
	count, err := agg.Count()
	require.NoError(t, err)
	require.Equal(t, uint64(0), count)

	pts, err := agg.Points()
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))
}

func new2() (_, _ *Aggregator) {
	alloc := New(2)
	return &alloc[0], &alloc[1]
}

func new4() (_, _, _, _ *Aggregator) {
	alloc := New(4)
	return &alloc[0], &alloc[1], &alloc[2], &alloc[3]
}

func sumOf(samples []aggregation.Point, k number.Kind) number.Number {
	var n number.Number
	for _, s := range samples {
		n.AddNumber(k, s.Number)
	}
	return n
}

func (ut *updateTest) run(t *testing.T, profile aggregatortest.Profile) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)
	agg, ckpt := new2()

	all := aggregatortest.NewNumbers(profile.NumberKind)

	for i := 0; i < ut.count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		advance()
		aggregatortest.CheckedUpdate(t, agg, x, descriptor)

		y := profile.Random(-1)
		all.Append(y)
		advance()
		aggregatortest.CheckedUpdate(t, agg, y, descriptor)
	}

	err := agg.SynchronizedMove(ckpt, descriptor)
	require.NoError(t, err)

	checkZero(t, agg, descriptor)

	all.Sort()

	pts, err := ckpt.Points()
	require.Nil(t, err)
	sum := sumOf(pts, profile.NumberKind)
	allSum := all.Sum()
	require.InEpsilon(t,
		allSum.CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		0.0000001,
		"Same sum")
	count, err := ckpt.Count()
	require.Nil(t, err)
	require.Equal(t, all.Count(), count, "Same count")
}

func TestExactUpdate(t *testing.T) {
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

func advance() {
	time.Sleep(time.Nanosecond)
}

func (mt *mergeTest) run(t *testing.T, profile aggregatortest.Profile) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)
	agg1, agg2, ckpt1, ckpt2 := new4()

	all := aggregatortest.NewNumbers(profile.NumberKind)

	for i := 0; i < mt.count; i++ {
		x1 := profile.Random(+1)
		all.Append(x1)
		advance()
		aggregatortest.CheckedUpdate(t, agg1, x1, descriptor)

		x2 := profile.Random(+1)
		all.Append(x2)
		advance()
		aggregatortest.CheckedUpdate(t, agg2, x2, descriptor)

		if !mt.absolute {
			y1 := profile.Random(-1)
			all.Append(y1)
			advance()
			aggregatortest.CheckedUpdate(t, agg1, y1, descriptor)

			y2 := profile.Random(-1)
			all.Append(y2)
			advance()
			aggregatortest.CheckedUpdate(t, agg2, y2, descriptor)
		}
	}

	require.NoError(t, agg1.SynchronizedMove(ckpt1, descriptor))
	require.NoError(t, agg2.SynchronizedMove(ckpt2, descriptor))

	checkZero(t, agg1, descriptor)
	checkZero(t, agg2, descriptor)

	aggregatortest.CheckedMerge(t, ckpt1, ckpt2, descriptor)

	pts, err := ckpt1.Points()
	require.Nil(t, err)

	received := aggregatortest.NewNumbers(profile.NumberKind)
	for i, s := range pts {
		received.Append(s.Number)

		if i > 0 {
			requireNotAfter(t, pts[i-1].Time, pts[i].Time)
		}
	}

	allSum := all.Sum()
	sum := sumOf(pts, profile.NumberKind)
	require.InEpsilon(t,
		allSum.CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		0.0000001,
		"Same sum - absolute")
	count, err := ckpt1.Count()
	require.Nil(t, err)
	require.Equal(t, all.Count(), count, "Same count - absolute")
	require.Equal(t, all, received, "Same ordered contents")
}

func TestExactMerge(t *testing.T) {
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

func TestExactErrors(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		agg, ckpt := new2()

		descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)

		advance()
		aggregatortest.CheckedUpdate(t, agg, number.Number(0), descriptor)

		if profile.NumberKind == number.Float64Kind {
			advance()
			aggregatortest.CheckedUpdate(t, agg, number.NewFloat64Number(math.NaN()), descriptor)
		}
		require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

		count, err := ckpt.Count()
		require.Equal(t, uint64(1), count, "NaN value was not counted")
		require.Nil(t, err)
	})
}

func TestExactFloat64(t *testing.T) {
	descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, number.Float64Kind)

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

	all := aggregatortest.NewNumbers(number.Float64Kind)

	agg, ckpt := new2()

	startTime := time.Now()

	for _, f := range fpsf(1) {
		all.Append(number.NewFloat64Number(f))
		advance()
		aggregatortest.CheckedUpdate(t, agg, number.NewFloat64Number(f), descriptor)
	}

	for _, f := range fpsf(-1) {
		all.Append(number.NewFloat64Number(f))
		advance()
		aggregatortest.CheckedUpdate(t, agg, number.NewFloat64Number(f), descriptor)
	}

	endTime := time.Now()

	require.NoError(t, agg.SynchronizedMove(ckpt, descriptor))

	pts, err := ckpt.Points()
	require.Nil(t, err)

	allSum := all.Sum()
	sum := sumOf(pts, number.Float64Kind)
	require.InEpsilon(t, allSum.AsFloat64(), sum.AsFloat64(), 0.0000001, "Same sum")

	count, err := ckpt.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	po, err := ckpt.Points()
	require.Nil(t, err)
	require.Equal(t, all.Len(), len(po), "Points() must have same length of updates")
	for i := 0; i < len(po); i++ {
		require.Equal(t, all.Points()[i], po[i].Number, "Wrong point at position %d", i)
		if i > 0 {
			requireNotAfter(t, po[i-1].Time, po[i].Time)
		}
	}
	requireNotAfter(t, startTime, po[0].Time)
	requireNotAfter(t, po[len(po)-1].Time, endTime)
}

func TestSynchronizedMoveReset(t *testing.T) {
	aggregatortest.SynchronizedMoveResetTest(
		t,
		metric.ValueRecorderInstrumentKind,
		func(desc *metric.Descriptor) export.Aggregator {
			return &New(1)[0]
		},
	)
}

func TestMergeBehavior(t *testing.T) {
	aggregatortest.RunProfiles(t, func(t *testing.T, profile aggregatortest.Profile) {
		for _, forward := range []bool{false, true} {
			t.Run(fmt.Sprint("Forward=", forward), func(t *testing.T) {
				descriptor := aggregatortest.NewAggregatorTest(metric.ValueRecorderInstrumentKind, profile.NumberKind)
				agg1, agg2, ckpt, _ := new4()

				all := aggregatortest.NewNumbers(profile.NumberKind)

				for i := 0; i < 100; i++ {
					x1 := profile.Random(+1)
					all.Append(x1)
					advance()
					aggregatortest.CheckedUpdate(t, agg1, x1, descriptor)
				}

				for i := 0; i < 100; i++ {
					x2 := profile.Random(+1)
					all.Append(x2)
					advance()
					aggregatortest.CheckedUpdate(t, agg2, x2, descriptor)
				}

				if forward {
					aggregatortest.CheckedMerge(t, ckpt, agg1, descriptor)
					aggregatortest.CheckedMerge(t, ckpt, agg2, descriptor)
				} else {
					aggregatortest.CheckedMerge(t, ckpt, agg2, descriptor)
					aggregatortest.CheckedMerge(t, ckpt, agg1, descriptor)
				}

				pts, err := ckpt.Points()
				require.NoError(t, err)

				received := aggregatortest.NewNumbers(profile.NumberKind)
				for i, s := range pts {
					received.Append(s.Number)

					if i > 0 {
						requireNotAfter(t, pts[i-1].Time, pts[i].Time)
					}
				}

				allSum := all.Sum()
				sum := sumOf(pts, profile.NumberKind)
				require.InEpsilon(t,
					allSum.CoerceToFloat64(profile.NumberKind),
					sum.CoerceToFloat64(profile.NumberKind),
					0.0000001,
					"Same sum - absolute")
				count, err := ckpt.Count()
				require.NoError(t, err)
				require.Equal(t, all.Count(), count, "Same count - absolute")
				require.Equal(t, all, received, "Same ordered contents")
			})
		}
	})
}
