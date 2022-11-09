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

	mp.setDelegate(metric.NewNoopMeterProvider())
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
			_, _ = mtr.Float64ObservableCounter(name)
			_, _ = mtr.Float64ObservableUpDownCounter(name)
			_, _ = mtr.Float64ObservableGauge(name)
			_, _ = mtr.Int64ObservableCounter(name)
			_, _ = mtr.Int64ObservableUpDownCounter(name)
			_, _ = mtr.Int64ObservableGauge(name)
			_, _ = mtr.Float64Counter(name)
			_, _ = mtr.Float64UpDownCounter(name)
			_, _ = mtr.Float64Histogram(name)
			_, _ = mtr.Int64Counter(name)
			_, _ = mtr.Int64UpDownCounter(name)
			_, _ = mtr.Int64Histogram(name)
			_, _ = mtr.RegisterCallback(func(ctx context.Context) error {
				return nil
			}, nil)
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
	mtr.setDelegate(metric.NewNoopMeterProvider())
	close(finish)
}

func testSetupAllInstrumentTypes(t *testing.T, m metric.Meter) (metric.Float64Counter, metric.Float64ObservableCounter) {
	afcounter, err := m.Float64ObservableCounter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.Float64ObservableUpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.Float64ObservableGauge("test_Async_Gauge")
	assert.NoError(t, err)

	_, err = m.Int64ObservableCounter("test_Async_Counter")
	assert.NoError(t, err)
	_, err = m.Int64ObservableUpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.Int64ObservableGauge("test_Async_Gauge")
	assert.NoError(t, err)

	_, err = m.RegisterCallback(func(ctx context.Context) error {
		afcounter.Observe(ctx, 3)
		return nil
	}, afcounter)
	assert.NoError(t, err)

	sfcounter, err := m.Float64Counter("test_Async_Counter")
	require.NoError(t, err)
	_, err = m.Float64UpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.Float64Histogram("test_Async_Histogram")
	assert.NoError(t, err)

	_, err = m.Int64Counter("test_Async_Counter")
	assert.NoError(t, err)
	_, err = m.Int64UpDownCounter("test_Async_UpDownCounter")
	assert.NoError(t, err)
	_, err = m.Int64Histogram("test_Async_Histogram")
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

	assert.Equal(t, 1, tMeter.afCounter)
	assert.Equal(t, 1, tMeter.afUpDownCounter)
	assert.Equal(t, 1, tMeter.afGauge)
	assert.Equal(t, 1, tMeter.aiCounter)
	assert.Equal(t, 1, tMeter.aiUpDownCounter)
	assert.Equal(t, 1, tMeter.aiGauge)
	assert.Equal(t, 1, tMeter.sfCounter)
	assert.Equal(t, 1, tMeter.sfUpDownCounter)
	assert.Equal(t, 1, tMeter.sfHistogram)
	assert.Equal(t, 1, tMeter.siCounter)
	assert.Equal(t, 1, tMeter.siUpDownCounter)
	assert.Equal(t, 1, tMeter.siHistogram)
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
	assert.Equal(t, 1, tMeter.afCounter)
	assert.Equal(t, 1, tMeter.afUpDownCounter)
	assert.Equal(t, 1, tMeter.afGauge)
	assert.Equal(t, 1, tMeter.aiCounter)
	assert.Equal(t, 1, tMeter.aiUpDownCounter)
	assert.Equal(t, 1, tMeter.aiGauge)
	assert.Equal(t, 1, tMeter.sfCounter)
	assert.Equal(t, 1, tMeter.sfUpDownCounter)
	assert.Equal(t, 1, tMeter.sfHistogram)
	assert.Equal(t, 1, tMeter.siCounter)
	assert.Equal(t, 1, tMeter.siUpDownCounter)
	assert.Equal(t, 1, tMeter.siHistogram)

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
	assert.Equal(t, 1, tMeter.afCounter)
	assert.Equal(t, 1, tMeter.afUpDownCounter)
	assert.Equal(t, 1, tMeter.afGauge)
	assert.Equal(t, 1, tMeter.aiCounter)
	assert.Equal(t, 1, tMeter.aiUpDownCounter)
	assert.Equal(t, 1, tMeter.aiGauge)
	assert.Equal(t, 1, tMeter.sfCounter)
	assert.Equal(t, 1, tMeter.sfUpDownCounter)
	assert.Equal(t, 1, tMeter.sfHistogram)
	assert.Equal(t, 1, tMeter.siCounter)
	assert.Equal(t, 1, tMeter.siUpDownCounter)
	assert.Equal(t, 1, tMeter.siHistogram)

	// Because the Meter was a delegate it should return a delegated instrument

	assert.IsType(t, &sfCounter{}, ctr)
	assert.IsType(t, &afCounter{}, actr)
	assert.Equal(t, 1, mp.count)
}
