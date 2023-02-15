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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/internaltest"
)

func TestNewNoopMeterProvider(t *testing.T) {
	mp := NewNoopMeterProvider()
	assert.Equal(t, mp, noopMeterProvider{})
	meter := mp.Meter("")
	assert.Equal(t, meter, noopMeter{})
}

func TestSyncFloat64(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	meter := NewNoopMeterProvider().Meter("test instrumentation")
	assert.NotPanics(t, func() {
		inst := meter.Float64Counter("test instrument")
		eh.RequireNoErrors(t)
		inst.Add(context.Background(), 1.0, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst := meter.Float64UpDownCounter("test instrument")
		eh.RequireNoErrors(t)
		inst.Add(context.Background(), -1.0, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst := meter.Float64Histogram("test instrument")
		eh.RequireNoErrors(t)
		inst.Record(context.Background(), 1.0, attribute.String("key", "value"))
	})
}

func TestSyncInt64(t *testing.T) {
	eh := internaltest.NewErrorHandler()
	otel.SetErrorHandler(eh)

	meter := NewNoopMeterProvider().Meter("test instrumentation")
	assert.NotPanics(t, func() {
		inst := meter.Int64Counter("test instrument")
		eh.RequireNoErrors(t)
		inst.Add(context.Background(), 1, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst := meter.Int64UpDownCounter("test instrument")
		eh.RequireNoErrors(t)
		inst.Add(context.Background(), -1, attribute.String("key", "value"))
	})

	assert.NotPanics(t, func() {
		inst := meter.Int64Histogram("test instrument")
		eh.RequireNoErrors(t)
		inst.Record(context.Background(), 1, attribute.String("key", "value"))
	})
}
