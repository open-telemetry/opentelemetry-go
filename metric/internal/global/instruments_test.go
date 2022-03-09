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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/nonrecording"
)

func Test_asyncInstrument_setDelegate_race(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64 Instruments", func(t *testing.T) {
		t.Run("Async Counter", func(t *testing.T) {
			delegate := &afCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Observe(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Async UpDownCounter", func(t *testing.T) {
			delegate := &afUpDownCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Observe(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Async Gauge", func(t *testing.T) {
			delegate := &afGauge{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Observe(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})
	})

	// Int64 Instruments

	t.Run("int64 Instruments", func(t *testing.T) {
		t.Run("Async Counter", func(t *testing.T) {
			delegate := &aiCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Observe(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Async UpDownCounter", func(t *testing.T) {
			delegate := &aiUpDownCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Observe(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Async Gauge", func(t *testing.T) {
			delegate := &aiGauge{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Observe(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})
	})
}

func Test_syncInstrument_setDelegate_race(t *testing.T) {
	// Float64 Instruments
	// Float64 Instruments
	t.Run("Float64 Instruments", func(t *testing.T) {
		t.Run("Sync Counter", func(t *testing.T) {
			delegate := &sfCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Add(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Sync UpDownCounter", func(t *testing.T) {
			delegate := &sfUpDownCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Add(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Sync Histogram", func(t *testing.T) {
			delegate := &sfHistogram{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Record(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})
	})

	// Int64 Instruments

	t.Run("Int64 Instruments", func(t *testing.T) {
		t.Run("Sync Counter", func(t *testing.T) {
			delegate := &siCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Add(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Sync UpDownCounter", func(t *testing.T) {
			delegate := &siUpDownCounter{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Add(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
		})

		t.Run("Sync Histogram", func(t *testing.T) {
			delegate := &siHistogram{
				name: "testName",
				opts: []instrument.Option{},
			}
			finish := make(chan struct{})
			go func() {
				for {
					delegate.Record(context.Background(), 1)
					select {
					case <-finish:
						return
					default:
					}
				}
			}()

			delegate.setDelegate(nonrecording.NewNoopMeter())
			close(finish)
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
