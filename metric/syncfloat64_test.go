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
}

type float64Config interface {
	Description() string
	Unit() string
}

func TestFloat64ExplicitBucketHistogramConfiguration(t *testing.T) {
	bounds := []float64{0.1, 0.5, 1.0}
	got := NewFloat64HistogramConfig(WithExplicitBucketBoundaries(bounds...))
	assert.Equal(t, bounds, got.ExplicitBucketBoundaries(), "boundaries")
}
