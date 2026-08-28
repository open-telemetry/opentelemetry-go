// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

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
	"time"

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

func testSumAggregateOutput(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	*dest = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint[int64]{{Value: 1}},
	}
	return 1
}

func TestNewPipeline(t *testing.T) {
	pipe := newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter, 0)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(t.Context(), &output)
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

	err = pipe.produce(t.Context(), &output)
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	require.Len(t, output.ScopeMetrics, 1)
	require.Len(t, output.ScopeMetrics[0].Metrics, 1)
}

func TestPipelineUsesResource(t *testing.T) {
	res := resource.NewWithAttributes("noSchema", attribute.String("test", "resource"))
	pipe := newPipeline(res, nil, nil, exemplar.AlwaysOffFilter, 0)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(t.Context(), &output)
	assert.NoError(t, err)
	assert.Equal(t, res, output.Resource)
}

func TestPipelineConcurrentSafe(t *testing.T) {
	pipe := newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter, 0)
	ctx := t.Context()
	var output metricdata.ResourceMetrics

	var wg sync.WaitGroup
	const threads = 2
	for i := range threads {
		wg.Go(func() {
			_ = pipe.produce(ctx, &output)
		})

		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			name := fmt.Sprintf("name %d", n)
			sync := instrumentSync{name, "desc", "1", testSumAggregateOutput}
			pipe.addSync(instrumentation.Scope{}, sync)
		}(i)

		wg.Go(func() {
			pipe.addMultiCallback(func(context.Context) error { return nil })
		})

		wg.Go(func() {
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
		})
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
				pipe: newPipeline(nil, reader, nil, exemplar.AlwaysOffFilter, 0),
			},
			{
				name: "NoMatchingView",
				pipe: newPipeline(nil, reader, []View{
					NewView(Instrument{Name: "foo"}, Stream{Name: "bar"}),
				}, exemplar.AlwaysOffFilter, 0),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var c cache[string, instID]
				i := newInserter[N](test.pipe, &c)
				readerAggregation := i.readerDefaultAggregation(inst.Kind)
				got, err := i.Instrument(inst, nil, readerAggregation)
				require.NoError(t, err)
				assert.Len(t, got, 1, "default view not applied")
				for _, in := range got {
					in(t.Context(), 1, *attribute.EmptySet())
				}

				out := metricdata.ResourceMetrics{}
				err = test.pipe.produce(t.Context(), &out)
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

		i := newInserter[int64](newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter, 0), &vc)
		i.logConflict(instID{Name: tc.name})

		if tc.conflict {
			assert.Containsf(
				t, msg, "duplicate metric stream definitions",
				"warning not logged for conflicting names: %s, %s",
				tc.existing, tc.name,
			)
		} else {
			assert.Emptyf(
				t, msg,
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
	i := newInserter[int64](newPipeline(nil, nil, nil, exemplar.AlwaysOffFilter, 0), &vc)

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
	pipe := newPipeline(nil, NewManualReader(), nil, exemplar.AlwaysOffFilter, 0)
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
	nCPU := max(runtime.GOMAXPROCS(0), 1)
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
		require.NoError(t, r.Collect(t.Context(), rm))

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

	ctx := t.Context()
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		SpanID:     trace.SpanID{0o1},
		TraceID:    trace.TraceID{0o1},
		TraceFlags: trace.FlagsSampled,
	})
	sampled := trace.ContextWithSpanContext(t.Context(), sc)

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
		reservoirProviderSelector := func(Aggregation) exemplar.ReservoirProvider {
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
	wg.Go(func() {
		_, err := m.Int64ObservableCounter("int64-observable-counter-2")
		require.NoError(t, err)
	})

	wg.Go(func() {
		_, err := m.RegisterCallback(
			func(_ context.Context, o metric.Observer) error {
				o.ObserveInt64(oc1, 2)
				return nil
			}, oc1,
		)
		require.NoError(t, err)
	})

	wg.Go(func() {
		_ = mp.pipes[0].produce(t.Context(), &metricdata.ResourceMetrics{})
	})

	wg.Go(func() {
		_ = mp.pipes[1].produce(t.Context(), &metricdata.ResourceMetrics{})
	})

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
		}, oc,
	)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, reg.Unregister()) })
	ctx := t.Context()
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

