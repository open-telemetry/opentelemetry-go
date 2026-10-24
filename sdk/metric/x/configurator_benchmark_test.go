// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
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
	histogramObservations := make([]float64, len(boundaries))
	for i, bound := range boundaries {
		histogramObservations[i] = bound + 1
	}

	// As in benchSyncViews, each provider and its instruments are created once
	// and shared by every sub-benchmark that uses them.
	type instruments struct {
		counter   metric.Int64Counter
		histogram metric.Int64Histogram
	}
	newInstruments := func(opts ...sdkmetric.Option) instruments {
		opts = append(opts, sdkmetric.WithReader(sdkmetric.NewManualReader()))
		m := sdkmetric.NewMeterProvider(opts...).Meter("benchSyncViews")
		c, err := m.Int64Counter("int64-counter")
		if err != nil {
			b.Fatal(err)
		}
		h, err := m.Int64Histogram("int64-histogram")
		if err != nil {
			b.Fatal(err)
		}
		return instruments{counter: c, histogram: h}
	}

	handle := x.NewMeterConfiguratorHandle()
	handle.Set(func(instrumentation.Scope) x.MeterConfig {
		return x.NewMeterConfig(x.WithMeterEnabled(true))
	})
	providers := []struct {
		name string
		inst instruments
	}{
		{"stable", newInstruments()},
		{"configurator", newInstruments(x.WithMeterConfigurator(handle))},
	}

	ctx := trace.ContextWithSpanContext(b.Context(), notSampled)

	// measure returns a function that makes the i-th measurement with set.
	type measure func(inst instruments, set attribute.Set) func(i int)
	kinds := []struct {
		name    string
		measure measure
	}{
		{"Int64Counter", func(inst instruments, set attribute.Set) func(int) {
			o := []metric.AddOption{metric.WithAttributeSet(set)}
			return func(int) { inst.counter.Add(ctx, 1, o...) }
		}},
		{"Int64Histogram", func(inst instruments, set attribute.Set) func(int) {
			o := []metric.RecordOption{metric.WithAttributeSet(set)}
			return func(i int) {
				inst.histogram.Record(ctx, int64(histogramObservations[i%len(histogramObservations)]), o...)
			}
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

	for _, k := range kinds {
		for _, as := range attrSets {
			for _, p := range providers {
				b.Run(k.name+"/Attributes/"+as.name+"/provider="+p.name, func(b *testing.B) {
					f := k.measure(p.inst, as.set)
					b.RunParallel(func(pb *testing.PB) {
						i := 0
						for pb.Next() {
							f(i)
							i++
						}
					})
				})
			}
		}
	}
}
