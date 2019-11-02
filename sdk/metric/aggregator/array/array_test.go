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

package array

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/sdk/export"
	"go.opentelemetry.io/sdk/metric/aggregator/test"
)

func TestArrayAbsolute(t *testing.T) {
	// Test with an odd an even number of measurements
	for count := 999; count <= 1000; count++ {
		// Test absolute and non-absolute
		for _, absolute := range []bool{false, true} {
			t.Run(fmt.Sprint("Absolute=", absolute), func(t *testing.T) {
				// Test integer and floating point
				test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
					ctx := context.Background()

					batcher, record := test.NewAggregatorTest(export.MeasureMetricKind, profile.NumberKind, !absolute)

					agg := New()

					all := test.NewNumbers(profile.NumberKind)

					for i := 0; i < count; i++ {
						x := profile.Random(+1)
						all.Append(x)
						agg.Update(ctx, x, record)

						if !absolute {
							y := profile.Random(-1)
							all.Append(y)
							agg.Update(ctx, y, record)
						}
					}

					agg.Collect(ctx, record, batcher)

					all.Sort()

					require.InEpsilon(t,
						all.Sum().CoerceToFloat64(profile.NumberKind),
						agg.Sum().CoerceToFloat64(profile.NumberKind),
						0.0000001,
						"Same sum - absolute")
					require.Equal(t, all.Count(), agg.Count(), "Same count - absolute")

					min, err := agg.Min()
					require.Nil(t, err)
					require.Equal(t, all.Min(), min, "Same min - absolute")

					max, err := agg.Max()
					require.Nil(t, err)
					require.Equal(t, all.Max(), max, "Same max - absolute")

					qx, err := agg.Quantile(0.5)
					require.Nil(t, err)
					require.Equal(t, all.Median(), qx, "Same median - absolute")
				})
			})
		}
	}
}

// TODO: test empty, test small, test NaN and other stuff

// TODO: test merge
