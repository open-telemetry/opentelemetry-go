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

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/registry"
	"go.opentelemetry.io/otel/oteltest"
)

type (
	newFunc func(m metric.Meter, name string) (metric.InstrumentImpl, error)
)

var (
	allNew = map[string]newFunc{
		"counter.int64": func(m metric.Meter, name string) (metric.InstrumentImpl, error) {
			return unwrap(m.NewInt64Counter(name))
		},
		"counter.float64": func(m metric.Meter, name string) (metric.InstrumentImpl, error) {
			return unwrap(m.NewFloat64Counter(name))
		},
		"valuerecorder.int64": func(m metric.Meter, name string) (metric.InstrumentImpl, error) {
			return unwrap(m.NewInt64ValueRecorder(name))
		},
		"valuerecorder.float64": func(m metric.Meter, name string) (metric.InstrumentImpl, error) {
			return unwrap(m.NewFloat64ValueRecorder(name))
		},
		"valueobserver.int64": func(m metric.Meter, name string) (metric.InstrumentImpl, error) {
			return unwrap(m.NewInt64ValueObserver(name, func(context.Context, metric.Int64ObserverResult) {}))
		},
		"valueobserver.float64": func(m metric.Meter, name string) (metric.InstrumentImpl, error) {
			return unwrap(m.NewFloat64ValueObserver(name, func(context.Context, metric.Float64ObserverResult) {}))
		},
	}
)

func unwrap(impl interface{}, err error) (metric.InstrumentImpl, error) {
	if impl == nil {
		return nil, err
	}
	if s, ok := impl.(interface {
		SyncImpl() metric.SyncImpl
	}); ok {
		return s.SyncImpl(), err
	}
	if a, ok := impl.(interface {
		AsyncImpl() metric.AsyncImpl
	}); ok {
		return a.AsyncImpl(), err
	}
	return nil, err
}

func TestRegistrySameInstruments(t *testing.T) {
	for _, nf := range allNew {
		_, provider := oteltest.NewMeterProvider()

		meter := provider.Meter("meter")
		inst1, err1 := nf(meter, "this")
		inst2, err2 := nf(meter, "this")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Equal(t, inst1, inst2)
	}
}

func TestRegistryDifferentNamespace(t *testing.T) {
	for _, nf := range allNew {
		_, provider := oteltest.NewMeterProvider()

		meter1 := provider.Meter("meter1")
		meter2 := provider.Meter("meter2")
		inst1, err1 := nf(meter1, "this")
		inst2, err2 := nf(meter2, "this")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, inst1, inst2)
	}
}

func TestRegistryDiffInstruments(t *testing.T) {
	for origName, origf := range allNew {
		_, provider := oteltest.NewMeterProvider()
		meter := provider.Meter("meter")

		_, err := origf(meter, "this")
		require.NoError(t, err)

		for newName, nf := range allNew {
			if newName == origName {
				continue
			}

			other, err := nf(meter, "this")
			require.Error(t, err)
			require.NotNil(t, other)
			require.True(t, errors.Is(err, registry.ErrMetricKindMismatch))
		}
	}
}

func TestMeterProvider(t *testing.T) {
	impl, _ := oteltest.NewMeter()
	p := registry.NewMeterProvider(impl)
	m1 := p.Meter("m1")
	m1p := p.Meter("m1")
	m2 := p.Meter("m2")

	require.Equal(t, m1, m1p)
	require.NotEqual(t, m1, m2)
}
