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
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/internal/metric/registry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/metrictest"
)

type (
	newFunc func(name, libraryName string) (interface{}, error)
)

var (
	allNew = map[string]newFunc{
		"counter.int64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewInt64Counter(name)
		},
		"counter.float64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewFloat64Counter(name)
		},
		"up_down_counter.int64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewInt64UpDownCounter(name)
		},
		"up_down_counter.float64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewFloat64UpDownCounter(name)
		},
		"histogram.int64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewInt64Histogram(name)
		},
		"histogram.float64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewFloat64Histogram(name)
		},
		"gauge_observer.int64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewInt64GaugeObserver(name, func(context.Context, metric.Int64ObserverResult) {})
		},
		"gauge_observer.float64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewFloat64GaugeObserver(name, func(context.Context, metric.Float64ObserverResult) {})
		},
		"counter_observer.int64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewInt64CounterObserver(name, func(context.Context, metric.Int64ObserverResult) {})
		},
		"counter_observer.float64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewFloat64CounterObserver(name, func(context.Context, metric.Float64ObserverResult) {})
		},
		"up_down_counter_observer.int64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewInt64UpDownCounterObserver(name, func(context.Context, metric.Int64ObserverResult) {})
		},
		"up_down_counter_observer.float64": func(name, libraryName string) (interface{}, error) {
			return MeterProvider().Meter(libraryName).NewFloat64UpDownCounterObserver(name, func(context.Context, metric.Float64ObserverResult) {})
		},
	}
)

func TestRegistrySameInstruments(t *testing.T) {
	for _, nf := range allNew {
		ResetForTest()
		inst1, err1 := nf("this", "meter")
		inst2, err2 := nf("this", "meter")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Equal(t, inst1, inst2)

		SetMeterProvider(metrictest.NewMeterProvider())

		require.Equal(t, inst1, inst2)
	}
}

func TestRegistryDiffInstruments(t *testing.T) {
	for origName, origf := range allNew {
		ResetForTest()

		_, err := origf("this", "super")
		require.NoError(t, err)

		for newName, nf := range allNew {
			if newName == origName {
				continue
			}

			other, err := nf("this", "super")
			require.Error(t, err)
			require.Nil(t, other)
			require.True(t, errors.Is(err, registry.ErrMetricKindMismatch))
			require.Contains(t, err.Error(), "by this name with another kind or number type")
		}
	}
}
