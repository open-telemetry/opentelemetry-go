// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
)

// A meter should be able to make instruments concurrently.
func TestMeterInstrumentConcurrency(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	wg := &sync.WaitGroup{}
	wg.Add(12)

	m := NewMeterProvider().Meter("inst-concurrency")

	go func() {
		_ = m.Float64ObservableCounter("AFCounter")
		wg.Done()
	}()
	go func() {
		_ = m.Float64ObservableUpDownCounter("AFUpDownCounter")
		wg.Done()
	}()
	go func() {
		_ = m.Float64ObservableGauge("AFGauge")
		wg.Done()
	}()
	go func() {
		_ = m.Int64ObservableCounter("AICounter")
		wg.Done()
	}()
	go func() {
		_ = m.Int64ObservableUpDownCounter("AIUpDownCounter")
		wg.Done()
	}()
	go func() {
		_ = m.Int64ObservableGauge("AIGauge")
		wg.Done()
	}()
	go func() {
		_ = m.Float64Counter("SFCounter")
		wg.Done()
	}()
	go func() {
		_ = m.Float64UpDownCounter("SFUpDownCounter")
		wg.Done()
	}()
	go func() {
		_ = m.Float64Histogram("SFHistogram")
		wg.Done()
	}()
	go func() {
		_ = m.Int64Counter("SICounter")
		wg.Done()
	}()
	go func() {
		_ = m.Int64UpDownCounter("SIUpDownCounter")
		wg.Done()
	}()
	go func() {
		_ = m.Int64Histogram("SIHistogram")
		wg.Done()
	}()

	wg.Wait()
	eh.RequireNoErrors(t)
}

var emptyCallback metric.Callback = func(context.Context, metric.Observer) error { return nil }

// A Meter Should be able register Callbacks Concurrently.
func TestMeterCallbackCreationConcurrency(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	m := NewMeterProvider().Meter("callback-concurrency")

	go func() {
		_ = m.RegisterCallback(emptyCallback)
		wg.Done()
	}()
	go func() {
		_ = m.RegisterCallback(emptyCallback)
		wg.Done()
	}()
	wg.Wait()
	eh.RequireNoErrors(t)
}

func TestNoopCallbackUnregisterConcurrency(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	m := NewMeterProvider().Meter("noop-unregister-concurrency")
	reg := m.RegisterCallback(emptyCallback)
	eh.RequireNoErrors(t)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		reg.Unregister()
		wg.Done()
	}()
	go func() {
		reg.Unregister()
		wg.Done()
	}()
	wg.Wait()
	eh.AssertNoErrors(t)
}

func TestCallbackUnregisterConcurrency(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	reader := NewManualReader()
	provider := NewMeterProvider(WithReader(reader))
	meter := provider.Meter("unregister-concurrency")

	actr := meter.Float64ObservableCounter("counter")
	eh.RequireNoErrors(t)

	ag := meter.Int64ObservableGauge("gauge")
	eh.RequireNoErrors(t)

	regCtr := meter.RegisterCallback(emptyCallback, actr)
	eh.RequireNoErrors(t)

	regG := meter.RegisterCallback(emptyCallback, ag)
	eh.RequireNoErrors(t)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		regCtr.Unregister()
		regG.Unregister()
		wg.Done()
	}()
	go func() {
		regCtr.Unregister()
		regG.Unregister()
		wg.Done()
	}()
	wg.Wait()
	eh.AssertNoErrors(t)
}

