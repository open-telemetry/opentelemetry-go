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

package metric

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func BenchmarkInstrument(b *testing.B) {
	attr := func(id int) attribute.Set {
		return attribute.NewSet(
			attribute.String("user", "Alice"),
			attribute.Bool("admin", true),
			attribute.Int("id", id),
		)
	}

	b.Run("instrumentImpl/aggregate", func(b *testing.B) {
		build := aggregate.Builder[int64]{}
		var meas []aggregate.Measure[int64]

		in, _ := build.LastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.CumulativeTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		inst := int64Inst{measures: meas}
		ctx := context.Background()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inst.aggregate(ctx, int64(i), attr(i))
		}
	})

	b.Run("observable/observe", func(b *testing.B) {
		build := aggregate.Builder[int64]{}
		var meas []aggregate.Measure[int64]

		in, _ := build.LastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.CumulativeTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		o := observable[int64]{measures: meas}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			o.observe(int64(i), attr(i))
		}
	})
}
