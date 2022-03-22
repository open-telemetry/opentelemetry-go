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
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/nonrecording"
)

func TestMeterProviderRace(t *testing.T) {
	mp := &meterProvider{}
	finish := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			mp.Meter(fmt.Sprintf("a%d", i))
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	mp.setDelegate(nonrecording.NewNoopMeterProvider())
	close(finish)

}

func TestMeterRace(t *testing.T) {
	mtr := &meter{}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	finish := make(chan struct{})
	go func() {
		for i, once := 0, false; ; i++ {
			name := fmt.Sprintf("a%d", i)
			_, _ = mtr.AsyncFloat64().Counter(name)
			_, _ = mtr.AsyncFloat64().UpDownCounter(name)
			_, _ = mtr.AsyncFloat64().Gauge(name)
			_, _ = mtr.AsyncInt64().Counter(name)
			_, _ = mtr.AsyncInt64().UpDownCounter(name)
			_, _ = mtr.AsyncInt64().Gauge(name)
			_, _ = mtr.SyncFloat64().Counter(name)
			_, _ = mtr.SyncFloat64().UpDownCounter(name)
			_, _ = mtr.SyncFloat64().Histogram(name)
			_, _ = mtr.SyncInt64().Counter(name)
			_, _ = mtr.SyncInt64().UpDownCounter(name)
			_, _ = mtr.SyncInt64().Histogram(name)
			_ = mtr.RegisterCallback(nil, func(ctx context.Context) {})
			if !once {
				wg.Done()
				once = true
			}
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	wg.Wait()
	mtr.setDelegate(nonrecording.NewNoopMeterProvider())
	close(finish)
}

func testSetupAllInstrumentTypes(t *testing.T, m metric.Meter) (syncfloat64.Counter, asyncfloat64.Counter) {

	afcounter, err := m.AsyncFloat64().Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.AsyncFloat64().UpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.AsyncFloat64().Gauge("test_Async_Gauge")
	assert.NoError(t, err)

	_, err = m.AsyncInt64().Counter("test_Async_Counter")
	assert.NoError(t, err)
	_, err = m.AsyncInt64().UpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.AsyncInt64().Gauge("test_Async_Gauge")
	assert.NoError(t, err)

	require.NoError(t, m.RegisterCallback([]instrument.Asynchronous{afcounter}, func(ctx context.Context) {
		afcounter.Observe(ctx, 3)
	}))

	sfcounter, err := m.SyncFloat64().Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.SyncFloat64().UpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.SyncFloat64().Histogram("test_Async_Histogram")
	assert.NoError(t, err)

	_, err = m.SyncInt64().Counter("test_Async_Counter")
	assert.NoError(t, err)
	_, err = m.SyncInt64().UpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.SyncInt64().Histogram("test_Async_Histogram")
	assert.NoError(t, err)

	return sfcounter, afcounter
}

// This is to emulate a read from an exporter.
func testCollect(t *testing.T, m metric.Meter) {
	if tMeter, ok := m.(*meter); ok {
		m, ok = tMeter.delegate.Load().(metric.Meter)
		if !ok {
			t.Error("meter was not delegated")
			return
		}
	}
	tMeter, ok := m.(*testMeter)
	if !ok {
		t.Error("collect called on non-test Meter")
		return
	}
	tMeter.collect()
}

func TestMeterProviderDelegatesCalls(t *testing.T) {

	// The global MeterProvider should directly call the underlying MeterProvider
	// if it is set prior to Meter() being called.

	// globalMeterProvider := otel.GetMeterProvider
	globalMeterProvider := &meterProvider{}

	mp := &testMeterProvider{}

	// otel.SetMeterProvider(mp)
	globalMeterProvider.setDelegate(mp)

	assert.Equal(t, 0, mp.count)

	meter := globalMeterProvider.Meter("go.opentelemetry.io/otel/metric/internal/global/meter_test")

	ctr, actr := testSetupAllInstrumentTypes(t, meter)

	ctr.Add(context.Background(), 5)

	testCollect(t, meter) // This is a hacky way to emulate a read from an exporter

	// Calls to Meter() after setDelegate() should be executed by the delegate
	require.IsType(t, &testMeter{}, meter)
	tMeter := meter.(*testMeter)
	assert.Equal(t, 3, tMeter.afCount)
	assert.Equal(t, 3, tMeter.aiCount)
	assert.Equal(t, 3, tMeter.sfCount)
	assert.Equal(t, 3, tMeter.siCount)
	assert.Equal(t, 1, len(tMeter.callbacks))

	// Because the Meter was provided by testmeterProvider it should also return our test instrument
	require.IsType(t, &testCountingFloatInstrument{}, ctr, "the meter did not delegate calls to the meter")
	assert.Equal(t, 1, ctr.(*testCountingFloatInstrument).count)

	require.IsType(t, &testCountingFloatInstrument{}, actr, "the meter did not delegate calls to the meter")
	assert.Equal(t, 1, actr.(*testCountingFloatInstrument).count)

	assert.Equal(t, 1, mp.count)
}

func TestMeterDelegatesCalls(t *testing.T) {

	// The global MeterProvider should directly provide a Meter instance that
	// can be updated.  If the SetMeterProvider is called after a Meter was
	// obtained, but before instruments only the instrument should be generated
	// by the delegated type.

	globalMeterProvider := &meterProvider{}

	mp := &testMeterProvider{}

	assert.Equal(t, 0, mp.count)

	m := globalMeterProvider.Meter("go.opentelemetry.io/otel/metric/internal/global/meter_test")

	globalMeterProvider.setDelegate(mp)

	ctr, actr := testSetupAllInstrumentTypes(t, m)

	ctr.Add(context.Background(), 5)

	testCollect(t, m) // This is a hacky way to emulate a read from an exporter

	// Calls to Meter methods after setDelegate() should be executed by the delegate
	require.IsType(t, &meter{}, m)
	tMeter := m.(*meter).delegate.Load().(*testMeter)
	require.NotNil(t, tMeter)
	assert.Equal(t, 3, tMeter.afCount)
	assert.Equal(t, 3, tMeter.aiCount)
	assert.Equal(t, 3, tMeter.sfCount)
	assert.Equal(t, 3, tMeter.siCount)

	// Because the Meter was provided by testmeterProvider it should also return our test instrument
	require.IsType(t, &testCountingFloatInstrument{}, ctr, "the meter did not delegate calls to the meter")
	assert.Equal(t, 1, ctr.(*testCountingFloatInstrument).count)

	// Because the Meter was provided by testmeterProvider it should also return our test instrument
	require.IsType(t, &testCountingFloatInstrument{}, actr, "the meter did not delegate calls to the meter")
	assert.Equal(t, 1, actr.(*testCountingFloatInstrument).count)

	assert.Equal(t, 1, mp.count)
}

func TestMeterDefersDelegations(t *testing.T) {

	// If SetMeterProvider is called after instruments are registered, the
	// instruments should be recreated with the new meter.

	// globalMeterProvider := otel.GetMeterProvider
	globalMeterProvider := &meterProvider{}

	m := globalMeterProvider.Meter("go.opentelemetry.io/otel/metric/internal/global/meter_test")

	ctr, actr := testSetupAllInstrumentTypes(t, m)

	ctr.Add(context.Background(), 5)

	mp := &testMeterProvider{}

	// otel.SetMeterProvider(mp)
	globalMeterProvider.setDelegate(mp)

	testCollect(t, m) // This is a hacky way to emulate a read from an exporter

	// Calls to Meter() before setDelegate() should be the delegated type
	require.IsType(t, &meter{}, m)
	tMeter := m.(*meter).delegate.Load().(*testMeter)
	require.NotNil(t, tMeter)
	assert.Equal(t, 3, tMeter.afCount)
	assert.Equal(t, 3, tMeter.aiCount)
	assert.Equal(t, 3, tMeter.sfCount)
	assert.Equal(t, 3, tMeter.siCount)

	// Because the Meter was a delegate it should return a delegated instrument

	assert.IsType(t, &sfCounter{}, ctr)
	assert.IsType(t, &afCounter{}, actr)
	assert.Equal(t, 1, mp.count)
}
