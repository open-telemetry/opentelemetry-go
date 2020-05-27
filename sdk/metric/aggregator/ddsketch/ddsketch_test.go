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
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 1000

type updateTest struct {
}

func (ut *updateTest) run(t *testing.T, profile test.Profile) {
	ctx := context.Background()

	descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)
	agg := New(NewDefaultConfig(), descriptor)

	all := test.NewNumbers(profile.NumberKind)
	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		test.CheckedUpdate(t, agg, x, descriptor)

		y := profile.Random(-1)
		all.Append(y)
		test.CheckedUpdate(t, agg, y, descriptor)
	}

	agg.Checkpoint(ctx, descriptor)

	all.Sort()

	sum, err := agg.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InDelta(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		sum.CoerceToFloat64(profile.NumberKind),
		1,
		"Same sum")

	count, err := agg.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	max, err := agg.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max")

	median, err := agg.Quantile(0.5)
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
	test.RunProfiles(t, ut.run)
}

type mergeTest struct {
	absolute bool
}

func (mt *mergeTest) run(t *testing.T, profile test.Profile) {
	ctx := context.Background()
	descriptor := test.NewAggregatorTest(metric.ValueRecorderKind, profile.NumberKind)

	agg1 := New(NewDefaultConfig(), descriptor)
	agg2 := New(NewDefaultConfig(), descriptor)

	all := test.NewNumbers(profile.NumberKind)
	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		test.CheckedUpdate(t, agg1, x, descriptor)

		if !mt.absolute {
			y := profile.Random(-1)
			all.Append(y)
			test.CheckedUpdate(t, agg1, y, descriptor)
		}
	}

	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		test.CheckedUpdate(t, agg2, x, descriptor)

		if !mt.absolute {
			y := profile.Random(-1)
			all.Append(y)
			test.CheckedUpdate(t, agg2, y, descriptor)
		}
	}

	agg1.Checkpoint(ctx, descriptor)
	agg2.Checkpoint(ctx, descriptor)

	test.CheckedMerge(t, agg1, agg2, descriptor)

	all.Sort()

	aggSum, err := agg1.Sum()
	require.Nil(t, err)
	allSum := all.Sum()
	require.InDelta(t,
		(&allSum).CoerceToFloat64(profile.NumberKind),
		aggSum.CoerceToFloat64(profile.NumberKind),
		1,
		"Same sum")

	count, err := agg1.Count()
	require.Equal(t, all.Count(), count, "Same count")
	require.Nil(t, err)

	max, err := agg1.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max")

	median, err := agg1.Quantile(0.5)
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
			test.RunProfiles(t, mt.run)
		})
	}
}
