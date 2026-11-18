// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64Configuration(t *testing.T) {
	const (
		token  float64 = 43
		desc           = "Instrument description."
		uBytes         = "By"
	)

	run := func(got float64Config) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")
		}
	}

	t.Run("Float64Counter", run(
		NewFloat64CounterConfig(WithDescription(desc), WithUnit(uBytes)),
	))

	t.Run("Float64UpDownCounter", run(
		NewFloat64UpDownCounterConfig(WithDescription(desc), WithUnit(uBytes)),
	))

	t.Run("Float64Histogram", run(
		NewFloat64HistogramConfig(WithDescription(desc), WithUnit(uBytes)),
	))

	t.Run("Float64Gauge", run(
		NewFloat64GaugeConfig(WithDescription(desc), WithUnit(uBytes)),
	))
}

type float64Config interface {
	Description() string
	Unit() string
}

func TestFloat64ExplicitBucketHistogramConfiguration(t *testing.T) {
	bounds := []float64{0.1, 0.5, 1.0}
	for _, tt := range []struct {
		name string
		opts []Float64HistogramOption
		want []float64
	}{
		{name: "omitted"},
		{
			name: "no arguments",
			opts: []Float64HistogramOption{WithExplicitBucketBoundaries()},
			want: []float64{},
		},
		{
			name: "nil slice",
			opts: []Float64HistogramOption{WithExplicitBucketBoundaries(nil...)},
			want: []float64{},
		},
		{
			name: "empty slice",
			opts: []Float64HistogramOption{WithExplicitBucketBoundaries([]float64{}...)},
			want: []float64{},
		},
		{
			name: "boundaries",
			opts: []Float64HistogramOption{WithExplicitBucketBoundaries(bounds...)},
			want: bounds,
		},
		{
			name: "last option wins",
			opts: []Float64HistogramOption{WithExplicitBucketBoundaries(bounds...), WithExplicitBucketBoundaries()},
			want: []float64{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFloat64HistogramConfig(tt.opts...)
			assert.Equal(t, tt.want, got.ExplicitBucketBoundaries(), "boundaries")
		})
	}
}
