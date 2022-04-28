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
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/registry"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	newFunc func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error)
)

var (
	allNew = map[string]newFunc{
		"counter.int64": func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error) {
			return unwrap(m.SyncInt64().Counter(name))
		},
		"counter.float64": func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error) {
			return unwrap(m.SyncFloat64().Counter(name))
		},
		"histogram.int64": func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error) {
			return unwrap(m.SyncInt64().Histogram(name))
		},
		"histogram.float64": func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error) {
			return unwrap(m.SyncFloat64().Histogram(name))
		},
		"gaugeobserver.int64": func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error) {
			return unwrap(m.AsyncInt64().Gauge(name))
		},
		"gaugeobserver.float64": func(m metric.Meter, name string) (sdkapi.InstrumentImpl, error) {
			return unwrap(m.AsyncFloat64().Gauge(name))
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

// TODO Replace with controller.
func testMeterWithRegistry(name string) metric.Meter {
	return sdkapi.WrapMeterImpl(
		registry.NewUniqueInstrumentMeterImpl(
			metricsdk.NewAccumulator(nil),
		),
	)
}

func TestRegistrySameInstruments(t *testing.T) {
	for _, nf := range allNew {
		meter := testMeterWithRegistry("meter")
		inst1, err1 := nf(meter, "this")
		inst2, err2 := nf(meter, "this")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Equal(t, inst1, inst2)
	}
}

func TestRegistryDiffInstruments(t *testing.T) {
	for origName, origf := range allNew {
		meter := testMeterWithRegistry("meter")

		_, err := origf(meter, "this")
		require.NoError(t, err)

		for newName, nf := range allNew {
			if newName == origName {
				continue
			}

			other, err := nf(meter, "this")
			require.Error(t, err)
			require.Nil(t, other)
			require.True(t, errors.Is(err, registry.ErrMetricKindMismatch))
		}
	}
}
