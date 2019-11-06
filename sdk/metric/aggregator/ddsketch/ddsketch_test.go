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

package ddsketch

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const count = 1000

type updateTest struct {
	absolute bool
}

func (ut *updateTest) run(t *testing.T, profile test.Profile) {
	ctx := context.Background()

	descriptor := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, !ut.absolute)
	agg := New(NewDefaultConfig(), descriptor)

	all := test.NewNumbers(profile.NumberKind)
	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		agg.Update(ctx, x, descriptor)

		if !ut.absolute {
			y := profile.Random(-1)
			all.Append(y)
			agg.Update(ctx, y, descriptor)
		}
	}

	agg.Checkpoint(ctx, descriptor)

	all.Sort()

	require.InDelta(t,
		all.Sum().CoerceToFloat64(profile.NumberKind),
		agg.Sum().CoerceToFloat64(profile.NumberKind),
		1,
		"Same sum - absolute")
	require.Equal(t, all.Count(), agg.Count(), "Same count - absolute")

	max, err := agg.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max - absolute")

	median, err := agg.Quantile(0.5)
	require.Nil(t, err)
	require.InDelta(t,
		all.Median().CoerceToFloat64(profile.NumberKind),
		median.CoerceToFloat64(profile.NumberKind),
		10,
		"Same median - absolute")
}

func TestDDSketchUpdate(t *testing.T) {
	// Test absolute and non-absolute
	for _, absolute := range []bool{false, true} {
		t.Run(fmt.Sprint("Absolute=", absolute), func(t *testing.T) {
			ut := updateTest{
				absolute: absolute,
			}
			// Test integer and floating point
			test.RunProfiles(t, ut.run)
		})
	}
}

type mergeTest struct {
	absolute bool
}

func (mt *mergeTest) run(t *testing.T, profile test.Profile) {
	ctx := context.Background()
	descriptor := test.NewAggregatorTest(export.MeasureKind, profile.NumberKind, !mt.absolute)

	agg1 := New(NewDefaultConfig(), descriptor)
	agg2 := New(NewDefaultConfig(), descriptor)

	all := test.NewNumbers(profile.NumberKind)
	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		agg1.Update(ctx, x, descriptor)

		if !mt.absolute {
			y := profile.Random(-1)
			all.Append(y)
			agg1.Update(ctx, y, descriptor)
		}
	}

	for i := 0; i < count; i++ {
		x := profile.Random(+1)
		all.Append(x)
		agg2.Update(ctx, x, descriptor)

		if !mt.absolute {
			y := profile.Random(-1)
			all.Append(y)
			agg2.Update(ctx, y, descriptor)
		}
	}

	agg1.Checkpoint(ctx, descriptor)
	agg2.Checkpoint(ctx, descriptor)

	agg1.Merge(agg2, descriptor)

	all.Sort()

	require.InDelta(t,
		all.Sum().CoerceToFloat64(profile.NumberKind),
		agg1.Sum().CoerceToFloat64(profile.NumberKind),
		1,
		"Same sum - absolute")
	require.Equal(t, all.Count(), agg1.Count(), "Same count - absolute")

	max, err := agg1.Max()
	require.Nil(t, err)
	require.Equal(t,
		all.Max(),
		max,
		"Same max - absolute")

	median, err := agg1.Quantile(0.5)
	require.Nil(t, err)
	require.InDelta(t,
		all.Median().CoerceToFloat64(profile.NumberKind),
		median.CoerceToFloat64(profile.NumberKind),
		10,
		"Same median - absolute")
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
