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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
)

// A meter should be able to make instruments concurrently.
func TestMeterInstrumentConcurrency(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(12)

	m := NewMeterProvider().Meter("inst-concurrency")

	go func() {
		_, _ = m.AsyncFloat64().Counter("AFCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.AsyncFloat64().UpDownCounter("AFUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.AsyncFloat64().Gauge("AFGauge")
		wg.Done()
	}()
	go func() {
		_, _ = m.AsyncInt64().Counter("AICounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.AsyncInt64().UpDownCounter("AIUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.AsyncInt64().Gauge("AIGauge")
		wg.Done()
	}()
	go func() {
		_, _ = m.SyncFloat64().Counter("SFCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.SyncFloat64().UpDownCounter("SFUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.SyncFloat64().Histogram("SFHistogram")
		wg.Done()
	}()
	go func() {
		_, _ = m.SyncInt64().Counter("SICounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.SyncInt64().UpDownCounter("SIUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.SyncInt64().Histogram("SIHistogram")
		wg.Done()
	}()

	wg.Wait()
}

// A Meter Should be able register Callbacks Concurrently.
func TestMeterCallbackCreationConcurrency(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	m := NewMeterProvider().Meter("callback-concurrency")

	go func() {
		_, _ = m.RegisterCallback([]instrument.Asynchronous{}, func(ctx context.Context) {})
		wg.Done()
	}()
	go func() {
		_, _ = m.RegisterCallback([]instrument.Asynchronous{}, func(ctx context.Context) {})
		wg.Done()
	}()
	wg.Wait()
}

func TestNoopCallbackUnregisterConcurrency(t *testing.T) {
	m := NewMeterProvider().Meter("noop-unregister-concurrency")
	reg, err := m.RegisterCallback(nil, func(ctx context.Context) {})
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		_ = reg.Unregister()
		wg.Done()
	}()
	go func() {
		_ = reg.Unregister()
		wg.Done()
	}()
	wg.Wait()
}

func TestCallbackUnregisterConcurrency(t *testing.T) {
	reader := NewManualReader()
	provider := NewMeterProvider(WithReader(reader))
	meter := provider.Meter("unregister-concurrency")

	actr, err := meter.AsyncFloat64().Counter("counter")
	require.NoError(t, err)

	ag, err := meter.AsyncInt64().Gauge("gauge")
	require.NoError(t, err)

	i := []instrument.Asynchronous{actr}
	regCtr, err := meter.RegisterCallback(i, func(ctx context.Context) {})
	require.NoError(t, err)

	i = []instrument.Asynchronous{ag}
	regG, err := meter.RegisterCallback(i, func(ctx context.Context) {})
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		_ = regCtr.Unregister()
		_ = regG.Unregister()
		wg.Done()
	}()
	go func() {
		_ = regCtr.Unregister()
		_ = regG.Unregister()
		wg.Done()
	}()
	wg.Wait()
}

// Instruments should produce correct ResourceMetrics.
func TestMeterCreatesInstruments(t *testing.T) {
	seven := 7.0
	testCases := []struct {
		name string
		fn   func(*testing.T, metric.Meter)
		want metricdata.Metrics
	}{
		{
			name: "AsyncInt64Count",
			fn: func(t *testing.T, m metric.Meter) {
				ctr, err := m.AsyncInt64().Counter("aint")
				assert.NoError(t, err)
				_, err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 3)
				})
				assert.NoError(t, err)

				// Observed outside of a callback, it should be ignored.
				ctr.Observe(context.Background(), 19)
			},
			want: metricdata.Metrics{
				Name: "aint",
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
			name: "AsyncInt64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				ctr, err := m.AsyncInt64().UpDownCounter("aint")
				assert.NoError(t, err)
				_, err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 11)
				})
				assert.NoError(t, err)

				// Observed outside of a callback, it should be ignored.
				ctr.Observe(context.Background(), 19)
			},
			want: metricdata.Metrics{
				Name: "aint",
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
			name: "AsyncInt64Gauge",
			fn: func(t *testing.T, m metric.Meter) {
				gauge, err := m.AsyncInt64().Gauge("agauge")
				assert.NoError(t, err)
				_, err = m.RegisterCallback([]instrument.Asynchronous{gauge}, func(ctx context.Context) {
					gauge.Observe(ctx, 11)
				})
				assert.NoError(t, err)

				// Observed outside of a callback, it should be ignored.
				gauge.Observe(context.Background(), 19)
			},
			want: metricdata.Metrics{
				Name: "agauge",
				Data: metricdata.Gauge[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Value: 11},
					},
				},
			},
		},
		{
			name: "AsyncFloat64Count",
			fn: func(t *testing.T, m metric.Meter) {
				ctr, err := m.AsyncFloat64().Counter("afloat")
				assert.NoError(t, err)
				_, err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 3)
				})
				assert.NoError(t, err)

				// Observed outside of a callback, it should be ignored.
				ctr.Observe(context.Background(), 19)
			},
			want: metricdata.Metrics{
				Name: "afloat",
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
			name: "AsyncFloat64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				ctr, err := m.AsyncFloat64().UpDownCounter("afloat")
				assert.NoError(t, err)
				_, err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 11)
				})
				assert.NoError(t, err)

				// Observed outside of a callback, it should be ignored.
				ctr.Observe(context.Background(), 19)
			},
			want: metricdata.Metrics{
				Name: "afloat",
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
			name: "AsyncFloat64Gauge",
			fn: func(t *testing.T, m metric.Meter) {
				gauge, err := m.AsyncFloat64().Gauge("agauge")
				assert.NoError(t, err)
				_, err = m.RegisterCallback([]instrument.Asynchronous{gauge}, func(ctx context.Context) {
					gauge.Observe(ctx, 11)
				})
				assert.NoError(t, err)

				// Observed outside of a callback, it should be ignored.
				gauge.Observe(context.Background(), 19)
			},
			want: metricdata.Metrics{
				Name: "agauge",
				Data: metricdata.Gauge[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Value: 11},
					},
				},
			},
		},

		{
			name: "SyncInt64Count",
			fn: func(t *testing.T, m metric.Meter) {
				ctr, err := m.SyncInt64().Counter("sint")
				assert.NoError(t, err)

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
				ctr, err := m.SyncInt64().UpDownCounter("sint")
				assert.NoError(t, err)

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
				gauge, err := m.SyncInt64().Histogram("histogram")
				assert.NoError(t, err)

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
							Min:          &seven,
							Max:          &seven,
							Sum:          7.0,
						},
					},
				},
			},
		},
		{
			name: "SyncFloat64Count",
			fn: func(t *testing.T, m metric.Meter) {
				ctr, err := m.SyncFloat64().Counter("sfloat")
				assert.NoError(t, err)

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
				ctr, err := m.SyncFloat64().UpDownCounter("sfloat")
				assert.NoError(t, err)

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
				gauge, err := m.SyncFloat64().Histogram("histogram")
				assert.NoError(t, err)

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
							Min:          &seven,
							Max:          &seven,
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

func TestMetersProvideScope(t *testing.T) {
	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	ctr1, err := m1.AsyncFloat64().Counter("ctr1")
	assert.NoError(t, err)
	_, err = m1.RegisterCallback([]instrument.Asynchronous{ctr1}, func(ctx context.Context) {
		ctr1.Observe(ctx, 5)
	})
	assert.NoError(t, err)

	m2 := mp.Meter("scope2")
	ctr2, err := m2.AsyncInt64().Counter("ctr2")
	assert.NoError(t, err)
	_, err = m1.RegisterCallback([]instrument.Asynchronous{ctr2}, func(ctx context.Context) {
		ctr2.Observe(ctx, 7)
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
	r := NewManualReader()
	mp := NewMeterProvider(WithReader(r))
	m := mp.Meter("TestUnregisterUnregisters")

	int64Counter, err := m.AsyncInt64().Counter("int64.counter")
	require.NoError(t, err)

	int64UpDownCounter, err := m.AsyncInt64().UpDownCounter("int64.up_down_counter")
	require.NoError(t, err)

	int64Gauge, err := m.AsyncInt64().Gauge("int64.gauge")
	require.NoError(t, err)

	floag64Counter, err := m.AsyncFloat64().Counter("floag64.counter")
	require.NoError(t, err)

	floag64UpDownCounter, err := m.AsyncFloat64().UpDownCounter("floag64.up_down_counter")
	require.NoError(t, err)

	floag64Gauge, err := m.AsyncFloat64().Gauge("floag64.gauge")
	require.NoError(t, err)

	var called bool
	reg, err := m.RegisterCallback([]instrument.Asynchronous{
		int64Counter,
		int64UpDownCounter,
		int64Gauge,
		floag64Counter,
		floag64UpDownCounter,
		floag64Gauge,
	}, func(context.Context) { called = true })
	require.NoError(t, err)

	ctx := context.Background()
	_, err = r.Collect(ctx)
	require.NoError(t, err)
	assert.True(t, called, "callback not called for registered callback")

	called = false
	require.NoError(t, reg.Unregister(), "unregister")

	_, err = r.Collect(ctx)
	require.NoError(t, err)
	assert.False(t, called, "callback called for unregistered callback")
}

func TestRegisterCallbackDropAggregations(t *testing.T) {
	aggFn := func(InstrumentKind) aggregation.Aggregation {
		return aggregation.Drop{}
	}
	r := NewManualReader(WithAggregationSelector(aggFn))
	mp := NewMeterProvider(WithReader(r))
	m := mp.Meter("testRegisterCallbackDropAggregations")

	int64Counter, err := m.AsyncInt64().Counter("int64.counter")
	require.NoError(t, err)

	int64UpDownCounter, err := m.AsyncInt64().UpDownCounter("int64.up_down_counter")
	require.NoError(t, err)

	int64Gauge, err := m.AsyncInt64().Gauge("int64.gauge")
	require.NoError(t, err)

	floag64Counter, err := m.AsyncFloat64().Counter("floag64.counter")
	require.NoError(t, err)

	floag64UpDownCounter, err := m.AsyncFloat64().UpDownCounter("floag64.up_down_counter")
	require.NoError(t, err)

	floag64Gauge, err := m.AsyncFloat64().Gauge("floag64.gauge")
	require.NoError(t, err)

	var called bool
	_, err = m.RegisterCallback([]instrument.Asynchronous{
		int64Counter,
		int64UpDownCounter,
		int64Gauge,
		floag64Counter,
		floag64UpDownCounter,
		floag64Gauge,
	}, func(context.Context) { called = true })
	require.NoError(t, err)

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
	one := 1.0
	two := 2.0
	testcases := []struct {
		name       string
		register   func(t *testing.T, mtr metric.Meter) error
		wantMetric metricdata.Metrics
	}{
		{
			name: "AsyncFloat64Counter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.AsyncFloat64().Counter("afcounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
				})
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "afcounter",
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
			name: "AsyncFloat64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.AsyncFloat64().UpDownCounter("afupdowncounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
				})
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "afupdowncounter",
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
			name: "AsyncFloat64Gauge",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.AsyncFloat64().Gauge("afgauge")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
				})
				return err
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
			name: "AsyncInt64Counter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.AsyncInt64().Counter("aicounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
				})
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aicounter",
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
			name: "AsyncInt64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.AsyncInt64().UpDownCounter("aiupdowncounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
				})
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aiupdowncounter",
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
			name: "AsyncInt64Gauge",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.AsyncInt64().Gauge("aigauge")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
					ctr.Observe(ctx, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
				})
				return err
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.SyncFloat64().Counter("sfcounter")
				if err != nil {
					return err
				}

				ctr.Add(context.Background(), 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
				return nil
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.SyncFloat64().UpDownCounter("sfupdowncounter")
				if err != nil {
					return err
				}

				ctr.Add(context.Background(), 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
				return nil
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.SyncFloat64().Histogram("sfhistogram")
				if err != nil {
					return err
				}

				ctr.Record(context.Background(), 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Record(context.Background(), 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
				return nil
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
							Min:          &one,
							Max:          &two,
							Sum:          3.0,
						},
					},
					Temporality: temporality,
				},
			},
		},
		{
			name: "SyncInt64Counter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.SyncInt64().Counter("sicounter")
				if err != nil {
					return err
				}

				ctr.Add(context.Background(), 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
				return nil
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.SyncInt64().UpDownCounter("siupdowncounter")
				if err != nil {
					return err
				}

				ctr.Add(context.Background(), 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Add(context.Background(), 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
				return nil
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.SyncInt64().Histogram("sihistogram")
				if err != nil {
					return err
				}

				ctr.Record(context.Background(), 1, attribute.String("foo", "bar"), attribute.Int("version", 1))
				ctr.Record(context.Background(), 2, attribute.String("foo", "bar"), attribute.Int("version", 2))
				return nil
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
							Min:          &one,
							Max:          &two,
							Sum:          3.0,
						},
					},
					Temporality: temporality,
				},
			},
		},
	}

	return func(t *testing.T) {
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
				require.NoError(t, tt.register(t, mtr))

				m, err := rdr.Collect(context.Background())
				assert.NoError(t, err)

				require.Len(t, m.ScopeMetrics, 1)
				require.Len(t, m.ScopeMetrics[0].Metrics, 1)

				metricdatatest.AssertEqual(t, tt.wantMetric, m.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
			})
		}
	}
}