// Instruments should produce correct ResourceMetrics.
func TestMeterCreatesInstruments(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	extrema := metricdata.NewExtrema(7.)
	attrs := []attribute.KeyValue{attribute.String("name", "alice")}
	testCases := []struct {
		name string
		fn   func(*testing.T, metric.Meter)
		want metricdata.Metrics
	}{
		{
			name: "ObservableInt64Count",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o instrument.Int64Observer) error {
					o.Observe(4, attrs...)
					return nil
				}
				ctr := m.Int64ObservableCounter("aint", instrument.WithInt64Callback(cback))
				eh.AssertNoErrors(t)
				_ = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 3)
					return nil
				}, ctr)
				eh.AssertNoErrors(t)
			},
			want: metricdata.Metrics{
				Name: "aint",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...), Value: 4},
						{Value: 3},
					},
				},
			},
		},
		{
			name: "ObservableInt64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o instrument.Int64Observer) error {
					o.Observe(4, attrs...)
					return nil
				}
				ctr := m.Int64ObservableUpDownCounter("aint", instrument.WithInt64Callback(cback))
				eh.AssertNoErrors(t)
				_ = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 11)
					return nil
				}, ctr)
				eh.AssertNoErrors(t)
			},
			want: metricdata.Metrics{
				Name: "aint",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...), Value: 4},
						{Value: 11},
					},
				},
			},
		},
		{
			name: "ObservableInt64Gauge",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o instrument.Int64Observer) error {
					o.Observe(4, attrs...)
					return nil
				}
				gauge := m.Int64ObservableGauge("agauge", instrument.WithInt64Callback(cback))
				eh.AssertNoErrors(t)
				_ = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(gauge, 11)
					return nil
				}, gauge)
				eh.AssertNoErrors(t)
			},
			want: metricdata.Metrics{
				Name: "agauge",
				Data: metricdata.Gauge[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attribute.NewSet(attrs...), Value: 4},
						{Value: 11},
					},
				},
			},
		},
		{
			name: "ObservableFloat64Count",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o instrument.Float64Observer) error {
					o.Observe(4, attrs...)
					return nil
				}
				ctr := m.Float64ObservableCounter("afloat", instrument.WithFloat64Callback(cback))
				eh.AssertNoErrors(t)
				_ = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 3)
					return nil
				}, ctr)
				eh.AssertNoErrors(t)
			},
			want: metricdata.Metrics{
				Name: "afloat",
				Data: metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: attribute.NewSet(attrs...), Value: 4},
						{Value: 3},
					},
				},
			},
		},
		{
			name: "ObservableFloat64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o instrument.Float64Observer) error {
					o.Observe(4, attrs...)
					return nil
				}
				ctr := m.Float64ObservableUpDownCounter("afloat", instrument.WithFloat64Callback(cback))
				eh.AssertNoErrors(t)
				_ = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 11)
					return nil
				}, ctr)
				eh.AssertNoErrors(t)
			},
			want: metricdata.Metrics{
				Name: "afloat",
				Data: metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: attribute.NewSet(attrs...), Value: 4},
						{Value: 11},
					},
				},
			},
		},
		{
			name: "ObservableFloat64Gauge",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o instrument.Float64Observer) error {
					o.Observe(4, attrs...)
					return nil
				}
				gauge := m.Float64ObservableGauge("agauge", instrument.WithFloat64Callback(cback))
				eh.AssertNoErrors(t)
				_ = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(gauge, 11)
					return nil
				}, gauge)
				eh.AssertNoErrors(t)
			},
			want: metricdata.Metrics{
				Name: "agauge",
				Data: metricdata.Gauge[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: attribute.NewSet(attrs...), Value: 4},
						{Value: 11},
					},
				},
			},
		},

		{
			name: "SyncInt64Count",
			fn: func(t *testing.T, m metric.Meter) {
				ctr := m.Int64Counter("sint")
				eh.AssertNoErrors(t)

				ctr.Add(context.Background(), 3)
			},
			want: metricdata.Metrics{
				Name: "sint",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{Value: 3},
					},
				},
			},
		},
		{
			name: "SyncInt64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				ctr := m.Int64UpDownCounter("sint")
				eh.AssertNoErrors(t)

				ctr.Add(context.Background(), 11)
			},
			want: metricdata.Metrics{
				Name: "sint",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[int64]{
						{Value: 11},
					},
				},
			},
		},
		{
			name: "SyncInt64Histogram",
			fn: func(t *testing.T, m metric.Meter) {
				gauge := m.Int64Histogram("histogram")
				eh.AssertNoErrors(t)

				gauge.Record(context.Background(), 7)
			},
			want: metricdata.Metrics{
				Name: "histogram",
				Data: metricdata.Histogram{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint{
						{
							Attributes:   attribute.Set{},
							Count:        1,
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Min:          extrema,
							Max:          extrema,
							Sum:          7.0,
						},
					},
				},
			},
		},
		{
			name: "SyncFloat64Count",
			fn: func(t *testing.T, m metric.Meter) {
				ctr := m.Float64Counter("sfloat")
				eh.AssertNoErrors(t)

				ctr.Add(context.Background(), 3)
			},
			want: metricdata.Metrics{
				Name: "sfloat",
				Data: metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[float64]{
						{Value: 3},
					},
				},
			},
		},
		{
			name: "SyncFloat64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				ctr := m.Float64UpDownCounter("sfloat")
				eh.AssertNoErrors(t)

				ctr.Add(context.Background(), 11)
			},
			want: metricdata.Metrics{
				Name: "sfloat",
				Data: metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[float64]{
						{Value: 11},
					},
				},
			},
		},
		{
			name: "SyncFloat64Histogram",
			fn: func(t *testing.T, m metric.Meter) {
				gauge := m.Float64Histogram("histogram")
				eh.AssertNoErrors(t)

				gauge.Record(context.Background(), 7)
			},
			want: metricdata.Metrics{
				Name: "histogram",
				Data: metricdata.Histogram{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint{
						{
							Attributes:   attribute.Set{},
							Count:        1,
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Min:          extrema,
							Max:          extrema,
							Sum:          7.0,
						},
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rdr := NewManualReader()
			m := NewMeterProvider(WithReader(rdr)).Meter("testInstruments")

			tt.fn(t, m)

			rm, err := rdr.Collect(context.Background())
			assert.NoError(t, err)

			require.Len(t, rm.ScopeMetrics, 1)
			sm := rm.ScopeMetrics[0]
			require.Len(t, sm.Metrics, 1)
			got := sm.Metrics[0]
			metricdatatest.AssertEqual(t, tt.want, got, metricdatatest.IgnoreTimestamp())
		})
	}
}

