// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"strconv"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func BenchmarkFinishLastValueMeasure(b *testing.B) {
	b.Run("Int64", benchmarkFinishLastValueMeasure[int64])
	b.Run("Float64", benchmarkFinishLastValueMeasure[float64])
}

func benchmarkFinishLastValueMeasure[N int64 | float64](b *testing.B) {
	for _, enabled := range []bool{false, true} {
		factory := func() Measure[N] {
			if enabled {
				return Builder[N]{}.FinishLastValue().Measure
			}
			measure, _ := Builder[N]{}.LastValue()
			return measure
		}
		b.Run("finish="+strconv.FormatBool(enabled), func(b *testing.B) {
			b.Run("mode=serial", func(b *testing.B) {
				measure := factory()
				b.ReportAllocs()
				for b.Loop() {
					measure(b.Context(), N(1), alice)
				}
			})
			b.Run("mode=parallel", func(b *testing.B) {
				measure := factory()
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						measure(b.Context(), N(1), alice)
					}
				})
			})
		})
	}
}

func BenchmarkFinishLastValueLifecycle(b *testing.B) {
	for _, name := range []string{"MeasureAndFinish", "RepeatedFinish", "UnknownFinish"} {
		b.Run(name, func(b *testing.B) {
			agg := Builder[int64]{}.FinishLastValue()
			agg.Measure(b.Context(), 1, alice)
			distinct := alice.Equivalent()
			if name == "UnknownFinish" {
				distinct = bob.Equivalent()
			}
			b.ReportAllocs()
			for b.Loop() {
				if name == "MeasureAndFinish" {
					agg.Measure(b.Context(), 1, alice)
				}
				agg.Finish(distinct, y2k)
			}
		})
	}
}

func BenchmarkFinishLastValueCycle(b *testing.B) {
	for _, temporality := range []metricdata.Temporality{
		metricdata.CumulativeTemporality, metricdata.DeltaTemporality,
	} {
		b.Run(temporality.String(), func(b *testing.B) {
			agg := Builder[int64]{Temporality: temporality}.FinishLastValue()
			var data metricdata.Aggregation
			b.ReportAllocs()
			for b.Loop() {
				agg.Measure(b.Context(), 1, alice)
				agg.ComputeAggregation(&data)
			}
		})
	}
}
