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

package global // import "go.opentelemetry.io/otel/metric/internal/global"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

func Test_MeterProvider_delegates_calls(t *testing.T) {

	// The global MeterProvider should directly call the underlying MeterProvider
	// if it is set prior to Meter() being called.

	// globalMeterProvider := otel.GetMeterProvider
	globalMeterProvider := &meterProvider{}

	mp := &test_MeterProvider{}

	// otel.SetMeterProvider(mp)
	globalMeterProvider.setDelegate(mp)

	require.Equal(t, 0, mp.count)

	meter := globalMeterProvider.Meter("go.opentelemetry.io/otel/metric/internal/global/meter_test")

	ctr, actr := test_setup_all_instrument_types(t, meter)

	ctr.Add(context.Background(), 5)

	test_collect(t, meter)

	// Calls to Meter() after setDelegate() should be executed by the delegate
	require.IsType(t, &test_Meter{}, meter)
	if t_meter, ok := meter.(*test_Meter); ok {
		require.Equal(t, 3, t_meter.afCount)
		require.Equal(t, 3, t_meter.aiCount)
		require.Equal(t, 3, t_meter.sfCount)
		require.Equal(t, 3, t_meter.siCount)
		require.Equal(t, 1, len(t_meter.callbacks))
	}

	// Because the Meter was provided by test_meterProvider it should also return our test instrument
	require.IsType(t, &test_counting_float_instrument{}, ctr, "the meter did not delegate calls to the meter")
	if test_ctr, ok := ctr.(*test_counting_float_instrument); ok {
		require.Equal(t, 1, test_ctr.count)
	}

	require.IsType(t, &test_counting_float_instrument{}, actr, "the meter did not delegate calls to the meter")
	if test_ctr, ok := actr.(*test_counting_float_instrument); ok {
		require.Equal(t, 1, test_ctr.count)
	}

	require.Equal(t, 1, mp.count)
}

func test_setup_all_instrument_types(t *testing.T, m metric.Meter) (syncfloat64.Counter, asyncfloat64.Counter) {

	afcounter, err := m.AsyncFloat64().Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.AsyncFloat64().UpDownCounter("test_Async_UpDownCounter")
	require.NoError(t, err)
	_, err = m.AsyncFloat64().Gauge("test_Async_Gauge")
	require.NoError(t, err)

	_, err = m.AsyncInt64().Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.AsyncInt64().UpDownCounter("test_Async_UpDownCounter")
	require.NoError(t, err)
	_, err = m.AsyncInt64().Gauge("test_Async_Gauge")
	require.NoError(t, err)

	m.RegisterCallback([]instrument.Asynchronous{afcounter}, func(ctx context.Context) {
		afcounter.Observe(ctx, 3)
	})

	sfcounter, err := m.SyncFloat64().Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.SyncFloat64().UpDownCounter("test_Async_UpDownCounter")
	require.NoError(t, err)
	_, err = m.SyncFloat64().Histogram("test_Async_Histogram")
	require.NoError(t, err)

	_, err = m.SyncInt64().Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.SyncInt64().UpDownCounter("test_Async_UpDownCounter")
	require.NoError(t, err)
	_, err = m.SyncInt64().Histogram("test_Async_Histogram")
	require.NoError(t, err)

	return sfcounter, afcounter
}

func test_collect(t *testing.T, m metric.Meter) {
	if t_meter, ok := m.(*meter); ok {
		m, ok = t_meter.delegate.Load().(metric.Meter)
		if !ok {
			t.Error("meter was not delegated")
			return
		}
	}
	t_meter, ok := m.(*test_Meter)
	if !ok {
		t.Error("collect called on non-test Meter")
		return
	}
	t_meter.collect()
}

func Test_Meter_delegates_calls(t *testing.T) {

	// The global MeterProvider should directly provide a Meter instance that
	// can be updated.  If the SetMeterProvider is called after a Meter was
	// obtained, but before instruments only the instrument should be generated
	// by the delegated type.

	// globalMeterProvider := otel.GetMeterProvider
	globalMeterProvider := &meterProvider{}

	mp := &test_MeterProvider{}

	require.Equal(t, 0, mp.count)

	m := globalMeterProvider.Meter("go.opentelemetry.io/otel/metric/internal/global/meter_test")

	// otel.SetMeterProvider(mp)
	globalMeterProvider.setDelegate(mp)

	ctr, actr := test_setup_all_instrument_types(t, m)

	ctr.Add(context.Background(), 5)

	test_collect(t, m)

	// Calls to Meter methods after setDelegate() should be executed by the delegate
	require.IsType(t, &meter{}, m)
	if d_meter, ok := m.(*meter); ok {
		m := d_meter.delegate.Load().(*test_Meter)
		require.NotNil(t, m)
		require.Equal(t, 3, m.afCount)
		require.Equal(t, 3, m.aiCount)
		require.Equal(t, 3, m.sfCount)
		require.Equal(t, 3, m.siCount)
	}

	// Because the Meter was provided by test_meterProvider it should also return our test instrument
	require.IsType(t, &test_counting_float_instrument{}, ctr, "the meter did not delegate calls to the meter")
	if test_ctr, ok := ctr.(*test_counting_float_instrument); ok {
		require.Equal(t, 1, test_ctr.count)
	}

	// Because the Meter was provided by test_meterProvider it should also return our test instrument
	require.IsType(t, &test_counting_float_instrument{}, actr, "the meter did not delegate calls to the meter")
	if test_ctr, ok := actr.(*test_counting_float_instrument); ok {
		require.Equal(t, 1, test_ctr.count)
	}

	require.Equal(t, 1, mp.count)
}

func Test_Meter_defers_delegations(t *testing.T) {

	// If SetMeterProvider is called after insturments are registered, the
	// instruments should be recreated with the new meter.

	// globalMeterProvider := otel.GetMeterProvider
	globalMeterProvider := &meterProvider{}

	m := globalMeterProvider.Meter("go.opentelemetry.io/otel/metric/internal/global/meter_test")

	ctr, actr := test_setup_all_instrument_types(t, m)

	ctr.Add(context.Background(), 5)

	mp := &test_MeterProvider{}

	// otel.SetMeterProvider(mp)
	globalMeterProvider.setDelegate(mp)

	// Calls to Meter() before setDelegate() should be the delegated type
	require.IsType(t, &meter{}, m)

	if d_meter, ok := m.(*meter); ok {
		m := d_meter.delegate.Load().(*test_Meter)
		require.NotNil(t, m)
		require.Equal(t, 1, m.afCount)
		require.Equal(t, 1, m.aiCount)
		require.Equal(t, 1, m.sfCount)
		require.Equal(t, 1, m.siCount)
	}

	// Because the Meter was a delegate it should return a delegated instrument

	require.IsType(t, &sfCounter{}, ctr)

	require.Equal(t, 1, mp.count)
}