func TestRegisterNonSDKObserverErrors(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))
	meter := mp.Meter("scope")

	type obsrv struct{ instrument.Asynchronous }
	o := obsrv{}

	_ = meter.RegisterCallback(
		func(context.Context, metric.Observer) error { return nil },
		o,
	)
	require.Equal(t, 1, eh.Len(), "external instrument registration should error")
	assert.ErrorContains(
		t,
		eh.Errors()[0],
		"invalid observable: from different implementation",
		"External instrument registred",
	)
}

func TestMeterMixingOnRegisterErrors(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	m2 := mp.Meter("scope2")
	iCtr := m2.Int64ObservableCounter("int64 ctr")
	eh.RequireNoErrors(t)
	fCtr := m2.Float64ObservableCounter("float64 ctr")
	eh.RequireNoErrors(t)
	_ = m1.RegisterCallback(
		func(context.Context, metric.Observer) error { return nil },
		iCtr, fCtr,
	)
	require.Equal(t, 2, eh.Len(), "instrument registration from alt meter should error")
	errs := eh.Errors()
	assert.ErrorContains(
		t,
		errs[0],
		`invalid registration: observable "int64 ctr" from Meter "scope2", registered with Meter "scope1"`,
		"Instrument registred with non-creation Meter",
	)
	assert.ErrorContains(
		t,
		errs[1],
		`invalid registration: observable "float64 ctr" from Meter "scope2", registered with Meter "scope1"`,
		"Instrument registred with non-creation Meter",
	)
}

func TestCallbackObserverNonRegistered(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	valid := m1.Int64ObservableCounter("ctr")
	eh.RequireNoErrors(t)

	m2 := mp.Meter("scope2")
	iCtr := m2.Int64ObservableCounter("int64 ctr")
	eh.RequireNoErrors(t)
	fCtr := m2.Float64ObservableCounter("float64 ctr")
	eh.RequireNoErrors(t)

	type int64Obsrv struct{ instrument.Int64Observable }
	int64Foreign := int64Obsrv{}
	type float64Obsrv struct{ instrument.Float64Observable }
	float64Foreign := float64Obsrv{}

	_ = m1.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(valid, 1)
			o.ObserveInt64(iCtr, 1)
			o.ObserveFloat64(fCtr, 1)
			o.ObserveInt64(int64Foreign, 1)
			o.ObserveFloat64(float64Foreign, 1)
			return nil
		},
		valid,
	)
	eh.RequireNoErrors(t)

	var got metricdata.ResourceMetrics
	var err error
	assert.NotPanics(t, func() {
		got, err = rdr.Collect(context.Background())
	})

	assert.NoError(t, err)
	want := metricdata.ResourceMetrics{
		Resource: resource.Default(),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: "scope1",
				},
				Metrics: []metricdata.Metrics{
					{
						Name: "ctr",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Value: 1,
								},
							},
						},
					},
				},
			},
		},
	}
	metricdatatest.AssertEqual(t, want, got, metricdatatest.IgnoreTimestamp())
}

type logSink struct {
	logr.LogSink

	messages []string
}

func newLogSink(t *testing.T) *logSink {
	return &logSink{LogSink: testr.New(t).GetSink()}
}

func (l *logSink) Info(level int, msg string, keysAndValues ...interface{}) {
	l.messages = append(l.messages, msg)
	l.LogSink.Info(level, msg, keysAndValues...)
}

func (l *logSink) Error(err error, msg string, keysAndValues ...interface{}) {
	l.messages = append(l.messages, fmt.Sprintf("%s: %s", err, msg))
	l.LogSink.Error(err, msg, keysAndValues...)
}

