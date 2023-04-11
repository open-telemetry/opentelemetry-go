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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/metricembed"
	"go.opentelemetry.io/otel/metric/noop"
)

func testFloat64Race(interact func(context.Context, float64, ...attribute.KeyValue), setDelegate func(metric.Meter)) {
	finish := make(chan struct{})
	go func() {
		for {
			interact(context.Background(), 1)
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

func testInt64Race(interact func(context.Context, int64, ...attribute.KeyValue), setDelegate func(metric.Meter)) {
	finish := make(chan struct{})
	go func() {
		for {
			interact(context.Background(), 1)
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

func TestAsyncInstrumentSetDelegateRace(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &afCounter{}
			f := func(context.Context, float64, ...attribute.KeyValue) { _ = delegate.Unwrap() }
			testFloat64Race(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &afUpDownCounter{}
			f := func(context.Context, float64, ...attribute.KeyValue) { _ = delegate.Unwrap() }
			testFloat64Race(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(t *testing.T) {
			delegate := &afGauge{}
			f := func(context.Context, float64, ...attribute.KeyValue) { _ = delegate.Unwrap() }
			testFloat64Race(f, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &aiCounter{}
			f := func(context.Context, int64, ...attribute.KeyValue) { _ = delegate.Unwrap() }
			testInt64Race(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &aiUpDownCounter{}
			f := func(context.Context, int64, ...attribute.KeyValue) { _ = delegate.Unwrap() }
			testInt64Race(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(t *testing.T) {
			delegate := &aiGauge{}
			f := func(context.Context, int64, ...attribute.KeyValue) { _ = delegate.Unwrap() }
			testInt64Race(f, delegate.setDelegate)
		})
	})
}

func TestSyncInstrumentSetDelegateRace(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &sfCounter{}
			testFloat64Race(delegate.Add, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &sfUpDownCounter{}
			testFloat64Race(delegate.Add, delegate.setDelegate)
		})

		t.Run("Histogram", func(t *testing.T) {
			delegate := &sfHistogram{}
			testFloat64Race(delegate.Record, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &siCounter{}
			testInt64Race(delegate.Add, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &siUpDownCounter{}
			testInt64Race(delegate.Add, delegate.setDelegate)
		})

		t.Run("Histogram", func(t *testing.T) {
			delegate := &siHistogram{}
			testInt64Race(delegate.Record, delegate.setDelegate)
		})
	})
}

type testCountingFloatInstrument struct {
	count int

	instrument.Float64Observable
	metricembed.Float64Counter
	metricembed.Float64UpDownCounter
	metricembed.Float64Histogram
	metricembed.Float64ObservableCounter
	metricembed.Float64ObservableUpDownCounter
	metricembed.Float64ObservableGauge
}

func (i *testCountingFloatInstrument) observe() {
	i.count++
}
func (i *testCountingFloatInstrument) Add(context.Context, float64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingFloatInstrument) Record(context.Context, float64, ...attribute.KeyValue) {
	i.count++
}

type testCountingIntInstrument struct {
	count int

	instrument.Int64Observable
	metricembed.Int64Counter
	metricembed.Int64UpDownCounter
	metricembed.Int64Histogram
	metricembed.Int64ObservableCounter
	metricembed.Int64ObservableUpDownCounter
	metricembed.Int64ObservableGauge
}

func (i *testCountingIntInstrument) observe() {
	i.count++
}
func (i *testCountingIntInstrument) Add(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingIntInstrument) Record(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
