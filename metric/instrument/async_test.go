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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/embedded"
)

func TestObservableConfiguration(t *testing.T) {
	t.Run("Int64", testObservableConfiguration[int64]())
	t.Run("Float64", testObservableConfiguration[float64]())
}

func testObservableConfiguration[N int64 | float64]() func(t *testing.T) {
	const (
		token  = 43
		desc   = "Instrument description."
		uBytes = "By"
	)

	run := func(got observableConfig[N]) func(*testing.T) {
		return func(t *testing.T) {
			assert.Equal(t, desc, got.Description(), "description")
			assert.Equal(t, uBytes, got.Unit(), "unit")

			// Functions are not comparable.
			cBacks := got.Callbacks()
			require.Len(t, cBacks, 1, "callbacks")
			o := &observer[N]{}
			err := cBacks[0](context.Background(), o)
			require.NoError(t, err)
			assert.Equal(t, N(token), o.got, "callback not set")
		}
	}

	cback := func(_ context.Context, obsrv ObserverT[N]) error {
		obsrv.Observe(token)
		return nil
	}

	return func(t *testing.T) {
		t.Run("ObservableCounter", run(
			NewObservableCounterConfig[N](
				WithDescription[N](desc),
				WithUnit[N](uBytes),
				WithCallback(cback),
			),
		))

		t.Run("ObservableUpDownCounter", run(
			NewObservableUpDownCounterConfig[N](
				WithDescription[N](desc),
				WithUnit[N](uBytes),
				WithCallback(cback),
			),
		))

		t.Run("ObservableGauge", run(
			NewObservableGaugeConfig[N](
				WithDescription[N](desc),
				WithUnit[N](uBytes),
				WithCallback(cback),
			),
		))
	}
}

type observableConfig[N int64 | float64] interface {
	Description() string
	Unit() string
	Callbacks() []Callback[N]
}

type observer[N int64 | float64] struct {
	embedded.ObserverT[N]
	got N
}

func (o *observer[N]) Observe(v N, _ ...attribute.KeyValue) {
	o.got = v
}