func (l *logSink) String() string {
	out := make([]string, len(l.messages))
	for i := range l.messages {
		out[i] = "\t-" + l.messages[i]
	}
	return strings.Join(out, "\n")
}

func TestGlobalInstRegisterCallback(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	l := newLogSink(t)
	otel.SetLogger(logr.New(l))

	const mtrName = "TestGlobalInstRegisterCallback"
	preMtr := global.Meter(mtrName)
	preInt64Ctr := preMtr.Int64ObservableCounter("pre.int64.counter")
	eh.RequireNoErrors(t)
	preFloat64Ctr := preMtr.Float64ObservableCounter("pre.float64.counter")
	eh.RequireNoErrors(t)

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), WithResource(resource.Empty()))
	global.SetMeterProvider(mp)

	postMtr := global.Meter(mtrName)
	postInt64Ctr := postMtr.Int64ObservableCounter("post.int64.counter")
	eh.RequireNoErrors(t)
	postFloat64Ctr := postMtr.Float64ObservableCounter("post.float64.counter")
	eh.RequireNoErrors(t)

	cb := func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(preInt64Ctr, 1)
		o.ObserveFloat64(preFloat64Ctr, 2)
		o.ObserveInt64(postInt64Ctr, 3)
		o.ObserveFloat64(postFloat64Ctr, 4)
		return nil
	}

	_ = preMtr.RegisterCallback(cb, preInt64Ctr, preFloat64Ctr, postInt64Ctr, postFloat64Ctr)
	eh.AssertNoErrors(t)

	_ = preMtr.RegisterCallback(cb, preInt64Ctr, preFloat64Ctr, postInt64Ctr, postFloat64Ctr)
	eh.AssertNoErrors(t)

	got, err := rdr.Collect(context.Background())
	assert.NoError(t, err)
	assert.Lenf(t, l.messages, 0, "Warnings and errors logged:\n%s", l)
	metricdatatest.AssertEqual(t, metricdata.ResourceMetrics{
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "TestGlobalInstRegisterCallback"},
				Metrics: []metricdata.Metrics{
					{
						Name: "pre.int64.counter",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints:  []metricdata.DataPoint[int64]{{Value: 1}},
						},
					},
					{
						Name: "pre.float64.counter",
						Data: metricdata.Sum[float64]{
							DataPoints:  []metricdata.DataPoint[float64]{{Value: 2}},
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
						},
					},
					{
						Name: "post.int64.counter",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints:  []metricdata.DataPoint[int64]{{Value: 3}},
						},
					},
					{
						Name: "post.float64.counter",
						Data: metricdata.Sum[float64]{
							DataPoints:  []metricdata.DataPoint[float64]{{Value: 4}},
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
						},
					},
				},
			},
		},
	}, got, metricdatatest.IgnoreTimestamp())
}

