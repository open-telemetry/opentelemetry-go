// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64Configuration(t *testing.T) {
	const (
		token  int64 = 43
		desc         = "Instrument description."
		uBytes       = "By"
	)

	run := func(got int64Config) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")
		}
	}

	t.Run("Int64Counter", run(
		NewInt64CounterConfig(WithDescription(desc), WithUnit(uBytes)),
	))

	t.Run("Int64UpDownCounter", run(
		NewInt64UpDownCounterConfig(WithDescription(desc), WithUnit(uBytes)),
	))

	t.Run("Int64Histogram", run(
		NewInt64HistogramConfig(WithDescription(desc), WithUnit(uBytes)),
	))

	t.Run("Int64Gauge", run(
		NewInt64GaugeConfig(WithDescription(desc), WithUnit(uBytes)),
	))
}

type int64Config interface {
	Description() string
	Unit() string
}

func TestInt64ExplicitBucketHistogramConfiguration(t *testing.T) {
	bounds := []float64{0.1, 0.5, 1.0}
	for _, tt := range []struct {
		name string
		opts []Int64HistogramOption
		want []float64
	}{
		{name: "omitted"},
		{
			name: "no arguments",
			opts: []Int64HistogramOption{WithExplicitBucketBoundaries()},
			want: []float64{},
		},
		{
			name: "nil slice",
			opts: []Int64HistogramOption{WithExplicitBucketBoundaries(nil...)},
			want: []float64{},
		},
		{
			name: "empty slice",
			opts: []Int64HistogramOption{WithExplicitBucketBoundaries([]float64{}...)},
			want: []float64{},
		},
		{
			name: "boundaries",
			opts: []Int64HistogramOption{WithExplicitBucketBoundaries(bounds...)},
			want: bounds,
		},
		{
			name: "last option wins",
			opts: []Int64HistogramOption{WithExplicitBucketBoundaries(bounds...), WithExplicitBucketBoundaries()},
			want: []float64{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInt64HistogramConfig(tt.opts...)
			assert.Equal(t, tt.want, got.ExplicitBucketBoundaries(), "boundaries")
		})
	}
}
