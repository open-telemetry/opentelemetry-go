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

package asyncfloat64 // import "go.opentelemetry.io/otel/metric/instrument/asyncfloat64"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/unit"
)

func TestOptions(t *testing.T) {
	const (
		token  = 43.0
		desc   = "Instrument description."
		uBytes = unit.Bytes
	)

	want := []Measurement{{Value: token}}
	got := NewConfig(
		WithDescription(desc),
		WithUnit(uBytes),
		WithCallback(func(context.Context) ([]Measurement, error) {
			return want, nil
		}),
	)
	assert.Equal(t, desc, got.Description(), "description")
	assert.Equal(t, uBytes, got.Unit(), "unit")

	// Functions are not comparable.
	cBacks := got.Callbacks()
	require.Len(t, cBacks, 1, "callbacks")
	val, err := cBacks[0](context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, val, "callback not set")
}