func TestMetersProvideScope(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	ctr1 := m1.Float64ObservableCounter("ctr1")
	eh.AssertNoErrors(t)
	_ = m1.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveFloat64(ctr1, 5)
		return nil
	}, ctr1)
	eh.AssertNoErrors(t)

	m2 := mp.Meter("scope2")
	ctr2 := m2.Int64ObservableCounter("ctr2")
	eh.AssertNoErrors(t)
	_ = m2.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(ctr2, 7)
		return nil
	}, ctr2)
	eh.AssertNoErrors(t)

	want := metricdata.ResourceMetrics{
		Resource: resource.Default(),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name: "scope1",
				},
				Metrics: []metricdata.Metrics{
					{
						Name: "ctr1",
						Data: metricdata.Sum[float64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Value: 5,
								},
							},
						},
					},
				},
			},
			{
				Scope: instrumentation.Scope{
					Name: "scope2",
				},
				Metrics: []metricdata.Metrics{
					{
						Name: "ctr2",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Value: 7,
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := rdr.Collect(context.Background())
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, want, got, metricdatatest.IgnoreTimestamp())
}

func TestUnregisterUnregisters(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	r := NewManualReader()
	mp := NewMeterProvider(WithReader(r))
	m := mp.Meter("TestUnregisterUnregisters")

	int64Counter := m.Int64ObservableCounter("int64.counter")
	eh.RequireNoErrors(t)

	int64UpDownCounter := m.Int64ObservableUpDownCounter("int64.up_down_counter")
	eh.RequireNoErrors(t)

	int64Gauge := m.Int64ObservableGauge("int64.gauge")
	eh.RequireNoErrors(t)

	floag64Counter := m.Float64ObservableCounter("floag64.counter")
	eh.RequireNoErrors(t)

	floag64UpDownCounter := m.Float64ObservableUpDownCounter("floag64.up_down_counter")
	eh.RequireNoErrors(t)

	floag64Gauge := m.Float64ObservableGauge("floag64.gauge")
	eh.RequireNoErrors(t)

	var called bool
	reg := m.RegisterCallback(
		func(context.Context, metric.Observer) error {
			called = true
			return nil
		},
		int64Counter,
		int64UpDownCounter,
		int64Gauge,
		floag64Counter,
		floag64UpDownCounter,
		floag64Gauge,
	)
	eh.RequireNoErrors(t)

	ctx := context.Background()
	_, err := r.Collect(ctx)
	require.NoError(t, err)
	assert.True(t, called, "callback not called for registered callback")

	called = false
	reg.Unregister()
	eh.RequireNoErrors(t, "unregister")

	_, err = r.Collect(ctx)
	require.NoError(t, err)
	assert.False(t, called, "callback called for unregistered callback")
}

func TestRegisterCallbackDropAggregations(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	aggFn := func(InstrumentKind) aggregation.Aggregation {
		return aggregation.Drop{}
	}
	r := NewManualReader(WithAggregationSelector(aggFn))
	mp := NewMeterProvider(WithReader(r))
	m := mp.Meter("testRegisterCallbackDropAggregations")

	int64Counter := m.Int64ObservableCounter("int64.counter")
	eh.RequireNoErrors(t)

	int64UpDownCounter := m.Int64ObservableUpDownCounter("int64.up_down_counter")
	eh.RequireNoErrors(t)

	int64Gauge := m.Int64ObservableGauge("int64.gauge")
	eh.RequireNoErrors(t)

	floag64Counter := m.Float64ObservableCounter("floag64.counter")
	eh.RequireNoErrors(t)

	floag64UpDownCounter := m.Float64ObservableUpDownCounter("floag64.up_down_counter")
	eh.RequireNoErrors(t)

	floag64Gauge := m.Float64ObservableGauge("floag64.gauge")
	eh.RequireNoErrors(t)

	var called bool
	_ = m.RegisterCallback(
		func(context.Context, metric.Observer) error {
			called = true
			return nil
		},
		int64Counter,
		int64UpDownCounter,
		int64Gauge,
		floag64Counter,
		floag64UpDownCounter,
		floag64Gauge,
	)
	eh.RequireNoErrors(t)

	data, err := r.Collect(context.Background())
	require.NoError(t, err)

	assert.False(t, called, "callback called for all drop instruments")
	assert.Len(t, data.ScopeMetrics, 0, "metrics exported for drop instruments")
}

func TestAttributeFilter(t *testing.T) {
	t.Run("Delta", testAttributeFilter(metricdata.DeltaTemporality))
	t.Run("Cumulative", testAttributeFilter(metricdata.CumulativeTemporality))
}

func testAttributeFilter(temporality metricdata.Temporality) func(*testing.T) {
	testcases := []struct {
		name       string
		register   func(t *testing.T, mtr metric.Meter)
		wantMetric metricdata.Metrics
	}{
		{
			name: "ObservableFloat64Counter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Float64ObservableCounter("afcounter")
				_ = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					o.ObserveFloat64(ctr, 2.0, attribute.String("foo", "bar"))
					o.ObserveFloat64(ctr, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
			},
			wantMetric: metricdata.Metrics{
				Name: "afcounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      4.0,
						},
					},
					Temporality: temporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "ObservableFloat64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Float64ObservableUpDownCounter("afupdowncounter")
				_ = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					o.ObserveFloat64(ctr, 2.0, attribute.String("foo", "bar"))
					o.ObserveFloat64(ctr, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
			},
			wantMetric: metricdata.Metrics{
				Name: "afupdowncounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      4.0,
						},
					},
					Temporality: temporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "ObservableFloat64Gauge",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Float64ObservableGauge("afgauge")
				_ = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					o.ObserveFloat64(ctr, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
			},
			wantMetric: metricdata.Metrics{
				Name: "afgauge",
				Data: metricdata.Gauge[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      2.0,
						},
					},
				},
			},
		},
		{
			name: "ObservableInt64Counter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Int64ObservableCounter("aicounter")
				_ = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					o.ObserveInt64(ctr, 20, attribute.String("foo", "bar"))
					o.ObserveInt64(ctr, 10, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
			},
			wantMetric: metricdata.Metrics{
				Name: "aicounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      40,
						},
					},
					Temporality: temporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "ObservableInt64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Int64ObservableUpDownCounter("aiupdowncounter")
				_ = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					o.ObserveInt64(ctr, 20, attribute.String("foo", "bar"))
					o.ObserveInt64(ctr, 10, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
			},
			wantMetric: metricdata.Metrics{
				Name: "aiupdowncounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      40,
						},
					},
					Temporality: temporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "ObservableInt64Gauge",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Int64ObservableGauge("aigauge")
				_ = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					o.ObserveInt64(ctr, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
			},
			wantMetric: metricdata.Metrics{
				Name: "aigauge",
				Data: metricdata.Gauge[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      20,
						},
					},
				},
			},
		},
		{
			name: "SyncFloat64Counter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Float64Counter("sfcounter")
				ctr.Add(context.Background(), 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
			},
			wantMetric: metricdata.Metrics{
				Name: "sfcounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      3.0,
						},
					},
					Temporality: temporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "SyncFloat64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Float64UpDownCounter("sfupdowncounter")
				ctr.Add(context.Background(), 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
			},
			wantMetric: metricdata.Metrics{
				Name: "sfupdowncounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      3.0,
						},
					},
					Temporality: temporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "SyncFloat64Histogram",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Float64Histogram("sfhistogram")
				ctr.Record(context.Background(), 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Record(context.Background(), 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
			},
			wantMetric: metricdata.Metrics{
				Name: "sfhistogram",
				Data: metricdata.Histogram{
					DataPoints: []metricdata.HistogramDataPoint{
						{
							Attributes:   attribute.NewSet(attribute.String("foo", "bar")),
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Count:        2,
							Min:          metricdata.NewExtrema(1.),
							Max:          metricdata.NewExtrema(2.),
							Sum:          3.0,
						},
					},
					Temporality: temporality,
				},
			},
		},
		{
			name: "SyncInt64Counter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Int64Counter("sicounter")
				ctr.Add(context.Background(), 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
			},
			wantMetric: metricdata.Metrics{
				Name: "sicounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      30,
						},
					},
					Temporality: temporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "SyncInt64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Int64UpDownCounter("siupdowncounter")
				ctr.Add(context.Background(), 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
			},
			wantMetric: metricdata.Metrics{
				Name: "siupdowncounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      30,
						},
					},
					Temporality: temporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "SyncInt64Histogram",
			register: func(t *testing.T, mtr metric.Meter) {
				ctr := mtr.Int64Histogram("sihistogram")
				ctr.Record(context.Background(), 1, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Record(context.Background(), 2, attribute.String("foo", "bar"), attribute.Int("version", 2))
			},
			wantMetric: metricdata.Metrics{
				Name: "sihistogram",
				Data: metricdata.Histogram{
					DataPoints: []metricdata.HistogramDataPoint{
						{
							Attributes:   attribute.NewSet(attribute.String("foo", "bar")),
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Count:        2,
							Min:          metricdata.NewExtrema(1.),
							Max:          metricdata.NewExtrema(2.),
							Sum:          3.0,
						},
					},
					Temporality: temporality,
				},
			},
		},
	}

	return func(t *testing.T) {
		eh := internaltest.NewErrorHandler()
		otel.SetErrorHandler(eh)
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				rdr := NewManualReader(WithTemporalitySelector(func(InstrumentKind) metricdata.Temporality {
					return temporality
				}))
				mtr := NewMeterProvider(
					WithReader(rdr),
					WithView(NewView(
						Instrument{Name: "*"},
						Stream{AttributeFilter: func(kv attribute.KeyValue) bool {
							return kv.Key == attribute.Key("foo")
						}},
					)),
				).Meter("TestAttributeFilter")
				tt.register(t, mtr)
				eh.RequireNoErrors(t)

				m, err := rdr.Collect(context.Background())
				assert.NoError(t, err)

				require.Len(t, m.ScopeMetrics, 1)
				require.Len(t, m.ScopeMetrics[0].Metrics, 1)

				metricdatatest.AssertEqual(t, tt.wantMetric, m.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
			})
		}
	}
}

