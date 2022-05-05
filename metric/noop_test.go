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

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewNoopMeterProvider(t *testing.T) {
	mp := NewNoopMeterProvider()
	assert.Equal(t, mp, noopMeterProvider{})
	meter := mp.Meter("")
	assert.Equal(t, meter, noopMeter{})
}

func TestSyncFloat64(t *testing.T) {
	meter := NewNoopMeterProvider().Meter("test instrumentation")
	assert.NotPanics(t, func() {
		inst, err := meter.SyncFloat64().Counter("test instrument")
		require.NoError(t, err)
		inst.Add(context.Background(), 1.0, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.SyncFloat64().UpDownCounter("test instrument")
		require.NoError(t, err)
		inst.Add(context.Background(), -1.0, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.SyncFloat64().Histogram("test instrument")
		require.NoError(t, err)
		inst.Record(context.Background(), 1.0, attribute.String("key", "value"))
	})
}

func TestSyncInt64(t *testing.T) {
	meter := NewNoopMeterProvider().Meter("test instrumentation")
	assert.NotPanics(t, func() {
		inst, err := meter.SyncInt64().Counter("test instrument")
		require.NoError(t, err)
		inst.Add(context.Background(), 1, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.SyncInt64().UpDownCounter("test instrument")
		require.NoError(t, err)
		inst.Add(context.Background(), -1, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.SyncInt64().Histogram("test instrument")
		require.NoError(t, err)
		inst.Record(context.Background(), 1, attribute.String("key", "value"))
	})
}

func TestAsyncFloat64(t *testing.T) {
	meter := NewNoopMeterProvider().Meter("test instrumentation")
	assert.NotPanics(t, func() {
		inst, err := meter.AsyncFloat64().Counter("test instrument")
		require.NoError(t, err)
		inst.Observe(context.Background(), 1.0, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.AsyncFloat64().UpDownCounter("test instrument")
		require.NoError(t, err)
		inst.Observe(context.Background(), -1.0, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.AsyncFloat64().Gauge("test instrument")
		require.NoError(t, err)
		inst.Observe(context.Background(), 1.0, attribute.String("key", "value"))
	})
}

func TestAsyncInt64(t *testing.T) {
	meter := NewNoopMeterProvider().Meter("test instrumentation")
	assert.NotPanics(t, func() {
		inst, err := meter.AsyncInt64().Counter("test instrument")
		require.NoError(t, err)
		inst.Observe(context.Background(), 1, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.AsyncInt64().UpDownCounter("test instrument")
		require.NoError(t, err)
		inst.Observe(context.Background(), -1, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst, err := meter.AsyncInt64().Gauge("test instrument")
		require.NoError(t, err)
		inst.Observe(context.Background(), 1, attribute.String("key", "value"))
	})
}
