// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64Configuration(t *testing.T) {
	const (
		token  float64 = 43
		desc           = "Instrument description."
		uBytes         = "By"
		optIn          = true
	)

	run := func(got float64Config) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")
			assert.Equal(t, optIn, got.OptIn(), "optIn")
		}
	}

	t.Run("Float64Counter", run(
		NewFloat64CounterConfig(WithDescription(desc), WithUnit(uBytes), WithOptIn()),
	))

	t.Run("Float64UpDownCounter", run(
		NewFloat64UpDownCounterConfig(WithDescription(desc), WithUnit(uBytes), WithOptIn()),
	))

	t.Run("Float64Histogram", run(
		NewFloat64HistogramConfig(WithDescription(desc), WithUnit(uBytes), WithOptIn()),
	))

	t.Run("Float64Gauge", run(
		NewFloat64GaugeConfig(WithDescription(desc), WithUnit(uBytes), WithOptIn()),
	))
}

type float64Config interface {
	Description() string
	Unit() string
	OptIn() bool
}

func TestFloat64ExplicitBucketHistogramConfiguration(t *testing.T) {
	bounds := []float64{0.1, 0.5, 1.0}
	got := NewFloat64HistogramConfig(WithExplicitBucketBoundaries(bounds...))
	assert.Equal(t, bounds, got.ExplicitBucketBoundaries(), "boundaries")
}
