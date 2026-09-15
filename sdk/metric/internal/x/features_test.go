// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParallelCallbacks(t *testing.T) {
	const key = "OTEL_GO_X_PARALLEL_CALLBACKS"
	require.Contains(t, ParallelCallbacks.Keys(), key)

	tests := []struct {
		name    string
		value   string
		enabled bool
	}{
		{name: "empty", value: "", enabled: false},
		{name: "false", value: "false", enabled: false},
		{name: "other", value: "1", enabled: false},
		{name: "true", value: "true", enabled: true},
		{name: "mixed case", value: "TrUe", enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(key, tt.value)
			assert.Equal(t, tt.enabled, ParallelCallbacks.Enabled())
			got, ok := ParallelCallbacks.Lookup()
			assert.Equal(t, tt.enabled, ok)
			assert.Equal(t, tt.enabled, got)
		})
	}
}
