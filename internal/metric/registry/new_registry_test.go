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

package registry_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/internal/metric/registry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/metrictest"
)

type (
	newRegisterFunc func(m metric.Meter, name string) (interface{}, error)
)

var (
	allNewRegisterFunc = map[string]newRegisterFunc{
		"counter.int64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewInt64Counter(name)
		},
		"counter.float64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewFloat64Counter(name)
		},
		"up_down_counter.int64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewInt64UpDownCounter(name)
		},
		"up_down_counter.float64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewFloat64UpDownCounter(name)
		},
		"histogram.int64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewInt64Histogram(name)
		},
		"histogram.float64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewFloat64Histogram(name)
		},
		"gauge_observer.int64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewInt64GaugeObserver(name, func(context.Context, metric.Int64ObserverResult) {})
		},
		"gauge_observer.float64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewFloat64GaugeObserver(name, func(context.Context, metric.Float64ObserverResult) {})
		},
		"counter_observer.int64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewInt64CounterObserver(name, func(context.Context, metric.Int64ObserverResult) {})
		},
		"counter_observer.float64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewFloat64CounterObserver(name, func(context.Context, metric.Float64ObserverResult) {})
		},
		"up_down_counter_observer.int64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewInt64UpDownCounterObserver(name, func(context.Context, metric.Int64ObserverResult) {})
		},
		"up_down_counter_observer.float64": func(m metric.Meter, name string) (interface{}, error) {
			return m.NewFloat64UpDownCounterObserver(name, func(context.Context, metric.Float64ObserverResult) {})
		},
	}
)

func testMeterWithNewRegistry(provider metric.MeterProvider, name string) metric.Meter {
	return registry.NewUniqueInstrumentMeter(provider.Meter(name))
}

func TestNewRegistrySameInstruments(t *testing.T) {
	for name, nf := range allNewRegisterFunc {
		t.Run(name, func(t *testing.T) {
			meter := testMeterWithNewRegistry(metrictest.NewMeterProvider(), "meter")
			inst1, err1 := nf(meter, "this")
			inst2, err2 := nf(meter, "this")

			require.NoError(t, err1)
			require.NoError(t, err2)
			require.Equal(t, inst1, inst2)
		})
	}
}

func TestNewRegistryDifferentNamespace(t *testing.T) {
	for name, nf := range allNewRegisterFunc {
		t.Run(name, func(t *testing.T) {
			provider := metrictest.NewMeterProvider()

			meter1 := testMeterWithNewRegistry(provider, "meter1")
			meter2 := testMeterWithNewRegistry(provider, "meter2")

			inst1, err1 := nf(meter1, "this")
			inst2, err2 := nf(meter2, "this")
			require.NoError(t, err1)
			require.NoError(t, err2)
			require.NotEqual(t, inst1, inst2)
		})
	}
}

func TestNewRegistryDiffInstruments(t *testing.T) {
	for origName, origF := range allNewRegisterFunc {
		t.Run(origName, func(t *testing.T) {
			meter := testMeterWithNewRegistry(metrictest.NewMeterProvider(), "meter")

			origI, err := origF(meter, "this")
			require.NoError(t, err)

			for newName, newF := range allNewRegisterFunc {
				newI, err := newF(meter, "this")
				if newName == origName {
					require.NoError(t, err)
					assert.Equal(t, origI, newI)
					continue
				}

				require.Error(t, err)
				require.Nil(t, newI)
				require.True(t, errors.Is(err, registry.ErrMetricKindMismatch))
			}
		})
	}
}