var (
	aiCounter       asyncint64.Counter
	aiUpDownCounter asyncint64.UpDownCounter
	aiGauge         asyncint64.Gauge

	afCounter       asyncfloat64.Counter
	afUpDownCounter asyncfloat64.UpDownCounter
	afGauge         asyncfloat64.Gauge

	siCounter       syncint64.Counter
	siUpDownCounter syncint64.UpDownCounter
	siHistogram     syncint64.Histogram

	sfCounter       syncfloat64.Counter
	sfUpDownCounter syncfloat64.UpDownCounter
	sfHistogram     syncfloat64.Histogram
)

func BenchmarkInstrumentCreation(b *testing.B) {
	provider := NewMeterProvider(WithReader(NewManualReader()))
	meter := provider.Meter("BenchmarkInstrumentCreation")

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		aiCounter, _ = meter.AsyncInt64().Counter("async.int64.counter")
		aiUpDownCounter, _ = meter.AsyncInt64().UpDownCounter("async.int64.up.down.counter")
		aiGauge, _ = meter.AsyncInt64().Gauge("async.int64.gauge")

		afCounter, _ = meter.AsyncFloat64().Counter("async.float64.counter")
		afUpDownCounter, _ = meter.AsyncFloat64().UpDownCounter("async.float64.up.down.counter")
		afGauge, _ = meter.AsyncFloat64().Gauge("async.float64.gauge")

		siCounter, _ = meter.SyncInt64().Counter("sync.int64.counter")
		siUpDownCounter, _ = meter.SyncInt64().UpDownCounter("sync.int64.up.down.counter")
		siHistogram, _ = meter.SyncInt64().Histogram("sync.int64.histogram")

		sfCounter, _ = meter.SyncFloat64().Counter("sync.float64.counter")
		sfUpDownCounter, _ = meter.SyncFloat64().UpDownCounter("sync.float64.up.down.counter")
		sfHistogram, _ = meter.SyncFloat64().Histogram("sync.float64.histogram")
	}
}
