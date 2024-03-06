// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"
)

var viewBenchmarks = []struct {
	Name  string
	Views []View
}{
	{"NoView", []View{}},
	{
		"DropView",
		[]View{NewView(
			Instrument{Name: "*"},
			Stream{Aggregation: AggregationDrop{}},
		)},
	},
	{
		"AttrFilterView",
		[]View{NewView(
			Instrument{Name: "*"},
			Stream{AttributeFilter: attribute.NewAllowKeysFilter("K")},
		)},
	},
}

func BenchmarkSyncMeasure(b *testing.B) {
	for _, bc := range viewBenchmarks {
		b.Run(bc.Name, benchSyncViews(bc.Views...))
	}
}

func benchSyncViews(views ...View) func(*testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr), WithView(views...))
	meter := provider.Meter("benchSyncViews")
	return func(b *testing.B) {
		iCtr, err := meter.Int64Counter("int64-counter")
		assert.NoError(b, err)
		b.Run("Int64Counter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func() {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func() { iCtr.Add(ctx, 1, o...) }
			}
		}()))

		fCtr, err := meter.Float64Counter("float64-counter")
		assert.NoError(b, err)
		b.Run("Float64Counter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func() {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func() { fCtr.Add(ctx, 1, o...) }
			}
		}()))

		iUDCtr, err := meter.Int64UpDownCounter("int64-up-down-counter")
		assert.NoError(b, err)
		b.Run("Int64UpDownCounter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func() {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func() { iUDCtr.Add(ctx, 1, o...) }
			}
		}()))

		fUDCtr, err := meter.Float64UpDownCounter("float64-up-down-counter")
		assert.NoError(b, err)
		b.Run("Float64UpDownCounter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func() {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func() { fUDCtr.Add(ctx, 1, o...) }
			}
		}()))

		iHist, err := meter.Int64Histogram("int64-histogram")
		assert.NoError(b, err)
		b.Run("Int64Histogram", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func() {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func() { iHist.Record(ctx, 1, o...) }
			}
		}()))

		fHist, err := meter.Float64Histogram("float64-histogram")
		assert.NoError(b, err)
		b.Run("Float64Histogram", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func() {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func() { fHist.Record(ctx, 1, o...) }
			}
		}()))
	}
}

type measF func(s attribute.Set) func()

func benchMeasAttrs(meas measF) func(*testing.B) {
	return func(b *testing.B) {
		b.Run("Attributes/0", func(b *testing.B) {
			f := meas(*attribute.EmptySet())
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				f()
			}
		})
		b.Run("Attributes/1", func(b *testing.B) {
			f := meas(attribute.NewSet(attribute.Bool("K", true)))
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				f()
			}
		})
		b.Run("Attributes/10", func(b *testing.B) {
			n := 10
			attrs := make([]attribute.KeyValue, 0)
			attrs = append(attrs, attribute.Bool("K", true))
			for i := 2; i < n; i++ {
				attrs = append(attrs, attribute.Int(strconv.Itoa(i), i))
			}
			f := meas(attribute.NewSet(attrs...))
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				f()
			}
		})
	}
}

func BenchmarkCollect(b *testing.B) {
	for _, bc := range viewBenchmarks {
		b.Run(bc.Name, benchCollectViews(bc.Views...))
	}
}

