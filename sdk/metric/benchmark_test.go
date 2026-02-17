// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
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

// BenchmarkEndToEndCounterAdd measures the performance of adding to a counter,
// but takes into account the costs of constructing options to pass attributes
// to the API in different user scenarios:
//   - In the "Precomputed" case, attributes are known ahead of time, and
//     options are not computed for each call.
//   - In the "Dynamic" case, attributes are not known ahead of time, and
//     options are computed for each counter increment. The "Dynamic" case
//     applies performance optimizations that are part of our contributor
//     guidelines.
//   - In the "Naive" case, the user uses the API and SDK in the simplest and
//     most obvious way without applying any performance optimizations.
func BenchmarkEndToEndCounterAdd(b *testing.B) {
	testCounter := func(b *testing.B, mp metric.MeterProvider) metric.Float64Counter {
		meter := mp.Meter("BenchmarkEndToEndCounterAdd")
		counter, err := meter.Float64Counter("test.counter")
		assert.NoError(b, err)
		return counter
	}
	var addOptPool = &sync.Pool{
		New: func() any {
			const n = 1 // WithAttributeSet
			o := make([]metric.AddOption, 0, n)
			// Return a pointer to avoid extra allocation on Put().
			return &o
		},
	}
	ctx := b.Context()
	for _, mp := range []struct {
		name     string
		provider func() metric.MeterProvider
	}{
		{
			name: "NoFilter",
			provider: func() metric.MeterProvider {
				return NewMeterProvider(
					WithReader(NewManualReader()),
					WithExemplarFilter(exemplar.AlwaysOffFilter),
				)
			},
		},
		{
			name: "Filtered",
			provider: func() metric.MeterProvider {
				view := NewView(
					Instrument{
						Name: "test.counter",
					},
					// Filter out one attribute from each call.
					Stream{AttributeFilter: attribute.NewDenyKeysFilter("a")},
				)
				return NewMeterProvider(
					WithView(view),
					WithReader(NewManualReader()),
					WithExemplarFilter(exemplar.AlwaysOffFilter),
				)
			},
		},
	} {
		b.Run(mp.name, func(b *testing.B) {
			for _, attrsLen := range []int{1, 5, 10} {
				attrPool := sync.Pool{
					New: func() any {
						// Pre-allocate common capacity
						s := make([]attribute.KeyValue, 0, attrsLen)
						// Return a pointer to avoid extra allocation on Put().
						return &s
					},
				}
				b.Run(fmt.Sprintf("Attributes/%d", attrsLen), func(b *testing.B) {
					// This case shows the performance of our API + SDK when
					// following our contributor guidance for recording
					// cached attributes by passing attribute.Set:
					// https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#cache-common-attribute-sets-for-repeated-measurements
					b.Run("Precomputed/WithAttributeSet", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						precomputedOpts := []metric.AddOption{
							metric.WithAttributeSet(attribute.NewSet(getAttributes(attrsLen)...)),
						}
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								counter.Add(ctx, 1, precomputedOpts...)
							}
						})
					})
					// This case shows the performance of our API + SDK when
					// following our contributor guidance for recording
					// cached attributes by passing []attribute.KeyValue:
					// https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#cache-common-attribute-sets-for-repeated-measurements
					b.Run("Precomputed/WithAttributes", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						precomputedOpts := []metric.AddOption{metric.WithAttributes(getAttributes(attrsLen)...)}
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								counter.Add(ctx, 1, precomputedOpts...)
							}
						})
					})
					// This case shows the performance of our API + SDK when
					// following our contributor guidance for recording
					// varying attributes by passing attribute.Set:
					// https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#attribute-and-option-allocation-management
					b.Run("Dynamic/WithAttributeSet", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								// Wrap in a function so we can use defer.
								func() {
									attrsSlice := attrPool.Get().(*[]attribute.KeyValue)
									defer func() {
										*attrsSlice = (*attrsSlice)[:0] // Reset.
										attrPool.Put(attrsSlice)
									}()
									appendAttributes(attrsLen, attrsSlice)
									addOpt := addOptPool.Get().(*[]metric.AddOption)
									defer func() {
										*addOpt = (*addOpt)[:0]
										addOptPool.Put(addOpt)
									}()
									set := attribute.NewSet(*attrsSlice...)
									*addOpt = append(*addOpt, metric.WithAttributeSet(set))
									counter.Add(ctx, 1, *addOpt...)
								}()
							}
						})
					})
					// This case shows the performance of our API + SDK when
					// following our contributor guidance for recording
					// varying attributes by passing []attribute.KeyValue:
					// https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#attribute-and-option-allocation-management
					b.Run("Dynamic/WithAttributes", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							i := 0
							for pb.Next() {
								// Wrap in a function so we can use defer.
								func() {
									attrsSlice := attrPool.Get().(*[]attribute.KeyValue)
									defer func() {
										*attrsSlice = (*attrsSlice)[:0] // Reset.
										attrPool.Put(attrsSlice)
									}()
									appendAttributes(attrsLen, attrsSlice)
									addOpt := addOptPool.Get().(*[]metric.AddOption)
									defer func() {
										*addOpt = (*addOpt)[:0]
										addOptPool.Put(addOpt)
									}()
									counter.Add(ctx, 1, metric.WithAttributes(*attrsSlice...))
								}()
								i++
							}
						})
					})
					// This case shows the performance of our API + SDK when
					// users use it in the "obvious" way, without explicitly
					// trying to optimize for performance.
					b.Run("Naive/WithAttributes", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								counter.Add(ctx, 1, metric.WithAttributes(getAttributes(attrsLen)...))
							}
						})
					})
				})
			}
		})
	}
}

func getAttributes(number int) []attribute.KeyValue {
	kvs := make([]attribute.KeyValue, 0, number)
	appendAttributes(number, &kvs)
	return kvs
}

func appendAttributes(number int, kvs *[]attribute.KeyValue) {
	switch number {
	case 1:
		*kvs = append(*kvs,
			attribute.String("a", "a"),
		)
	case 5:
		*kvs = append(*kvs,
			attribute.String("a", "a"),
			attribute.String("b", "b"),
			attribute.String("c", "c"),
			attribute.String("d", "d"),
			attribute.String("e", "e"),
		)
	case 10:
		*kvs = append(*kvs,
			attribute.String("a", "a"),
			attribute.String("b", "b"),
			attribute.String("c", "c"),
			attribute.String("d", "d"),
			attribute.String("e", "e"),
			attribute.String("f", "f"),
			attribute.String("g", "g"),
			attribute.String("h", "h"),
			attribute.String("i", "i"),
			attribute.String("j", "j"),
		)
	default:
		panic("unknown number of attributes")
	}
}
