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

var (
	sampledSpanContext = trace.NewSpanContext(trace.SpanContextConfig{
		SpanID:     trace.SpanID{0o1},
		TraceID:    trace.TraceID{0o1},
		TraceFlags: trace.FlagsSampled,
	})
	notSampledSpanContext = trace.NewSpanContext(trace.SpanContextConfig{
		SpanID:  trace.SpanID{0o1},
		TraceID: trace.TraceID{0o1},
	})
)

var exemplarBenchmarks = []struct {
	Name        string
	SpanContext trace.SpanContext
}{
	{"ExemplarsDisabled", notSampledSpanContext},
	{"ExemplarsEnabled", sampledSpanContext},
}

func BenchmarkSyncMeasure(b *testing.B) {
	for _, bc := range viewBenchmarks {
		for _, eb := range exemplarBenchmarks {
			b.Run(fmt.Sprintf("%s/%s", bc.Name, eb.Name), benchSyncViews(eb.SpanContext, bc.Views...))
		}
	}
}

func exponentialAggregationSelector(ik InstrumentKind) Aggregation {
	if ik == InstrumentKindHistogram {
		return AggregationBase2ExponentialHistogram{MaxScale: 20, MaxSize: 160}
	}
	return AggregationDefault{}
}

func benchSyncViews(sc trace.SpanContext, views ...View) func(*testing.B) {
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr), WithView(views...))
	meter := provider.Meter("benchSyncViews")
	expRdr := NewManualReader(WithAggregationSelector(exponentialAggregationSelector))
	expProvider := NewMeterProvider(WithReader(expRdr), WithView(views...))
	expMeter := expProvider.Meter("benchSyncViews")
	// Precompute histogram values so they are distributed equally to buckets.
	histogramBuckets := DefaultAggregationSelector(InstrumentKindHistogram).(AggregationExplicitBucketHistogram).Boundaries
	histogramObservations := make([]float64, len(histogramBuckets))
	for i, bucket := range histogramBuckets {
		histogramObservations[i] = bucket + 1
	}
	return func(b *testing.B) {
		ctx := trace.ContextWithSpanContext(b.Context(), sc)
		iCtr, err := meter.Int64Counter("int64-counter")
		assert.NoError(b, err)
		b.Run("Int64Counter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func(int) { iCtr.Add(ctx, 1, o...) }
			}
		}()))

		fCtr, err := meter.Float64Counter("float64-counter")
		assert.NoError(b, err)
		b.Run("Float64Counter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func(int) { fCtr.Add(ctx, 1, o...) }
			}
		}()))

		iUDCtr, err := meter.Int64UpDownCounter("int64-up-down-counter")
		assert.NoError(b, err)
		b.Run("Int64UpDownCounter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func(int) { iUDCtr.Add(ctx, 1, o...) }
			}
		}()))

		fUDCtr, err := meter.Float64UpDownCounter("float64-up-down-counter")
		assert.NoError(b, err)
		b.Run("Float64UpDownCounter", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.AddOption{metric.WithAttributeSet(s)}
				return func(int) { fUDCtr.Add(ctx, 1, o...) }
			}
		}()))

		iGauge, err := meter.Int64Gauge("int64-gauge")
		assert.NoError(b, err)
		b.Run("Int64Gauge", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func(int) { iGauge.Record(ctx, 1, o...) }
			}
		}()))

		fGauge, err := meter.Float64Gauge("float64-gauge")
		assert.NoError(b, err)
		b.Run("Float64Gauge", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func(int) { fGauge.Record(ctx, 1, o...) }
			}
		}()))

		iHist, err := meter.Int64Histogram("int64-histogram")
		assert.NoError(b, err)
		b.Run("Int64Histogram", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func(i int) { iHist.Record(ctx, int64(histogramObservations[i%len(histogramObservations)]), o...) }
			}
		}()))

		fHist, err := meter.Float64Histogram("float64-histogram")
		assert.NoError(b, err)
		b.Run("Float64Histogram", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(i int) {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func(i int) { fHist.Record(ctx, histogramObservations[i%len(histogramObservations)], o...) }
			}
		}()))

		expIHist, err := expMeter.Int64Histogram("exponential-int64-histogram")
		assert.NoError(b, err)
		b.Run("ExponentialInt64Histogram", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func(int) { expIHist.Record(ctx, 1, o...) }
			}
		}()))

		expFHist, err := expMeter.Float64Histogram("exponential-float64-histogram")
		assert.NoError(b, err)
		b.Run("ExponentialFloat64Histogram", benchMeasAttrs(func() measF {
			return func(s attribute.Set) func(int) {
				o := []metric.RecordOption{metric.WithAttributeSet(s)}
				return func(int) { expFHist.Record(ctx, 1, o...) }
			}
		}()))
	}
}

