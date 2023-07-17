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
	"go.opentelemetry.io/otel/metric"
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

var emptyCallback metric.Callback = func(context.Context, metric.Observer) error { return nil }

// A Meter Should be able register Callbacks Concurrently.
func TestMeterCallbackCreationConcurrency(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	m := NewMeterProvider().Meter("callback-concurrency")

	go func() {
		_, _ = m.RegisterCallback(emptyCallback)
		wg.Done()
	}()
	go func() {
		_, _ = m.RegisterCallback(emptyCallback)
		wg.Done()
	}()
	wg.Wait()
}

func TestNoopCallbackUnregisterConcurrency(t *testing.T) {
	m := NewMeterProvider().Meter("noop-unregister-concurrency")
	reg, err := m.RegisterCallback(emptyCallback)
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

	actr, err := meter.Float64ObservableCounter("counter")
	require.NoError(t, err)

	ag, err := meter.Int64ObservableGauge("gauge")
	require.NoError(t, err)

	regCtr, err := meter.RegisterCallback(emptyCallback, actr)
	require.NoError(t, err)

	regG, err := meter.RegisterCallback(emptyCallback, ag)
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
	attrs := attribute.NewSet(attribute.String("name", "alice"))
	opt := metric.WithAttributeSet(attrs)
	testCases := []struct {
		name string
		fn   func(*testing.T, metric.Meter)
		want metricdata.Metrics
	}{
		{
			name: "ObservableInt64Count",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o metric.Int64Observer) error {
					o.Observe(4, opt)
					return nil
				}
				ctr, err := m.Int64ObservableCounter("aint", metric.WithInt64Callback(cback))
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 3)
					return nil
				}, ctr)
				assert.NoError(t, err)
			},
			want: metricdata.Metrics{
				Name: "aint",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attrs, Value: 4},
						{Value: 3},
					},
				},
			},
		},
		{
			name: "ObservableInt64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o metric.Int64Observer) error {
					o.Observe(4, opt)
					return nil
				}
				ctr, err := m.Int64ObservableUpDownCounter("aint", metric.WithInt64Callback(cback))
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 11)
					return nil
				}, ctr)
				assert.NoError(t, err)
			},
			want: metricdata.Metrics{
				Name: "aint",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attrs, Value: 4},
						{Value: 11},
					},
				},
			},
		},
		{
			name: "ObservableInt64Gauge",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o metric.Int64Observer) error {
					o.Observe(4, opt)
					return nil
				}
				gauge, err := m.Int64ObservableGauge("agauge", metric.WithInt64Callback(cback))
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(gauge, 11)
					return nil
				}, gauge)
				assert.NoError(t, err)
			},
			want: metricdata.Metrics{
				Name: "agauge",
				Data: metricdata.Gauge[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: attrs, Value: 4},
						{Value: 11},
					},
				},
			},
		},
		{
			name: "ObservableFloat64Count",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o metric.Float64Observer) error {
					o.Observe(4, opt)
					return nil
				}
				ctr, err := m.Float64ObservableCounter("afloat", metric.WithFloat64Callback(cback))
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 3)
					return nil
				}, ctr)
				assert.NoError(t, err)
			},
			want: metricdata.Metrics{
				Name: "afloat",
				Data: metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: attrs, Value: 4},
						{Value: 3},
					},
				},
			},
		},
		{
			name: "ObservableFloat64UpDownCount",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o metric.Float64Observer) error {
					o.Observe(4, opt)
					return nil
				}
				ctr, err := m.Float64ObservableUpDownCounter("afloat", metric.WithFloat64Callback(cback))
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 11)
					return nil
				}, ctr)
				assert.NoError(t, err)
			},
			want: metricdata.Metrics{
				Name: "afloat",
				Data: metricdata.Sum[float64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: attrs, Value: 4},
						{Value: 11},
					},
				},
			},
		},
		{
			name: "ObservableFloat64Gauge",
			fn: func(t *testing.T, m metric.Meter) {
				cback := func(_ context.Context, o metric.Float64Observer) error {
					o.Observe(4, opt)
					return nil
				}
				gauge, err := m.Float64ObservableGauge("agauge", metric.WithFloat64Callback(cback))
				assert.NoError(t, err)
				_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(gauge, 11)
					return nil
				}, gauge)
				assert.NoError(t, err)
			},
			want: metricdata.Metrics{
				Name: "agauge",
				Data: metricdata.Gauge[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: attrs, Value: 4},
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
				Data: metricdata.Histogram[int64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[int64]{
						{
							Attributes:   attribute.Set{},
							Count:        1,
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Min:          metricdata.NewExtrema[int64](7),
							Max:          metricdata.NewExtrema[int64](7),
							Sum:          7,
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
				Data: metricdata.Histogram[float64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Attributes:   attribute.Set{},
							Count:        1,
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Min:          metricdata.NewExtrema[float64](7.),
							Max:          metricdata.NewExtrema[float64](7.),
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

			rm := metricdata.ResourceMetrics{}
			err := rdr.Collect(context.Background(), &rm)
			assert.NoError(t, err)

			require.Len(t, rm.ScopeMetrics, 1)
			sm := rm.ScopeMetrics[0]
			require.Len(t, sm.Metrics, 1)
			got := sm.Metrics[0]
			metricdatatest.AssertEqual(t, tt.want, got, metricdatatest.IgnoreTimestamp())
		})
	}
}

func TestMeterCreatesInstrumentsValidations(t *testing.T) {
	testCases := []struct {
		name string
		fn   func(*testing.T, metric.Meter) error

		wantErr error
	}{
		{
			name: "Int64Counter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64Counter("counter")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Int64Counter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64Counter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Int64UpDownCounter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64UpDownCounter("upDownCounter")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Int64UpDownCounter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64UpDownCounter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Int64Histogram with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64Histogram("histogram")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Int64Histogram with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64Histogram("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Int64ObservableCounter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableCounter("aint")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Int64ObservableCounter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableCounter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Int64ObservableUpDownCounter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableUpDownCounter("aint")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Int64ObservableUpDownCounter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableUpDownCounter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Int64ObservableGauge with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableGauge("aint")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Int64ObservableGauge with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableGauge("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Float64Counter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64Counter("counter")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Float64Counter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64Counter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Float64UpDownCounter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64UpDownCounter("upDownCounter")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Float64UpDownCounter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64UpDownCounter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Float64Histogram with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64Histogram("histogram")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Float64Histogram with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64Histogram("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Float64ObservableCounter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64ObservableCounter("aint")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Float64ObservableCounter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Int64ObservableCounter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Float64ObservableUpDownCounter with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64ObservableUpDownCounter("aint")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Float64ObservableUpDownCounter with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64ObservableUpDownCounter("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
		{
			name: "Float64ObservableGauge with no validation issues",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64ObservableGauge("aint")
				assert.NotNil(t, i)
				return err
			},
		},
		{
			name: "Float64ObservableGauge with an invalid name",

			fn: func(t *testing.T, m metric.Meter) error {
				i, err := m.Float64ObservableGauge("_")
				assert.NotNil(t, i)
				return err
			},

			wantErr: fmt.Errorf("%w: _: must start with a letter", ErrInstrumentName),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMeterProvider().Meter("testInstruments")
			err := tt.fn(t, m)
			assert.Equal(t, err, tt.wantErr)
		})
	}
}

func TestValidateInstrumentName(t *testing.T) {
	testCases := []struct {
		name string

		wantErr error
	}{
		{
			name:    "",
			wantErr: fmt.Errorf("%w: : is empty", ErrInstrumentName),
		},
		{
			name:    "1",
			wantErr: fmt.Errorf("%w: 1: must start with a letter", ErrInstrumentName),
		},
		{
			name: "a",
		},
		{
			name: "n4me",
		},
		{
			name: "n-me",
		},
		{
			name: "na_e",
		},
		{
			name: "nam.",
		},
		{
			name:    "name!",
			wantErr: fmt.Errorf("%w: name!: must only contain [A-Za-z0-9_.-]", ErrInstrumentName),
		},
		{
			name:    "someverylongnamewhichisover63charactersbutallofwhichmatchtheregexp",
			wantErr: fmt.Errorf("%w: someverylongnamewhichisover63charactersbutallofwhichmatchtheregexp: longer than 63 characters", ErrInstrumentName),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, validateInstrumentName(tt.name))
		})
	}
}

