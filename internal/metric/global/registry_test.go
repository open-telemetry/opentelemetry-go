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
	"go.opentelemetry.io/otel/metric/sdkapi"
)

type (
	newFunc func(name, libraryName string) (sdkapi.InstrumentImpl, error)
)

var (
	allNew = map[string]newFunc{
		"counter.int64": func(name, libraryName string) (sdkapi.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64Counter(name))
		},
		"counter.float64": func(name, libraryName string) (sdkapi.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64Counter(name))
		},
		"histogram.int64": func(name, libraryName string) (sdkapi.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64Histogram(name))
		},
		"histogram.float64": func(name, libraryName string) (sdkapi.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64Histogram(name))
		},
		"gauge.int64": func(name, libraryName string) (sdkapi.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64GaugeObserver(name, func(context.Context, metric.Int64ObserverResult) {}))
		},
		"gauge.float64": func(name, libraryName string) (sdkapi.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64GaugeObserver(name, func(context.Context, metric.Float64ObserverResult) {}))
		},
	}
)

func unwrap(impl interface{}, err error) (sdkapi.InstrumentImpl, error) {
	if impl == nil {
		return nil, err
	}
	if s, ok := impl.(interface {
		SyncImpl() sdkapi.SyncImpl
	}); ok {
		return s.SyncImpl(), err
	}
	if a, ok := impl.(interface {
		AsyncImpl() sdkapi.AsyncImpl
	}); ok {
		return a.AsyncImpl(), err
	}
	return nil, err
}

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

func TestRegistryDifferentNamespace(t *testing.T) {
	for _, nf := range allNew {
		ResetForTest()
		inst1, err1 := nf("this", "meter1")
		inst2, err2 := nf("this", "meter2")

		require.NoError(t, err1)
		require.NoError(t, err2)

		if inst1.Descriptor().InstrumentKind().Synchronous() {
			// They're equal because of a `nil` pointer at this point.
			// (Only for synchronous instruments, which lack callacks.)
			require.EqualValues(t, inst1, inst2)
		}

		SetMeterProvider(metrictest.NewMeterProvider())

		// They're different after the deferred setup.
		require.NotEqual(t, inst1, inst2)
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
			require.NotNil(t, other)
			require.True(t, errors.Is(err, registry.ErrMetricKindMismatch))
			require.Contains(t, err.Error(), "by this name with another kind or number type")
		}
	}
}
