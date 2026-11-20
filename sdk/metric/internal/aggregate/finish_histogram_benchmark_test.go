// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"strconv"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func BenchmarkFinishHistogramMeasure(b *testing.B) {
	b.Run("Int64", benchmarkFinishHistogramMeasure[int64])
	b.Run("Float64", benchmarkFinishHistogramMeasure[float64])
}

func benchmarkFinishHistogramMeasure[N int64 | float64](b *testing.B) {
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		b.Run(temporality.String(), func(b *testing.B) {
			for _, enabled := range []bool{false, true} {
				factory := func() Measure[N] {
					builder := Builder[N]{Temporality: temporality}
					if enabled {
						return builder.FinishExplicitBucketHistogram(
							[]float64{0, 5, 10, 25, 50, 75, 100},
							false,
							false,
						).Measure
					}
					measure, _ := builder.ExplicitBucketHistogram([]float64{0, 5, 10, 25, 50, 75, 100}, false, false)
					return measure
				}
				b.Run("finish="+strconv.FormatBool(enabled), func(b *testing.B) {
					b.Run("mode=serial", func(b *testing.B) {
						measure := factory()
						measure(b.Context(), N(7), alice)
						b.ReportAllocs()
						for b.Loop() {
							measure(b.Context(), N(7), alice)
						}
					})
					b.Run("mode=parallel", func(b *testing.B) {
						measure := factory()
						measure(b.Context(), N(7), alice)
						b.ReportAllocs()
						b.ResetTimer()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								measure(b.Context(), N(7), alice)
							}
						})
					})
				})
			}
		})
	}
}

func BenchmarkFinishHistogramLifecycle(b *testing.B) {
	b.Run("Int64", benchmarkFinishHistogramLifecycle[int64])
	b.Run("Float64", benchmarkFinishHistogramLifecycle[float64])
}

func benchmarkFinishHistogramLifecycle[N int64 | float64](b *testing.B) {
	for _, name := range []string{"MeasureAndFinish", "RepeatedFinish", "UnknownFinish"} {
		b.Run(name, func(b *testing.B) {
			agg := Builder[N]{}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
			agg.Measure(b.Context(), N(7), alice)
			finishTime := now()
			distinct := alice.Equivalent()
			if name == "UnknownFinish" {
				distinct = bob.Equivalent()
			}
			b.ReportAllocs()
			for b.Loop() {
				if name == "MeasureAndFinish" {
					agg.Measure(b.Context(), N(7), alice)
				}
				agg.Finish(distinct, finishTime)
			}
		})
	}
}

func BenchmarkFinishHistogramCycle(b *testing.B) {
	b.Run("Int64", benchmarkFinishHistogramCycle[int64])
	b.Run("Float64", benchmarkFinishHistogramCycle[float64])
}

func benchmarkFinishHistogramCycle[N int64 | float64](b *testing.B) {
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		b.Run(temporality.String(), func(b *testing.B) {
			agg := Builder[N]{
				Temporality: temporality,
			}.FinishExplicitBucketHistogram(
				finishHistogramBounds,
				false,
				false,
			)
			var data metricdata.Aggregation
			b.ReportAllocs()
			for b.Loop() {
				agg.Measure(b.Context(), N(7), alice)
				agg.Finish(alice.Equivalent(), now())
				agg.ComputeAggregation(&data)
			}
		})
	}
}