func TestCallbackCanUnregisterDuringCollect(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	var called atomic.Int64
	var reg metric.Registration
	reg, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		called.Add(1)
		o.ObserveInt64(counter, 1)
		return reg.Unregister()
	}, counter)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.Equal(t, int64(1), called.Load())
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

	rm = metricdata.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.Equal(t, int64(1), called.Load())
	assert.Empty(t, rm.ScopeMetrics)
}

// TestCallbackCanUnregisterDuringCollectFromDeepStack checks that a callback
// unregistering itself is recognized even when it does so from a call stack
// deeper than inCallback's initial frame buffer.
func TestCallbackCanUnregisterDuringCollectFromDeepStack(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	var reg metric.Registration
	var unregisterAtDepth func(depth int) error
	unregisterAtDepth = func(depth int) error {
		if depth == 0 {
			return reg.Unregister()
		}
		return unregisterAtDepth(depth - 1)
	}

	var called atomic.Int64
	reg, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		called.Add(1)
		o.ObserveInt64(counter, 1)
		return unregisterAtDepth(200)
	}, counter)
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() {
		var rm metricdata.ResourceMetrics
		done <- reader.Collect(t.Context(), &rm)
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Unregister blocked inside the callback it unregisters")
	}
	assert.Equal(t, int64(1), called.Load())

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.Equal(t, int64(1), called.Load(), "callback ran after unregistering itself")
}

func TestUnregisterRemovesCallbackWithoutCollect(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	for range 10 {
		reg, err := m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(counter, 1)
			return nil
		}, counter)
		require.NoError(t, err)
		require.NoError(t, reg.Unregister())
	}

	pipe := mp.pipes[0]
	pipe.Lock()
	defer pipe.Unlock()
	assert.Empty(t, pipe.multiCallbacks)
}

func TestUnregisterDuringCollectRemovesCallback(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	var calls atomic.Int64
	started := make(chan struct{})
	release := make(chan struct{})
	reg, err := m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		if calls.Add(1) == 1 {
			close(started)
			<-release
		}
		o.ObserveInt64(counter, 1)
		return nil
	}, counter)
	require.NoError(t, err)

	collectDone := make(chan error, 1)
	go func() {
		var rm metricdata.ResourceMetrics
		collectDone <- reader.Collect(t.Context(), &rm)
	}()
	<-started

	unregDone := make(chan error, 1)
	go func() { unregDone <- reg.Unregister() }()

	// The pipeline releases the callback before Unregister waits for the
	// in-flight call to return.
	pipe := mp.pipes[0]
	assert.Eventually(t, func() bool {
		pipe.Lock()
		defer pipe.Unlock()
		return len(pipe.multiCallbacks) == 0
	}, 5*time.Second, time.Millisecond, "pipeline retained the callback while Unregister waited")

	close(release)
	select {
	case err := <-unregDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Unregister did not return after the callback returned")
	}
	require.NoError(t, <-collectDone)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.Equal(t, int64(1), calls.Load(), "callback ran after Unregister returned")
}

func TestCallbackCanUnregisterPeerDuringCollect(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	var aCalls, bCalls atomic.Int64
	var regB metric.Registration
	_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		aCalls.Add(1)
		o.ObserveInt64(counter, 1)
		return regB.Unregister()
	}, counter)
	require.NoError(t, err)
	regB, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		bCalls.Add(1)
		o.ObserveInt64(counter, 1)
		return nil
	}, counter)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.Equal(t, int64(1), aCalls.Load())
	assert.Equal(
		t, int64(0), bCalls.Load(),
		"callback ran in the same cycle after a peer callback unregistered it",
	)

	rm = metricdata.ResourceMetrics{}
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.Equal(t, int64(2), aCalls.Load())
	assert.Equal(t, int64(0), bCalls.Load())
}

