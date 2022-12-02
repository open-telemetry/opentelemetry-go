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
	"go.opentelemetry.io/otel/metric/unit"
)

func TestSynchronousOptions(t *testing.T) {
	const (
		desc   = "Instrument description."
		uBytes = unit.Bytes
	)

	got := NewSynchronousConfig(WithDescription(desc), WithUnit(uBytes))
	assert.Equal(t, desc, got.Description(), "description")
	assert.Equal(t, uBytes, got.Unit(), uBytes)
}

func TestAsynchronousOptions(t *testing.T) {
	const (
		token  = "token"
		desc   = "Instrument description."
		uBytes = unit.Bytes
	)

	var received string
	cBack := Callback(func(context.Context, Asynchronous) error {
		received = token
		return nil
	})

	got := NewAsynchronousConfig(
		WithDescription(desc),
		WithUnit(uBytes),
		WithCallback(cBack),
	)
	assert.Equal(t, desc, got.Description(), "description")
	assert.Equal(t, uBytes, got.Unit(), uBytes)

	// Functions are not comparable.
	cBacks := got.Callbacks()
	require.Len(t, cBacks, 1, "callbacks")
	_ = cBacks[0](context.Background(), nil)
	assert.Equal(t, token, received, "callback not set")
}
