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

package global

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/noop"
)

func testRace(interact func(), setDelegate func()) {
	finish := make(chan struct{})
	go func() {
		for {
			interact()
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	setDelegate()
	close(finish)
}

func must[T any, O any](t *testing.T, f func(string, ...O) (T, error)) T {
	v, err := f("")
	require.NoError(t, err)
	return v
}

func TestInstrumentSetDelegateRace(t *testing.T) {
	meter := meter{}
	altM := noop.NewMeterProvider().Meter("")

	t.Run("Float64Counter", func(t *testing.T) {
		i := must(t, meter.Float64Counter)
		interact := func() { i.Add(context.Background(), 1) }
		setDelegate := func() { i.(storer).Store(must(t, altM.Float64Counter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Float64UpDownCounter", func(t *testing.T) {
		i := must(t, meter.Float64UpDownCounter)
		interact := func() { i.Add(context.Background(), 1) }
		setDelegate := func() { i.(storer).Store(must(t, altM.Float64UpDownCounter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Float64Histogram", func(t *testing.T) {
		i := must(t, meter.Float64Histogram)
		interact := func() { i.Record(context.Background(), 1) }
		setDelegate := func() { i.(storer).Store(must(t, altM.Float64Histogram)) }
		testRace(interact, setDelegate)
	})

	t.Run("Float64ObservableCounter", func(t *testing.T) {
		i := must(t, meter.Float64ObservableCounter)
		interact := func() { _ = i.(unwrapper).Unwrap() }
		setDelegate := func() { i.(storer).Store(must(t, altM.Float64ObservableCounter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Float64ObservableUpDownCounter", func(t *testing.T) {
		i := must(t, meter.Float64ObservableUpDownCounter)
		interact := func() { _ = i.(unwrapper).Unwrap() }
		setDelegate := func() { i.(storer).Store(must(t, altM.Float64ObservableUpDownCounter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Float64ObservableGauge", func(t *testing.T) {
		i := must(t, meter.Float64ObservableGauge)
		interact := func() { _ = i.(unwrapper).Unwrap() }
		setDelegate := func() { i.(storer).Store(must(t, altM.Float64ObservableGauge)) }
		testRace(interact, setDelegate)
	})

	t.Run("Int64Counter", func(t *testing.T) {
		i := must(t, meter.Int64Counter)
		interact := func() { i.Add(context.Background(), 1) }
		setDelegate := func() { i.(storer).Store(must(t, altM.Int64Counter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Int64UpDownCounter", func(t *testing.T) {
		i := must(t, meter.Int64UpDownCounter)
		interact := func() { i.Add(context.Background(), 1) }
		setDelegate := func() { i.(storer).Store(must(t, altM.Int64UpDownCounter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Int64Histogram", func(t *testing.T) {
		i := must(t, meter.Int64Histogram)
		interact := func() { i.Record(context.Background(), 1) }
		setDelegate := func() { i.(storer).Store(must(t, altM.Int64Histogram)) }
		testRace(interact, setDelegate)
	})

	t.Run("Int64ObservableCounter", func(t *testing.T) {
		i := must(t, meter.Int64ObservableCounter)
		interact := func() { _ = i.(unwrapper).Unwrap() }
		setDelegate := func() { i.(storer).Store(must(t, altM.Int64ObservableCounter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Int64ObservableUpDownCounter", func(t *testing.T) {
		i := must(t, meter.Int64ObservableUpDownCounter)
		interact := func() { _ = i.(unwrapper).Unwrap() }
		setDelegate := func() { i.(storer).Store(must(t, altM.Int64ObservableUpDownCounter)) }
		testRace(interact, setDelegate)
	})

	t.Run("Int64ObservableGauge", func(t *testing.T) {
		i := must(t, meter.Int64ObservableGauge)
		interact := func() { _ = i.(unwrapper).Unwrap() }
		setDelegate := func() { i.(storer).Store(must(t, altM.Int64ObservableGauge)) }
		testRace(interact, setDelegate)
	})
}

type testCountingInstrument[N int64 | float64] struct {
	count int

	instrument.ObservableT[N]
	embedded.Counter[N]
	embedded.UpDownCounter[N]
	embedded.Histogram[N]
	embedded.ObservableCounter[N]
	embedded.ObservableUpDownCounter[N]
	embedded.ObservableGauge[N]
}

func (i *testCountingInstrument[N]) observe() {
	i.count++
}
func (i *testCountingInstrument[N]) Add(context.Context, N, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingInstrument[N]) Record(context.Context, N, ...attribute.KeyValue) {
	i.count++
}
