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
	"go.opentelemetry.io/otel/metric/nonrecording"
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

	setDelegate(nonrecording.NewNoopMeter())
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

	setDelegate(nonrecording.NewNoopMeter())
	close(finish)
}

func TestAsyncInstrumentSetDelegateRace(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &afCounter{}
			testFloat64Race(delegate.Observe, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &afUpDownCounter{}
			testFloat64Race(delegate.Observe, delegate.setDelegate)
		})

		t.Run("Gauge", func(t *testing.T) {
			delegate := &afGauge{}
			testFloat64Race(delegate.Observe, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(t *testing.T) {
		t.Run("Counter", func(t *testing.T) {
			delegate := &aiCounter{}
			testInt64Race(delegate.Observe, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(t *testing.T) {
			delegate := &aiUpDownCounter{}
			testInt64Race(delegate.Observe, delegate.setDelegate)
		})

		t.Run("Gauge", func(t *testing.T) {
			delegate := &aiGauge{}
			testInt64Race(delegate.Observe, delegate.setDelegate)
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

	instrument.Asynchronous
	instrument.Synchronous
}

func (i *testCountingFloatInstrument) Observe(context.Context, float64, ...attribute.KeyValue) {
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

	instrument.Asynchronous
	instrument.Synchronous
}

func (i *testCountingIntInstrument) Observe(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingIntInstrument) Add(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
func (i *testCountingIntInstrument) Record(context.Context, int64, ...attribute.KeyValue) {
	i.count++
}
