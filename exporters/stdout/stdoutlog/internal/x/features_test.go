// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObservability(t *testing.T) {
	const key = "OTEL_GO_X_OBSERVABILITY"
	require.Contains(t, Observability.Keys(), key)

	t.Run("100", run(setenv(key, "100"), assertDisabled(Observability)))
	t.Run("true", run(setenv(key, "true"), assertEnabled(Observability, "true")))
	t.Run("True", run(setenv(key, "True"), assertEnabled(Observability, "True")))
	t.Run("false", run(setenv(key, "false"), assertDisabled(Observability)))
	t.Run("empty", run(assertDisabled(Observability)))
}
