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
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
	"go.opentelemetry.io/sdk/metric/aggregator/test"
)

const count = 100

// N.B. DDSketch only supports absolute measures

func TestMeasureAbsolute(t *testing.T) {
	ctx := context.Background()

	test.RunProfiles(t, func(t *testing.T, profile test.Profile) {
		batcher, record := test.NewAggregatorTest(export.MeasureMetricKind, profile.NumberKind, false)

		agg := New(NewDefaultConfig(), record.Descriptor())

		var all test.Numbers
		for i := 0; i < count; i++ {
			x := profile.Random(+1)
			all = append(all, x)
			agg.Update(ctx, x, record)
		}

		agg.Collect(ctx, record, batcher)

		all.Sort()

		require.InEpsilon(t,
			asFloat(all.Sum(profile.NumberKind), profile.NumberKind),
			agg.Sum(),
			0.0000001,
			"Same sum - monotonic")
		require.Equal(t, all.Count(), agg.Count(), "Same sum - monotonic")
		require.Equal(t,
			asFloat(all[len(all)-1], profile.NumberKind),
			agg.Max(),
			"Same sum - monotonic")
		// Median
		require.InEpsilon(t,
			asFloat(all[len(all)/2], profile.NumberKind),
			agg.Quantile(0.5),
			0.1,
			"Same sum - monotonic")
	})
}

func asFloat(n core.Number, kind core.NumberKind) float64 {
	switch kind {
	case core.Float64NumberKind:
		return n.AsFloat64()
	case core.Int64NumberKind:
		return float64(n.AsInt64())
	default:
		panic("Not tested")
	}
}
