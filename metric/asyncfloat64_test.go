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

func TestFloat64ObservableConfiguration(t *testing.T) {
	const (
		token  float64 = 43
		desc           = "Instrument description."
		uBytes         = "By"
	)
	cback := func(_ context.Context, obsrv Float64Observer) error {
		obsrv.Observe(token)
		return nil
	}
	run := func(got float64ObservableConfig, wantCallbacks int) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")
			callbacks := got.Callbacks()
			require.Len(t, callbacks, wantCallbacks, "callbacks")
			for _, callback := range callbacks {
				o := &float64Observer{}
				require.NoError(t, callback(t.Context(), o))
				assert.Equal(t, token, o.got, "callback not set")
			}
		}
	}
	for _, tc := range []struct {
		name          string
		callbacks     []Float64Callback
		wantCallbacks int
	}{
		{name: "NoCallbacks"},
		{name: "NilCallback", callbacks: []Float64Callback{nil}},
		{name: "Callback", callbacks: []Float64Callback{cback}, wantCallbacks: 1},
		{name: "NilBeforeAndAfterCallback", callbacks: []Float64Callback{nil, cback, nil}, wantCallbacks: 1},
		{name: "CallbacksWithNil", callbacks: []Float64Callback{cback, nil, cback}, wantCallbacks: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counterOpts := []Float64ObservableCounterOption{WithDescription(desc), WithUnit(uBytes)}
			upDownOpts := []Float64ObservableUpDownCounterOption{WithDescription(desc), WithUnit(uBytes)}
			gaugeOpts := []Float64ObservableGaugeOption{WithDescription(desc), WithUnit(uBytes)}
			for _, callback := range tc.callbacks {
				option := WithFloat64Callback(callback)
				counterOpts = append(counterOpts, option)
				upDownOpts = append(upDownOpts, option)
				gaugeOpts = append(gaugeOpts, option)
			}
			t.Run("Float64ObservableCounter", run(NewFloat64ObservableCounterConfig(counterOpts...), tc.wantCallbacks))
			t.Run(
				"Float64ObservableUpDownCounter",
				run(NewFloat64ObservableUpDownCounterConfig(upDownOpts...), tc.wantCallbacks),
			)
			t.Run("Float64ObservableGauge", run(NewFloat64ObservableGaugeConfig(gaugeOpts...), tc.wantCallbacks))
		})
	}
}

type float64ObservableConfig interface {
	Description() string
	Unit() string
	Callbacks() []Float64Callback
}

type float64Observer struct {
	embedded.Float64Observer
	Observable
	got float64
}

func (o *float64Observer) Observe(v float64, _ ...ObserveOption) {
	o.got = v
}
