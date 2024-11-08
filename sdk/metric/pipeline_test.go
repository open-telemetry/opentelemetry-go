// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func testSumAggregateOutput(dest *metricdata.Aggregation) int {
	*dest = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint[int64]{{Value: 1}},
	}
	return 1
}

func TestNewPipeline(t *testing.T) {
	pipe := newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(context.Background(), &output)
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	assert.Empty(t, output.ScopeMetrics)

	iSync := instrumentSync{"name", "desc", "1", testSumAggregateOutput}
	assert.NotPanics(t, func() {
		pipe.addSync(instrumentation.Scope{}, iSync)
	})

	require.NotPanics(t, func() {
		pipe.addMultiCallback(func(context.Context) error { return nil })
	})

	err = pipe.produce(context.Background(), &output)
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	require.Len(t, output.ScopeMetrics, 1)
	require.Len(t, output.ScopeMetrics[0].Metrics, 1)
}

func TestPipelineUsesResource(t *testing.T) {
	res := resource.NewWithAttributes("noSchema", attribute.String("test", "resource"))
	pipe := newPipeline(res, nil, nil, exemplar.AlwaysOffFilter)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(context.Background(), &output)
	assert.NoError(t, err)
	assert.Equal(t, res, output.Resource)
}

func TestPipelineConcurrentSafe(t *testing.T) {
	pipe := newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter)
	ctx := context.Background()
	var output metricdata.ResourceMetrics

	var wg sync.WaitGroup
	const threads = 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = pipe.produce(ctx, &output)
		}()

		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			name := fmt.Sprintf("name %d", n)
			sync := instrumentSync{name, "desc", "1", testSumAggregateOutput}
			pipe.addSync(instrumentation.Scope{}, sync)
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			pipe.addMultiCallback(func(context.Context) error { return nil })
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			b := aggregate.Builder[int64]{
				Temporality:      metricdata.CumulativeTemporality,
				ReservoirFunc:    nil,
				AggregationLimit: 0,
			}
			var oID observableID[int64]
			m, _ := b.PrecomputedSum(false)
			measures := []aggregate.Measure[int64]{}
			measures = append(measures, m)
			pipe.addInt64Measure(oID, measures)
		}()
	}
	wg.Wait()
}

func TestDefaultViewImplicit(t *testing.T) {
	t.Run("Int64", testDefaultViewImplicit[int64]())
	t.Run("Float64", testDefaultViewImplicit[float64]())
}

func testDefaultViewImplicit[N int64 | float64]() func(t *testing.T) {
	inst := Instrument{
		Name:        "requests",
		Description: "count of requests received",
		Kind:        InstrumentKindCounter,
		Unit:        "1",
	}
	return func(t *testing.T) {
		reader := NewManualReader()
		tests := []struct {
			name string
			pipe *pipeline
		}{
			{
				name: "NoView",
				pipe: newPipeline(nil, reader, nil, exemplar.AlwaysOffFilter),
			},
			{
				name: "NoMatchingView",
				pipe: newPipeline(nil, reader, []View{
					NewView(Instrument{Name: "foo"}, Stream{Name: "bar"}),
				}, exemplar.AlwaysOffFilter),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var c cache[string, instID]
				i := newInserter[N](test.pipe, &c)
				readerAggregation := i.readerDefaultAggregation(inst.Kind)
				got, err := i.Instrument(inst, readerAggregation)
				require.NoError(t, err)
				assert.Len(t, got, 1, "default view not applied")
				for _, in := range got {
					in(context.Background(), 1, *attribute.EmptySet())
				}

				out := metricdata.ResourceMetrics{}
				err = test.pipe.produce(context.Background(), &out)
				require.NoError(t, err)
				require.Len(t, out.ScopeMetrics, 1, "Aggregator not registered with pipeline")
				sm := out.ScopeMetrics[0]
				require.Len(t, sm.Metrics, 1, "metrics not produced from default view")
				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        inst.Name,
					Description: inst.Description,
					Unit:        "1",
					Data: metricdata.Sum[N]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints:  []metricdata.DataPoint[N]{{Value: N(1)}},
					},
				}, sm.Metrics[0], metricdatatest.IgnoreTimestamp())
			})
		}
	}
}