// TestUnregisterWaitsForInFlightCallback covers the metric.Meter.RegisterCallback
// contract: once Unregister returns, f is neither running nor called again,
// regardless of how many readers were collecting or how many goroutines
// called Unregister.
func TestUnregisterWaitsForInFlightCallback(t *testing.T) {
	tests := []struct {
		name    string
		readers int
		callers int
	}{
		{name: "OneReader", readers: 1, callers: 1},
		{name: "TwoReaders", readers: 2, callers: 1},
		{name: "ConcurrentCallers", readers: 1, callers: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readers := make([]*ManualReader, tt.readers)
			opts := make([]Option, tt.readers)
			for i := range readers {
				readers[i] = NewManualReader()
				opts[i] = WithReader(readers[i])
			}
			mp := NewMeterProvider(opts...)
			m := mp.Meter("test")

			counter, err := m.Int64ObservableCounter("int64-observable-counter")
			require.NoError(t, err)

			started := make(chan struct{}, tt.readers)
			release := make(chan struct{})
			releaseOnce := sync.OnceFunc(func() { close(release) })
			defer releaseOnce()
			var returned atomic.Int64
			reg, err := m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
				started <- struct{}{}
				<-release
				o.ObserveInt64(counter, 1)
				returned.Add(1)
				return nil
			}, counter)
			require.NoError(t, err)

			// Block every reader's collection inside f.
			collectDone := make(chan error, tt.readers)
			for _, r := range readers {
				go func() {
					var rm metricdata.ResourceMetrics
					collectDone <- r.Collect(t.Context(), &rm)
				}()
			}
			for range tt.readers {
				<-started
			}

			// Unregister from outside f while f is in flight.
			type result struct {
				err      error
				returned int64 // invocations of f that had returned when Unregister returned.
			}
			results := make(chan result, tt.callers)
			for range tt.callers {
				go func() {
					err := reg.Unregister()
					results <- result{err: err, returned: returned.Load()}
				}()
			}

			select {
			case <-results:
				t.Fatal("Unregister returned while f was in flight")
			case <-time.After(100 * time.Millisecond):
			}

			releaseOnce()
			for range tt.callers {
				select {
				case res := <-results:
					require.NoError(t, res.err)
					assert.Equal(t, int64(tt.readers), res.returned, "Unregister returned before f returned")
				case <-time.After(5 * time.Second):
					t.Fatal("Unregister did not return after f returned")
				}
			}
			for range tt.readers {
				require.NoError(t, <-collectDone)
			}

			// f is not called by any later collection.
			for _, r := range readers {
				var rm metricdata.ResourceMetrics
				require.NoError(t, r.Collect(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics)
			}
			assert.Equal(t, int64(tt.readers), returned.Load(), "f ran after Unregister returned")
		})
	}
}

// TestUnregisterFromGoroutineSpawnedByCallbackWaits checks that a goroutine
// started by f is treated as being outside f: its call to Unregister waits
// for f to return.
func TestUnregisterFromGoroutineSpawnedByCallbackWaits(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	release := make(chan struct{})
	releaseOnce := sync.OnceFunc(func() { close(release) })
	defer releaseOnce()
	var returned atomic.Bool
	// unregDone receives whether f had returned when Unregister returned.
	unregDone := make(chan bool, 1)
	var reg metric.Registration
	reg, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		go func() {
			_ = reg.Unregister()
			unregDone <- returned.Load()
		}()
		<-release
		o.ObserveInt64(counter, 1)
		returned.Store(true)
		return nil
	}, counter)
	require.NoError(t, err)

	collectDone := make(chan error, 1)
	go func() {
		var rm metricdata.ResourceMetrics
		collectDone <- reader.Collect(t.Context(), &rm)
	}()

	// Wait until the spawned goroutine has removed the callback from the
	// pipeline; it then has to wait for f to return.
	pipe := mp.pipes[0]
	require.Eventually(t, func() bool {
		pipe.Lock()
		defer pipe.Unlock()
		return len(pipe.multiCallbacks) == 0
	}, 5*time.Second, time.Millisecond)
	select {
	case <-unregDone:
		t.Fatal("Unregister returned while f was in flight")
	case <-time.After(100 * time.Millisecond):
	}

	releaseOnce()
	assert.True(t, <-unregDone, "Unregister returned before f returned")
	require.NoError(t, <-collectDone)
}

