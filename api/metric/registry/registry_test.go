// Copyright 2020, OpenTelemetry Authors
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
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/registry"
	mockTest "go.opentelemetry.io/otel/internal/metric"
)

type (
	syncImpler interface {
		SyncImpl() metric.SyncImpl
	}

	newSyncFunc func(m metric.Meter, name string) (metric.SyncImpl, error)
)

var (
	allNewSync = map[string]newSyncFunc{
		"counter.int64": func(m metric.Meter, name string) (metric.SyncImpl, error) {
			return unwrapImpl(m.NewInt64Counter(name))
		},
		"counter.float64": func(m metric.Meter, name string) (metric.SyncImpl, error) {
			return unwrapImpl(m.NewFloat64Counter(name))
		},
		"measure.int64": func(m metric.Meter, name string) (metric.SyncImpl, error) {
			return unwrapImpl(m.NewInt64Measure(name))
		},
		"measure.float64": func(m metric.Meter, name string) (metric.SyncImpl, error) {
			return unwrapImpl(m.NewFloat64Measure(name))
		},
	}

	// HERE @@@
	// allNewAsync = map[string]newSyncFunc{
	// 	"observer.int64": func(m metric.Meter, name string) (metric.SyncImpl, error) {
	// 		return unwrapImpl(m.NewInt64Counter(name))
	// 	},
	// 	"observer.float64": func(m metric.Meter, name string) (metric.SyncImpl, error) {
	// 		return unwrapImpl(m.NewFloat64Counter(name))
	// 	},
	// }
)

func unwrapImpl(impl syncImpler, err error) (metric.SyncImpl, error) {
	if impl == nil {
		return nil, err
	}
	return impl.SyncImpl(), err
}

func TestRegistrySameSyncInstruments(t *testing.T) {
	for _, nf := range allNewSync {
		_, provider := mockTest.NewProvider()

		meter := provider.Meter("meter")
		inst1, err1 := nf(meter, "this")
		inst2, err2 := nf(meter, "this")

		require.Nil(t, err1)
		require.Nil(t, err2)
		require.Equal(t, inst1, inst2)
	}
}

func TestRegistryDiffSyncInstruments(t *testing.T) {
	for origName, origf := range allNewSync {
		_, provider := mockTest.NewProvider()
		meter := provider.Meter("meter")

		_, err := origf(meter, "this")
		require.Nil(t, err)

		for newName, nf := range allNewSync {
			if newName == origName {
				continue
			}

			other, err := nf(meter, "this")
			require.NotNil(t, err)
			require.NotNil(t, other)
			require.True(t, errors.Is(err, registry.ErrMetricTypeMismatch))
		}
	}
}
