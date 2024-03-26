// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSetting(t *testing.T) {
	const val int = 1
	s := newSetting(val)
	assert.True(t, s.Set, "returned unset value")
	assert.Equal(t, val, s.Value, "value not set")
}

func TestSettingResolve(t *testing.T) {
	t.Run("clearLessThanOne", func(t *testing.T) {
		var s setting[int]
		s.Value = -10
		s = s.Resolve(clearLessThanOne[int]())
		assert.False(t, s.Set)
		assert.Equal(t, 0, s.Value)

		s = newSetting[int](1)
		s = s.Resolve(clearLessThanOne[int]())
		assert.True(t, s.Set)
		assert.Equal(t, 1, s.Value)
	})

	t.Run("getenv", func(t *testing.T) {
		const key = "key"
		t.Setenv(key, "10")

		var s setting[int]
		s = s.Resolve(getenv[int](key))
		assert.True(t, s.Set)
		assert.Equal(t, 10, s.Value)

		t.Setenv(key, "20")
		s = s.Resolve(getenv[int](key))
		assert.Equal(t, 10, s.Value, "set setting overridden")
	})

	t.Run("fallback", func(t *testing.T) {
		var s setting[int]
		s = s.Resolve(fallback[int](10))
		assert.True(t, s.Set)
		assert.Equal(t, 10, s.Value)
	})

	t.Run("Precedence", func(t *testing.T) {
		const key = "key"

		var s setting[int]
		s = s.Resolve(
			clearLessThanOne[int](),
			getenv[int](key), // Unset.
			fallback[int](10),
		)
		assert.True(t, s.Set)
		assert.Equal(t, 10, s.Value)

		t.Setenv(key, "20")
		s = s.Resolve(
			clearLessThanOne[int](),
			getenv[int](key),  // Should not apply, already set.
			fallback[int](15), // Should not apply, already set.
		)
		assert.True(t, s.Set)
		assert.Equal(t, 10, s.Value)

		s = setting[int]{}
		s = s.Resolve(
			getenv[int](key),
			fallback[int](15), // Should not apply, already set.
		)
		assert.True(t, s.Set)
		assert.Equal(t, 20, s.Value)
	})
}
