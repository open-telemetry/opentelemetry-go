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
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/metric/view"
	"go.opentelemetry.io/otel/sdk/resource"
)

// A meter should be able to make instruments concurrently.
func TestMeterInstrumentConcurrency(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(12)

	m := NewMeterProvider().Meter("inst-concurrency")

	go func() {
		_, _ = m.Float64ObservableCounter("AFCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Float64ObservableUpDownCounter("AFUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Float64ObservableGauge("AFGauge")
		wg.Done()
	}()
	go func() {
		_, _ = m.Int64ObservableCounter("AICounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Int64ObservableUpDownCounter("AIUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Int64ObservableGauge("AIGauge")
		wg.Done()
	}()
	go func() {
		_, _ = m.Float64Counter("SFCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Float64UpDownCounter("SFUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Float64Histogram("SFHistogram")
		wg.Done()
	}()
	go func() {
		_, _ = m.Int64Counter("SICounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Int64UpDownCounter("SIUpDownCounter")
		wg.Done()
	}()
	go func() {
		_, _ = m.Int64Histogram("SIHistogram")
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
		_, _ = m.RegisterCallback(func(ctx context.Context) error {
			return nil
		}, nil)
		wg.Done()
	}()
	go func() {
		_, _ = m.RegisterCallback(func(ctx context.Context) error {
			return nil
		}, nil)
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
				ctr, err := m.Int64ObservableCounter("aint")
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 3)
					return nil
				}, ctr)
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
				ctr, err := m.Int64ObservableUpDownCounter("aint")
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 11)
					return nil
				}, ctr)
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
				gauge, err := m.Int64ObservableGauge("agauge")
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(ctx context.Context) error {
					gauge.Observe(ctx, 11)
					return nil
				}, gauge)
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
				ctr, err := m.Float64ObservableCounter("afloat")
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 3)
					return nil
				}, ctr)
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
				ctr, err := m.Float64ObservableUpDownCounter("afloat")
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 11)
					return nil
				}, ctr)
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
				gauge, err := m.Float64ObservableGauge("agauge")
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(ctx context.Context) error {
					gauge.Observe(ctx, 11)
					return nil
				}, gauge)
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
				ctr, err := m.Int64Counter("sint")
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
				ctr, err := m.Int64UpDownCounter("sint")
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
				gauge, err := m.Int64Histogram("histogram")
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
				ctr, err := m.Float64Counter("sfloat")
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
				ctr, err := m.Float64UpDownCounter("sfloat")
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
				gauge, err := m.Float64Histogram("histogram")
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
	ctr1, err := m1.Float64ObservableCounter("ctr1")
	assert.NoError(t, err)
	_, err = m1.RegisterCallback(func(ctx context.Context) error {
		ctr1.Observe(ctx, 5)
		return nil
	}, ctr1)
	assert.NoError(t, err)

	m2 := mp.Meter("scope2")
	ctr2, err := m2.Int64ObservableCounter("ctr2")
	assert.NoError(t, err)
	_, err = m1.RegisterCallback(func(ctx context.Context) error {
		ctr2.Observe(ctx, 7)
		return nil
	}, ctr2)
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

func TestAttributeFilter(t *testing.T) {
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
				ctr, err := mtr.Float64ObservableCounter("afcounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "afcounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      2.0, // TODO (#3439): This should be 3.0.
						},
					},
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "AsyncFloat64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64ObservableUpDownCounter("afupdowncounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "afupdowncounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      2.0, // TODO (#3439): This should be 3.0.
						},
					},
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "AsyncFloat64Gauge",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64ObservableGauge("afgauge")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 1.0, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 2.0, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
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
				ctr, err := mtr.Int64ObservableCounter("aicounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aicounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      20, // TODO (#3439): This should be 30.
						},
					},
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "AsyncInt64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64ObservableUpDownCounter("aiupdowncounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aiupdowncounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(attribute.String("foo", "bar")),
							Value:      20, // TODO (#3439): This should be 30.
						},
					},
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "AsyncInt64Gauge",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64ObservableGauge("aigauge")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(ctx context.Context) error {
					ctr.Observe(ctx, 10, attribute.String("foo", "bar"), attribute.Int("version", 1))
					ctr.Observe(ctx, 20, attribute.String("foo", "bar"), attribute.Int("version", 2))
					return nil
				}, ctr)
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
				ctr, err := mtr.Float64Counter("sfcounter")
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
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "SyncFloat64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64UpDownCounter("sfupdowncounter")
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
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "SyncFloat64Histogram",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64Histogram("sfhistogram")
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
					Temporality: metricdata.CumulativeTemporality,
				},
			},
		},
		{
			name: "SyncInt64Counter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64Counter("sicounter")
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
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "SyncInt64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64UpDownCounter("siupdowncounter")
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
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "SyncInt64Histogram",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64Histogram("sihistogram")
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
					Temporality: metricdata.CumulativeTemporality,
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			v, err := view.New(
				view.MatchInstrumentName("*"),
				view.WithFilterAttributes(attribute.Key("foo")),
			)
			require.NoError(t, err)
			rdr := NewManualReader()
			mtr := NewMeterProvider(
				WithReader(rdr),
				WithView(v),
			).Meter("TestAttributeFilter")

			err = tt.register(t, mtr)
			require.NoError(t, err)

			m, err := rdr.Collect(context.Background())
			assert.NoError(t, err)

			require.Len(t, m.ScopeMetrics, 1)
			require.Len(t, m.ScopeMetrics[0].Metrics, 1)

			metricdatatest.AssertEqual(t, tt.wantMetric, m.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
		})
	}
}