func TestLogConflictName(t *testing.T) {
	testcases := []struct {
		existing, name string
		conflict       bool
	}{
		{
			existing: "requestCount",
			name:     "requestCount",
			conflict: false,
		},
		{
			existing: "requestCount",
			name:     "requestDuration",
			conflict: false,
		},
		{
			existing: "requestCount",
			name:     "requestcount",
			conflict: true,
		},
		{
			existing: "requestCount",
			name:     "REQUESTCOUNT",
			conflict: true,
		},
		{
			existing: "requestCount",
			name:     "rEqUeStCoUnT",
			conflict: true,
		},
	}

	var msg string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) {
			msg = args
		}, funcr.Options{Verbosity: 20}))
		return func() { otel.SetLogger(orig) }
	}(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))))

	for _, tc := range testcases {
		var vc cache[string, instID]

		name := strings.ToLower(tc.existing)
		_ = vc.Lookup(name, func() instID {
			return instID{Name: tc.existing}
		})

		i := newInserter[int64](newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter), &vc)
		i.logConflict(instID{Name: tc.name})

		if tc.conflict {
			assert.Containsf(
				t, msg, "duplicate metric stream definitions",
				"warning not logged for conflicting names: %s, %s",
				tc.existing, tc.name,
			)
		} else {
			assert.Equalf(
				t, "", msg,
				"warning logged for non-conflicting names: %s, %s",
				tc.existing, tc.name,
			)
		}

		// Reset.
		msg = ""
	}
}

