// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelfObservability(t *testing.T) {
	const key = "OTEL_GO_X_SELF_OBSERVABILITY"
	require.Equal(t, key, SelfObservability.Key())

	t.Run("100", run(setenv(key, "100"), assertDisabled(SelfObservability)))
	t.Run("true", run(setenv(key, "true"), assertEnabled(SelfObservability, "true")))
	t.Run("True", run(setenv(key, "True"), assertEnabled(SelfObservability, "True")))
	t.Run("false", run(setenv(key, "false"), assertDisabled(SelfObservability)))
	t.Run("empty", run(assertDisabled(SelfObservability)))
}

func run(steps ...func(*testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		for _, step := range steps {
			step(t)
		}
	}
}

func setenv(k, v string) func(t *testing.T) { //nolint:unparam // This is a reusable test utility function.
	return func(t *testing.T) { t.Setenv(k, v) }
}

func assertEnabled[T any](f Feature[T], want T) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		assert.True(t, f.Enabled(), "not enabled")

		v, ok := f.Lookup()
		assert.True(t, ok, "Lookup state")
		assert.Equal(t, want, v, "Lookup value")
	}
}

func assertDisabled[T any](f Feature[T]) func(*testing.T) {
	var zero T
	return func(t *testing.T) {
		t.Helper()

		assert.False(t, f.Enabled(), "enabled")

		v, ok := f.Lookup()
		assert.False(t, ok, "Lookup state")
		assert.Equal(t, zero, v, "Lookup value")
	}
}
