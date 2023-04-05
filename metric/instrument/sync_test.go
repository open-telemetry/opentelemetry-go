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

package instrument // import "go.opentelemetry.io/otel/metric/instrument"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64Configuration(t *testing.T) {
	t.Run("Int64", testConfiguration[int64]())
	t.Run("Float64", testConfiguration[float64]())
}

func testConfiguration[N int64 | float64]() func(t *testing.T) {
	const (
		token  int64 = 43
		desc         = "Instrument description."
		uBytes       = "By"
	)

	run := func(got config) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")
		}
	}

	return func(t *testing.T) {
		t.Run("Counter", run(
			NewCounterConfig[N](WithDescription[N](desc), WithUnit[N](uBytes)),
		))

		t.Run("UpDownCounter", run(
			NewUpDownCounterConfig[N](WithDescription[N](desc), WithUnit[N](uBytes)),
		))

		t.Run("Histogram", run(
			NewHistogramConfig[N](WithDescription[N](desc), WithUnit[N](uBytes)),
		))
	}
}

type config interface {
	Description() string
	Unit() string
}