func benchCollectViews(views ...View) func(*testing.B) {
	setup := func(name string) (metric.Meter, Reader) {
		r := NewManualReader()
		mp := NewMeterProvider(WithReader(r), WithView(views...))
		return mp.Meter(name), r
	}
	ctx := context.Background()
	return func(b *testing.B) {
		b.Run("Int64Counter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Counter")
			i, err := m.Int64Counter("int64-counter")
			assert.NoError(b, err)
			i.Add(ctx, 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Int64Counter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Counter")
			i, err := m.Int64Counter("int64-counter")
			assert.NoError(b, err)
			for n := 0; n < 10; n++ {
				i.Add(ctx, 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Float64Counter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Counter")
			i, err := m.Float64Counter("float64-counter")
			assert.NoError(b, err)
			i.Add(ctx, 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Float64Counter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Counter")
			i, err := m.Float64Counter("float64-counter")
			assert.NoError(b, err)
			for n := 0; n < 10; n++ {
				i.Add(ctx, 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Int64UpDownCounter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64UpDownCounter")
			i, err := m.Int64UpDownCounter("int64-up-down-counter")
			assert.NoError(b, err)
			i.Add(ctx, 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Int64UpDownCounter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64UpDownCounter")
			i, err := m.Int64UpDownCounter("int64-up-down-counter")
			assert.NoError(b, err)
			for n := 0; n < 10; n++ {
				i.Add(ctx, 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Float64UpDownCounter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64UpDownCounter")
			i, err := m.Float64UpDownCounter("float64-up-down-counter")
			assert.NoError(b, err)
			i.Add(ctx, 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Float64UpDownCounter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64UpDownCounter")
			i, err := m.Float64UpDownCounter("float64-up-down-counter")
			assert.NoError(b, err)
			for n := 0; n < 10; n++ {
				i.Add(ctx, 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Int64Histogram/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Histogram")
			i, err := m.Int64Histogram("int64-histogram")
			assert.NoError(b, err)
			i.Record(ctx, 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Int64Histogram/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Histogram")
			i, err := m.Int64Histogram("int64-histogram")
			assert.NoError(b, err)
			for n := 0; n < 10; n++ {
				i.Record(ctx, 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Float64Histogram/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Histogram")
			i, err := m.Float64Histogram("float64-histogram")
			assert.NoError(b, err)
			i.Record(ctx, 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Float64Histogram/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Histogram")
			i, err := m.Float64Histogram("float64-histogram")
			assert.NoError(b, err)
			for n := 0; n < 10; n++ {
				i.Record(ctx, 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Int64ObservableCounter", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64ObservableCounter")
			_, err := m.Int64ObservableCounter(
				"int64-observable-counter",
				metric.WithInt64Callback(int64Cback(s)),
			)
			assert.NoError(b, err)
			return r
		}))

		b.Run("Float64ObservableCounter", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64ObservableCounter")
			_, err := m.Float64ObservableCounter(
				"float64-observable-counter",
				metric.WithFloat64Callback(float64Cback(s)),
			)
			assert.NoError(b, err)
			return r
		}))

		b.Run("Int64ObservableUpDownCounter", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64ObservableUpDownCounter")
			_, err := m.Int64ObservableUpDownCounter(
				"int64-observable-up-down-counter",
				metric.WithInt64Callback(int64Cback(s)),
			)
			assert.NoError(b, err)
			return r
		}))

		b.Run("Float64ObservableUpDownCounter", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64ObservableUpDownCounter")
			_, err := m.Float64ObservableUpDownCounter(
				"float64-observable-up-down-counter",
				metric.WithFloat64Callback(float64Cback(s)),
			)
			assert.NoError(b, err)
			return r
		}))

		b.Run("Int64ObservableGauge", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64ObservableGauge")
			_, err := m.Int64ObservableGauge(
				"int64-observable-gauge",
				metric.WithInt64Callback(int64Cback(s)),
			)
			assert.NoError(b, err)
			return r
		}))

		b.Run("Float64ObservableGauge", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64ObservableGauge")
			_, err := m.Float64ObservableGauge(
				"float64-observable-gauge",
				metric.WithFloat64Callback(float64Cback(s)),
			)
			assert.NoError(b, err)
			return r
		}))
	}
}

func int64Cback(s attribute.Set) metric.Int64Callback {
	opt := []metric.ObserveOption{metric.WithAttributeSet(s)}
	return func(_ context.Context, o metric.Int64Observer) error {
		o.Observe(1, opt...)
		return nil
	}
}

func float64Cback(s attribute.Set) metric.Float64Callback {
	opt := []metric.ObserveOption{metric.WithAttributeSet(s)}
	return func(_ context.Context, o metric.Float64Observer) error {
		o.Observe(1, opt...)
		return nil
	}
}

func benchCollectAttrs(setup func(attribute.Set) Reader) func(*testing.B) {
	ctx := context.Background()
	out := new(metricdata.ResourceMetrics)
	run := func(reader Reader) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_ = reader.Collect(ctx, out)
			}
		}
	}
	return func(b *testing.B) {
		b.Run("Attributes/0", run(setup(*attribute.EmptySet())))

		attrs := []attribute.KeyValue{attribute.Bool("K", true)}
		b.Run("Attributes/1", run(setup(attribute.NewSet(attrs...))))

		for i := 2; i < 10; i++ {
			attrs = append(attrs, attribute.Int(strconv.Itoa(i), i))
		}
		b.Run("Attributes/10", run(setup(attribute.NewSet(attrs...))))
	}
}

func BenchmarkExemplars(b *testing.B) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		SpanID:     trace.SpanID{0o1},
		TraceID:    trace.TraceID{0o1},
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)

	attr := attribute.NewSet(
		attribute.String("user", "Alice"),
		attribute.Bool("admin", true),
	)

	setup := func(name string) (metric.Meter, Reader) {
		r := NewManualReader()
		v := NewView(Instrument{Name: "*"}, Stream{
			AttributeFilter: func(kv attribute.KeyValue) bool {
				return kv.Key == attribute.Key("user")
			},
		})
		mp := NewMeterProvider(WithReader(r), WithView(v))
		return mp.Meter(name), r
	}
	nCPU := runtime.NumCPU() // Size of the fixed reservoir used.

	b.Setenv("OTEL_GO_X_EXEMPLAR", "true")

	name := fmt.Sprintf("Int64Counter/%d", nCPU)
	b.Run(name, func(b *testing.B) {
		m, r := setup("Int64Counter")
		i, err := m.Int64Counter("int64-counter")
		assert.NoError(b, err)

		rm := newRM(metricdata.Sum[int64]{
			DataPoints: []metricdata.DataPoint[int64]{
				{Exemplars: make([]metricdata.Exemplar[int64], 0, nCPU)},
			},
		})
		e := &(rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64]).DataPoints[0].Exemplars)

		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for j := 0; j < 2*nCPU; j++ {
				i.Add(ctx, 1, metric.WithAttributeSet(attr))
			}

			_ = r.Collect(ctx, rm)
			assert.Len(b, *e, nCPU)
		}
	})

	name = fmt.Sprintf("Int64Histogram/%d", nCPU)
	b.Run(name, func(b *testing.B) {
		m, r := setup("Int64Counter")
		i, err := m.Int64Histogram("int64-histogram")
		assert.NoError(b, err)

		rm := newRM(metricdata.Histogram[int64]{
			DataPoints: []metricdata.HistogramDataPoint[int64]{
				{Exemplars: make([]metricdata.Exemplar[int64], 0, 1)},
			},
		})
		e := &(rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[int64]).DataPoints[0].Exemplars)

		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for j := 0; j < 2*nCPU; j++ {
				i.Record(ctx, 1, metric.WithAttributeSet(attr))
			}

			_ = r.Collect(ctx, rm)
			assert.Len(b, *e, 1)
		}
	})
}

func newRM(a metricdata.Aggregation) *metricdata.ResourceMetrics {
	return &metricdata.ResourceMetrics{
		ScopeMetrics: []metricdata.ScopeMetrics{
			{Metrics: []metricdata.Metrics{{Data: a}}},
		},
	}
}