type measF func(s attribute.Set) func(i int)

func benchMeasAttrs(meas measF) func(*testing.B) {
	return func(b *testing.B) {
		b.Run("Attributes/0", func(b *testing.B) {
			f := meas(*attribute.EmptySet())
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					f(i)
					i++
				}
			})
		})
		b.Run("Attributes/1", func(b *testing.B) {
			f := meas(attribute.NewSet(attribute.Bool("K", true)))
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					f(i)
					i++
				}
			})
		})
		b.Run("Attributes/10", func(b *testing.B) {
			n := 10
			attrs := make([]attribute.KeyValue, 0)
			attrs = append(attrs, attribute.Bool("K", true))
			for i := 2; i < n; i++ {
				attrs = append(attrs, attribute.Int(strconv.Itoa(i), i))
			}
			f := meas(attribute.NewSet(attrs...))
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
	return func(b *testing.B) {
		b.Run("Int64Counter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Counter")
			i, err := m.Int64Counter("int64-counter")
			assert.NoError(b, err)
			i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Int64Counter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Counter")
			i, err := m.Int64Counter("int64-counter")
			assert.NoError(b, err)
			for range 10 {
				i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Float64Counter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Counter")
			i, err := m.Float64Counter("float64-counter")
			assert.NoError(b, err)
			i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Float64Counter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Counter")
			i, err := m.Float64Counter("float64-counter")
			assert.NoError(b, err)
			for range 10 {
				i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Int64UpDownCounter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64UpDownCounter")
			i, err := m.Int64UpDownCounter("int64-up-down-counter")
			assert.NoError(b, err)
			i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Int64UpDownCounter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64UpDownCounter")
			i, err := m.Int64UpDownCounter("int64-up-down-counter")
			assert.NoError(b, err)
			for range 10 {
				i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Float64UpDownCounter/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64UpDownCounter")
			i, err := m.Float64UpDownCounter("float64-up-down-counter")
			assert.NoError(b, err)
			i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Float64UpDownCounter/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64UpDownCounter")
			i, err := m.Float64UpDownCounter("float64-up-down-counter")
			assert.NoError(b, err)
			for range 10 {
				i.Add(b.Context(), 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Int64Histogram/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Histogram")
			i, err := m.Int64Histogram("int64-histogram")
			assert.NoError(b, err)
			i.Record(b.Context(), 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Int64Histogram/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Int64Histogram")
			i, err := m.Int64Histogram("int64-histogram")
			assert.NoError(b, err)
			for range 10 {
				i.Record(b.Context(), 1, metric.WithAttributeSet(s))
			}
			return r
		}))

		b.Run("Float64Histogram/1", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Histogram")
			i, err := m.Float64Histogram("float64-histogram")
			assert.NoError(b, err)
			i.Record(b.Context(), 1, metric.WithAttributeSet(s))
			return r
		}))
		b.Run("Float64Histogram/10", benchCollectAttrs(func(s attribute.Set) Reader {
			m, r := setup("benchCollectViews/Float64Histogram")
			i, err := m.Float64Histogram("float64-histogram")
			assert.NoError(b, err)
			for range 10 {
				i.Record(b.Context(), 1, metric.WithAttributeSet(s))
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
	out := new(metricdata.ResourceMetrics)
	run := func(reader Reader) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_ = reader.Collect(b.Context(), out)
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
	ctx := trace.ContextWithSpanContext(b.Context(), sampledSpanContext)

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
