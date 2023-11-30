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

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExemplars(t *testing.T) {
	const key = "OTEL_GO_X_EXEMPLAR"
	require.Equal(t, key, Exemplars.Key())

	t.Run("true", func(t *testing.T) {
		t.Setenv(key, "true")
		assert.True(t, Exemplars.Enabled())

		v, ok := Exemplars.Lookup()
		assert.True(t, ok)
		assert.Equal(t, "true", v)
	})

	t.Run("false", func(t *testing.T) {
		t.Setenv(key, "false")
		assert.False(t, Exemplars.Enabled())

		v, ok := Exemplars.Lookup()
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})

	t.Run("empty", func(t *testing.T) {
		assert.False(t, Exemplars.Enabled())

		v, ok := Exemplars.Lookup()
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})
}

func TestCardinalityLimit(t *testing.T) {
	const key = "OTEL_GO_X_CARDINALITY_LIMIT"
	require.Equal(t, key, CardinalityLimit.Key())

	t.Run("100", func(t *testing.T) {
		t.Setenv(key, "100")
		assert.True(t, CardinalityLimit.Enabled())

		v, ok := CardinalityLimit.Lookup()
		assert.True(t, ok)
		assert.Equal(t, 100, v)
	})

	t.Run("-1", func(t *testing.T) {
		t.Setenv(key, "-1")
		assert.True(t, CardinalityLimit.Enabled())

		v, ok := CardinalityLimit.Lookup()
		assert.True(t, ok)
		assert.Equal(t, -1, v)
	})

	t.Run("false", func(t *testing.T) {
		t.Setenv(key, "false")
		assert.False(t, CardinalityLimit.Enabled())

		v, ok := CardinalityLimit.Lookup()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv(key, "false")
		assert.False(t, CardinalityLimit.Enabled())

		v, ok := CardinalityLimit.Lookup()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})
}
