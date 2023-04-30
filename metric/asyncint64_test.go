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
			err := cBacks[0](context.Background(), o)
			require.NoError(t, err)
			assert.Equal(t, token, o.got, "callback not set")
		}
	}

	cback := func(ctx context.Context, obsrv Int64Observer) error {
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