func TestAsynchronousExample(t *testing.T) {
	// This example can be found:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/8b91585e6175dd52b51e1d60bea105041225e35d/specification/metrics/supplementary-guidelines.md#asynchronous-example
	var (
		threadID1 = attribute.Int("tid", 1)
		threadID2 = attribute.Int("tid", 2)
		threadID3 = attribute.Int("tid", 3)

		processID1001 = attribute.String("pid", "1001")

		thread1 = attribute.NewSet(processID1001, threadID1)
		thread2 = attribute.NewSet(processID1001, threadID2)
		thread3 = attribute.NewSet(processID1001, threadID3)

		process1001 = attribute.NewSet(processID1001)
	)

	setup := func(t *testing.T, temp metricdata.Temporality) (map[attribute.Set]int64, func(*testing.T), *metricdata.ScopeMetrics, *int64, *int64, *int64) {
		t.Helper()

		eh := internaltest.NewErrorHandler()
		otel.SetErrorHandler(eh)

		const (
			instName       = "pageFaults"
			filteredStream = "filteredPageFaults"
			scopeName      = "AsynchronousExample"
		)

		selector := func(InstrumentKind) metricdata.Temporality { return temp }
		reader := NewManualReader(WithTemporalitySelector(selector))

		noopFilter := func(kv attribute.KeyValue) bool { return true }
		noFiltered := NewView(Instrument{Name: instName}, Stream{Name: instName, AttributeFilter: noopFilter})

		filter := func(kv attribute.KeyValue) bool { return kv.Key != "tid" }
		filtered := NewView(Instrument{Name: instName}, Stream{Name: filteredStream, AttributeFilter: filter})

		mp := NewMeterProvider(WithReader(reader), WithView(noFiltered, filtered))
		meter := mp.Meter(scopeName)

		observations := make(map[attribute.Set]int64)
		_ = meter.Int64ObservableCounter(instName, instrument.WithInt64Callback(
			func(_ context.Context, o instrument.Int64Observer) error {
				for attrSet, val := range observations {
					o.Observe(val, attrSet.ToSlice()...)
				}
				return nil
			},
		))
		eh.RequireNoErrors(t)

		want := &metricdata.ScopeMetrics{
			Scope: instrumentation.Scope{Name: scopeName},
			Metrics: []metricdata.Metrics{
				{
					Name: filteredStream,
					Data: metricdata.Sum[int64]{
						Temporality: temp,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{Attributes: process1001},
						},
					},
				},
				{
					Name: instName,
					Data: metricdata.Sum[int64]{
						Temporality: temp,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{Attributes: thread1},
							{Attributes: thread2},
						},
					},
				},
			},
		}
		wantFiltered := &want.Metrics[0].Data.(metricdata.Sum[int64]).DataPoints[0].Value
		wantThread1 := &want.Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[0].Value
		wantThread2 := &want.Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[1].Value

		collect := func(t *testing.T) {
			t.Helper()
			got, err := reader.Collect(context.Background())
			require.NoError(t, err)
			require.Len(t, got.ScopeMetrics, 1)
			metricdatatest.AssertEqual(t, *want, got.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())
		}

		return observations, collect, want, wantFiltered, wantThread1, wantThread2
	}

	t.Run("Cumulative", func(t *testing.T) {
		temporality := metricdata.CumulativeTemporality
		observations, verify, want, wantFiltered, wantThread1, wantThread2 := setup(t, temporality)

		// During the time range (T0, T1]:
		//     pid = 1001, tid = 1, #PF = 50
		//     pid = 1001, tid = 2, #PF = 30
		observations[thread1] = 50
		observations[thread2] = 30

		*wantFiltered = 80
		*wantThread1 = 50
		*wantThread2 = 30

		verify(t)

		// During the time range (T1, T2]:
		//     pid = 1001, tid = 1, #PF = 53
		//     pid = 1001, tid = 2, #PF = 38
		observations[thread1] = 53
		observations[thread2] = 38

		*wantFiltered = 91
		*wantThread1 = 53
		*wantThread2 = 38

		verify(t)

		// During the time range (T2, T3]
		//     pid = 1001, tid = 1, #PF = 56
		//     pid = 1001, tid = 2, #PF = 42
		observations[thread1] = 56
		observations[thread2] = 42

		*wantFiltered = 98
		*wantThread1 = 56
		*wantThread2 = 42

		verify(t)

		// During the time range (T3, T4]:
		//     pid = 1001, tid = 1, #PF = 60
		//     pid = 1001, tid = 2, #PF = 47
		observations[thread1] = 60
		observations[thread2] = 47

		*wantFiltered = 107
		*wantThread1 = 60
		*wantThread2 = 47

		verify(t)

		// During the time range (T4, T5]:
		//     thread 1 died, thread 3 started
		//     pid = 1001, tid = 2, #PF = 53
		//     pid = 1001, tid = 3, #PF = 5
		delete(observations, thread1)
		observations[thread2] = 53
		observations[thread3] = 5

		*wantFiltered = 58
		want.Metrics[1].Data = metricdata.Sum[int64]{
			Temporality: temporality,
			IsMonotonic: true,
			DataPoints: []metricdata.DataPoint[int64]{
				// Thread 1 remains at last measured value.
				{Attributes: thread1, Value: 60},
				{Attributes: thread2, Value: 53},
				{Attributes: thread3, Value: 5},
			},
		}

		verify(t)
	})

	t.Run("Delta", func(t *testing.T) {
		temporality := metricdata.DeltaTemporality
		observations, verify, want, wantFiltered, wantThread1, wantThread2 := setup(t, temporality)

		// During the time range (T0, T1]:
		//     pid = 1001, tid = 1, #PF = 50
		//     pid = 1001, tid = 2, #PF = 30
		observations[thread1] = 50
		observations[thread2] = 30

		*wantFiltered = 80
		*wantThread1 = 50
		*wantThread2 = 30

		verify(t)

		// During the time range (T1, T2]:
		//     pid = 1001, tid = 1, #PF = 53
		//     pid = 1001, tid = 2, #PF = 38
		observations[thread1] = 53
		observations[thread2] = 38

		*wantFiltered = 11
		*wantThread1 = 3
		*wantThread2 = 8

		verify(t)

		// During the time range (T2, T3]
		//     pid = 1001, tid = 1, #PF = 56
		//     pid = 1001, tid = 2, #PF = 42
		observations[thread1] = 56
		observations[thread2] = 42

		*wantFiltered = 7
		*wantThread1 = 3
		*wantThread2 = 4

		verify(t)

		// During the time range (T3, T4]:
		//     pid = 1001, tid = 1, #PF = 60
		//     pid = 1001, tid = 2, #PF = 47
		observations[thread1] = 60
		observations[thread2] = 47

		*wantFiltered = 9
		*wantThread1 = 4
		*wantThread2 = 5

		verify(t)

		// During the time range (T4, T5]:
		//     thread 1 died, thread 3 started
		//     pid = 1001, tid = 2, #PF = 53
		//     pid = 1001, tid = 3, #PF = 5
		delete(observations, thread1)
		observations[thread2] = 53
		observations[thread3] = 5

		*wantFiltered = -49
		want.Metrics[1].Data = metricdata.Sum[int64]{
			Temporality: temporality,
			IsMonotonic: true,
			DataPoints: []metricdata.DataPoint[int64]{
				// Thread 1 remains at last measured value.
				{Attributes: thread1, Value: 0},
				{Attributes: thread2, Value: 6},
				{Attributes: thread3, Value: 5},
			},
		}

		verify(t)
	})
}

