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
		_ = m.RegisterCallback([]instrument.Asynchronous{}, func(ctx context.Context) {})
		wg.Done()
	}()
	go func() {
		_ = m.RegisterCallback([]instrument.Asynchronous{}, func(ctx context.Context) {})
		wg.Done()
	}()
	wg.Wait()
}

// Instruments should produce correct ResourceMetrics.
func TestMeterCreatesInstruments(t *testing.T) {
	var seven float64 = 7.0
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
				err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
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
				err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
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
				err = m.RegisterCallback([]instrument.Asynchronous{gauge}, func(ctx context.Context) {
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
				err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
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
				err = m.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
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
				err = m.RegisterCallback([]instrument.Asynchronous{gauge}, func(ctx context.Context) {
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
	err = m1.RegisterCallback([]instrument.Asynchronous{ctr1}, func(ctx context.Context) {
		ctr1.Observe(ctx, 5)
	})
	assert.NoError(t, err)

	m2 := mp.Meter("scope2")
	ctr2, err := m2.AsyncInt64().Counter("ctr2")
	assert.NoError(t, err)
	err = m1.RegisterCallback([]instrument.Asynchronous{ctr2}, func(ctx context.Context) {
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