func TestRegisterNonSDKObserverErrors(t *testing.T) {
	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))
	meter := mp.Meter("scope")

	type obsrv struct{ metric.Observable }
	o := obsrv{}

	_, err := meter.RegisterCallback(
		func(context.Context, metric.Observer) error { return nil },
		o,
	)
	assert.ErrorContains(
		t,
		err,
		"invalid observable: from different implementation",
		"External instrument registered",
	)
}

func TestMeterMixingOnRegisterErrors(t *testing.T) {
	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	m2 := mp.Meter("scope2")
	iCtr, err := m2.Int64ObservableCounter("int64ctr")
	require.NoError(t, err)
	fCtr, err := m2.Float64ObservableCounter("float64ctr")
	require.NoError(t, err)
	_, err = m1.RegisterCallback(
		func(context.Context, metric.Observer) error { return nil },
		iCtr, fCtr,
	)
	assert.ErrorContains(
		t,
		err,
		`invalid registration: observable "int64ctr" from Meter "scope2", registered with Meter "scope1"`,
		"Instrument registered with non-creation Meter",
	)
	assert.ErrorContains(
		t,
		err,
		`invalid registration: observable "float64ctr" from Meter "scope2", registered with Meter "scope1"`,
		"Instrument registered with non-creation Meter",
	)
}