var (
	aiCounter       instrument.Int64ObservableCounter
	aiUpDownCounter instrument.Int64ObservableUpDownCounter
	aiGauge         instrument.Int64ObservableGauge

	afCounter       instrument.Float64ObservableCounter
	afUpDownCounter instrument.Float64ObservableUpDownCounter
	afGauge         instrument.Float64ObservableGauge

	siCounter       instrument.Int64Counter
	siUpDownCounter instrument.Int64UpDownCounter
	siHistogram     instrument.Int64Histogram

	sfCounter       instrument.Float64Counter
	sfUpDownCounter instrument.Float64UpDownCounter
	sfHistogram     instrument.Float64Histogram
)

func BenchmarkInstrumentCreation(b *testing.B) {
	provider := NewMeterProvider(WithReader(NewManualReader()))
	meter := provider.Meter("BenchmarkInstrumentCreation")

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		aiCounter = meter.Int64ObservableCounter("observable.int64.counter")
		aiUpDownCounter = meter.Int64ObservableUpDownCounter("observable.int64.up.down.counter")
		aiGauge = meter.Int64ObservableGauge("observable.int64.gauge")

		afCounter = meter.Float64ObservableCounter("observable.float64.counter")
		afUpDownCounter = meter.Float64ObservableUpDownCounter("observable.float64.up.down.counter")
		afGauge = meter.Float64ObservableGauge("observable.float64.gauge")

		siCounter = meter.Int64Counter("sync.int64.counter")
		siUpDownCounter = meter.Int64UpDownCounter("sync.int64.up.down.counter")
		siHistogram = meter.Int64Histogram("sync.int64.histogram")

		sfCounter = meter.Float64Counter("sync.float64.counter")
		sfUpDownCounter = meter.Float64UpDownCounter("sync.float64.up.down.counter")
		sfHistogram = meter.Float64Histogram("sync.float64.histogram")
	}
}