// TestUnregisterFromCallbackDoesNotWaitForInFlightPeer checks that a callback
// can unregister a peer that is in flight on another reader without blocking.
// Waiting inside a callback could deadlock, and callbacks are expected not to
// block.
func TestUnregisterFromCallbackDoesNotWaitForInFlightPeer(t *testing.T) {
	tests := []struct {
		name string
		// register registers a callback that unregisters reg the first time
		// it runs while inFlight is set, reporting the result on done.
		register func(m metric.Meter, reg *metric.Registration, inFlight *atomic.Bool, done chan<- error) error
	}{
		{
			name: "RegisterCallback",
			register: func(m metric.Meter, reg *metric.Registration, inFlight *atomic.Bool, done chan<- error) error {
				gauge, err := m.Int64ObservableGauge("int64-observable-gauge")
				if err != nil {
					return err
				}
				var once sync.Once
				_, err = m.RegisterCallback(func(context.Context, metric.Observer) error {
					if inFlight.Load() {
						once.Do(func() { done <- (*reg).Unregister() })
					}
					return nil
				}, gauge)
				return err
			},
		},
		{
			name: "InstrumentCallback",
			register: func(m metric.Meter, reg *metric.Registration, inFlight *atomic.Bool, done chan<- error) error {
				var once sync.Once
				_, err := m.Int64ObservableGauge("int64-observable-gauge", metric.WithInt64Callback(
					func(context.Context, metric.Int64Observer) error {
						if inFlight.Load() {
							once.Do(func() { done <- (*reg).Unregister() })
						}
						return nil
					},
				))
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readerA, readerB := NewManualReader(), NewManualReader()
			mp := NewMeterProvider(WithReader(readerA), WithReader(readerB))
			m := mp.Meter("test")

			counter, err := m.Int64ObservableCounter("int64-observable-counter")
			require.NoError(t, err)

			var fCalls atomic.Int64
			var fInFlight atomic.Bool
			fStarted := make(chan struct{})
			releaseF := make(chan struct{})
			releaseFOnce := sync.OnceFunc(func() { close(releaseF) })
			defer releaseFOnce()
			regF, err := m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
				if fCalls.Add(1) == 1 {
					fInFlight.Store(true)
					close(fStarted)
					<-releaseF
				}
				o.ObserveInt64(counter, 1)
				return nil
			}, counter)
			require.NoError(t, err)

			unregDone := make(chan error, 1)
			require.NoError(t, tt.register(m, &regF, &fInFlight, unregDone))

			// Block reader B's collection inside f.
			collectB := make(chan error, 1)
			go func() {
				var rm metricdata.ResourceMetrics
				collectB <- readerB.Collect(t.Context(), &rm)
			}()
			<-fStarted

			// Reader A's collection runs the peer, which unregisters f.
			collectA := make(chan error, 1)
			go func() {
				var rm metricdata.ResourceMetrics
				collectA <- readerA.Collect(t.Context(), &rm)
			}()
			select {
			case err := <-collectA:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("Unregister blocked inside a callback waiting for a peer in flight on another reader")
			}
			require.NoError(t, <-unregDone)

			releaseFOnce()
			require.NoError(t, <-collectB)

			// f is not called again.
			calls := fCalls.Load()
			for _, r := range []*ManualReader{readerA, readerB} {
				var rm metricdata.ResourceMetrics
				require.NoError(t, r.Collect(t.Context(), &rm))
			}
			assert.Equal(t, calls, fCalls.Load(), "f ran after Unregister returned")
		})
	}
}

// TestCallbacksCanUnregisterEachOtherAcrossReaders checks that two callbacks
// in flight on different readers can unregister each other without
// deadlocking, which is why Unregister does not wait when called from within
// a callback.
func TestCallbacksCanUnregisterEachOtherAcrossReaders(t *testing.T) {
	readerA, readerB := NewManualReader(), NewManualReader()
	mp := NewMeterProvider(WithReader(readerA), WithReader(readerB))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	started := make(chan struct{}, 2)
	bothStarted := make(chan struct{})
	bothStartedOnce := sync.OnceFunc(func() { close(bothStarted) })
	defer bothStartedOnce()

	var regF, regG metric.Registration
	var fCalls, gCalls atomic.Int64
	regF, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(counter, 1)
		if fCalls.Add(1) != 1 {
			return nil
		}
		started <- struct{}{}
		<-bothStarted
		return regG.Unregister()
	}, counter)
	require.NoError(t, err)
	regG, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(counter, 1)
		if gCalls.Add(1) != 1 {
			return nil
		}
		started <- struct{}{}
		<-bothStarted
		return regF.Unregister()
	}, counter)
	require.NoError(t, err)

	// Reader A blocks in f. Reader B then runs f, which returns at once,
	// and blocks in g.
	collectA := make(chan error, 1)
	go func() {
		var rm metricdata.ResourceMetrics
		collectA <- readerA.Collect(t.Context(), &rm)
	}()
	<-started
	collectB := make(chan error, 1)
	go func() {
		var rm metricdata.ResourceMetrics
		collectB <- readerB.Collect(t.Context(), &rm)
	}()
	<-started
	bothStartedOnce()

	for _, done := range []chan error{collectA, collectB} {
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("callbacks unregistering each other across readers deadlocked")
		}
	}

	// Neither callback runs again.
	f, g := fCalls.Load(), gCalls.Load()
	for _, r := range []*ManualReader{readerA, readerB} {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))
		assert.Empty(t, rm.ScopeMetrics)
	}
	assert.Equal(t, f, fCalls.Load())
	assert.Equal(t, g, gCalls.Load())
}

