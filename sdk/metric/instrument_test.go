// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type attributesProvider struct {
	set attribute.Set
}

func (p *attributesProvider) Attributes() attribute.Set {
	return p.set
}

func (p *attributesProvider) AttributesDistinct() attribute.Distinct {
	return p.set.Equivalent()
}

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

		build.Temporality = metricdata.CumulativeTemporality
		in, _ := build.LastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.LastValue()
		meas = append(meas, in)

		build.Temporality = metricdata.CumulativeTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		build.Temporality = metricdata.DeltaTemporality
		in, _ = build.Sum(true)
		meas = append(meas, in)

		inst := int64Inst{measures: meas}
		ctx := b.Context()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inst.aggregate(ctx, int64(i), &attributesProvider{set: attr(i)})
		}
	})

	b.Run("observable/observe", func(b *testing.B) {
		build := aggregate.Builder[int64]{}
		var meas []aggregate.Measure[int64]

		in, _ := build.PrecomputedLastValue()
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
			o.observe(int64(i), &attributesProvider{set: attr(i)})
		}
	})
}
