// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"strconv"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/x"
	"go.opentelemetry.io/otel/trace"
)

// BenchmarkConfiguratorSyncMeasure compares synchronous recording on a
// MeterProvider without a MeterConfigurator (stable) against one with a
// configurator that enables every scope (configurator). The option is the only
// difference, so the gap is what configured providers pay per measurement.
// Settings mirror BenchmarkSyncMeasure's NoView/ExemplarsDisabled case.
func BenchmarkConfiguratorSyncMeasure(b *testing.B) {
	// An unsampled span context with the default trace-based exemplar
	// filter, as in BenchmarkSyncMeasure's ExemplarsDisabled case.
	notSampled := trace.NewSpanContext(trace.SpanContextConfig{
		SpanID:  trace.SpanID{0o1},
		TraceID: trace.TraceID{0o1},
	})

	boundaries := sdkmetric.DefaultAggregationSelector(sdkmetric.InstrumentKindHistogram).(sdkmetric.AggregationExplicitBucketHistogram).Boundaries
	histogramObservations := make([]int64, len(boundaries))
	for i, bound := range boundaries {
		histogramObservations[i] = int64(bound + 1)
	}

	// record returns a function that makes the i-th measurement with set.
	type record func(set attribute.Set) func(ctx context.Context, i int)
	instruments := []struct {
		name  string
		build func(metric.Meter) (record, error)
	}{
		{"Int64Counter", func(m metric.Meter) (record, error) {
			c, err := m.Int64Counter("int64-counter")
			return func(set attribute.Set) func(context.Context, int) {
				o := []metric.AddOption{metric.WithAttributeSet(set)}
				return func(ctx context.Context, _ int) { c.Add(ctx, 1, o...) }
			}, err
		}},
		{"Int64Histogram", func(m metric.Meter) (record, error) {
			h, err := m.Int64Histogram("int64-histogram")
			return func(set attribute.Set) func(context.Context, int) {
				o := []metric.RecordOption{metric.WithAttributeSet(set)}
				return func(ctx context.Context, i int) {
					h.Record(ctx, histogramObservations[i%len(histogramObservations)], o...)
				}
			}, err
		}},
	}

	attrs := []attribute.KeyValue{attribute.Bool("K", true)}
	for i := 2; i < 10; i++ {
		attrs = append(attrs, attribute.Int(strconv.Itoa(i), i))
	}
	attrSets := []struct {
		name string
		set  attribute.Set
	}{
		{"0", *attribute.EmptySet()},
		{"10", attribute.NewSet(attrs...)},
	}

	providers := []struct {
		name string
		opts func() []sdkmetric.Option
	}{
		{"stable", func() []sdkmetric.Option { return nil }},
		{"configurator", func() []sdkmetric.Option {
			h := x.NewMeterConfiguratorHandle()
			h.Set(func(instrumentation.Scope) x.MeterConfig {
				return x.NewMeterConfig(x.WithMeterEnabled(true))
			})
			return []sdkmetric.Option{x.WithMeterConfigurator(h)}
		}},
	}

	for _, inst := range instruments {
		for _, as := range attrSets {
			for _, p := range providers {
				b.Run(inst.name+"/Attributes/"+as.name+"/provider="+p.name, func(b *testing.B) {
					opts := append(p.opts(), sdkmetric.WithReader(sdkmetric.NewManualReader()))
					rec, err := inst.build(sdkmetric.NewMeterProvider(opts...).Meter("bench"))
					if err != nil {
						b.Fatal(err)
					}
					f := rec(as.set)
					ctx := trace.ContextWithSpanContext(b.Context(), notSampled)
					b.ReportAllocs()
					b.ResetTimer()
					b.RunParallel(func(pb *testing.PB) {
						i := 0
						for pb.Next() {
							f(ctx, i)
							i++
						}
					})
				})
			}
		}
	}
}