// TestUnregisterAfterCallbackPanic checks that a callback that panics during
// collection does not leave its registration in a state that blocks
// Unregister.
func TestUnregisterAfterCallbackPanic(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	reg, err := m.RegisterCallback(func(context.Context, metric.Observer) error {
		panic("callback panic")
	}, counter)
	require.NoError(t, err)

	func() {
		defer func() { assert.Equal(t, "callback panic", recover()) }()
		var rm metricdata.ResourceMetrics
		_ = reader.Collect(t.Context(), &rm)
	}()

	done := make(chan error, 1)
	go func() { done <- reg.Unregister() }()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Unregister blocked after the callback panicked")
	}
}

// TestCallbackCanUnregisterDuringConcurrentCollects checks that a callback
// can unregister itself while it is in flight on two readers at once.
func TestCallbackCanUnregisterDuringConcurrentCollects(t *testing.T) {
	readerA, readerB := NewManualReader(), NewManualReader()
	mp := NewMeterProvider(WithReader(readerA), WithReader(readerB))
	m := mp.Meter("test")

	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)

	started := make(chan struct{}, 2)
	bothStarted := make(chan struct{})
	bothStartedOnce := sync.OnceFunc(func() { close(bothStarted) })
	defer bothStartedOnce()

	var calls atomic.Int64
	var reg metric.Registration
	reg, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		calls.Add(1)
		o.ObserveInt64(counter, 1)
		started <- struct{}{}
		<-bothStarted
		return reg.Unregister()
	}, counter)
	require.NoError(t, err)

	collectDone := make(chan error, 2)
	for _, r := range []*ManualReader{readerA, readerB} {
		go func() {
			var rm metricdata.ResourceMetrics
			collectDone <- r.Collect(t.Context(), &rm)
		}()
	}
	<-started
	<-started
	bothStartedOnce()

	for range 2 {
		select {
		case err := <-collectDone:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("Unregister blocked inside the callback it unregisters")
		}
	}
	assert.Equal(t, int64(2), calls.Load())

	for _, r := range []*ManualReader{readerA, readerB} {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))
		assert.Empty(t, rm.ScopeMetrics)
	}
	assert.Equal(t, int64(2), calls.Load(), "callback ran after unregistering itself")
}

func TestInCallback(t *testing.T) {
	assert.False(t, inCallback(), "outside any callback")

	var inside bool
	require.NoError(t, runCallback(t.Context(), func(context.Context) error {
		inside = inCallback()
		return nil
	}))
	assert.True(t, inside, "within a callback run by runCallback")

	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	var inMulti, inInstrument atomic.Bool
	counter, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)
	_, err = m.RegisterCallback(func(context.Context, metric.Observer) error {
		inMulti.Store(inCallback())
		return nil
	}, counter)
	require.NoError(t, err)
	_, err = m.Int64ObservableGauge("int64-observable-gauge", metric.WithInt64Callback(
		func(context.Context, metric.Int64Observer) error {
			inInstrument.Store(inCallback())
			return nil
		},
	))
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.True(t, inMulti.Load(), "within a RegisterCallback callback")
	assert.True(t, inInstrument.Load(), "within an instrument callback")
	assert.False(t, inCallback(), "after collection")
}

