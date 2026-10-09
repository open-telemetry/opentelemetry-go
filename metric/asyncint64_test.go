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
	cback := func(_ context.Context, obsrv Int64Observer) error {
		obsrv.Observe(token)
		return nil
	}
	run := func(got int64ObservableConfig, wantCallbacks int) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")
			callbacks := got.Callbacks()
			require.Len(t, callbacks, wantCallbacks, "callbacks")
			for _, callback := range callbacks {
				o := &int64Observer{}
				require.NoError(t, callback(t.Context(), o))
				assert.Equal(t, token, o.got, "callback not set")
			}
		}
	}
	for _, tc := range []struct {
		name          string
		callbacks     []Int64Callback
		wantCallbacks int
	}{
		{name: "NoCallbacks"},
		{name: "NilCallback", callbacks: []Int64Callback{nil}},
		{name: "Callback", callbacks: []Int64Callback{cback}, wantCallbacks: 1},
		{name: "NilBeforeAndAfterCallback", callbacks: []Int64Callback{nil, cback, nil}, wantCallbacks: 1},
		{name: "CallbacksWithNil", callbacks: []Int64Callback{cback, nil, cback}, wantCallbacks: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counterOpts := []Int64ObservableCounterOption{WithDescription(desc), WithUnit(uBytes)}
			upDownOpts := []Int64ObservableUpDownCounterOption{WithDescription(desc), WithUnit(uBytes)}
			gaugeOpts := []Int64ObservableGaugeOption{WithDescription(desc), WithUnit(uBytes)}
			for _, callback := range tc.callbacks {
				option := WithInt64Callback(callback)
				counterOpts = append(counterOpts, option)
				upDownOpts = append(upDownOpts, option)
				gaugeOpts = append(gaugeOpts, option)
			}
			t.Run("Int64ObservableCounter", run(NewInt64ObservableCounterConfig(counterOpts...), tc.wantCallbacks))
			t.Run(
				"Int64ObservableUpDownCounter",
				run(NewInt64ObservableUpDownCounterConfig(upDownOpts...), tc.wantCallbacks),
			)
			t.Run("Int64ObservableGauge", run(NewInt64ObservableGaugeConfig(gaugeOpts...), tc.wantCallbacks))
		})
	}
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