func TestCallbackObserverNonRegistered(t *testing.T) {
	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	valid, err := m1.Int64ObservableCounter("ctr")
	require.NoError(t, err)

	m2 := mp.Meter("scope2")
	iCtr, err := m2.Int64ObservableCounter("int64ctr")
	require.NoError(t, err)
	fCtr, err := m2.Float64ObservableCounter("float64ctr")
	require.NoError(t, err)

	type int64Obsrv struct{ metric.Int64Observable }
	int64Foreign := int64Obsrv{}
	type float64Obsrv struct{ metric.Float64Observable }
	float64Foreign := float64Obsrv{}

	_, err = m1.RegisterCallback(
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
	require.NoError(t, err)

	var got metricdata.ResourceMetrics
	assert.NotPanics(t, func() {
		err = rdr.Collect(context.Background(), &got)
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
	l := newLogSink(t)
	otel.SetLogger(logr.New(l))

	const mtrName = "TestGlobalInstRegisterCallback"
	preMtr := otel.Meter(mtrName)
	preInt64Ctr, err := preMtr.Int64ObservableCounter("pre.int64.counter")
	require.NoError(t, err)
	preFloat64Ctr, err := preMtr.Float64ObservableCounter("pre.float64.counter")
	require.NoError(t, err)

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), WithResource(resource.Empty()))
	otel.SetMeterProvider(mp)

	postMtr := otel.Meter(mtrName)
	postInt64Ctr, err := postMtr.Int64ObservableCounter("post.int64.counter")
	require.NoError(t, err)
	postFloat64Ctr, err := postMtr.Float64ObservableCounter("post.float64.counter")
	require.NoError(t, err)

	cb := func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(preInt64Ctr, 1)
		o.ObserveFloat64(preFloat64Ctr, 2)
		o.ObserveInt64(postInt64Ctr, 3)
		o.ObserveFloat64(postFloat64Ctr, 4)
		return nil
	}

	_, err = preMtr.RegisterCallback(cb, preInt64Ctr, preFloat64Ctr, postInt64Ctr, postFloat64Ctr)
	assert.NoError(t, err)

	_, err = preMtr.RegisterCallback(cb, preInt64Ctr, preFloat64Ctr, postInt64Ctr, postFloat64Ctr)
	assert.NoError(t, err)

	got := metricdata.ResourceMetrics{}
	err = rdr.Collect(context.Background(), &got)
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
	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr))

	m1 := mp.Meter("scope1")
	ctr1, err := m1.Float64ObservableCounter("ctr1")
	assert.NoError(t, err)
	_, err = m1.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveFloat64(ctr1, 5)
		return nil
	}, ctr1)
	assert.NoError(t, err)

	m2 := mp.Meter("scope2")
	ctr2, err := m2.Int64ObservableCounter("ctr2")
	assert.NoError(t, err)
	_, err = m2.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(ctr2, 7)
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

	got := metricdata.ResourceMetrics{}
	err = rdr.Collect(context.Background(), &got)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, want, got, metricdatatest.IgnoreTimestamp())
}