func TestCallbackCanCallMeterMethods(t *testing.T) {
	reader := NewManualReader()
	mp := NewMeterProvider(WithReader(reader))
	m := mp.Meter("test")

	int64Observable, err := m.Int64ObservableCounter("int64-observable-counter")
	require.NoError(t, err)
	float64Observable, err := m.Float64ObservableCounter("float64-observable-counter")
	require.NoError(t, err)

	var called atomic.Bool
	var reg metric.Registration
	reg, err = m.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		called.Store(true)
		o.ObserveInt64(int64Observable, 1)
		o.ObserveFloat64(float64Observable, 1)

		int64Counter, err := m.Int64Counter("int64-counter")
		if err != nil {
			return err
		}
		int64Counter.Add(ctx, 1)
		_ = int64Counter.Enabled(ctx)

		int64UpDownCounter, err := m.Int64UpDownCounter("int64-up-down-counter")
		if err != nil {
			return err
		}
		int64UpDownCounter.Add(ctx, 1)
		_ = int64UpDownCounter.Enabled(ctx)

		int64Histogram, err := m.Int64Histogram("int64-histogram")
		if err != nil {
			return err
		}
		int64Histogram.Record(ctx, 1)
		_ = int64Histogram.Enabled(ctx)

		int64Gauge, err := m.Int64Gauge("int64-gauge")
		if err != nil {
			return err
		}
		int64Gauge.Record(ctx, 1)
		_ = int64Gauge.Enabled(ctx)

		float64Counter, err := m.Float64Counter("float64-counter")
		if err != nil {
			return err
		}
		float64Counter.Add(ctx, 1)
		_ = float64Counter.Enabled(ctx)

		float64UpDownCounter, err := m.Float64UpDownCounter("float64-up-down-counter")
		if err != nil {
			return err
		}
		float64UpDownCounter.Add(ctx, 1)
		_ = float64UpDownCounter.Enabled(ctx)

		float64Histogram, err := m.Float64Histogram("float64-histogram")
		if err != nil {
			return err
		}
		float64Histogram.Record(ctx, 1)
		_ = float64Histogram.Enabled(ctx)

		float64Gauge, err := m.Float64Gauge("float64-gauge")
		if err != nil {
			return err
		}
		float64Gauge.Record(ctx, 1)
		_ = float64Gauge.Enabled(ctx)

		_, err = m.Int64ObservableCounter(
			"int64-observable-counter-from-callback",
			metric.WithInt64Callback(func(context.Context, metric.Int64Observer) error { return nil }),
		)
		if err != nil {
			return err
		}
		_, err = m.Int64ObservableUpDownCounter(
			"int64-observable-up-down-counter-from-callback",
			metric.WithInt64Callback(func(context.Context, metric.Int64Observer) error { return nil }),
		)
		if err != nil {
			return err
		}
		_, err = m.Int64ObservableGauge(
			"int64-observable-gauge-from-callback",
			metric.WithInt64Callback(func(context.Context, metric.Int64Observer) error { return nil }),
		)
		if err != nil {
			return err
		}
		_, err = m.Float64ObservableCounter(
			"float64-observable-counter-from-callback",
			metric.WithFloat64Callback(func(context.Context, metric.Float64Observer) error { return nil }),
		)
		if err != nil {
			return err
		}
		_, err = m.Float64ObservableUpDownCounter(
			"float64-observable-up-down-counter-from-callback",
			metric.WithFloat64Callback(func(context.Context, metric.Float64Observer) error { return nil }),
		)
		if err != nil {
			return err
		}
		observable, err := m.Float64ObservableGauge("registered-observable-gauge")
		if err != nil {
			return err
		}
		innerReg, err := m.RegisterCallback(
			func(context.Context, metric.Observer) error { return nil },
			observable,
		)
		if err != nil {
			return err
		}
		if err := innerReg.Unregister(); err != nil {
			return err
		}
		return reg.Unregister()
	}, int64Observable, float64Observable)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	assert.True(t, called.Load())
}

