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

package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/registry"
)

type (
	newFunc func(name, libraryName string) (metric.InstrumentImpl, error)
)

var (
	allNew = map[string]newFunc{
		"counter.int64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64Counter(name))
		},
		"counter.float64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64Counter(name))
		},
		"valuerecorder.int64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64ValueRecorder(name))
		},
		"valuerecorder.float64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64ValueRecorder(name))
		},
		"valueobserver.int64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewInt64ValueObserver(name, func(context.Context, metric.Int64ObserverResult) {}))
		},
		"valueobserver.float64": func(name, libraryName string) (metric.InstrumentImpl, error) {
			return unwrap(MeterProvider().Meter(libraryName).NewFloat64ValueObserver(name, func(context.Context, metric.Float64ObserverResult) {}))
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
		ResetForTest()
		inst1, err1 := nf("this", "meter")
		inst2, err2 := nf("this", "meter")

		require.NoError(t, err1)
		require.NoError(t, err2)
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
			require.Contains(t, err.Error(), "super")
		}
	}
}