func TestUnregisterUnregisters(t *testing.T) {
	r := NewManualReader()
	mp := NewMeterProvider(WithReader(r))
	m := mp.Meter("TestUnregisterUnregisters")

	int64Counter, err := m.Int64ObservableCounter("int64.counter")
	require.NoError(t, err)

	int64UpDownCounter, err := m.Int64ObservableUpDownCounter("int64.up_down_counter")
	require.NoError(t, err)

	int64Gauge, err := m.Int64ObservableGauge("int64.gauge")
	require.NoError(t, err)

	floag64Counter, err := m.Float64ObservableCounter("floag64.counter")
	require.NoError(t, err)

	floag64UpDownCounter, err := m.Float64ObservableUpDownCounter("floag64.up_down_counter")
	require.NoError(t, err)

	floag64Gauge, err := m.Float64ObservableGauge("floag64.gauge")
	require.NoError(t, err)

	var called bool
	reg, err := m.RegisterCallback(
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
	require.NoError(t, err)

	ctx := context.Background()
	err = r.Collect(ctx, &metricdata.ResourceMetrics{})
	require.NoError(t, err)
	assert.True(t, called, "callback not called for registered callback")

	called = false
	require.NoError(t, reg.Unregister(), "unregister")

	err = r.Collect(ctx, &metricdata.ResourceMetrics{})
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

	int64Counter, err := m.Int64ObservableCounter("int64.counter")
	require.NoError(t, err)

	int64UpDownCounter, err := m.Int64ObservableUpDownCounter("int64.up_down_counter")
	require.NoError(t, err)

	int64Gauge, err := m.Int64ObservableGauge("int64.gauge")
	require.NoError(t, err)

	floag64Counter, err := m.Float64ObservableCounter("floag64.counter")
	require.NoError(t, err)

	floag64UpDownCounter, err := m.Float64ObservableUpDownCounter("floag64.up_down_counter")
	require.NoError(t, err)

	floag64Gauge, err := m.Float64ObservableGauge("floag64.gauge")
	require.NoError(t, err)

	var called bool
	_, err = m.RegisterCallback(
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
	require.NoError(t, err)

	data := metricdata.ResourceMetrics{}
	err = r.Collect(context.Background(), &data)
	require.NoError(t, err)

	assert.False(t, called, "callback called for all drop instruments")
	assert.Len(t, data.ScopeMetrics, 0, "metrics exported for drop instruments")
}

func TestAttributeFilter(t *testing.T) {
	t.Run("Delta", testAttributeFilter(metricdata.DeltaTemporality))
	t.Run("Cumulative", testAttributeFilter(metricdata.CumulativeTemporality))
}

func testAttributeFilter(temporality metricdata.Temporality) func(*testing.T) {
	fooBar := attribute.NewSet(attribute.String("foo", "bar"))
	withFooBar := metric.WithAttributeSet(fooBar)
	v1 := attribute.NewSet(attribute.String("foo", "bar"), attribute.Int("version", 1))
	withV1 := metric.WithAttributeSet(v1)
	v2 := attribute.NewSet(attribute.String("foo", "bar"), attribute.Int("version", 2))
	withV2 := metric.WithAttributeSet(v2)
	testcases := []struct {
		name       string
		register   func(t *testing.T, mtr metric.Meter) error
		wantMetric metricdata.Metrics
	}{
		{
			name: "ObservableFloat64Counter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64ObservableCounter("afcounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 1.0, withV1)
					o.ObserveFloat64(ctr, 2.0, withFooBar)
					o.ObserveFloat64(ctr, 1.0, withV2)
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "afcounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: fooBar, Value: 4.0},
					},
					Temporality: temporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "ObservableFloat64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64ObservableUpDownCounter("afupdowncounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 1.0, withV1)
					o.ObserveFloat64(ctr, 2.0, withFooBar)
					o.ObserveFloat64(ctr, 1.0, withV2)
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Float64ObservableGauge("afgauge")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveFloat64(ctr, 1.0, withV1)
					o.ObserveFloat64(ctr, 2.0, withV2)
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "afgauge",
				Data: metricdata.Gauge[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: fooBar, Value: 2.0},
					},
				},
			},
		},
		{
			name: "ObservableInt64Counter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64ObservableCounter("aicounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 10, withV1)
					o.ObserveInt64(ctr, 20, withFooBar)
					o.ObserveInt64(ctr, 10, withV2)
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aicounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: fooBar, Value: 40},
					},
					Temporality: temporality,
					IsMonotonic: true,
				},
			},
		},
		{
			name: "ObservableInt64UpDownCounter",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64ObservableUpDownCounter("aiupdowncounter")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 10, withV1)
					o.ObserveInt64(ctr, 20, withFooBar)
					o.ObserveInt64(ctr, 10, withV2)
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aiupdowncounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: fooBar, Value: 40},
					},
					Temporality: temporality,
					IsMonotonic: false,
				},
			},
		},
		{
			name: "ObservableInt64Gauge",
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64ObservableGauge("aigauge")
				if err != nil {
					return err
				}
				_, err = mtr.RegisterCallback(func(_ context.Context, o metric.Observer) error {
					o.ObserveInt64(ctr, 10, withV1)
					o.ObserveInt64(ctr, 20, withV2)
					return nil
				}, ctr)
				return err
			},
			wantMetric: metricdata.Metrics{
				Name: "aigauge",
				Data: metricdata.Gauge[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: fooBar, Value: 20},
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

				ctr.Add(context.Background(), 1.0, withV1)
				ctr.Add(context.Background(), 2.0, withV2)
				return nil
			},
			wantMetric: metricdata.Metrics{
				Name: "sfcounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: fooBar, Value: 3.0},
					},
					Temporality: temporality,
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

				ctr.Add(context.Background(), 1.0, withV1)
				ctr.Add(context.Background(), 2.0, withV2)
				return nil
			},
			wantMetric: metricdata.Metrics{
				Name: "sfupdowncounter",
				Data: metricdata.Sum[float64]{
					DataPoints: []metricdata.DataPoint[float64]{
						{Attributes: fooBar, Value: 3.0},
					},
					Temporality: temporality,
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

				ctr.Record(context.Background(), 1.0, withV1)
				ctr.Record(context.Background(), 2.0, withV2)
				return nil
			},
			wantMetric: metricdata.Metrics{
				Name: "sfhistogram",
				Data: metricdata.Histogram[float64]{
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Attributes:   fooBar,
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
			register: func(t *testing.T, mtr metric.Meter) error {
				ctr, err := mtr.Int64Counter("sicounter")
				if err != nil {
					return err
				}

				ctr.Add(context.Background(), 10, withV1)
				ctr.Add(context.Background(), 20, withV2)
				return nil
			},
			wantMetric: metricdata.Metrics{
				Name: "sicounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: fooBar, Value: 30},
					},
					Temporality: temporality,
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

				ctr.Add(context.Background(), 10, withV1)
				ctr.Add(context.Background(), 20, withV2)
				return nil
			},
			wantMetric: metricdata.Metrics{
				Name: "siupdowncounter",
				Data: metricdata.Sum[int64]{
					DataPoints: []metricdata.DataPoint[int64]{
						{Attributes: fooBar, Value: 30},
					},
					Temporality: temporality,
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

				ctr.Record(context.Background(), 1, withV1)
				ctr.Record(context.Background(), 2, withV2)
				return nil
			},
			wantMetric: metricdata.Metrics{
				Name: "sihistogram",
				Data: metricdata.Histogram[int64]{
					DataPoints: []metricdata.HistogramDataPoint[int64]{
						{
							Attributes:   fooBar,
							Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
							BucketCounts: []uint64{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Count:        2,
							Min:          metricdata.NewExtrema[int64](1),
							Max:          metricdata.NewExtrema[int64](2),
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
						Stream{AllowAttributeKeys: []attribute.Key{"foo"}},
					)),
				).Meter("TestAttributeFilter")
				require.NoError(t, tt.register(t, mtr))

				m := metricdata.ResourceMetrics{}
				err := rdr.Collect(context.Background(), &m)
				assert.NoError(t, err)

				require.Len(t, m.ScopeMetrics, 1)
				require.Len(t, m.ScopeMetrics[0].Metrics, 1)

				metricdatatest.AssertEqual(t, tt.wantMetric, m.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
			})
		}
	}
}

func TestObservableExample(t *testing.T) {
	// This example can be found:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/metrics/supplementary-guidelines.md#asynchronous-example
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

		const (
			instName       = "pageFaults"
			filteredStream = "filteredPageFaults"
			scopeName      = "ObservableExample"
		)

		selector := func(InstrumentKind) metricdata.Temporality { return temp }
		reader := NewManualReader(WithTemporalitySelector(selector))

		noFiltered := NewView(Instrument{Name: instName}, Stream{Name: instName})
		filtered := NewView(Instrument{Name: instName}, Stream{Name: filteredStream, AllowAttributeKeys: []attribute.Key{"pid"}})

		mp := NewMeterProvider(WithReader(reader), WithView(noFiltered, filtered))
		meter := mp.Meter(scopeName)

		observations := make(map[attribute.Set]int64)
		_, err := meter.Int64ObservableCounter(instName, metric.WithInt64Callback(
			func(_ context.Context, o metric.Int64Observer) error {
				for attrSet, val := range observations {
					o.Observe(val, metric.WithAttributeSet(attrSet))
				}
				return nil
			},
		))
		require.NoError(t, err)

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
			got := metricdata.ResourceMetrics{}
			err := reader.Collect(context.Background(), &got)
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
				// Thread 1 is no longer exported.
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
				// Thread 1 is no longer exported.
				{Attributes: thread2, Value: 6},
				{Attributes: thread3, Value: 5},
			},
		}

		verify(t)
	})
}

