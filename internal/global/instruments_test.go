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
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/noop"
)

func testRace[N int64 | float64](interact func(N), setDelegate func(metric.Meter) error) {
	finish := make(chan struct{})
	go func() {
		for {
			interact(1)
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	setDelegate(noop.NewMeterProvider().Meter(""))
	close(finish)
}

func TestInstrumentSetDelegateRace(t *testing.T) {
	meter := meter{}
	t.Run("Float64Counter", func(t *testing.T) {
		i, err := meter.Float64Counter("")
		require.NoError(t, err)
		interact := func(v float64) { i.Add(context.Background(), v) }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Float64UpDownCounter", func(t *testing.T) {
		i, err := meter.Float64UpDownCounter("")
		require.NoError(t, err)
		interact := func(v float64) { i.Add(context.Background(), v) }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Float64Histogram", func(t *testing.T) {
		i, err := meter.Float64Histogram("")
		require.NoError(t, err)
		interact := func(v float64) { i.Record(context.Background(), v) }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Float64ObservableCounter", func(t *testing.T) {
		i, err := meter.Float64ObservableCounter("")
		require.NoError(t, err)
		require.Implements(t, (*unwrapper)(nil), i)
		interact := func(float64) { _ = i.(unwrapper).Unwrap() }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Float64ObservableUpDownCounter", func(t *testing.T) {
		i, err := meter.Float64ObservableUpDownCounter("")
		require.NoError(t, err)
		require.Implements(t, (*unwrapper)(nil), i)
		interact := func(float64) { _ = i.(unwrapper).Unwrap() }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Float64ObservableGauge", func(t *testing.T) {
		i, err := meter.Float64ObservableGauge("")
		require.NoError(t, err)
		require.Implements(t, (*unwrapper)(nil), i)
		interact := func(float64) { _ = i.(unwrapper).Unwrap() }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Int64Counter", func(t *testing.T) {
		i, err := meter.Int64Counter("")
		require.NoError(t, err)
		interact := func(v int64) { i.Add(context.Background(), v) }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Int64UpDownCounter", func(t *testing.T) {
		i, err := meter.Int64UpDownCounter("")
		require.NoError(t, err)
		interact := func(v int64) { i.Add(context.Background(), v) }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Int64Histogram", func(t *testing.T) {
		i, err := meter.Int64Histogram("")
		require.NoError(t, err)
		interact := func(v int64) { i.Record(context.Background(), v) }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Int64ObservableCounter", func(t *testing.T) {
		i, err := meter.Int64ObservableCounter("")
		require.NoError(t, err)
		require.Implements(t, (*unwrapper)(nil), i)
		interact := func(int64) { _ = i.(unwrapper).Unwrap() }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Int64ObservableUpDownCounter", func(t *testing.T) {
		i, err := meter.Int64ObservableUpDownCounter("")
		require.NoError(t, err)
		require.Implements(t, (*unwrapper)(nil), i)
		interact := func(int64) { _ = i.(unwrapper).Unwrap() }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
	})

	t.Run("Int64ObservableGauge", func(t *testing.T) {
		i, err := meter.Int64ObservableGauge("")
		require.NoError(t, err)
		require.Implements(t, (*unwrapper)(nil), i)
		interact := func(int64) { _ = i.(unwrapper).Unwrap() }
		require.Implements(t, (*delegatedInstrument)(nil), i)
		testRace(interact, i.(delegatedInstrument).SetDelegate)
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
