// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/noop"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func TestWithViewCombinesWithMeterOptions(t *testing.T) {
	view := sdkmetricx.WithView(sdkmetricx.NewView(
		sdkmetric.Instrument{Name: "requests"},
		sdkmetricx.Stream{},
	))

	t.Run("Finish", func(t *testing.T) {
		mp := sdkmetric.NewMeterProvider(view, sdkmetricx.WithFinish())
		t.Cleanup(func() { assert.NoError(t, mp.Shutdown(t.Context())) })

		counter, err := mp.Meter("test").Int64Counter("requests")
		require.NoError(t, err)
		_, ok := counter.(metricx.Finisher)
		assert.True(t, ok)
	})

	t.Run("MeterConfigurator", func(t *testing.T) {
		handle := sdkmetricx.NewMeterConfiguratorHandle()
		mp := sdkmetric.NewMeterProvider(view, sdkmetricx.WithMeterConfigurator(handle))
		t.Cleanup(func() { assert.NoError(t, mp.Shutdown(t.Context())) })

		assert.NotEqual(t, noop.Meter{}, mp.Meter("test"))
	})
}

func TestMeterOptionsAreRejectedTogether(t *testing.T) {
	mp := sdkmetric.NewMeterProvider(
		sdkmetricx.WithFinish(),
		sdkmetricx.WithMeterConfigurator(sdkmetricx.NewMeterConfiguratorHandle()),
	)
	t.Cleanup(func() { assert.NoError(t, mp.Shutdown(t.Context())) })

	assert.Equal(t, noop.Meter{}, mp.Meter("test"))
}
