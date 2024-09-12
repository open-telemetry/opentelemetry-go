// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResource(t *testing.T) {
	const key = "OTEL_GO_X_RESOURCE"
	require.Equal(t, key, Resource.Key())

	t.Run("true", run(setenv(key, "true"), assertEnabled(Resource, "true")))
	t.Run("True", run(setenv(key, "True"), assertEnabled(Resource, "True")))
	t.Run("TRUE", run(setenv(key, "TRUE"), assertEnabled(Resource, "TRUE")))
	t.Run("false", run(setenv(key, "false"), assertDisabled(Resource)))
	t.Run("1", run(setenv(key, "1"), assertDisabled(Resource)))
	t.Run("empty", run(assertDisabled(Resource)))
}

func run(steps ...func(*testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		for _, step := range steps {
			step(t)
		}
	}
}

func setenv(k, v string) func(t *testing.T) { //nolint:unparam
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