// TestPipelineProduceErrors tests the issue described in https://github.com/open-telemetry/opentelemetry-go/issues/6344.
// Earlier implementations of the pipeline produce method could corrupt metric data point state when the passed context
// was canceled during execution of callbacks. In this case, corroption was the result of some or all callbacks being
// invoked without instrument compAgg functions called.
func TestPipelineProduceErrors(t *testing.T) {
	// Create a test pipeline with aggregations
	pipeReader := NewManualReader()
	pipe := newPipeline(nil, pipeReader, nil, exemplar.AlwaysOffFilter, 0)

	// Set up an observable with callbacks
	var testObsID observableID[int64]
	aggBuilder := aggregate.Builder[int64]{Temporality: metricdata.CumulativeTemporality}
	measure, _ := aggBuilder.Sum(true)
	pipe.addInt64Measure(testObsID, []aggregate.Measure[int64]{measure})

	// Add an aggregation that just sets the data point value to the number of times the aggregation is invoked
	aggCallCount := 0
	inst := instrumentSync{
		name:        "test-metric",
		description: "test description",
		unit:        "test unit",
		compAgg: func(dest *metricdata.Aggregation) int {
			aggCallCount++

			*dest = metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: false,
				DataPoints:  []metricdata.DataPoint[int64]{{Value: int64(aggCallCount)}},
			}
			return aggCallCount
		},
	}
	pipe.addSync(instrumentation.Scope{Name: "test"}, inst)

	ctx, cancelCtx := context.WithCancel(t.Context())
	var shouldCancelContext bool // When true, the second callback cancels ctx
	var shouldReturnError bool   // When true, the third callback returns an error
	var callbackCounts [3]int

	pipe.callbacks = append(pipe.callbacks,
		// Callback 1: cancels the context during execution but continues to populate data
		func(ctx context.Context) error {
			callbackCounts[0]++
			for _, m := range pipe.int64Measures[testObsID] {
				m(ctx, 123, *attribute.EmptySet())
			}
			return nil
		},
		// Callback 2: populates int64 observable data
		func(context.Context) error {
			callbackCounts[1]++
			if shouldCancelContext {
				cancelCtx()
			}
			return nil
		},
		// Callback 3: return an error
		func(context.Context) error {
			callbackCounts[2]++
			if shouldReturnError {
				return fmt.Errorf("test callback error")
			}
			return nil
		})

	assertMetrics := func(rm *metricdata.ResourceMetrics, expectVal int64) {
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		metricdatatest.AssertEqual(t, metricdata.Metrics{
			Name:        inst.name,
			Description: inst.description,
			Unit:        inst.unit,
			Data: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: false,
				DataPoints:  []metricdata.DataPoint[int64]{{Value: expectVal}},
			},
		}, rm.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
	}

	t.Run("no errors", func(t *testing.T) {
		var rm metricdata.ResourceMetrics
		err := pipe.produce(ctx, &rm)
		require.NoError(t, err)

		assert.Equal(t, [3]int{1, 1, 1}, callbackCounts)
		assert.Equal(t, 1, aggCallCount)

		assertMetrics(&rm, 1)
	})

	t.Run("callback returns error", func(t *testing.T) {
		shouldReturnError = true

		var rm metricdata.ResourceMetrics
		err := pipe.produce(ctx, &rm)
		require.ErrorContains(t, err, "test callback error")

		// Even though a callback returned an error, the agg function is still called
		assert.Equal(t, [3]int{2, 2, 2}, callbackCounts)
		assert.Equal(t, 2, aggCallCount)

		assertMetrics(&rm, 2)
	})

	t.Run("context canceled during produce", func(t *testing.T) {
		shouldCancelContext = true

		var rm metricdata.ResourceMetrics
		err := pipe.produce(ctx, &rm)
		require.ErrorContains(t, err, "test callback error")

		// Even though the context was canceled midway through invoking callbacks,
		// all remaining callbacks and agg functions are still called
		assert.Equal(t, [3]int{3, 3, 3}, callbackCounts)
		assert.Equal(t, 3, aggCallCount)
	})

	t.Run("context already cancelled", func(t *testing.T) {
		var output metricdata.ResourceMetrics
		err := pipe.produce(ctx, &output)
		require.ErrorIs(t, err, context.Canceled)

		// No callbacks or agg functions are called since the context was canceled prior to invoking
		// the produce method
		assert.Equal(t, [3]int{3, 3, 3}, callbackCounts)
		assert.Equal(t, 3, aggCallCount)
	})
}