var (
	aiCounter       metric.Int64ObservableCounter
	aiUpDownCounter metric.Int64ObservableUpDownCounter
	aiGauge         metric.Int64ObservableGauge

	afCounter       metric.Float64ObservableCounter
	afUpDownCounter metric.Float64ObservableUpDownCounter
	afGauge         metric.Float64ObservableGauge

	siCounter       metric.Int64Counter
	siUpDownCounter metric.Int64UpDownCounter
	siHistogram     metric.Int64Histogram

	sfCounter       metric.Float64Counter
	sfUpDownCounter metric.Float64UpDownCounter
	sfHistogram     metric.Float64Histogram
)

func BenchmarkInstrumentCreation(b *testing.B) {
	provider := NewMeterProvider(WithReader(NewManualReader()))
	meter := provider.Meter("BenchmarkInstrumentCreation")

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		aiCounter, _ = meter.Int64ObservableCounter("observable.int64.counter")
		aiUpDownCounter, _ = meter.Int64ObservableUpDownCounter("observable.int64.up.down.counter")
		aiGauge, _ = meter.Int64ObservableGauge("observable.int64.gauge")

		afCounter, _ = meter.Float64ObservableCounter("observable.float64.counter")
		afUpDownCounter, _ = meter.Float64ObservableUpDownCounter("observable.float64.up.down.counter")
		afGauge, _ = meter.Float64ObservableGauge("observable.float64.gauge")

		siCounter, _ = meter.Int64Counter("sync.int64.counter")
		siUpDownCounter, _ = meter.Int64UpDownCounter("sync.int64.up.down.counter")
		siHistogram, _ = meter.Int64Histogram("sync.int64.histogram")

		sfCounter, _ = meter.Float64Counter("sync.float64.counter")
		sfUpDownCounter, _ = meter.Float64UpDownCounter("sync.float64.up.down.counter")
		sfHistogram, _ = meter.Float64Histogram("sync.float64.histogram")
	}
}