func TestLogConflictSuggestView(t *testing.T) {
	var msg string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) {
			msg = args
		}, funcr.Options{Verbosity: 20}))
		return func() { otel.SetLogger(orig) }
	}(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))))

	orig := instID{
		Name:        "requestCount",
		Description: "number of requests",
		Kind:        InstrumentKindCounter,
		Unit:        "1",
		Number:      "int64",
	}

	var vc cache[string, instID]
	name := strings.ToLower(orig.Name)
	_ = vc.Lookup(name, func() instID { return orig })
	i := newInserter[int64](newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter), &vc)

	viewSuggestion := func(inst instID, stream string) string {
		return `"NewView(Instrument{` +
			`Name: \"` + inst.Name +
			`\", Description: \"` + inst.Description +
			`\", Kind: \"InstrumentKind` + inst.Kind.String() +
			`\", Unit: \"` + inst.Unit +
			`\"}, ` +
			stream +
			`)"`
	}

	t.Run("Name", func(t *testing.T) {
		inst := instID{
			Name:        "requestcount",
			Description: orig.Description,
			Kind:        orig.Kind,
			Unit:        orig.Unit,
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.Containsf(t, msg, viewSuggestion(
			inst, `Stream{Name: \"{{NEW_NAME}}\"}`,
		), "no suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Description", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: "alt",
			Kind:        orig.Kind,
			Unit:        orig.Unit,
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.Containsf(t, msg, viewSuggestion(
			inst, `Stream{Description: \"`+orig.Description+`\"}`,
		), "no suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Kind", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: orig.Description,
			Kind:        InstrumentKindHistogram,
			Unit:        orig.Unit,
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.Containsf(t, msg, viewSuggestion(
			inst, `Stream{Name: \"{{NEW_NAME}}\"}`,
		), "no suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Unit", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: orig.Description,
			Kind:        orig.Kind,
			Unit:        "ms",
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.NotContains(t, msg, "NewView", "suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Number", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: orig.Description,
			Kind:        orig.Kind,
			Unit:        orig.Unit,
			Number:      "float64",
		}
		i.logConflict(inst)
		assert.NotContains(t, msg, "NewView", "suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})
}

func TestInserterCachedAggregatorNameConflict(t *testing.T) {
	const name = "requestCount"
	scope := instrumentation.Scope{Name: "pipeline_test"}
	kind := InstrumentKindCounter
	stream := Stream{
		Name:        name,
		Aggregation: AggregationSum{},
	}

	var vc cache[string, instID]
	pipe := newPipeline(nil, NewManualReader(), nil, exemplar.AlwaysOffFilter)
	i := newInserter[int64](pipe, &vc)

	readerAggregation := i.readerDefaultAggregation(kind)
	_, origID, err := i.cachedAggregator(scope, kind, stream, readerAggregation)
	require.NoError(t, err)

	require.Len(t, pipe.aggregations, 1)
	require.Contains(t, pipe.aggregations, scope)
	iSync := pipe.aggregations[scope]
	require.Len(t, iSync, 1)
	require.Equal(t, name, iSync[0].name)

	stream.Name = "RequestCount"
	_, id, err := i.cachedAggregator(scope, kind, stream, readerAggregation)
	require.NoError(t, err)
	assert.Equal(t, origID, id, "multiple aggregators for equivalent name")

	assert.Len(t, pipe.aggregations, 1, "additional scope added")
	require.Contains(t, pipe.aggregations, scope, "original scope removed")
	iSync = pipe.aggregations[scope]
	require.Len(t, iSync, 1, "registered instrumentSync changed")
	assert.Equal(t, name, iSync[0].name, "stream name changed")
}

func TestExemplars(t *testing.T) {
	nCPU := runtime.NumCPU()
	setup := func(name string) (metric.Meter, Reader) {
		r := NewManualReader()
		v := NewView(Instrument{Name: "int64-expo-histogram"}, Stream{
			Aggregation: AggregationBase2ExponentialHistogram{
				MaxSize:  160, // > 20, reservoir size should default to 20.
				MaxScale: 20,
			},
		})
		return NewMeterProvider(WithReader(r), WithView(v)).Meter(name), r
	}

	measure := func(ctx context.Context, m metric.Meter) {
		i, err := m.Int64Counter("int64-counter")
		require.NoError(t, err)

		h, err := m.Int64Histogram("int64-histogram")
		require.NoError(t, err)

		e, err := m.Int64Histogram("int64-expo-histogram")
		require.NoError(t, err)

		for j := 0; j < 20*nCPU; j++ { // will be >= 20 and > nCPU
			i.Add(ctx, 1)
			h.Record(ctx, 1)
			e.Record(ctx, 1)
		}
	}

	check := func(t *testing.T, r Reader, nSum, nHist, nExpo int) {
		t.Helper()

		rm := new(metricdata.ResourceMetrics)
		require.NoError(t, r.Collect(context.Background(), rm))

		require.Len(t, rm.ScopeMetrics, 1, "ScopeMetrics")
		sm := rm.ScopeMetrics[0]
		require.Len(t, sm.Metrics, 3, "Metrics")

		require.IsType(t, metricdata.Sum[int64]{}, sm.Metrics[0].Data, sm.Metrics[0].Name)
		sum := sm.Metrics[0].Data.(metricdata.Sum[int64])
		assert.Len(t, sum.DataPoints[0].Exemplars, nSum)

		require.IsType(t, metricdata.Histogram[int64]{}, sm.Metrics[1].Data, sm.Metrics[1].Name)
		hist := sm.Metrics[1].Data.(metricdata.Histogram[int64])
		assert.Len(t, hist.DataPoints[0].Exemplars, nHist)

		require.IsType(t, metricdata.ExponentialHistogram[int64]{}, sm.Metrics[2].Data, sm.Metrics[2].Name)
		expo := sm.Metrics[2].Data.(metricdata.ExponentialHistogram[int64])
		assert.Len(t, expo.DataPoints[0].Exemplars, nExpo)
	}

	ctx := context.Background()
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		SpanID:     trace.SpanID{0o1},
		TraceID:    trace.TraceID{0o1},
		TraceFlags: trace.FlagsSampled,
	})
	sampled := trace.ContextWithSpanContext(context.Background(), sc)

	t.Run("Default", func(t *testing.T) {
		m, r := setup("default")
		measure(ctx, m)
		check(t, r, 0, 0, 0)

		measure(sampled, m)
		check(t, r, nCPU, 1, 20)
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Setenv("OTEL_METRICS_EXEMPLAR_FILTER", "unrecognized")
		m, r := setup("default")
		measure(ctx, m)
		check(t, r, 0, 0, 0)

		measure(sampled, m)
		check(t, r, nCPU, 1, 20)
	})

	t.Run("always_on", func(t *testing.T) {
		t.Setenv("OTEL_METRICS_EXEMPLAR_FILTER", "always_on")
		m, r := setup("always_on")
		measure(ctx, m)
		check(t, r, nCPU, 1, 20)
	})

	t.Run("always_off", func(t *testing.T) {
		t.Setenv("OTEL_METRICS_EXEMPLAR_FILTER", "always_off")
		m, r := setup("always_off")
		measure(ctx, m)
		check(t, r, 0, 0, 0)
	})

	t.Run("trace_based", func(t *testing.T) {
		t.Setenv("OTEL_METRICS_EXEMPLAR_FILTER", "trace_based")
		m, r := setup("trace_based")
		measure(ctx, m)
		check(t, r, 0, 0, 0)

		measure(sampled, m)
		check(t, r, nCPU, 1, 20)
	})

	t.Run("Custom reservoir", func(t *testing.T) {
		r := NewManualReader()
		reservoirProviderSelector := func(agg Aggregation) exemplar.ReservoirProvider {
			return exemplar.FixedSizeReservoirProvider(2)
		}
		v1 := NewView(Instrument{Name: "int64-expo-histogram"}, Stream{
			Aggregation: AggregationBase2ExponentialHistogram{
				MaxSize:  160, // > 20, reservoir size should default to 20.
				MaxScale: 20,
			},
			ExemplarReservoirProviderSelector: reservoirProviderSelector,
		})
		v2 := NewView(Instrument{Name: "int64-counter"}, Stream{
			ExemplarReservoirProviderSelector: reservoirProviderSelector,
		})
		v3 := NewView(Instrument{Name: "int64-histogram"}, Stream{
			ExemplarReservoirProviderSelector: reservoirProviderSelector,
		})
		m := NewMeterProvider(WithReader(r), WithView(v1, v2, v3)).Meter("custom-reservoir")
		measure(ctx, m)
		check(t, r, 0, 0, 0)

		measure(sampled, m)
		check(t, r, 2, 2, 2)
	})
}

func TestAddingAndObservingMeasureConcurrentSafe(t *testing.T) {
	r1 := NewManualReader()
	r2 := NewManualReader()

	mp := NewMeterProvider(WithReader(r1), WithReader(r2))
	m := mp.Meter("test")

	oc1, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := m.Int64ObservableCounter("int64-observable-counter-2")
		require.NoError(t, err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := m.RegisterCallback(
			func(_ context.Context, o metric.Observer) error {
				o.ObserveInt64(oc1, 2)
				return nil
			}, oc1)
		require.NoError(t, err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = mp.pipes[0].produce(context.Background(), &metricdata.ResourceMetrics{})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = mp.pipes[1].produce(context.Background(), &metricdata.ResourceMetrics{})
	}()

	wg.Wait()
}

func TestPipelineWithMultipleReaders(t *testing.T) {
	r1 := NewManualReader()
	r2 := NewManualReader()
	mp := NewMeterProvider(WithReader(r1), WithReader(r2))
	m := mp.Meter("test")
	var val atomic.Int64
	oc, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)
	reg, err := m.RegisterCallback(
		// SDK calls this function when collecting data.
		func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(oc, val.Load())
			return nil
		}, oc)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, reg.Unregister()) })
	ctx := context.Background()
	rm := new(metricdata.ResourceMetrics)
	val.Add(1)
	err = r1.Collect(ctx, rm)
	require.NoError(t, err)
	if assert.Len(t, rm.ScopeMetrics, 1) &&
		assert.Len(t, rm.ScopeMetrics[0].Metrics, 1) {
		assert.Equal(t, int64(1), rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64]).DataPoints[0].Value)
	}
	val.Add(1)
	err = r2.Collect(ctx, rm)
	require.NoError(t, err)
	if assert.Len(t, rm.ScopeMetrics, 1) &&
		assert.Len(t, rm.ScopeMetrics[0].Metrics, 1) {
		assert.Equal(t, int64(2), rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64]).DataPoints[0].Value)
	}
}
