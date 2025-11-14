// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResource(t *testing.T) {
	const key = "OTEL_GO_X_RESOURCE"
	require.Contains(t, Resource.Keys(), key)

	t.Run("100", run(setenv(key, "100"), assertDisabled(Resource)))
	t.Run("true", run(setenv(key, "true"), assertEnabled(Resource, "true")))
	t.Run("True", run(setenv(key, "True"), assertEnabled(Resource, "True")))
	t.Run("false", run(setenv(key, "false"), assertDisabled(Resource)))
	t.Run("empty", run(assertDisabled(Resource)))
}
