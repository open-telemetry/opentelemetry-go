// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/embedded"
)

func TestInt64ObservableConfiguration(t *testing.T) {
	const (
		token  int64 = 43
		desc         = "Instrument description."
		uBytes       = "By"
	)

	run := func(got int64ObservableConfig) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")

			// Functions are not comparable.
			cBacks := got.Callbacks()
			require.Len(t, cBacks, 1, "callbacks")
			o := &int64Observer{}
			err := cBacks[0](t.Context(), o)
			require.NoError(t, err)
			assert.Equal(t, token, o.got, "callback not set")
		}
	}

	cback := func(_ context.Context, obsrv Int64Observer) error {
		obsrv.Observe(token)
		return nil
	}

	t.Run("Int64ObservableCounter", run(
		NewInt64ObservableCounterConfig(
			WithDescription(desc),
			WithUnit(uBytes),
			WithInt64Callback(cback),
		),
	))

	t.Run("Int64ObservableUpDownCounter", run(
		NewInt64ObservableUpDownCounterConfig(
			WithDescription(desc),
			WithUnit(uBytes),
			WithInt64Callback(cback),
		),
	))

	t.Run("Int64ObservableGauge", run(
		NewInt64ObservableGaugeConfig(
			WithDescription(desc),
			WithUnit(uBytes),
			WithInt64Callback(cback),
		),
	))
}

func TestWithInt64CallbackIgnoresNil(t *testing.T) {
	nilOption := WithInt64Callback(nil)
	assert.Empty(t, NewInt64ObservableCounterConfig(nilOption).Callbacks())
	assert.Empty(t, NewInt64ObservableUpDownCounterConfig(nilOption).Callbacks())
	assert.Empty(t, NewInt64ObservableGaugeConfig(nilOption).Callbacks())

	callback := WithInt64Callback(func(_ context.Context, observer Int64Observer) error {
		observer.Observe(43)
		return nil
	})
	checkCallbacks := func(callbacks []Int64Callback) {
		t.Helper()
		require.Len(t, callbacks, 1)
		observer := &int64Observer{}
		require.NoError(t, callbacks[0](t.Context(), observer))
		assert.Equal(t, int64(43), observer.got)
	}
	checkCallbacks(NewInt64ObservableCounterConfig(callback, nilOption).Callbacks())
	checkCallbacks(NewInt64ObservableUpDownCounterConfig(callback, nilOption).Callbacks())
	checkCallbacks(NewInt64ObservableGaugeConfig(callback, nilOption).Callbacks())
}

type int64ObservableConfig interface {
	Description() string
	Unit() string
	Callbacks() []Int64Callback
}

type int64Observer struct {
	embedded.Int64Observer
	Observable
	got int64
}

func (o *int64Observer) Observe(v int64, _ ...ObserveOption) {
	o.got = v
}
